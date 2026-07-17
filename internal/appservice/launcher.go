package appservice

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/Enigamitsuj/among-us-mod-launcher/internal/fsutil"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/game"
	githubapi "github.com/Enigamitsuj/among-us-mod-launcher/internal/githubapi"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/installer"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/installrecord"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/mods"
)

// Launcher is the Wails-bound service the frontend talks to.
type Launcher struct {
	mu         sync.Mutex
	github     *githubapi.Client
	installer  *installer.Service
	window     *application.WebviewWindow
	installing bool
	lastDetect mods.GameStatus
}

func New() *Launcher {
	return &Launcher{
		github:    githubapi.NewClient(),
		installer: installer.New(),
	}
}

// SetWindow stores the main window handle for minimize/close.
func (l *Launcher) SetWindow(w *application.WebviewWindow) {
	l.window = w
}

func (l *Launcher) GetMods() []mods.Mod {
	return mods.All()
}

func (l *Launcher) GetDefaultMod() mods.Mod {
	return mods.Default()
}

func (l *Launcher) DetectGame() mods.GameStatus {
	d := game.Detect()
	status := toGameStatus(d)
	l.mu.Lock()
	l.lastDetect = status
	l.mu.Unlock()
	return status
}

func (l *Launcher) IsGameRunning() bool {
	return game.IsRunning()
}

func (l *Launcher) GetReleases(modID string) ([]mods.Release, error) {
	mod, ok := mods.Get(modID)
	if !ok {
		return nil, errors.New("unknown mod")
	}
	if !mod.Enabled {
		return nil, errors.New("this mod is coming soon")
	}

	platform := game.PlatformUnknown
	l.mu.Lock()
	if l.lastDetect.Found && l.lastDetect.Platform != "" {
		platform = game.Platform(l.lastDetect.Platform)
	}
	l.mu.Unlock()

	// Refresh detection if we have no platform yet.
	var gameVersion string
	l.mu.Lock()
	gameVersion = l.lastDetect.Version
	l.mu.Unlock()

	if platform == game.PlatformUnknown || platform == "" {
		status := l.DetectGame()
		platform = game.Platform(status.Platform)
		gameVersion = status.Version
	}

	releases, err := l.github.ListReleases(mod.GitHubOwner, mod.GitHubRepo, platform, 20)
	if err != nil {
		return nil, err
	}

	for i := range releases {
		recommended := mod.IsRecommendedTag(releases[i].TagName)
		compat := game.EvaluateRelease(
			platform,
			gameVersion,
			releases[i].AssetMatched,
			recommended,
			mod.RequiredGameVersion,
		)
		releases[i].CompatLevel = compat.Level
		releases[i].CompatReason = compat.Reason
		releases[i].Installable = compat.Installable
		releases[i].Recommended = recommended
	}
	return releases, nil
}

func (l *Launcher) GetInstallState(modID, amongUsPath string) mods.InstallState {
	mod, ok := mods.Get(modID)
	if !ok || amongUsPath == "" {
		return mods.InstallState{}
	}
	path := fsutil.JoinInstallPath(filepath.Dir(filepath.Clean(amongUsPath)), mod.FolderName)
	exists := fsutil.DirExists(path)
	exe := filepath.Join(path, "Among Us.exe")
	installed := exists && fileExists(exe)
	version := ""
	managed := false
	if record, ok := installrecord.Load(path); ok && record.ModID == modID {
		version = record.Version
		managed = true
	}
	return mods.InstallState{
		Installed: installed,
		Exists:    exists,
		Path:      path,
		Version:   version,
		Managed:   managed,
	}
}

func (l *Launcher) DestinationExists(modID, amongUsPath string) bool {
	mod, ok := mods.Get(modID)
	if !ok || amongUsPath == "" {
		return false
	}
	return l.installer.DestinationExists(filepath.Dir(filepath.Clean(amongUsPath)), mod.FolderName)
}

