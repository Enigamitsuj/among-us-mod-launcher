package installer

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Enigamitsuj/among-us-mod-launcher/internal/fsutil"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/installrecord"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/mods"
)

const (
	stagingSuffix = ".__aumod-new"
	backupSuffix  = ".__aumod-old"
)

// ProgressFunc receives install progress updates.
type ProgressFunc func(mods.InstallProgress)

// Service performs download + extract + copy installs.
type Service struct {
	HTTP *http.Client
}

func New() *Service {
	return &Service{
		HTTP: &http.Client{Timeout: 10 * time.Minute},
	}
}

// DestinationExists reports whether the mod folder already exists.
func (s *Service) DestinationExists(installRoot, folderName string) bool {
	return fsutil.DirExists(fsutil.JoinInstallPath(installRoot, folderName))
}

// Install downloads a release ZIP and creates a sibling modded Among Us copy.
// The original game directory is never modified. Reinstalls build into a
// staging folder first, then swap into place so a failed update keeps the old copy.
func (s *Service) Install(opts mods.InstallOptions, mod mods.Mod, downloadURL string, onProgress ProgressFunc) (string, error) {
	emit := func(stage, message string, percent float64, done bool, errMsg, installDir string) {
		if onProgress != nil {
			onProgress(mods.InstallProgress{
				Stage:      stage,
				Message:    message,
				Percent:    percent,
				Done:       done,
				Error:      errMsg,
				InstallDir: installDir,
			})
		}
	}

	if opts.AmongUsPath == "" {
		return "", fmt.Errorf("Among Us not found")
	}
	if !fileExists(filepath.Join(opts.AmongUsPath, "Among Us.exe")) {
		return "", fmt.Errorf("Among Us not found")
	}
	if downloadURL == "" {
		return "", fmt.Errorf("no download available for this version")
	}

	destRoot := filepath.Dir(filepath.Clean(opts.AmongUsPath))
	dest := fsutil.JoinInstallPath(destRoot, mod.FolderName)
	staging := dest + stagingSuffix
	backup := dest + backupSuffix

	if !fsutil.DirIsWritable(destRoot) {
		return "", ErrProtectedInstallFolder
	}

	if fsutil.DirExists(dest) && !opts.ForceReinstall {
		return "", ErrAlreadyExists
	}

	// Clear leftover staging/backup from a previous interrupted update.
	_ = os.RemoveAll(staging)
	_ = os.RemoveAll(backup)

	tmpDir, err := os.MkdirTemp("", "aumod-install-*")
	if err != nil {
		return "", fmt.Errorf("could not create temporary folder")
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "mod.zip")
	emit("downloading", "Downloading...", 5, false, "", dest)
	if err := s.downloadFile(downloadURL, zipPath, func(p float64) {
		emit("downloading", "Downloading...", 5+p*40, false, "", dest)
	}); err != nil {
		emit("error", "Download failed", 0, true, "GitHub unavailable.", dest)
		return "", fmt.Errorf("GitHub unavailable")
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	emit("extracting", "Extracting...", 50, false, "", dest)
	if err := unzip(zipPath, extractDir); err != nil {
		return "", fmt.Errorf("failed to extract the download")
	}

	emit("installing", "Preparing new installation...", 70, false, "", dest)
	if err := copyDir(opts.AmongUsPath, staging); err != nil {
		_ = os.RemoveAll(staging)
		if isPermissionError(err) {
			return "", ErrProtectedInstallFolder
		}
		return "", fmt.Errorf("failed to copy Among Us files")
	}

	emit("installing", "Applying mod files...", 88, false, "", dest)
	modRoot := findModRoot(extractDir)
	if err := mergeDir(modRoot, staging); err != nil {
		_ = os.RemoveAll(staging)
		return "", fmt.Errorf("failed to apply mod files")
	}

	if err := installrecord.Save(staging, installrecord.Record{
		ModID:    opts.ModID,
		Path:     dest,
		Version:  opts.VersionTag,
		Platform: opts.Platform,
	}); err != nil {
		_ = os.RemoveAll(staging)
		return "", fmt.Errorf("failed to finalize the mod installation")
	}

	if !fileExists(filepath.Join(staging, "Among Us.exe")) {
		_ = os.RemoveAll(staging)
		return "", fmt.Errorf("failed to finalize the mod installation")
	}

	emit("installing", "Updating installation...", 95, false, "", dest)
	if err := promoteStaging(staging, dest, backup); err != nil {
		return "", err
	}

	emit("finished", "Finished.", 100, true, "", dest)
	return dest, nil
}

// promoteStaging moves a fully prepared staging directory into dest.
// If dest already exists, it is renamed aside first and restored if the swap fails.
func promoteStaging(staging, dest, backup string) error {
	if fsutil.DirExists(dest) {
		if err := os.Rename(dest, backup); err != nil {
			_ = os.RemoveAll(staging)
			return fmt.Errorf("could not update existing installation")
		}
		if err := os.Rename(staging, dest); err != nil {
			_ = os.Rename(backup, dest)
			_ = os.RemoveAll(staging)
			return fmt.Errorf("could not update existing installation")
		}
		_ = os.RemoveAll(backup)
		return nil
	}

	if err := os.Rename(staging, dest); err != nil {
		_ = os.RemoveAll(staging)
		return fmt.Errorf("could not finalize the mod installation")
	}
	return nil
}

func (s *Service) downloadFile(url, dest string, onPercent func(float64)) error {
	resp, err := s.HTTP.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	var written int64
	total := resp.ContentLength
	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
			written += int64(n)
			if total > 0 && onPercent != nil {
				onPercent(float64(written) / float64(total) * 100)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	if onPercent != nil {
		onPercent(100)
	}
	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}

	for _, f := range r.File {
		if err := extractZipFile(f, dest); err != nil {
			return err
		}
	}
	return nil
}

func extractZipFile(f *zip.File, dest string) error {
	target, err := safeJoin(dest, f.Name)
	if err != nil {
		return err
	}
	if f.FileInfo().IsDir() {
		return os.MkdirAll(target, 0o755)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, rc)
	return err
}

func safeJoin(base, name string) (string, error) {
	clean := filepath.Clean(name)
	target := filepath.Join(base, clean)
	rel, err := filepath.Rel(base, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("invalid archive path")
	}
	return target, nil
}

func findModRoot(extractDir string) string {
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return extractDir
	}
	// If the ZIP has a single top-level folder, use that.
	var dirs []string
	var files int
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join(extractDir, e.Name()))
		} else {
			files++
		}
	}
	if files == 0 && len(dirs) == 1 {
		return dirs[0]
	}
	return extractDir
}

func copyDir(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target, info.Mode())
	})
}

func mergeDir(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dest, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dest string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isPermissionError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, os.ErrPermission) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "access is denied") ||
		strings.Contains(msg, "permission denied") ||
		strings.Contains(msg, "denied")
}
