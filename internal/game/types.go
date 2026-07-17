package game

// Platform identifies where Among Us was installed from.
type Platform string

const (
	PlatformUnknown Platform = "unknown"
	PlatformSteam   Platform = "steam"
	PlatformEpic    Platform = "epic"
	PlatformItch    Platform = "itch"
	PlatformXbox    Platform = "xbox"
)

// DisplayName returns a user-facing storefront label.
func (p Platform) DisplayName() string {
	switch p {
	case PlatformSteam:
		return "Steam"
	case PlatformEpic:
		return "Epic Games"
	case PlatformItch:
		return "Itch.io"
	case PlatformXbox:
		return "Microsoft Store / Xbox"
	default:
		return "Unknown"
	}
}

// Install is one detected Among Us installation.
type Install struct {
	Platform    Platform `json:"platform"`
	Path        string   `json:"path"`
	Version     string   `json:"version"`
	BuildID     string   `json:"buildId"`
	Branch      string   `json:"branch"`
	Supported   bool     `json:"supported"`
	CanInstall  bool     `json:"canInstall"`
	AssetHint   string   `json:"assetHint"`
	Message     string   `json:"message"`
	LaunchHint  string   `json:"launchHint"`
}

// Detection is the full scan result shown in the UI.
type Detection struct {
	Found      bool      `json:"found"`
	Primary    Install   `json:"primary"`
	Installs   []Install `json:"installs"`
	Running    bool      `json:"running"`
	Message    string    `json:"message"`
	SteamFound bool      `json:"steamFound"`
	SteamPath  string    `json:"steamPath"`
	EpicFound  bool      `json:"epicFound"`
	ItchFound  bool      `json:"itchFound"`
	XboxFound  bool      `json:"xboxFound"`
}
