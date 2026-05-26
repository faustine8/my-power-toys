# AGENTS.md

Rules for AI agents (Codex, OpenCode, Claude Code, etc.) working on this project.

## Language and Framework

- Use Go.
- Use Cobra as the CLI entry point framework.
- Bubble Tea is only allowed in `internal/tui/`. Do not import it elsewhere.

## Architecture

- Business logic must NOT live in Cobra command handlers or Bubble Tea update functions.
- Core logic belongs in the `internal/` packages:
  - `internal/project/` - project management logic
  - `internal/opener/` - opener execution logic
  - `internal/config/` - configuration loading and saving
  - `internal/tui/` - Bubble Tea TUI (when implemented)
- The `cmd/` package is only for CLI wiring (parsing flags, calling internal packages).

## Code Quality

- Run `gofmt` after every change.
- Run `go test ./...` after every task is completed. All tests must pass.
- Do not introduce unnecessary dependencies.

## Scope Discipline

- Do not implement features outside the scope of the current task.
- If you notice something that should be done later, leave a TODO comment or note it in your output, but do not implement it.

## Output Requirements

- When completing a task, state:
  - Which files were added or modified
  - What tests were run
  - Whether tests passed or failed
