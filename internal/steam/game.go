package steam

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// IsGameRunning reports whether Among Us.exe is currently running.
func IsGameRunning() bool {
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

// ReadGameVersion attempts to read the game version from amongus_data or similar.
// Falls back to unknown when unavailable.
func ReadGameVersion(gamePath string) string {
	candidates := []string{
		filepath.Join(gamePath, "Among Us_Data", "StreamingAssets", "aa", "catalog.json"),
		filepath.Join(gamePath, "version.txt"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			// Version parsing is refined later; presence alone is enough for now.
			break
		}
	}
	return "unknown"
}
