# my-power-toys

A tiny project opener for developers.

## Status

First phase: `po`, a short-lived CLI for Linux / WSL style shell workflows.

This phase does not include a resident TUI. The CLI stores projects, prints a selected project path for shell functions, and can start OpenCode in a selected project directory.

## Commands

```bash
# Register the current directory using the directory name
po add

# Register the current directory with an alias
po add --name my-project

# List registered projects
po list

# Remove a project by name
po remove my-project

# Pick a project interactively and print only its path
po pick

# Search for a project and print only its path
po pick my-project

# Pick a project and start opencode in that directory
po oc

# Print version
po version
```

## Shell Integration

`po pick` writes prompts to stderr and prints only the selected path to stdout, so it can be used from a zsh/bash shell function:

```sh
pcd() {
  local dir
  dir="$(po pick "$@")" || return
  [ -n "$dir" ] && cd "$dir"
}
```

For example, `pcd my-project` searches registered projects and changes into the matched directory.

A Go binary cannot change the current directory of its parent shell directly. `po pick` prints the selected path, and the shell function performs `cd "$dir"` in the current shell process.

## Storage

Projects are stored in:

```text
~/.my-power-toys/projects.json
```

Example:

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

Project names must be unique. Project paths are stored as absolute paths and are de-duplicated.

## Build

```bash
go build -o po .
./po version
```

## Test

```bash
go test ./...
```
