package fsutil

import (
	"os"
	"path/filepath"
)

// LauncherDir returns the directory containing the running executable.
func LauncherDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		resolved = exe
	}
	return filepath.Dir(resolved), nil
}

// DefaultInstallRoot is next to the launcher executable (not AppData/Temp).
func DefaultInstallRoot() (string, error) {
	return LauncherDir()
}

// JoinInstallPath builds "<root>/<mod folder name>".
func JoinInstallPath(root, folderName string) string {
	return filepath.Join(root, folderName)
}

// DirExists reports whether path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// EnsureDir creates a directory if missing.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}
