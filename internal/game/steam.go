package game

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const steamAppID = "945360"

var steamPathRe = regexp.MustCompile(`"path"\s+"([^"]+)"`)

func findSteamInstalls() (steamPath string, installs []Install) {
	if runtime.GOOS != "windows" {
		return "", nil
	}
	steamPath = steamRoot()
	if steamPath == "" {
		return "", nil
	}

	libs := []string{filepath.Join(steamPath, "steamapps")}
	extra, _ := parseSteamLibraries(filepath.Join(steamPath, "steamapps", "libraryfolders.vdf"))
	for _, lib := range extra {
		apps := filepath.Join(lib, "steamapps")
		if !containsFold(libs, apps) {
			libs = append(libs, apps)
		}
	}

	for _, lib := range libs {
		candidate := filepath.Join(lib, "common", "Among Us")
		if !fileExists(filepath.Join(candidate, "Among Us.exe")) {
			continue
		}
		inst := Install{
			Platform: PlatformSteam,
			Path:     candidate,
			Version:  ReadVersion(candidate),
		}
		build, branch := readSteamManifest(lib)
		inst.BuildID = build
		inst.Branch = branch
		inst.Version = RefineSteamVersion(inst.Version, branch, build)
		EvaluateCompatibility(&inst)
		installs = append(installs, inst)
	}
	return steamPath, installs
}

func steamRoot() string {
	if p, err := steamPathFromRegistry(); err == nil && p != "" {
		return p
	}
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

func parseSteamLibraries(vdfPath string) ([]string, error) {
	f, err := os.Open(vdfPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var libs []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if m := steamPathRe.FindStringSubmatch(scanner.Text()); len(m) == 2 {
			p := strings.ReplaceAll(m[1], `\\`, `\`)
			libs = append(libs, filepath.Clean(p))
		}
	}
	return libs, scanner.Err()
}

func readSteamManifest(steamapps string) (buildID, branch string) {
	manifest := filepath.Join(steamapps, "appmanifest_"+steamAppID+".acf")
	data, err := os.ReadFile(manifest)
	if err != nil {
		return "", ""
	}
	text := string(data)
	buildID = vdfValue(text, "buildid")
	branch = vdfValue(text, "BetaKey")
	if branch == "" {
		branch = "public"
	}
	return buildID, branch
}

func vdfValue(text, key string) string {
	re := regexp.MustCompile(`"` + regexp.QuoteMeta(key) + `"\s+"([^"]*)"`)
	if m := re.FindStringSubmatch(text); len(m) == 2 {
		return m[1]
	}
	return ""
}

func containsFold(list []string, target string) bool {
	for _, p := range list {
		if strings.EqualFold(p, target) {
			return true
		}
	}
	return false
}
