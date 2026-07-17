package mods

import "sync"

var (
	mu       sync.RWMutex
	registry []Mod
)

func init() {
	registry = []Mod{
		{
			ID:          "tou-mira",
			Name:        "Town of Us: Mira",
			ShortName:   "TOU Mira",
			Description: "The community-favorite Town of Us experience, rebuilt for modern Among Us. Crewmates, impostors, and dozens of unique roles — ready in one click.",
			Author:      "AU-Avengers",
			GitHubOwner: "AU-Avengers",
			GitHubRepo:  "TOU-Mira",
			FolderName:  "Among Us - TOU Mira",
			AccentFrom:  "#a855f7",
			AccentTo:    "#ef4444",
			Enabled:     true,
			ComingSoon:  false,

			RequiredGameVersion: "17.3",
			RecommendedTags:     []string{"1.6.2", "1.6.3-beta2"},
		},
		{
			ID:          "better-crewlink",
			Name:        "BetterCrewLink",
			ShortName:   "BCL",
			Description: "Proximity voice chat for Among Us. Coming soon to the launcher.",
			Author:      "Community",
			FolderName:  "Among Us - BetterCrewLink",
			AccentFrom:  "#38bdf8",
			AccentTo:    "#6366f1",
			Enabled:     false,
			ComingSoon:  true,
		},
	}
}

// All returns every registered mod.
func All() []Mod {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Mod, len(registry))
	copy(out, registry)
	return out
}

// Enabled returns mods that are available to install.
func Enabled() []Mod {
	mu.RLock()
	defer mu.RUnlock()
	var out []Mod
	for _, m := range registry {
		if m.Enabled {
			out = append(out, m)
		}
	}
	return out
}

// Get returns a mod by ID.
func Get(id string) (Mod, bool) {
	mu.RLock()
	defer mu.RUnlock()
	for _, m := range registry {
		if m.ID == id {
			return m, true
		}
	}
	return Mod{}, false
}

// Default returns the primary mod (Town of Us: Mira).
func Default() Mod {
	if m, ok := Get("tou-mira"); ok {
		return m
	}
	all := All()
	if len(all) > 0 {
		return all[0]
	}
	return Mod{}
}
