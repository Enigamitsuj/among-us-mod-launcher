package shortcut

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// CreateDesktopShortcut creates a .lnk on the desktop that launches targetExe,
// using iconExe (typically the original Among Us.exe) for the icon.
func CreateDesktopShortcut(name, targetExe, iconExe, workingDir string) error {
	desktop, err := desktopDir()
	if err != nil {
		return err
	}
	lnk := filepath.Join(desktop, name+".lnk")
	if iconExe == "" {
		iconExe = targetExe
	}
	if workingDir == "" {
		workingDir = filepath.Dir(targetExe)
	}

	ps := fmt.Sprintf(`
$ws = New-Object -ComObject WScript.Shell
$s = $ws.CreateShortcut(%s)
$s.TargetPath = %s
$s.WorkingDirectory = %s
$s.IconLocation = %s
$s.Save()
`, psQuote(lnk), psQuote(targetExe), psQuote(workingDir), psQuote(iconExe+",0"))

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("could not create desktop shortcut: %s", string(out))
	}
	return nil
}

func desktopDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Prefer the public Desktop path via USERPROFILE\Desktop
	desktop := filepath.Join(home, "Desktop")
	if info, err := os.Stat(desktop); err == nil && info.IsDir() {
		return desktop, nil
	}
	// OneDrive Desktop fallback
	onedrive := filepath.Join(home, "OneDrive", "Desktop")
	if info, err := os.Stat(onedrive); err == nil && info.IsDir() {
		return onedrive, nil
	}
	return desktop, nil
}

func psQuote(s string) string {
	return "'" + filepath.Clean(s) + "'"
}
