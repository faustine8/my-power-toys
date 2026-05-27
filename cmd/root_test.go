package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/faustine8/my-power-toys/internal/config"
)

func TestRootCommand(t *testing.T) {
	stdout, _, err := executeTestCommand(t, rootOptions{}, "", []string{}...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(stdout, "po") {
		t.Errorf("expected output to contain 'po', got: %s", stdout)
	}
}

func TestVersionCommand(t *testing.T) {
	stdout, _, err := executeTestCommand(t, rootOptions{}, "", "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(stdout, "po dev") {
		t.Errorf("expected output to contain 'po dev', got: %s", stdout)
	}
}

func TestAddAndListCommandsUseProjectStore(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	projectDir := filepath.Join(t.TempDir(), "repo")
	options := rootOptions{
		StorePath: storePath,
		Getwd: func() (string, error) {
			return projectDir, nil
		},
	}

	if _, _, err := executeTestCommand(t, options, "", "add", "--name", "tools"); err != nil {
		t.Fatalf("add command: %v", err)
	}

	stdout, _, err := executeTestCommand(t, options, "", "list")
	if err != nil {
		t.Fatalf("list command: %v", err)
	}

	wantPath, err := filepath.Abs(projectDir)
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	want := "tools\t" + wantPath
	if !strings.Contains(stdout, want) {
		t.Fatalf("expected list output to contain %q, got %q", want, stdout)
	}
}

func TestRemoveCommandDeletesByName(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	store := config.Store{Path: storePath}
	if err := store.Save(config.File{
		Version: 1,
		Projects: []config.Project{
			{Name: "one", Path: "/tmp/one"},
			{Name: "two", Path: "/tmp/two"},
		},
	}); err != nil {
		t.Fatalf("save fixture: %v", err)
	}

	stdout, _, err := executeTestCommand(t, rootOptions{StorePath: storePath}, "", "remove", "one")
	if err != nil {
		t.Fatalf("remove command: %v", err)
	}
	if !strings.Contains(stdout, "removed\tone") {
		t.Fatalf("expected remove output, got %q", stdout)
	}

	file, err := store.Load()
	if err != nil {
		t.Fatalf("load projects: %v", err)
	}
	if len(file.Projects) != 1 || file.Projects[0].Name != "two" {
		t.Fatalf("expected only two to remain, got %#v", file.Projects)
	}
}

func TestPickCommandPrintsOnlySelectedPathToStdout(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	store := config.Store{Path: storePath}
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}
	if err := store.Save(config.File{Version: 1, Projects: projects}); err != nil {
		t.Fatalf("save fixture: %v", err)
	}

	stdout, stderr, err := executeTestCommand(t, rootOptions{StorePath: storePath}, "\x1b[B\r", "pick")
	if err != nil {
		t.Fatalf("pick command: %v", err)
	}

	if stdout != "/tmp/two\n" {
		t.Fatalf("expected stdout to contain only selected path, got %q", stdout)
	}
	if !strings.Contains(stderr, "Select project") {
		t.Fatalf("expected prompt on stderr, got %q", stderr)
	}
	if strings.Contains(stderr, "Enter number") {
		t.Fatalf("expected arrow-key prompt, got %q", stderr)
	}
}

func TestPickCommandWithoutQueryPromptsForOnlyProject(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	store := config.Store{Path: storePath}
	if err := store.Save(config.File{
		Version: 1,
		Projects: []config.Project{
			{Name: "one", Path: "/tmp/one"},
		},
	}); err != nil {
		t.Fatalf("save fixture: %v", err)
	}

	stdout, stderr, err := executeTestCommand(t, rootOptions{StorePath: storePath}, "\r", "pick")
	if err != nil {
		t.Fatalf("pick command: %v", err)
	}
	if stdout != "/tmp/one\n" {
		t.Fatalf("expected stdout to contain only selected path, got %q", stdout)
	}
	if !strings.Contains(stderr, "Select project") {
		t.Fatalf("expected prompt on stderr, got %q", stderr)
	}
}

