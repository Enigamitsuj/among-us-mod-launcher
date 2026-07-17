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

// DirIsWritable reports whether a new file can be created in dir.
// Used to detect UAC-protected game folders such as Program Files.
func DirIsWritable(dir string) bool {
	if dir == "" {
		return false
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}

	f, err := os.CreateTemp(dir, ".aumod-write-*")
	if err != nil {
		return false
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	return true
}
