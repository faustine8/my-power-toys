# my-power-toys 🧰✨

A tiny cross-platform TUI toolbox for personal developer workflows.

First toy inside the box: **a project opener**.

Because sometimes your company splits one system into 12 repositories,  
but your brain does not have to split with it. 🫠

---

## What is this? 🤔

`my-power-toys` is a personal terminal toolbox.

The first module is a fast TUI project launcher:

```bash
mpt
```

Type a few letters, pick a project, press Enter, and open it with your favorite tool:

- IntelliJ IDEA
- OpenCode
- Codex CLI
- Claude Code
- VS Code
- or any custom command you like

No Electron.  
No WebView.  
No hidden Chromium pretending to be a tiny utility. 🐘🚫

Just a small terminal tool.

---

## Why? 💡

Modern development often means jumping between many related repositories:

```text
backend-service
scheduler
common-lib
webhook-service
data-builder
docs
notes
```

Existing options are often annoying:

- IDE recent projects are tied to one IDE
- GUI launchers are heavier than they need to be
- taskbar jump lists may not work reliably
- some tools run a whole browser engine just to show a list
- switching between coding agents is becoming normal

So this project tries to be a simple, portable command center for opening workspaces.

---

## Core idea 🧠

You open a directory with `mpt`.

`mpt` remembers it.

Next time, you just search and open it again.

```text
my-power-toys / projects

Search: common_

★ iot_common_data       IOT     IDEA       D:\Dev\Code\iot\iot_common_data
★ iot_common_lib        IOT     IDEA       D:\Dev\Code\iot\iot_common_lib
  iot_scheduler         IOT     Codex      D:\Dev\Code\iot\iot_scheduler
  work_note             Notes   VS Code    D:\Dev\Code\github\work_note_2024

↑/↓    Select
Enter  Open
Tab    Change opener
e      Edit project
f      Toggle favorite
q      Quit
```

---

## Status 🚧

This project is currently in early development.

Planned MVP:

- [ ] Cross-platform Go CLI/TUI
- [ ] Project list
- [ ] Search box
- [ ] Keyboard selection
- [ ] Project opener selection
- [ ] Config file storage
- [ ] Open count and last-opened tracking
- [ ] Windows / macOS / Linux support

---

## Install 📦

Not released yet.

Eventually:

```bash
go install github.com/YOUR_USERNAME/my-power-toys/cmd/mpt@latest
```

Or download a binary from GitHub Releases.

---

## Usage 🚀

### Open the project picker

```bash
mpt
```

Equivalent to:

```bash
mpt projects
```

---

### Add current directory

```bash
mpt add .
```

Example:

```bash
cd D:\Dev\Code\iot\iot_common_data
mpt add .
```

`mpt` will remember this directory and use the folder name as the default project name:

```text
iot_common_data
```

---

### Open current directory

```bash
mpt open .
```

If this directory is not recorded yet, `mpt` will add it first.

If no default opener is configured, it will ask:

```text
Choose opener for iot_common_data:

> IntelliJ IDEA
  OpenCode
  Codex CLI
  Claude Code
  VS Code
  File Manager
```

---

### Open project by name

```bash
mpt open iot_common_data
```

If there is one clear match, it opens directly.

If there are multiple matches, it opens the TUI picker.

---

### Edit config

```bash
mpt config
```

This opens the config file in your default editor.

---

## Openers 🪄

An opener is a tool used to open a project.

Some openers receive the project path as an argument:

```json
{
  "id": "idea",
  "name": "IntelliJ IDEA",
  "command": "idea",
  "args": ["{{path}}"]
}
```

Some openers run inside the project directory:

```json
{
  "id": "codex",
  "name": "Codex CLI",
  "command": "codex",
  "args": [],
  "workingDir": "{{path}}"
}
```

This makes it possible to support both IDEs and coding agents.

---

## Example config ⚙️

`my-power-toys` stores config in the standard user config directory:

