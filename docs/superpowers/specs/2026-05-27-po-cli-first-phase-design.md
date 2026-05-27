# po CLI First Phase Design

## Goal

Build the first phase of `my-power-toys` as a short-lived Linux/WSL-friendly CLI named `po`. It manages project directories and either prints a selected path for shell functions or starts `opencode` in a selected project directory.

## Scope

This phase implements:

- `po add`
- `po add --name <alias>`
- `po list`
- `po remove <name>`
- `po pick`
- `po oc`

This phase explicitly does not implement a resident TUI, Bubble Tea interface, opener configuration, project groups, favorites, fuzzy scoring, or scan/import flows.

## Architecture

The Cobra command layer only parses arguments, flags, and wires stdout/stderr/stdin to internal packages.

- `internal/config` owns the storage file path, JSON load, and JSON save behavior.
- `internal/project` owns project records, uniqueness rules, add/list/remove behavior, and text selection.
- `internal/opener` owns launching external tools with the selected project directory as working directory.
- `cmd` owns the `po` command tree and delegates behavior to internal packages.

Bubble Tea remains disallowed outside `internal/tui/` and is not used in this phase.

## Storage

Project data is stored at:

```text
~/.my-power-toys/projects.json
```

The JSON shape is intentionally small:

```json
{
  "version": 1,
  "projects": [
    {
      "name": "my-power-toys",
      "path": "/home/me/dev/my-power-toys"
    }
  ]
}
```

The config directory is created on first save. A missing config file loads as an empty project list. Invalid JSON returns an error and is not overwritten.

## Project Rules

Each project has:

- `name`: unique project alias used by commands.
- `path`: absolute cleaned directory path.

`po add` uses the current directory as the project path. Without `--name`, the default name is the base directory name. With `--name`, the provided value is trimmed and used as the project name.

Validation:

- name cannot be empty.
- path cannot be empty.
- path is converted to an absolute path.
- duplicate paths are rejected.
- duplicate names are rejected.

## Commands

`po add`

Adds the current working directory using its base directory name.

`po add --name <alias>`

Adds the current working directory using `<alias>` as the project name.

`po list`

Prints registered projects as one project per line:

```text
name<TAB>path
```

`po remove <name>`

Deletes the project with the exact name. Missing names return a clear error.

`po pick`

Loads projects, prompts the user to choose one, and prints only the chosen path to stdout. Cancelled or empty selection prints nothing. This keeps it suitable for shell functions such as:

```sh
p() {
  local dir
  dir="$(po pick)"
  if [ -n "$dir" ]; then
    cd "$dir"
  fi
}
```

Selection UI writes prompts to stderr so stdout remains machine-readable.

`po oc`

Uses the same selection flow as `po pick`, then starts `opencode` with the selected project path as the process working directory. It does not print the selected path for shell consumption.

## Selection Behavior

The first version uses a simple numbered text selector:

```text
Select project:
1) my-power-toys    /home/me/dev/my-power-toys
2) notes            /home/me/dev/notes
Enter number:
```

The selector reads from stdin. Empty input means cancel. Invalid numbers return a clear error. If there is exactly one project, the selector returns that project without prompting.

## Error Handling

Errors are returned to Cobra and printed through Cobra's normal error path:

- duplicate project name
- duplicate project path
- unknown project for remove
- missing config cannot happen; it loads as empty
- invalid JSON
- invalid selection input
- `opencode` start failure

## Testing

Tests cover internal behavior instead of only command wiring:

- config load for missing files and saved files
- project add default name
- project add custom name
- duplicate name rejection
- duplicate path rejection
- remove existing project
- remove missing project
- selector stdout cleanliness through command tests
- `po pick` prints only the path on stdout

The required final verification is:

```bash
gofmt
go test ./...
```
