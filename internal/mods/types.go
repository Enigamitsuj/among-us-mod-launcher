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

	// RequiredGameVersion is the Among Us version this mod targets (e.g. "17.3").
	RequiredGameVersion string `json:"requiredGameVersion"`
	// RecommendedTags are curated release tags known to be stable.
	RecommendedTags []string `json:"recommendedTags"`
}

// IsRecommendedTag reports whether a release tag is in the curated list.
func (m Mod) IsRecommendedTag(tag string) bool {
	for _, t := range m.RecommendedTags {
		if t == tag {
			return true
		}
	}
	return false
}

// Release is a simplified GitHub release shown in the UI.
type Release struct {
	TagName      string `json:"tagName"`
	Name           string `json:"name"`
	Body           string `json:"body"`
	PublishedAt    string `json:"publishedAt"`
	Prerelease     bool   `json:"prerelease"`
	DownloadURL    string `json:"downloadUrl"`
	DownloadName   string `json:"downloadName"`
	DownloadSize   int64  `json:"downloadSize"`
	HTMLURL        string `json:"htmlUrl"`
	AssetMatched   bool   `json:"assetMatched"`
	AssetHint      string `json:"assetHint"`

	// Compatibility of this release against the currently detected game.
	CompatLevel  string `json:"compatLevel"`
	CompatReason string `json:"compatReason"`
	Installable  bool   `json:"installable"`
	Recommended  bool   `json:"recommended"`
}

// GameInstall is one detected Among Us copy (mirrors game.Install for bindings).
type GameInstall struct {
	Platform   string `json:"platform"`
	Path       string `json:"path"`
	Version    string `json:"version"`
	BuildID    string `json:"buildId"`
	Branch     string `json:"branch"`
	Supported  bool   `json:"supported"`
	CanInstall bool   `json:"canInstall"`
	AssetHint  string `json:"assetHint"`
	Message    string `json:"message"`
	LaunchHint string `json:"launchHint"`
}

// GameStatus describes the detected Among Us install(s).
type GameStatus struct {
	Found      bool          `json:"found"`
	Path       string        `json:"path"`
	Version    string        `json:"version"`
	Platform   string        `json:"platform"`
	PlatformLabel string     `json:"platformLabel"`
	BuildID    string        `json:"buildId"`
	Branch     string        `json:"branch"`
	Supported  bool          `json:"supported"`
	CanInstall bool          `json:"canInstall"`
	Running    bool          `json:"running"`
	AssetHint  string        `json:"assetHint"`
	SteamFound bool          `json:"steamFound"`
	SteamPath  string        `json:"steamPath"`
	EpicFound  bool          `json:"epicFound"`
	ItchFound  bool          `json:"itchFound"`
	XboxFound  bool          `json:"xboxFound"`
	Message    string        `json:"message"`
	Installs   []GameInstall `json:"installs"`
}

// InstallOptions are user choices for an install run.
type InstallOptions struct {
	ModID              string `json:"modId"`
	VersionTag         string `json:"versionTag"`
	AmongUsPath        string `json:"amongUsPath"`
	Platform           string `json:"platform"`
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
	Installed bool   `json:"installed"`
	Path      string `json:"path"`
	Version   string `json:"version"`
	Exists    bool   `json:"exists"`
}
