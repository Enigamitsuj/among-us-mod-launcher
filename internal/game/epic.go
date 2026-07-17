package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type epicManifest struct {
	DisplayName     string `json:"DisplayName"`
	InstallLocation string `json:"InstallLocation"`
	AppName         string `json:"AppName"`
	CatalogItemId   string `json:"CatalogItemId"`
}

func findEpicInstalls() (epicFound bool, installs []Install) {
	roots := epicManifestDirs()
	seen := map[string]bool{}

	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		epicFound = true
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".item") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(root, e.Name()))
			if err != nil {
				continue
			}
			var m epicManifest
			if err := json.Unmarshal(data, &m); err != nil {
				continue
			}
			if !looksLikeAmongUs(m.DisplayName, m.AppName, m.InstallLocation) {
				continue
			}
			path := filepath.Clean(m.InstallLocation)
			if !isAmongUsDir(path) {
				// Sometimes InstallLocation is the parent.
				alt := filepath.Join(path, "AmongUs")
				if isAmongUsDir(alt) {
					path = alt
				} else {
					alt = filepath.Join(path, "Among Us")
					if isAmongUsDir(alt) {
						path = alt
					} else {
						continue
					}
				}
			}
			key := strings.ToLower(path)
			if seen[key] {
				continue
			}
			seen[key] = true
			inst := Install{
				Platform: PlatformEpic,
				Path:     path,
				Version:  ReadVersion(path),
			}
			EvaluateCompatibility(&inst)
			installs = append(installs, inst)
		}
	}

	// Fallback common path
	for _, candidate := range []string{
		`C:\Program Files\Epic Games\AmongUs`,
		`C:\Program Files\Epic Games\Among Us`,
	} {
		if !isAmongUsDir(candidate) {
			continue
		}
		key := strings.ToLower(candidate)
		if seen[key] {
			continue
		}
		epicFound = true
		seen[key] = true
		inst := Install{
			Platform: PlatformEpic,
			Path:     candidate,
			Version:  ReadVersion(candidate),
		}
		EvaluateCompatibility(&inst)
		installs = append(installs, inst)
	}

	return epicFound, installs
}

func epicManifestDirs() []string {
	programData := os.Getenv("ProgramData")
	if programData == "" {
		programData = `C:\ProgramData`
	}
	return []string{
		filepath.Join(programData, "Epic", "EpicGamesLauncher", "Data", "Manifests"),
	}
}

func looksLikeAmongUs(parts ...string) bool {
	for _, p := range parts {
		l := strings.ToLower(p)
		if strings.Contains(l, "among us") || strings.Contains(l, "amongus") {
			return true
		}
	}
	return false
}

func isAmongUsDir(path string) bool {
	return fileExists(filepath.Join(path, "Among Us.exe"))
}
