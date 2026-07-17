package installer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPromoteStagingFreshInstall(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "Among Us - TOU Mira")
	staging := dest + stagingSuffix
	backup := dest + backupSuffix

	if err := os.MkdirAll(staging, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staging, "Among Us.exe"), []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := promoteStaging(staging, dest, backup); err != nil {
		t.Fatalf("promoteStaging() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "Among Us.exe")); err != nil {
		t.Fatalf("dest missing Among Us.exe: %v", err)
	}
	if _, err := os.Stat(staging); !os.IsNotExist(err) {
		t.Fatal("staging should be gone")
	}
}

func TestPromoteStagingKeepsOldOnSwapFailure(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "Among Us - TOU Mira")
	staging := dest + stagingSuffix
	backup := dest + backupSuffix

	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dest, "Among Us.exe"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Leave staging missing to force rename failure after backup succeeds.
	if err := promoteStaging(staging, dest, backup); err == nil {
		t.Fatal("expected promoteStaging to fail")
	}
	data, err := os.ReadFile(filepath.Join(dest, "Among Us.exe"))
	if err != nil {
		t.Fatalf("old install should be restored: %v", err)
	}
	if string(data) != "old" {
		t.Fatalf("got %q, want old", data)
	}
}

func TestPromoteStagingReplacesExisting(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "Among Us - TOU Mira")
	staging := dest + stagingSuffix
	backup := dest + backupSuffix

	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dest, "Among Us.exe"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(staging, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staging, "Among Us.exe"), []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := promoteStaging(staging, dest, backup); err != nil {
		t.Fatalf("promoteStaging() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "Among Us.exe"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new" {
		t.Fatalf("got %q, want new", data)
	}
	if _, err := os.Stat(backup); !os.IsNotExist(err) {
		t.Fatal("backup should be cleaned up")
	}
}
