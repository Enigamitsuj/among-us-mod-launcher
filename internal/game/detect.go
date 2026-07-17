package game

import (
	"path/filepath"
	"sort"

	"github.com/Enigamitsuj/among-us-mod-launcher/internal/fsutil"
)

// Detect scans Steam, Epic, Itch, and Xbox for Among Us installs.
func Detect() Detection {
	d := Detection{
		Message: "Among Us not found.",
	}

	steamPath, steamInstalls := findSteamInstalls()
	if steamPath != "" {
		d.SteamFound = true
		d.SteamPath = steamPath
	}
	epicFound, epicInstalls := findEpicInstalls()
	d.EpicFound = epicFound
	itchFound, itchInstalls := findItchInstalls()
	d.ItchFound = itchFound
	xboxFound, xboxInstalls := findXboxInstalls()
	d.XboxFound = xboxFound

	all := append([]Install{}, steamInstalls...)
	all = append(all, itchInstalls...)
	all = append(all, epicInstalls...)
	all = append(all, xboxInstalls...)

	for i := range all {
		applyInstallWritability(&all[i])
	}

	d.Running = IsRunning()
	if len(all) == 0 {
		if !d.SteamFound && !d.EpicFound && !d.ItchFound && !d.XboxFound {
			d.Message = "No game stores with Among Us were detected."
		} else {
			d.Message = "Among Us not found."
		}
		return d
	}

	sort.SliceStable(all, func(i, j int) bool {
		// Prefer installable platforms, then preference order.
		if all[i].CanInstall != all[j].CanInstall {
			return all[i].CanInstall
		}
		return PreferOrder(all[i].Platform) < PreferOrder(all[j].Platform)
	})

	d.Found = true
	d.Installs = all
	d.Primary = all[0]

	if d.Running {
		d.Message = "Game is currently running. Close Among Us before installing."
		return d
	}
	d.Message = d.Primary.Message
	return d
}

func applyInstallWritability(inst *Install) {
	if inst == nil || !inst.CanInstall || inst.Path == "" {
		return
	}
	parent := filepath.Dir(filepath.Clean(inst.Path))
	if fsutil.DirIsWritable(parent) {
		return
	}
	inst.CanInstall = false
	inst.Message = "Among Us is in a protected folder (often Program Files). Move the game to a user library folder, or run the launcher as administrator."
}
