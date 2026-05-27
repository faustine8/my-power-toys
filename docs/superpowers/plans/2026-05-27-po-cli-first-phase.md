# po CLI First Phase Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a short-lived `po` CLI that registers, lists, removes, picks, and opens project directories with OpenCode.

**Architecture:** Cobra remains the CLI wiring layer. Configuration, project logic, and opener execution are isolated in `internal/config`, `internal/project`, and `internal/opener` so business logic does not live in command handlers.

**Tech Stack:** Go, Cobra, standard library JSON/filesystem/process APIs.

---

## File Structure

- Create `internal/config/store.go`: storage path, load, save, and JSON model.
- Create `internal/config/store_test.go`: config load/save tests.
- Create `internal/project/project.go`: project model, add/remove/list service, and text selector.
- Create `internal/project/project_test.go`: project validation and selector tests.
- Create `internal/opener/opener.go`: `opencode` process execution with selected project as working directory.
- Modify `cmd/root.go`: change root command to `po` and add `add`, `list`, `remove`, `pick`, and `oc` commands.
- Modify `cmd/root_test.go`: reset command state between tests and cover command stdout behavior.
- Modify `README.md`: document first-phase `po` usage and shell function.

## Tasks

### Task 1: Config Store

- [ ] Write failing tests in `internal/config/store_test.go` for missing-file load, save, and reload.
- [ ] Run `go test ./internal/config` and confirm it fails because the package does not exist.
- [ ] Implement `internal/config/store.go` with `Project`, `File`, `Store`, `DefaultPath`, `Load`, and `Save`.
- [ ] Run `go test ./internal/config` and confirm it passes.

### Task 2: Project Service

- [ ] Write failing tests in `internal/project/project_test.go` for add default name, add custom name, duplicate path, duplicate name, remove existing, and remove missing.
- [ ] Run `go test ./internal/project` and confirm it fails because the package does not exist.
- [ ] Implement `internal/project/project.go` service methods.
- [ ] Run `go test ./internal/project` and confirm it passes.

### Task 3: Selector

- [ ] Add failing selector tests for single-project auto-pick, empty input cancel, valid numeric choice, and invalid numeric choice.
- [ ] Run `go test ./internal/project` and confirm the new selector tests fail.
- [ ] Implement the selector in `internal/project/project.go`, writing prompts to an injected writer.
- [ ] Run `go test ./internal/project` and confirm it passes.

### Task 4: Cobra Commands

- [ ] Write failing command tests in `cmd/root_test.go` for `po pick` stdout path output and root/version text using `po`.
- [ ] Run `go test ./cmd` and confirm the command tests fail.
- [ ] Modify `cmd/root.go` to build a fresh Cobra command with injectable dependencies and add `add`, `list`, `remove`, `pick`, and `oc`.
- [ ] Run `go test ./cmd` and confirm it passes.

### Task 5: OpenCode Launcher

- [ ] Write a focused test or command injection test proving `po oc` passes the selected project path as the working directory.
- [ ] Run the relevant test and confirm it fails.
- [ ] Implement `internal/opener.RunOpenCode` and wire `po oc` to it.
- [ ] Run the relevant test and confirm it passes.

### Task 6: Docs and Verification

- [ ] Update `README.md` to show `po` usage, storage path, and shell function integration.
- [ ] Run `gofmt` on Go files.
- [ ] Run `go test ./...`.
- [ ] Review `git diff` to confirm no unrelated changes were reverted.
