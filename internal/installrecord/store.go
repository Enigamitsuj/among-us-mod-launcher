package installrecord

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const markerName = ".among-us-mod-launcher.json"

type Record struct {
	ModID    string `json:"modId"`
	Path     string `json:"path"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
}

// Load reads the ownership marker from an installed mod directory.
func Load(installPath string) (Record, bool) {
	data, err := os.ReadFile(filepath.Join(installPath, markerName))
	if err != nil {
		return Record{}, false
	}
	var record Record
	if err := json.Unmarshal(data, &record); err != nil {
		return Record{}, false
	}
	return record, true
}

// Save marks a directory as a launcher-managed mod installation.
func Save(installPath string, record Record) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(installPath, markerName), data, 0o600)
}
