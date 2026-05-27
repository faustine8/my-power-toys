package opener

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Name string
	Args []string
	Dir  string
}

type Executor func(Command) error

type Runner struct {
	Execute Executor
}

func DefaultRunner() Runner {
	return Runner{Execute: ExecuteCommand}
}

func (r Runner) RunOpenCode(dir string) error {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return fmt.Errorf("project path cannot be empty")
	}
	return r.execute(Command{Name: "opencode", Dir: dir})
}

func ExecuteCommand(command Command) error {
	cmd := exec.Command(command.Name, command.Args...)
	cmd.Dir = command.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r Runner) execute(command Command) error {
	execute := r.Execute
	if execute == nil {
		execute = ExecuteCommand
	}
	return execute(command)
}
