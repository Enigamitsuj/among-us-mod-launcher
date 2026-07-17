package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirIsWritableTempDir(t *testing.T) {
	dir := t.TempDir()
	if !DirIsWritable(dir) {
		t.Fatal("expected temp dir to be writable")
	}
}

func TestDirIsWritableMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	if DirIsWritable(missing) {
		t.Fatal("expected missing dir to be unwritable")
	}
}

func TestDirIsWritableFilePath(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if DirIsWritable(file) {
		t.Fatal("expected file path to be unwritable as a directory")
	}
}
