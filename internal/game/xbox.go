package game

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func findXboxInstalls() (xboxFound bool, installs []Install) {
	seen := map[string]bool{}

	// Modern Xbox PC Game Pass / Xbox app path.
	for _, candidate := range []string{
		`C:\XboxGames\Among Us\Content`,
		filepath.Join(os.Getenv("SystemDrive")+`\`, "XboxGames", "Among Us", "Content"),
	} {
		if isAmongUsDir(candidate) {
			xboxFound = true
			addXbox(&installs, seen, candidate)
		}
	}

	// Appx package install location (when readable).
	ps := `Get-AppxPackage -Name "*AmongUs*" -ErrorAction SilentlyContinue | Select-Object -ExpandProperty InstallLocation`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			root := strings.TrimSpace(line)
			if root == "" {
				continue
			}
			xboxFound = true
			for _, candidate := range []string{
				filepath.Join(root, "Content"),
				root,
			} {
				if isAmongUsDir(candidate) {
					addXbox(&installs, seen, candidate)
				}
			}
		}
	}

	return xboxFound, installs
}

func addXbox(installs *[]Install, seen map[string]bool, path string) {
	key := strings.ToLower(path)
	if seen[key] {
		return
	}
	seen[key] = true
	inst := Install{
		Platform: PlatformXbox,
		Path:     path,
		Version:  ReadVersion(path),
	}
	EvaluateCompatibility(&inst)
	*installs = append(*installs, inst)
}
