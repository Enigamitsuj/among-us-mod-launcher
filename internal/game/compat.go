package game

import (
	"strings"
)

// AssetPattern returns the release ZIP filename substring for a platform.
func AssetPattern(p Platform) string {
	switch p {
	case PlatformSteam, PlatformItch:
		return "steam-itch"
	case PlatformEpic, PlatformXbox:
		return "epic-msstore"
	default:
		return ""
	}
}

// EvaluateCompatibility applies current Town of Us: Mira platform rules.
// Docs (2026): Among Us v17.4 is unstable with mods; Steam/Itch should be on v17.3.
// Microsoft Store cannot downgrade and is currently unsupported.
func EvaluateCompatibility(inst *Install) {
	if inst == nil {
		return
	}
	inst.AssetHint = AssetPattern(inst.Platform)

	switch inst.Platform {
	case PlatformXbox:
		inst.Supported = false
		inst.CanInstall = false
		inst.Message = "Microsoft Store / Xbox cannot be downgraded. Town of Us: Mira is currently unsupported on this platform."
		inst.LaunchHint = "Launch from the Xbox app after a manual install (when support returns)."
		return

	case PlatformSteam:
		inst.LaunchHint = "Launch Among Us.exe from the modded folder."
		if isUnsupportedLatest(inst.Version) {
			inst.Supported = false
			inst.CanInstall = false
			inst.Message = "Among Us v" + displayVer(inst.Version) + " detected. Switch Steam beta to “public previous” (v17.3), then come back."
			return
		}
		inst.Supported = true
		inst.CanInstall = true
		if inst.Version == "unknown" || inst.Version == "" {
			inst.Message = "Steam Among Us found. Version could not be confirmed — use beta “public previous” (v17.3) if install fails."
			return
		}
		inst.Message = "Steam · Among Us v" + displayVer(inst.Version) + " · Compatible"
		return

	case PlatformItch:
		inst.LaunchHint = "Launch Among Us.exe from the modded folder."
		if isUnsupportedLatest(inst.Version) {
			inst.Supported = false
			inst.CanInstall = false
			inst.Message = "Among Us v" + displayVer(inst.Version) + " detected. In Itch, switch to version 113 (v17.3), then try again."
			return
		}
		inst.Supported = true
		inst.CanInstall = true
		if inst.Version == "unknown" || inst.Version == "" {
			inst.Message = "Itch.io Among Us found. Confirm you are on version 113 (v17.3)."
			return
		}
		inst.Message = "Itch.io · Among Us v" + displayVer(inst.Version) + " · Compatible"
		return

	case PlatformEpic:
		inst.LaunchHint = "Epic installs should be launched with EpicGamesStarter.exe (outside Program Files)."
		if isUnsupportedLatest(inst.Version) {
			inst.Supported = false
			inst.CanInstall = false
			inst.Message = "Among Us v" + displayVer(inst.Version) + " on Epic is currently unstable with mods."
			return
		}
		inst.Supported = true
		inst.CanInstall = true
		if inst.Version == "unknown" || inst.Version == "" {
			inst.Message = "Epic Among Us found. The modded copy will be installed beside the game (required for Epic)."
			return
		}
		inst.Message = "Epic Games · Among Us v" + displayVer(inst.Version) + " · Compatible"
		return

	default:
		inst.Supported = false
		inst.CanInstall = false
		inst.Message = "Unknown Among Us install."
	}
}

func isUnsupportedLatest(version string) bool {
	v := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(version)), "v")
	return strings.HasPrefix(v, "17.4") || v == "174" || strings.HasPrefix(v, "17.4.")
}

func isCompatiblePrevious(version string) bool {
	v := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(version)), "v")
	return strings.HasPrefix(v, "17.3") || v == "173"
}

func displayVer(version string) string {
	v := strings.TrimSpace(version)
	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return "unknown"
	}
	return v
}

// PreferOrder ranks platforms for auto-selection (Steam first for modding).
func PreferOrder(p Platform) int {
	switch p {
	case PlatformSteam:
		return 0
	case PlatformItch:
		return 1
	case PlatformEpic:
		return 2
	case PlatformXbox:
		return 3
	default:
		return 99
	}
}