```text
Windows: %AppData%\my-power-toys\config.json
macOS:   ~/Library/Application Support/my-power-toys/config.json
Linux:   ~/.config/my-power-toys/config.json
```

Example:

```json
{
  "version": 1,
  "defaultModule": "projects",
  "projects": [
    {
      "id": "iot_common_data",
      "name": "iot_common_data",
      "alias": ["common_data", "common"],
      "path": "D:\\Dev\\Code\\iot\\iot_common_data",
      "group": "IOT",
      "favorite": true,
      "defaultOpener": "idea",
      "lastOpenedAt": "2026-05-26T10:30:00+08:00",
      "openCount": 12
    }
  ],
  "openers": [
    {
      "id": "idea",
      "name": "IntelliJ IDEA",
      "command": "idea",
      "args": ["{{path}}"]
    },
    {
      "id": "opencode",
      "name": "OpenCode",
      "command": "opencode",
      "args": ["{{path}}"]
    },
    {
      "id": "codex",
      "name": "Codex CLI",
      "command": "codex",
      "args": [],
      "workingDir": "{{path}}"
    },
    {
      "id": "claude",
      "name": "Claude Code",
      "command": "claude",
      "args": [],
      "workingDir": "{{path}}"
    },
    {
      "id": "vscode",
      "name": "VS Code",
      "command": "code",
      "args": ["{{path}}"]
    }
  ]
}
```

---

## Template variables 🧩

Openers support simple template variables:

| Variable    | Meaning       |
|-------------|---------------|
| `{{path}}`  | Project path  |
| `{{name}}`  | Project name  |
| `{{group}}` | Project group |

Example:

```json
{
  "id": "idea",
  "name": "IntelliJ IDEA",
  "command": "idea",
  "args": ["{{path}}"]
}
```

---

## Design goals 🎯

- Fast startup
- Keyboard-first workflow
- Cross-platform
- No Electron
- No WebView
- No database
- Human-readable config
- Easy to extend
- Good enough to use every day

---

## Non-goals 🙅

This project is not trying to be:

- a full IDE manager
- a JetBrains Toolbox clone
- a desktop GUI app
- a cloud-synced workspace platform
- a team collaboration tool
- another giant productivity suite

It is just a small personal toolbox.

Small tools are allowed to stay small. 🌱

---

## Planned roadmap 🗺️

### MVP

- [ ] `mpt`
- [ ] `mpt projects`
- [ ] `mpt add .`
- [ ] `mpt open .`
- [ ] `mpt open <name>`
- [ ] `mpt config`
- [ ] Project search
- [ ] Opener selection
- [ ] Last-opened tracking
- [ ] Open-count tracking

### Later

- [ ] `mpt scan <dir>` to scan Git repositories
- [ ] Favorite-only view
- [ ] Group filtering
- [ ] Open all projects in a group
- [ ] Fuzzy search scoring
- [ ] Import JetBrains recent projects
- [ ] More coding-agent templates
- [ ] Snippet launcher
- [ ] Port killer
- [ ] Env switcher
- [ ] Git branch cleaner

---

## Development 🛠️

Recommended stack:

- Go
- Bubble Tea
- Bubbles
- Lip Gloss

Suggested structure:

```text
my-power-toys/
├── cmd/
│   └── mpt/
│       └── main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── projects/
│   ├── tui/
│   └── platform/
├── go.mod
└── README.md
```

Run locally:

```bash
go run ./cmd/mpt
```

Build:

```bash
go build -o mpt ./cmd/mpt
```

On Windows:

```powershell
go build -o mpt.exe ./cmd/mpt
```

---

## Philosophy 🐢

A useful dev tool does not need to be huge.

It should:

- do one thing well
- start quickly
- stay out of the way
- be easy to understand
- be easy to modify
- make the daily workflow slightly less annoying

That is the whole point.

---

## License 📄

MIT License.

Use it, fork it, break it, fix it, and make it yours.
