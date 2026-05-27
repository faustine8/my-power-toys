package project

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/faustine8/my-power-toys/internal/config"
)

func TestAddUsesDirectoryBaseNameByDefault(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "my-power-toys")
	file := config.File{Version: 1}

	got, added, err := Add(file, dir, "")
	if err != nil {
		t.Fatalf("add project: %v", err)
	}

	wantPath, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	if added.Name != "my-power-toys" {
		t.Fatalf("expected name my-power-toys, got %q", added.Name)
	}
	if added.Path != wantPath {
		t.Fatalf("expected path %q, got %q", wantPath, added.Path)
	}
	if len(got.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(got.Projects))
	}
}

func TestAddUsesCustomName(t *testing.T) {
	dir := t.TempDir()
	file := config.File{Version: 1}

	_, added, err := Add(file, dir, "tools")
	if err != nil {
		t.Fatalf("add project: %v", err)
	}

	if added.Name != "tools" {
		t.Fatalf("expected custom name tools, got %q", added.Name)
	}
}

func TestAddRejectsDuplicateName(t *testing.T) {
	existingPath := filepath.Join(t.TempDir(), "existing")
	newPath := filepath.Join(t.TempDir(), "new")
	file := config.File{
		Version:  1,
		Projects: []config.Project{{Name: "tools", Path: existingPath}},
	}

	if _, _, err := Add(file, newPath, "tools"); err == nil {
		t.Fatal("expected duplicate name error")
	}
}

func TestAddRejectsDuplicatePath(t *testing.T) {
	dir := t.TempDir()
	absDir, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	file := config.File{
		Version:  1,
		Projects: []config.Project{{Name: "existing", Path: absDir}},
	}

	if _, _, err := Add(file, dir, "other"); err == nil {
		t.Fatal("expected duplicate path error")
	}
}

func TestRemoveDeletesExistingProject(t *testing.T) {
	file := config.File{
		Version: 1,
		Projects: []config.Project{
			{Name: "keep", Path: "/tmp/keep"},
			{Name: "remove", Path: "/tmp/remove"},
		},
	}

	got, removed, err := Remove(file, "remove")
	if err != nil {
		t.Fatalf("remove project: %v", err)
	}

	if removed.Name != "remove" {
		t.Fatalf("expected removed project remove, got %q", removed.Name)
	}
	if len(got.Projects) != 1 {
		t.Fatalf("expected 1 remaining project, got %d", len(got.Projects))
	}
	if got.Projects[0].Name != "keep" {
		t.Fatalf("expected keep project to remain, got %q", got.Projects[0].Name)
	}
}

func TestRemoveRejectsMissingProject(t *testing.T) {
	file := config.File{Version: 1}

	if _, _, err := Remove(file, "missing"); err == nil {
		t.Fatal("expected missing project error")
	}
}

func TestSearchMatchesProjectName(t *testing.T) {
	projects := []config.Project{
		{Name: "my-power-toys", Path: "/tmp/tools"},
		{Name: "other", Path: "/tmp/other"},
	}

	got := Search(projects, "POWER")

	if len(got) != 1 {
		t.Fatalf("expected 1 match, got %d", len(got))
	}
	if got[0].Name != "my-power-toys" {
		t.Fatalf("expected my-power-toys match, got %#v", got[0])
	}
}

func TestSearchMatchesProjectPath(t *testing.T) {
	projects := []config.Project{
		{Name: "tools", Path: "/home/me/dev/my-power-toys"},
		{Name: "other", Path: "/home/me/dev/other"},
	}

	got := Search(projects, "power-toys")

	if len(got) != 1 {
		t.Fatalf("expected 1 match, got %d", len(got))
	}
	if got[0].Path != "/home/me/dev/my-power-toys" {
		t.Fatalf("expected path match, got %#v", got[0])
	}
}

func TestSearchTrimsQuery(t *testing.T) {
	projects := []config.Project{{Name: "tools", Path: "/tmp/tools"}}

	got := Search(projects, "  tools  ")

	if len(got) != 1 {
		t.Fatalf("expected trimmed query to match, got %d", len(got))
	}
}

func TestSelectReturnsOnlyProjectWithoutPrompt(t *testing.T) {
	projects := []config.Project{{Name: "only", Path: "/tmp/only"}}
	var prompt bytes.Buffer

	got, ok, err := Select(projects, strings.NewReader(""), &prompt)
	if err != nil {
		t.Fatalf("select project: %v", err)
	}

	if !ok {
		t.Fatal("expected selection")
	}
	if got != projects[0] {
		t.Fatalf("expected %#v, got %#v", projects[0], got)
	}
	if prompt.Len() != 0 {
		t.Fatalf("expected no prompt, got %q", prompt.String())
	}
}

func TestSelectCancelsOnEmptyInput(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}

	_, ok, err := Select(projects, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("select project: %v", err)
	}
	if ok {
		t.Fatal("expected cancelled selection")
	}
}

func TestSelectReturnsNumberedChoice(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}

	got, ok, err := Select(projects, strings.NewReader("2\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("select project: %v", err)
	}

	if !ok {
		t.Fatal("expected selection")
	}
	if got != projects[1] {
		t.Fatalf("expected %#v, got %#v", projects[1], got)
	}
}

func TestSelectRejectsInvalidNumber(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}

	if _, _, err := Select(projects, strings.NewReader("3\n"), &bytes.Buffer{}); err == nil {
		t.Fatal("expected invalid selection error")
	}
}
