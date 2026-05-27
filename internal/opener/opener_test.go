package opener

import "testing"

func TestRunnerRunsOpenCodeInProjectDirectory(t *testing.T) {
	var got Command
	runner := Runner{
		Execute: func(command Command) error {
			got = command
			return nil
		},
	}

	if err := runner.RunOpenCode("/tmp/project"); err != nil {
		t.Fatalf("run opencode: %v", err)
	}

	if got.Name != "opencode" {
		t.Fatalf("expected command opencode, got %q", got.Name)
	}
	if got.Dir != "/tmp/project" {
		t.Fatalf("expected dir /tmp/project, got %q", got.Dir)
	}
	if len(got.Args) != 0 {
		t.Fatalf("expected no args, got %#v", got.Args)
	}
}
