package fsutil

import (
	"os"
	"path/filepath"
)

// JoinInstallPath builds "<root>/<mod folder name>".
func JoinInstallPath(root, folderName string) string {
	return filepath.Join(root, folderName)
}

// DirExists reports whether path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// EnsureDir creates a directory if missing.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}
