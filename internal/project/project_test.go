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

func TestSelectPromptsForOnlyProject(t *testing.T) {
	projects := []config.Project{{Name: "only", Path: "/tmp/only"}}
	var prompt bytes.Buffer

	got, ok, err := Select(projects, strings.NewReader("\r"), &prompt)
	if err != nil {
		t.Fatalf("select project: %v", err)
	}

	if !ok {
		t.Fatal("expected selection")
	}
	if got != projects[0] {
		t.Fatalf("expected %#v, got %#v", projects[0], got)
	}
	if !strings.Contains(prompt.String(), "Select project") {
		t.Fatalf("expected prompt, got %q", prompt.String())
	}
}

func TestSelectReturnsDefaultChoiceOnEnter(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}

	got, ok, err := Select(projects, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("select project: %v", err)
	}
	if !ok {
		t.Fatal("expected selection")
	}
	if got != projects[0] {
		t.Fatalf("expected %#v, got %#v", projects[0], got)
	}
}

func TestSelectReturnsChoiceWithArrowKeys(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}

	got, ok, err := Select(projects, strings.NewReader("\x1b[B\r"), &bytes.Buffer{})
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

func TestSelectMovesSelectionUpAndDown(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
		{Name: "three", Path: "/tmp/three"},
	}

	got, ok, err := Select(projects, strings.NewReader("\x1b[B\x1b[B\x1b[A\r"), &bytes.Buffer{})
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

func TestSelectCancelsOnEscape(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}

	_, ok, err := Select(projects, strings.NewReader("\x1b"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("select project: %v", err)
	}
	if ok {
		t.Fatal("expected cancelled selection")
	}
}

func TestSelectCancelsOnCtrlC(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}

	_, ok, err := Select(projects, strings.NewReader("\x03"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("select project: %v", err)
	}
	if ok {
		t.Fatal("expected cancelled selection")
	}
}

func TestRenderSelectionRerenderReturnsToLineStartAndClearsEveryLine(t *testing.T) {
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}

	var prompt bytes.Buffer
	lines := renderSelection(&prompt, projects, 1, 4)

	if lines != 4 {
		t.Fatalf("expected stable rendered line count 4, got %d", lines)
	}

	want := "\r\x1b[4A" +
		"\r\x1b[2KSelect project:\r\n" +
		"\r\x1b[2K  one  /tmp/one\r\n" +
		"\r\x1b[2K> two  /tmp/two\r\n" +
		"\r\x1b[2KUse Up/Down to move, Enter to select, Esc/Ctrl+C to cancel.\r\n"
	if prompt.String() != want {
		t.Fatalf("unexpected rerender output:\nwant %q\ngot  %q", want, prompt.String())
	}
}

func TestRenderProjectRowsAlignsPaths(t *testing.T) {
	projects := []config.Project{
		{Name: "a", Path: "/tmp/a"},
		{Name: "long-project-name", Path: "/tmp/long"},
		{Name: "mid", Path: "/tmp/mid"},
	}

	rows := renderProjectRows(projects, 1)

	if len(rows) != len(projects) {
		t.Fatalf("expected %d rows, got %d", len(projects), len(rows))
	}
	wantPathColumn := strings.Index(rows[0], projects[0].Path)
	if wantPathColumn < 0 {
		t.Fatalf("expected first row to contain path, got %q", rows[0])
	}
	for i, row := range rows {
		gotPathColumn := strings.Index(row, projects[i].Path)
		if gotPathColumn != wantPathColumn {
			t.Fatalf("expected path column %d for row %q, got %d", wantPathColumn, row, gotPathColumn)
		}
	}
	if !strings.HasPrefix(rows[1], "> long-project-name") {
		t.Fatalf("expected selected row to keep cursor, got %q", rows[1])
	}
}

func TestFitTerminalLineTruncatesLongLines(t *testing.T) {
	got := fitTerminalLine("> project-name\t/some/very/long/path", 16)
	want := "> project-nam..."

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	if len(got) > 16 {
		t.Fatalf("expected truncated line to fit width, got length %d", len(got))
	}
}
