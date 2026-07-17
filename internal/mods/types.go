package mods

// Channel represents a release channel for a mod.
type Channel string

const (
	ChannelLatest Channel = "latest"
	ChannelStable Channel = "stable"
	ChannelBeta   Channel = "beta"
)

// Mod describes a community mod the launcher can install.
type Mod struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Description string `json:"description"`
	Author      string `json:"author"`
	GitHubOwner string `json:"githubOwner"`
	GitHubRepo  string `json:"githubRepo"`
	FolderName  string `json:"folderName"`
	AccentFrom  string `json:"accentFrom"`
	AccentTo    string `json:"accentTo"`
	Enabled     bool   `json:"enabled"`
	ComingSoon  bool   `json:"comingSoon"`
}

// Release is a simplified GitHub release shown in the UI.
type Release struct {
	TagName     string `json:"tagName"`
	Name          string `json:"name"`
	Body          string `json:"body"`
	PublishedAt   string `json:"publishedAt"`
	Prerelease    bool   `json:"prerelease"`
	DownloadURL   string `json:"downloadUrl"`
	DownloadName  string `json:"downloadName"`
	DownloadSize  int64  `json:"downloadSize"`
	HTMLURL       string `json:"htmlUrl"`
}

// GameStatus describes the detected Among Us install.
type GameStatus struct {
	Found         bool   `json:"found"`
	Path          string `json:"path"`
	Version       string `json:"version"`
	Supported     bool   `json:"supported"`
	Running       bool   `json:"running"`
	SteamFound    bool   `json:"steamFound"`
	SteamPath     string `json:"steamPath"`
	Message       string `json:"message"`
}

// InstallOptions are user choices for an install run.
type InstallOptions struct {
	ModID              string `json:"modId"`
	VersionTag         string `json:"versionTag"`
	AmongUsPath        string `json:"amongUsPath"`
	InstallLocation    string `json:"installLocation"`
	CreateShortcut     bool   `json:"createShortcut"`
	LaunchAfterInstall bool   `json:"launchAfterInstall"`
	ForceReinstall     bool   `json:"forceReinstall"`
}

// InstallProgress is emitted during installation.
type InstallProgress struct {
	Stage      string  `json:"stage"`
	Message    string  `json:"message"`
	Percent    float64 `json:"percent"`
	Done       bool    `json:"done"`
	Error      string  `json:"error"`
	InstallDir string  `json:"installDir"`
}

// InstallState describes whether a mod is already installed.
type InstallState struct {
	Installed  bool   `json:"installed"`
	Path       string `json:"path"`
	Version    string `json:"version"`
	Exists     bool   `json:"exists"`
}