func TestPickCommandWithQueryPrintsOnlyMatchedPathToStdout(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	store := config.Store{Path: storePath}
	projects := []config.Project{
		{Name: "my-power-toys", Path: "/tmp/my-power-toys"},
		{Name: "other", Path: "/tmp/other"},
	}
	if err := store.Save(config.File{Version: 1, Projects: projects}); err != nil {
		t.Fatalf("save fixture: %v", err)
	}

	stdout, stderr, err := executeTestCommand(t, rootOptions{StorePath: storePath}, "", "pick", "my-power-toys")
	if err != nil {
		t.Fatalf("pick command: %v", err)
	}

	if stdout != "/tmp/my-power-toys\n" {
		t.Fatalf("expected stdout to contain only matched path, got %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("expected no stderr for single query match, got %q", stderr)
	}
}

func TestPickCommandWithQueryUsesFilteredSelectorForMultipleMatches(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	store := config.Store{Path: storePath}
	projects := []config.Project{
		{Name: "alpha-api", Path: "/tmp/alpha-api"},
		{Name: "alpha-web", Path: "/tmp/alpha-web"},
		{Name: "beta", Path: "/tmp/beta"},
	}
	if err := store.Save(config.File{Version: 1, Projects: projects}); err != nil {
		t.Fatalf("save fixture: %v", err)
	}

	stdout, stderr, err := executeTestCommand(t, rootOptions{StorePath: storePath}, "\x1b[B\r", "pick", "alpha")
	if err != nil {
		t.Fatalf("pick command: %v", err)
	}

	if stdout != "/tmp/alpha-web\n" {
		t.Fatalf("expected stdout to contain only selected filtered path, got %q", stdout)
	}
	if !strings.Contains(stderr, "Select project") {
		t.Fatalf("expected prompt on stderr, got %q", stderr)
	}
	if strings.Contains(stderr, "beta") {
		t.Fatalf("expected filtered prompt to exclude beta, got %q", stderr)
	}
}

func TestPickCommandCancelsWithoutStdout(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	store := config.Store{Path: storePath}
	if err := store.Save(config.File{
		Version: 1,
		Projects: []config.Project{
			{Name: "one", Path: "/tmp/one"},
			{Name: "two", Path: "/tmp/two"},
		},
	}); err != nil {
		t.Fatalf("save fixture: %v", err)
	}

	stdout, _, err := executeTestCommand(t, rootOptions{StorePath: storePath}, "\x1b", "pick")
	if err != nil {
		t.Fatalf("pick command: %v", err)
	}
	if stdout != "" {
		t.Fatalf("expected no stdout on cancel, got %q", stdout)
	}
}

func TestPickCommandWithQueryReturnsClearErrorWhenNoProjectsMatch(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	store := config.Store{Path: storePath}
	if err := store.Save(config.File{
		Version: 1,
		Projects: []config.Project{
			{Name: "one", Path: "/tmp/one"},
		},
	}); err != nil {
		t.Fatalf("save fixture: %v", err)
	}

	stdout, _, err := executeTestCommand(t, rootOptions{StorePath: storePath}, "", "pick", "missing")
	if err == nil {
		t.Fatal("expected missing project error")
	}
	if stdout != "" {
		t.Fatalf("expected no stdout on error, got %q", stdout)
	}
	if !strings.Contains(err.Error(), "no project matches query: missing") {
		t.Fatalf("expected clear no-match error, got %v", err)
	}
}

func TestPickCommandAcceptsAtMostOneQueryArgument(t *testing.T) {
	_, _, err := executeTestCommand(t, rootOptions{}, "", "pick", "one", "two")
	if err == nil {
		t.Fatal("expected maximum args error")
	}
}

func TestOpenCodeCommandRunsWithSelectedProjectPath(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "projects.json")
	store := config.Store{Path: storePath}
	projects := []config.Project{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}
	if err := store.Save(config.File{Version: 1, Projects: projects}); err != nil {
		t.Fatalf("save fixture: %v", err)
	}

	var gotDir string
	options := rootOptions{
		StorePath: storePath,
		RunOpenCode: func(dir string) error {
			gotDir = dir
			return nil
		},
	}
	stdout, _, err := executeTestCommand(t, options, "\r", "oc")
	if err != nil {
		t.Fatalf("oc command: %v", err)
	}

	if stdout != "" {
		t.Fatalf("expected no stdout, got %q", stdout)
	}
	if gotDir != "/tmp/one" {
		t.Fatalf("expected opencode dir /tmp/one, got %q", gotDir)
	}
}

func executeTestCommand(t *testing.T, options rootOptions, input string, args ...string) (string, string, error) {
	t.Helper()

	cmd := newRootCommand(options)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs(args)

	err := cmd.Execute()
	return stdout.String(), stderr.String(), err
}
