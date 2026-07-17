package steam

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const amongUsAppID = "945360"

// Install holds a detected Steam root and its library folders.
type Install struct {
	Path      string
	Libraries []string
}

// DetectSteam finds the Steam installation on this machine.
func DetectSteam() (*Install, error) {
	if runtime.GOOS != "windows" {
		return nil, ErrNotFound
	}

	steamPath, err := steamPathFromRegistry()
	if err != nil || steamPath == "" {
		steamPath = steamPathFromCommonLocations()
	}
	if steamPath == "" {
		return nil, ErrNotFound
	}

	libs := []string{filepath.Join(steamPath, "steamapps")}
	extra, _ := parseLibraryFolders(filepath.Join(steamPath, "steamapps", "libraryfolders.vdf"))
	for _, lib := range extra {
		apps := filepath.Join(lib, "steamapps")
		if !containsPath(libs, apps) {
			libs = append(libs, apps)
		}
	}

	return &Install{Path: steamPath, Libraries: libs}, nil
}

// FindAmongUs searches every Steam library for Among Us.
func FindAmongUs() (string, error) {
	install, err := DetectSteam()
	if err != nil {
		return "", err
	}

	for _, lib := range install.Libraries {
		candidate := filepath.Join(lib, "common", "Among Us")
		exe := filepath.Join(candidate, "Among Us.exe")
		if fileExists(exe) {
			return candidate, nil
		}
	}
	return "", ErrGameNotFound
}

func steamPathFromRegistry() (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		key, err = registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, registry.QUERY_VALUE)
		if err != nil {
			return "", err
		}
	}
	defer key.Close()

	path, _, err := key.GetStringValue("SteamPath")
	if err != nil {
		path, _, err = key.GetStringValue("InstallPath")
	}
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), nil
}

func steamPathFromCommonLocations() string {
	candidates := []string{
		`C:\Program Files (x86)\Steam`,
		`C:\Program Files\Steam`,
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Steam"),
		filepath.Join(os.Getenv("ProgramFiles"), "Steam"),
	}
	for _, c := range candidates {
		if fileExists(filepath.Join(c, "steam.exe")) {
			return c
		}
	}
	return ""
}

var pathRe = regexp.MustCompile(`"path"\s+"([^"]+)"`)

func parseLibraryFolders(vdfPath string) ([]string, error) {
	f, err := os.Open(vdfPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var libs []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if m := pathRe.FindStringSubmatch(line); len(m) == 2 {
			p := strings.ReplaceAll(m[1], `\\`, `\`)
			libs = append(libs, filepath.Clean(p))
		}
	}
	_ = amongUsAppID // reserved for future manifest checks
	return libs, scanner.Err()
}

func containsPath(list []string, target string) bool {
	for _, p := range list {
		if strings.EqualFold(p, target) {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
