package config

import (
	"path/filepath"
	"testing"
)

func TestLoadMissingFileReturnsEmptyProjects(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), "projects.json")}

	file, err := store.Load()
	if err != nil {
		t.Fatalf("load missing file: %v", err)
	}

	if file.Version != 1 {
		t.Fatalf("expected version 1, got %d", file.Version)
	}
	if len(file.Projects) != 0 {
		t.Fatalf("expected no projects, got %d", len(file.Projects))
	}
}

func TestSaveAndLoadProjects(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), ".my-power-toys", "projects.json")}
	want := File{
		Version: 1,
		Projects: []Project{
			{Name: "my-power-toys", Path: "/home/me/dev/my-power-toys"},
		},
	}

	if err := store.Save(want); err != nil {
		t.Fatalf("save projects: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("load saved projects: %v", err)
	}

	if got.Version != want.Version {
		t.Fatalf("expected version %d, got %d", want.Version, got.Version)
	}
	if len(got.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(got.Projects))
	}
	if got.Projects[0] != want.Projects[0] {
		t.Fatalf("expected project %#v, got %#v", want.Projects[0], got.Projects[0])
	}
}
