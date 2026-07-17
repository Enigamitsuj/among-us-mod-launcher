package game

import "strings"

// Compatibility levels for a specific mod release against the detected game.
const (
	CompatCompatible  = "compatible"
	CompatCaution     = "caution"
	CompatUnsupported = "unsupported"
)

// GameVersionUnsupported reports whether the Among Us version is the mod-breaking
// latest line (currently v17.4).
func GameVersionUnsupported(version string) bool {
	return isUnsupportedLatest(version)
}

// GameVersionCompatible reports whether the Among Us version is the modding target
// (currently v17.3).
func GameVersionCompatible(version string) bool {
	return isCompatiblePrevious(version)
}

// ReleaseCompat is the result of evaluating one mod release against the game.
type ReleaseCompat struct {
	Level       string
	Reason      string
	Installable bool
}

// EvaluateRelease decides whether a mod release can be installed onto the
// detected game right now. requiredVersion is the Among Us version the mod
// targets (e.g. "17.3"); recommended marks curated stable tags.
func EvaluateRelease(platform Platform, gameVersion string, assetMatched bool, recommended bool, requiredVersion string) ReleaseCompat {
	// Hard block: platform not supported at all.
	if platform == PlatformXbox {
		return ReleaseCompat{
			Level:       CompatUnsupported,
			Reason:      "Microsoft Store / Xbox can't be downgraded yet, so this version can't be installed.",
			Installable: false,
		}
	}

	// Hard block: this release has no package for the detected storefront.
	if AssetPattern(platform) != "" && !assetMatched {
		return ReleaseCompat{
			Level:       CompatUnsupported,
			Reason:      "This version has no " + AssetPattern(platform) + " download for " + platform.DisplayName() + ".",
			Installable: false,
		}
	}

	// Hard block: game on the mod-breaking latest version.
	if GameVersionUnsupported(gameVersion) {
		hint := downgradeHint(platform, requiredVersion)
		return ReleaseCompat{
			Level:       CompatUnsupported,
			Reason:      "Needs Among Us v" + normalizeReq(requiredVersion) + " — you're on v" + displayVer(gameVersion) + ". " + hint,
			Installable: false,
		}
	}

	// Compatible game version.
	if GameVersionCompatible(gameVersion) {
		if recommended {
			return ReleaseCompat{
				Level:       CompatCompatible,
				Reason:      "Recommended for Among Us v" + displayVer(gameVersion) + ".",
				Installable: true,
			}
		}
		return ReleaseCompat{
			Level:       CompatCompatible,
			Reason:      "Compatible with Among Us v" + displayVer(gameVersion) + ".",
			Installable: true,
		}
	}

	// Unknown game version — allow but caution.
	return ReleaseCompat{
		Level:       CompatCaution,
		Reason:      "Among Us version couldn't be confirmed. This targets v" + normalizeReq(requiredVersion) + " — install only if you're on that version.",
		Installable: true,
	}
}

func downgradeHint(platform Platform, requiredVersion string) string {
	switch platform {
	case PlatformSteam:
		return "Set Steam beta to “public previous”, then reopen the launcher."
	case PlatformItch:
		return "In Itch, switch to version 113, then reopen the launcher."
	case PlatformEpic:
		return "Use an Among Us v" + normalizeReq(requiredVersion) + " build for Epic."
	default:
		return ""
	}
}

func normalizeReq(requiredVersion string) string {
	v := strings.TrimSpace(requiredVersion)
	if v == "" {
		return "17.3"
	}
	return strings.TrimPrefix(v, "v")
}
