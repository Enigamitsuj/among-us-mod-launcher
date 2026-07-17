package game

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

var (
	amongUsPublicVerRe = regexp.MustCompile(`(?i)\b(?:v)?(17\.\d+(?:\.\d+)?)\b`)
	unityYearRe        = regexp.MustCompile(`\b(20\d{2}\.\d+\.\d+)\b`)
)

// IsRunning reports whether Among Us.exe is currently running.
func IsRunning() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq Among Us.exe", "/NH")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), "among us.exe")
}

// ReadVersion tries several strategies to determine a user-facing Among Us version.
// Prefer community labels like "17.3" over Unity engine ProductVersion strings.
func ReadVersion(gamePath string) string {
	if v := versionFromGlobalGameManagers(gamePath); v != "" {
		return v
	}
	if v := versionFromVersionFile(gamePath); v != "" {
		return v
	}
	if v := versionFromExe(gamePath); v != "" {
		// Only keep exe metadata if it looks like an Among Us public version.
		if m := amongUsPublicVerRe.FindStringSubmatch(v); len(m) == 2 {
			return m[1]
		}
	}
	return "unknown"
}

// RefineSteamVersion uses branch/build metadata when file probing is inconclusive.
func RefineSteamVersion(version, branch, buildID string) string {
	b := strings.ToLower(strings.TrimSpace(branch))
	switch {
	case strings.Contains(b, "previous"):
		return "17.3"
	case b == "public" || b == "":
		// Latest public branch is currently the unstable 17.4 line per TOU docs.
		if version == "unknown" || version == "" {
			return "17.4"
		}
	}
	_ = buildID
	if version != "" && version != "unknown" {
		return version
	}
	return "unknown"
}

func versionFromExe(gamePath string) string {
	exe := filepath.Join(gamePath, "Among Us.exe")
	if !fileExists(exe) {
		return ""
	}
	escaped := strings.ReplaceAll(exe, "'", "''")
	ps := `(Get-Item -LiteralPath '` + escaped + `').VersionInfo.ProductVersion`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err == nil {
		if v := strings.TrimSpace(string(out)); v != "" {
			return v
		}
	}
	ps = `(Get-Item -LiteralPath '` + escaped + `').VersionInfo.FileVersion`
	cmd = exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err = cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func versionFromGlobalGameManagers(gamePath string) string {
	path := filepath.Join(gamePath, "Among Us_Data", "globalgamemanagers")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	text := string(filterPrintable(data))
	if m := amongUsPublicVerRe.FindStringSubmatch(text); len(m) == 2 {
		return m[1]
	}
	// Fall back to year.build style if present (e.g. 2024.8.13).
	if m := unityYearRe.FindStringSubmatch(text); len(m) == 2 {
		return mapYearBuildToPublic(m[1])
	}
	return ""
}

func versionFromVersionFile(gamePath string) string {
	path := filepath.Join(gamePath, "version.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	raw := strings.TrimSpace(string(data))
	if m := amongUsPublicVerRe.FindStringSubmatch(raw); len(m) == 2 {
		return m[1]
	}
	return raw
}

func mapYearBuildToPublic(yearBuild string) string {
	// Keep the raw year.build for display when we cannot map confidently.
	// Known recent lines used by the community docs:
	switch {
	case strings.HasPrefix(yearBuild, "2024.9"), strings.HasPrefix(yearBuild, "2025."):
		return "17.4"
	case strings.HasPrefix(yearBuild, "2024.8"), strings.HasPrefix(yearBuild, "2024.6"), strings.HasPrefix(yearBuild, "2024.3"):
		return "17.3"
	default:
		return yearBuild
	}
}

func filterPrintable(in []byte) []byte {
	out := make([]byte, 0, len(in)/8)
	for _, b := range in {
		if b >= 32 && b < 127 {
			out = append(out, b)
		} else {
			out = append(out, ' ')
		}
	}
	return out
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
