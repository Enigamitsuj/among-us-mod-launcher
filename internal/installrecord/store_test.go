package installrecord

import "testing"

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	want := Record{
		ModID:    "tou-mira",
		Path:     dir,
		Version:  "1.6.2",
		Platform: "steam",
	}

	if err := Save(dir, want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, ok := Load(dir)
	if !ok {
		t.Fatal("Load() did not find saved marker")
	}
	if got != want {
		t.Fatalf("Load() = %#v, want %#v", got, want)
	}
}

func TestLoadRejectsMissingMarker(t *testing.T) {
	if _, ok := Load(t.TempDir()); ok {
		t.Fatal("Load() found a marker that does not exist")
	}
}