// Install starts a background install and emits "install:progress" events.
func (l *Launcher) Install(opts mods.InstallOptions) error {
	l.mu.Lock()
	if l.installing {
		l.mu.Unlock()
		return errors.New("an installation is already in progress")
	}
	l.installing = true
	l.mu.Unlock()

	mod, ok := mods.Get(opts.ModID)
	if !ok || !mod.Enabled {
		l.finishInstall()
		return errors.New("this mod is not available")
	}

	status := l.DetectGame()
	if !status.Found {
		l.finishInstall()
		return errors.New("Among Us not found.")
	}
	if status.Running {
		l.finishInstall()
		return errors.New("Game is currently running.")
	}
	if !status.CanInstall {
		l.finishInstall()
		if status.Message != "" {
			return errors.New(status.Message)
		}
		return errors.New("Unsupported version.")
	}

	platform := game.Platform(status.Platform)
	opts.AmongUsPath = status.Path
	opts.Platform = string(platform)

	releases, err := l.github.ListReleases(mod.GitHubOwner, mod.GitHubRepo, platform, 30)
	if err != nil {
		l.finishInstall()
		return errors.New("GitHub unavailable.")
	}

	var selected *mods.Release
	for i := range releases {
		if releases[i].TagName == opts.VersionTag {
			selected = &releases[i]
			break
		}
	}
	if selected == nil && len(releases) > 0 && opts.VersionTag == "" {
		selected = &releases[0]
		opts.VersionTag = releases[0].TagName
	}
	if selected == nil {
		l.finishInstall()
		return errors.New("selected version was not found")
	}

	compat := game.EvaluateRelease(
		platform,
		status.Version,
		selected.AssetMatched,
		mod.IsRecommendedTag(selected.TagName),
		mod.RequiredGameVersion,
	)
	if !compat.Installable {
		l.finishInstall()
		if compat.Reason != "" {
			return errors.New(compat.Reason)
		}
		return errors.New("This version is not compatible with your Among Us install.")
	}

	installRoot := filepath.Dir(filepath.Clean(opts.AmongUsPath))
	if l.installer.DestinationExists(installRoot, mod.FolderName) && !opts.ForceReinstall {
		l.finishInstall()
		return installer.ErrAlreadyExists
	}

	downloadURL := selected.DownloadURL
	go func() {
		defer l.finishInstall()
		app := application.Get()
		_, err := l.installer.Install(opts, mod, downloadURL, func(p mods.InstallProgress) {
			if app != nil {
				app.Event.Emit("install:progress", p)
			}
		})
		if err != nil {
			if app != nil {
				app.Event.Emit("install:progress", mods.InstallProgress{
					Stage:   "error",
					Message: friendlyError(err),
					Done:    true,
					Error:   friendlyError(err),
				})
			}
			return
		}
	}()

	return nil
}

func (l *Launcher) Uninstall(modID, amongUsPath string) error {
	mod, ok := mods.Get(modID)
	if !ok || amongUsPath == "" {
		return errors.New("installation not found")
	}
	if game.IsRunning() {
		return errors.New("Game is currently running.")
	}

	path := fsutil.JoinInstallPath(filepath.Dir(filepath.Clean(amongUsPath)), mod.FolderName)
	if !fsutil.DirExists(path) {
		return nil
	}
	record, managed := installrecord.Load(path)
	if !managed || record.ModID != modID {
		return errors.New("This folder was not created by the launcher, so it was not removed.")
	}
	if err := os.RemoveAll(path); err != nil {
		return errors.New("could not remove the mod installation")
	}
	return nil
}

func (l *Launcher) LaunchMod(installPath string) error {
	// Prefer Epic starter when present.
	starter := filepath.Join(installPath, "EpicGamesStarter.exe")
	exe := filepath.Join(installPath, "Among Us.exe")
	target := exe
	if fileExists(starter) {
		target = starter
	}
	if !fileExists(target) {
		return errors.New("installation not found")
	}
	cmd := exec.Command(target)
	cmd.Dir = installPath
	return cmd.Start()
}

func (l *Launcher) MinimizeWindow() {
	if l.window != nil {
		l.window.Minimise()
	}
}

func (l *Launcher) CloseWindow() {
	if l.window != nil {
		l.window.Close()
	}
}

func (l *Launcher) finishInstall() {
	l.mu.Lock()
	l.installing = false
	l.mu.Unlock()
}

func toGameStatus(d game.Detection) mods.GameStatus {
	status := mods.GameStatus{
		Found:      d.Found,
		Running:    d.Running,
		Message:    d.Message,
		SteamFound: d.SteamFound,
		SteamPath:  d.SteamPath,
		EpicFound:  d.EpicFound,
		ItchFound:  d.ItchFound,
		XboxFound:  d.XboxFound,
	}
	if !d.Found {
		status.Supported = false
		status.CanInstall = false
		return status
	}

	p := d.Primary
	status.Path = p.Path
	status.Version = p.Version
	status.Platform = string(p.Platform)
	status.PlatformLabel = p.Platform.DisplayName()
	status.BuildID = p.BuildID
	status.Branch = p.Branch
	status.Supported = p.Supported
	status.CanInstall = p.CanInstall
	status.AssetHint = p.AssetHint
	status.Message = d.Message

	status.Installs = make([]mods.GameInstall, 0, len(d.Installs))
	for _, inst := range d.Installs {
		status.Installs = append(status.Installs, mods.GameInstall{
			Platform:   string(inst.Platform),
			Path:       inst.Path,
			Version:    inst.Version,
			BuildID:    inst.BuildID,
			Branch:     inst.Branch,
			Supported:  inst.Supported,
			CanInstall: inst.CanInstall,
			AssetHint:  inst.AssetHint,
			Message:    inst.Message,
			LaunchHint: inst.LaunchHint,
		})
	}
	return status
}

func friendlyError(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, installer.ErrAlreadyExists) {
		return "Town of Us already exists."
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "Among Us not found"):
		return "Among Us not found."
	case strings.Contains(msg, "GitHub"):
		return "GitHub unavailable."
	case strings.Contains(msg, "running"):
		return "Game is currently running."
	case strings.Contains(msg, "Unsupported") || strings.Contains(msg, "downgrade"):
		return msg
	default:
		return msg
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
