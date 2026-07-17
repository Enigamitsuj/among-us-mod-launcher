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
	githubapi "github.com/Enigamitsuj/among-us-mod-launcher/internal/githubapi"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/installer"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/mods"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/shortcut"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/steam"
)

// Launcher is the Wails-bound service the frontend talks to.
type Launcher struct {
	mu         sync.Mutex
	github     *githubapi.Client
	installer  *installer.Service
	window     *application.WebviewWindow
	installing bool
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

func (l *Launcher) GetDefaultInstallLocation() (string, error) {
	return fsutil.DefaultInstallRoot()
}

func (l *Launcher) DetectGame() mods.GameStatus {
	status := mods.GameStatus{
		Supported: true,
		Message:   "Ready to install",
	}

	steamInstall, err := steam.DetectSteam()
	if err != nil {
		status.Message = "Steam not installed."
		return status
	}
	status.SteamFound = true
	status.SteamPath = steamInstall.Path

	gamePath, err := steam.FindAmongUs()
	if err != nil {
		status.Message = "Among Us not found."
		return status
	}

	status.Found = true
	status.Path = gamePath
	status.Version = steam.ReadGameVersion(gamePath)
	status.Running = steam.IsGameRunning()
	if status.Running {
		status.Message = "Game is currently running. Close Among Us before installing."
	}
	return status
}

func (l *Launcher) GetReleases(modID string) ([]mods.Release, error) {
	mod, ok := mods.Get(modID)
	if !ok {
		return nil, errors.New("unknown mod")
	}
	if !mod.Enabled {
		return nil, errors.New("this mod is coming soon")
	}
	return l.github.ListReleases(mod.GitHubOwner, mod.GitHubRepo, 20)
}

func (l *Launcher) GetInstallState(modID, installLocation string) mods.InstallState {
	mod, ok := mods.Get(modID)
	if !ok {
		return mods.InstallState{}
	}
	if installLocation == "" {
		installLocation, _ = fsutil.DefaultInstallRoot()
	}
	path := fsutil.JoinInstallPath(installLocation, mod.FolderName)
	exists := fsutil.DirExists(path)
	exe := filepath.Join(path, "Among Us.exe")
	installed := exists && fileExists(exe)
	return mods.InstallState{
		Installed: installed,
		Exists:    exists,
		Path:      path,
	}
}

func (l *Launcher) DestinationExists(modID, installLocation string) bool {
	mod, ok := mods.Get(modID)
	if !ok {
		return false
	}
	if installLocation == "" {
		installLocation, _ = fsutil.DefaultInstallRoot()
	}
	return l.installer.DestinationExists(installLocation, mod.FolderName)
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

	releases, err := l.github.ListReleases(mod.GitHubOwner, mod.GitHubRepo, 30)
	if err != nil {
		l.finishInstall()
		return errors.New("GitHub unavailable.")
	}

	downloadURL := ""
	for _, r := range releases {
		if r.TagName == opts.VersionTag {
			downloadURL = r.DownloadURL
			break
		}
	}
	if downloadURL == "" && len(releases) > 0 && opts.VersionTag == "" {
		downloadURL = releases[0].DownloadURL
		opts.VersionTag = releases[0].TagName
	}
	if downloadURL == "" {
		l.finishInstall()
		return errors.New("selected version was not found")
	}

	if opts.AmongUsPath == "" {
		path, err := steam.FindAmongUs()
		if err != nil {
			l.finishInstall()
			return errors.New("Among Us not found.")
		}
		opts.AmongUsPath = path
	}

	if steam.IsGameRunning() {
		l.finishInstall()
		return errors.New("Game is currently running.")
	}

	if l.installer.DestinationExists(opts.InstallLocation, mod.FolderName) && !opts.ForceReinstall {
		l.finishInstall()
		return installer.ErrAlreadyExists
	}

	go func() {
		defer l.finishInstall()
		app := application.Get()
		dest, err := l.installer.Install(opts, mod, downloadURL, func(p mods.InstallProgress) {
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

		if opts.CreateShortcut {
			_ = shortcut.CreateDesktopShortcut(
				mod.FolderName,
				filepath.Join(dest, "Among Us.exe"),
				filepath.Join(opts.AmongUsPath, "Among Us.exe"),
				dest,
			)
		}

		if opts.LaunchAfterInstall {
			_ = l.LaunchMod(dest)
		}
	}()

	return nil
}

func (l *Launcher) LaunchMod(installPath string) error {
	exe := filepath.Join(installPath, "Among Us.exe")
	if !fileExists(exe) {
		return errors.New("installation not found")
	}
	cmd := exec.Command(exe)
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
	case strings.Contains(msg, "Steam"):
		return "Steam not installed."
	case strings.Contains(msg, "GitHub"):
		return "GitHub unavailable."
	case strings.Contains(msg, "running"):
		return "Game is currently running."
	default:
		return msg
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
