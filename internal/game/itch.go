package game

import (
	"os"
	"path/filepath"
	"strings"
)

func findItchInstalls() (itchFound bool, installs []Install) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return false, nil
	}
	itchRoot := filepath.Join(appData, "itch")
	if !dirExists(itchRoot) {
		return false, nil
	}
	itchFound = true
	seen := map[string]bool{}

	// Default itch apps folder + a shallow walk for Among Us.exe.
	candidates := []string{
		filepath.Join(itchRoot, "apps"),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "itch", "apps"),
	}

	for _, root := range candidates {
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil || d == nil {
				return nil
			}
			// Keep walk shallow-ish for speed.
			rel, relErr := filepath.Rel(root, path)
			if relErr == nil && strings.Count(rel, string(os.PathSeparator)) > 4 {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.IsDir() {
				return nil
			}
			if !strings.EqualFold(d.Name(), "Among Us.exe") {
				return nil
			}
			dir := filepath.Dir(path)
			key := strings.ToLower(dir)
			if seen[key] {
				return nil
			}
			seen[key] = true
			inst := Install{
				Platform: PlatformItch,
				Path:     dir,
				Version:  ReadVersion(dir),
			}
			EvaluateCompatibility(&inst)
			installs = append(installs, inst)
			return nil
		})
	}

	return itchFound, installs
}
