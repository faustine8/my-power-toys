# my-power-toys 需求文档：跨平台 TUI 项目打开器

## 1. 项目背景

当前日常开发中存在多个强关联代码仓库，需要频繁在不同项目之间切换。JetBrains Toolbox、Windows 任务栏 Jump List、Electron 类启动器等方案存在以下问题：

- JetBrains Toolbox 更像托盘工具，不适合作为稳定的项目工作台；
- Windows 任务栏 Jump List 可能被系统设置或公司策略禁用；
- Electron/WebView 类工具资源占用偏高，不符合轻量工具诉求；
- GUI 工具需要常驻窗口，而实际需求更接近“命令面板”；
- 希望后续可以持续扩展为个人开发者工具箱。

因此准备开发一个跨平台 TUI 工具箱，暂定名称为 `my-power-toys`，第一个模块是“项目打开器”。

---

## 2. 项目名称

### 2.1 仓库名

`my-power-toys`

### 2.2 命令名

建议命令名：

```bash
mpt
```

`mpt` 是 `my-power-toys` 的缩写，适合高频输入。

---

## 3. 项目定位

`my-power-toys` 是一个跨平台的个人开发者 TUI 工具箱。

第一阶段实现 `projects` 模块，用于：

- 记录打开过的代码目录；
- 为目录生成默认项目名；
- 支持用户自定义项目名、别名、分组；
- 在 TUI 输入框中模糊搜索项目；
- 支持上下键选择候选项目；
- 选中项目后使用默认打开方式打开；
- 第一次打开项目时选择默认打开方式；
- 支持 IntelliJ IDEA、OpenCode、Codex CLI、Claude Code、VS Code 等 opener；
- 后续可以持续扩展更多个人小工具。

---

## 4. 技术选型

### 4.1 语言

使用 Go。

原因：

- 跨平台支持 Windows、macOS、Linux；
- 单二进制发布，部署简单；
- 启动速度快；
- 适合开发 CLI/TUI 工具；
- 不依赖 Electron、WebView、Python 运行时；
- 适合后续扩展为个人工具箱。

### 4.2 TUI 框架

建议使用：

```text
Bubble Tea
Lip Gloss
Bubbles
```

用途：

- Bubble Tea：TUI 主框架；
- Lip Gloss：样式；
- Bubbles：输入框、列表等基础组件。

### 4.3 配置存储

使用 JSON 文件作为第一版配置存储。

不使用数据库。

配置文件路径使用 Go 的 `os.UserConfigDir()` 获取。

不同系统下大致路径：

```text
Windows: %AppData%\my-power-toys\config.json
macOS:   ~/Library/Application Support/my-power-toys/config.json
Linux:   ~/.config/my-power-toys/config.json
```

---

## 5. 第一阶段目标

第一阶段只实现 `projects` 模块，也就是“项目打开器”。

核心目标：

```text
启动 mpt -> 进入项目搜索 TUI -> 输入关键词 -> 上下选择项目 -> 回车打开
```

如果项目还没有默认打开方式，则第一次打开时要求用户选择 opener，并保存为该项目的默认 opener。

---

## 6. MVP 功能范围

### 6.1 必须实现

1. `mpt` 启动 TUI 项目选择器；
2. 从配置文件读取项目列表和 opener 列表；
3. 显示项目列表；
4. 支持输入关键词实时过滤项目；
5. 支持上下键选择项目；
6. 支持 Enter 打开选中项目；
7. 支持项目默认 opener；
8. 如果项目没有默认 opener，第一次打开时弹出 opener 选择界面；
9. 打开成功后更新：
    - `lastOpenedAt`
    - `openCount`
10. 支持 `mpt add .` 将当前目录加入项目列表；
11. 自动根据目录名生成默认项目名；
12. 支持自定义项目名；
13. 支持 alias 别名；
14. 支持 group 分组；
15. 支持跨 Windows、macOS、Linux 运行。

### 6.2 暂不实现

以下功能可以后续迭代，不进入 MVP：

- 插件系统；
- GUI；
- 云同步；
- 数据库；
- 读取 JetBrains 历史项目；
- 自动识别所有 IDE；
- 复杂命令面板；
- 项目图标；
- 多窗口常驻；
- WebView/Electron 界面；
- 项目标签系统；
- 复杂权限管理；
- 远程 SSH 项目打开；
- 与 GitHub/GitLab API 集成。

---

## 7. 命令设计

### 7.1 打开默认 TUI

```bash
mpt
```

行为：

- 打开默认模块；
- 第一版默认模块为 `projects`；
- 展示项目搜索列表。

等价于：

```bash
mpt projects
```

---

### 7.2 打开项目模块

```bash
mpt projects
```

行为：

- 打开项目选择 TUI；
- 读取配置文件；
- 显示项目列表；
- 支持搜索、选择、打开。

---

### 7.3 添加当前目录

```bash
mpt add .
```

行为：

- 获取当前目录绝对路径；
- 如果该路径已存在于配置中，则提示已存在；
- 如果不存在，则添加为新项目；
- 默认项目名为当前目录名；
- `defaultOpener` 为空；
- 第一次打开时再选择 opener。

示例：

当前目录：

```text
D:\Dev\Code\ciot\ciot_common_data
```

执行：

```bash
mpt add .
```

自动生成项目：

```json
{
  "name": "ciot_common_data",
  "path": "D:\\Dev\\Code\\ciot\\ciot_common_data"
}
```

---

### 7.4 打开当前目录

```bash
mpt open .
```

行为：

- 如果当前目录已记录，则使用该项目默认 opener 打开；
- 如果当前目录未记录，则自动添加；
- 如果没有默认 opener，则要求用户选择 opener；
- 打开后更新 `lastOpenedAt` 和 `openCount`。

---

### 7.5 通过名称打开项目

```bash
mpt open ciot_common_data
```

行为：

- 根据项目 `name`、`alias`、`path` 做匹配；
- 如果唯一匹配，则直接打开；
- 如果有多个匹配，则进入 TUI 候选列表；
- 如果没有匹配，则提示找不到项目。

---

### 7.6 打开配置文件

```bash
mpt config
```

行为：

- 使用系统默认编辑器，或用户配置的编辑器打开配置文件；
- 如果配置文件不存在，则创建默认配置。

---

### 7.7 扫描目录

第二版再实现。

```bash
mpt scan D:\Dev\Code
```

预期行为：

- 扫描指定目录下的 Git 仓库；
- 将包含 `.git` 的目录识别为项目；
- 批量加入项目列表；
- 已存在路径不重复添加。

---

## 8. TUI 交互设计

### 8.1 项目选择界面

示例：

```text
my-power-toys / projects

Search: common_

★ ciot_common_data       CIOT     IDEA       D:\Dev\Code\ciot\ciot_common_data
★ ciot_common_lib        CIOT     IDEA       D:\Dev\Code\ciot\ciot_common_lib
  ciot_scheduler         CIOT     Codex      D:\Dev\Code\ciot\ciot_scheduler
  work_note_2024         Notes    VS Code    D:\Dev\Code\github\work_note_2024

↑/↓    Select
Enter  Open
Tab    Change opener
e      Edit project
f      Toggle favorite
a      Add current directory
q      Quit
```

### 8.2 搜索行为

搜索字段应支持以下匹配来源：

- 项目名 `name`
- 项目别名 `alias`
- 分组 `group`
- 路径 `path`

第一版可以使用简单匹配：

- 忽略大小写；
- `contains` 匹配；
- favorite 项目优先；
- 最近打开项目优先；
- 打开次数多的项目优先。

后续可以升级为 fuzzy score。

排序建议：

1. favorite 优先；
2. 最近打开优先；
3. openCount 高的优先；
4. name 字母顺序。

---

## 9. 第一次选择 opener 的交互

如果项目没有 `defaultOpener`，用户按 Enter 打开时，展示 opener 选择界面。

示例：

```text
No default opener for ciot_common_data.

Choose opener:

> IntelliJ IDEA
  OpenCode
  Codex CLI
  Claude Code
  VS Code
  File Manager

Enter  Select
Esc    Cancel
```

选择后：

- 将 opener id 保存到该项目的 `defaultOpener`；
- 立即用该 opener 打开项目；
- 更新项目的 `lastOpenedAt` 和 `openCount`。

---

## 10. 修改 opener 的交互

在项目选择界面选中项目后，按 `Tab` 可以修改项目默认 opener。

示例：

```text
Change opener for ciot_common_data:

Current: IntelliJ IDEA

> IntelliJ IDEA
  OpenCode
  Codex CLI
  Claude Code
  VS Code
  File Manager
```

选择后更新配置文件。

---

## 11. 编辑项目的交互

选中项目后按 `e` 进入编辑界面。

第一版至少支持编辑：

- `name`
- `alias`
- `group`
- `defaultOpener`

示例：

```text
Edit project

Name:           ciot_common_data
Alias:          common_data, common
Group:          CIOT
Path:           D:\Dev\Code\ciot\ciot_common_data
Default opener: IntelliJ IDEA

Enter Save
Esc   Cancel
```

`path` 第一版可以只读，避免误改导致项目无法打开。

---

## 12. 数据模型

### 12.1 配置文件结构

配置文件：`config.json`

示例：

```json
{
  "version": 1,
  "defaultModule": "projects",
  "projects": [
    {
      "id": "ciot_common_data",
      "name": "ciot_common_data",
      "alias": ["common_data", "common"],
      "path": "D:\\Dev\\Code\\ciot\\ciot_common_data",
      "group": "CIOT",
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
      "command": "idea64.exe",
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

### 12.2 Project 字段说明

| 字段 | 类型 | 必填 | 说明 |
|---|---:|---:|---|
| `id` | string | 是 | 项目唯一 ID，可以由 name/path 生成 |
| `name` | string | 是 | 项目显示名 |
| `alias` | string[] | 否 | 项目别名，用于搜索 |
| `path` | string | 是 | 项目本地路径 |
| `group` | string | 否 | 项目分组 |
| `favorite` | bool | 否 | 是否收藏 |
| `defaultOpener` | string | 否 | 默认打开方式，对应 opener id |
| `lastOpenedAt` | string | 否 | 最近打开时间 |
| `openCount` | number | 否 | 打开次数 |

---

### 12.3 Opener 字段说明

| 字段 | 类型 | 必填 | 说明 |
|---|---:|---:|---|
| `id` | string | 是 | opener 唯一 ID |
| `name` | string | 是 | opener 显示名 |
| `command` | string | 是 | 执行命令 |
| `args` | string[] | 否 | 命令参数 |
| `workingDir` | string | 否 | 执行命令时的工作目录 |

---

## 13. Opener 类型

opener 分为两类。

### 13.1 path 参数型 opener

将项目路径作为参数传给命令。

例如：

```json
{
  "id": "idea",
  "name": "IntelliJ IDEA",
  "command": "idea64.exe",
  "args": ["{{path}}"]
}
```

适用于：

- IntelliJ IDEA
- VS Code
- OpenCode
- Cursor
- 文件管理器

---

### 13.2 working directory 型 opener

先进入项目目录，再执行命令。

例如：

```json
{
  "id": "codex",
  "name": "Codex CLI",
  "command": "codex",
  "args": [],
  "workingDir": "{{path}}"
}
```

适用于：

- Codex CLI
- Claude Code
- 某些需要在项目根目录启动的 CLI Agent

---

## 14. 模板变量

opener 的 `args` 和 `workingDir` 支持模板变量。

第一版至少支持：

| 变量 | 说明 |
|---|---|
| `{{path}}` | 项目路径 |
| `{{name}}` | 项目名称 |
| `{{group}}` | 项目分组 |

示例：

```json
{
  "id": "idea",
  "name": "IntelliJ IDEA",
  "command": "idea64.exe",
  "args": ["{{path}}"]
}
```

---

## 15. 跨平台要求

### 15.1 Windows

需要支持：

- PowerShell；
- Windows Terminal；
- `.exe` 单文件运行；
- 路径格式如 `D:\Dev\Code\xxx`；
- 使用 `idea64.exe`、`code`、`codex`、`claude` 等命令。

### 15.2 macOS

需要支持：

- zsh/bash；
- 路径格式如 `/Users/name/Dev/xxx`；
- 使用 `idea`、`code`、`codex`、`claude` 等命令；
- 配置文件放在 `~/Library/Application Support/my-power-toys/config.json`。

### 15.3 Linux

需要支持：

- bash/zsh；
- 路径格式如 `/home/name/dev/xxx`；
- 使用 `idea`、`code`、`codex`、`claude` 等命令；
- 配置文件放在 `~/.config/my-power-toys/config.json`。

---

## 16. 默认配置初始化

当第一次运行 `mpt`，如果配置文件不存在，应自动创建默认配置。

默认配置至少包含：

```json
{
  "version": 1,
  "defaultModule": "projects",
  "projects": [],
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

注意：

- Windows 用户可能需要手动把 `idea` 改成 `idea64.exe`；
- 或者后续版本自动检测常见命令是否存在。

---

## 17. 错误处理

### 17.1 项目路径不存在

如果项目 path 不存在：

- 在列表中标记为 invalid；
- 打开时提示：

```text
Project path does not exist:
D:\Dev\Code\ciot\ciot_common_data
```

### 17.2 opener 命令不存在

如果 opener command 不存在：

- 打开失败；
- 提示用户检查 command 是否在 PATH 中；
- 不更新 `lastOpenedAt` 和 `openCount`。

示例：

```text
Failed to open project.

Opener command not found:
idea64.exe

Please check your PATH or edit config:
%AppData%\my-power-toys\config.json
```

### 17.3 配置文件 JSON 格式错误

如果配置文件无法解析：

- 不要覆盖原配置；
- 显示错误位置；
- 提示用户运行 `mpt config` 手动修复。

---

## 18. 打开项目后的行为

当用户选择项目并成功打开后：

1. 启动外部命令；
2. 更新该项目：
    - `lastOpenedAt = now`
    - `openCount = openCount + 1`
3. 保存配置文件；
4. TUI 可以自动退出。

第一版建议打开后自动退出，因为这个工具定位是“瞬时命令面板”，不是常驻软件。

后续可以增加设置：

```json
{
  "exitAfterOpen": true
}
```

---

## 19. 项目结构建议

建议 Go 项目结构：

```text
my-power-toys/
├── cmd/
│   └── mpt/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   ├── config.go
│   │   └── store.go
│   ├── projects/
│   │   ├── project.go
│   │   ├── service.go
│   │   └── opener.go
│   ├── tui/
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   └── platform/
│       ├── path.go
│       └── command.go
├── go.mod
├── go.sum
└── README.md
```

---

## 20. 模块边界

### 20.1 config 模块

职责：

- 获取配置文件路径；
- 初始化默认配置；
- 读取配置；
- 保存配置；
- 处理配置版本。

### 20.2 projects 模块

职责：

- 项目增删改查；
- 项目搜索；
- 项目排序；
- 项目打开；
- opener 选择；
- 更新打开记录。

### 20.3 tui 模块

职责：

- TUI 状态管理；
- 输入框；
- 列表；
- 键盘事件；
- 页面切换；
- 渲染。

### 20.4 platform 模块

职责：

- 跨平台路径处理；
- 命令查找；
- 启动外部进程；
- 打开系统默认编辑器；
- 获取用户配置目录。

---

## 21. 验收标准

### 21.1 基础验收

- 在 Windows、macOS、Linux 至少可以编译；
- 运行 `mpt` 可以进入 TUI；
- 首次运行可以自动创建配置文件；
- 可以通过 `mpt add .` 添加当前目录；
- 添加后重新运行 `mpt` 可以看到该项目；
- 输入关键词可以过滤项目；
- 上下键可以选择项目；
- Enter 可以打开项目；
- 第一次打开时可以选择 opener；
- 选择 opener 后可以保存为默认 opener；
- 后续再次打开同项目不再询问 opener；
- 成功打开后更新 `lastOpenedAt` 和 `openCount`。

### 21.2 错误验收

- 项目路径不存在时有明确提示；
- opener 命令不存在时有明确提示；
- 配置 JSON 错误时不会覆盖原文件；
- 空项目列表时界面不崩溃；
- 用户按 Esc/q 可以退出。

---

## 22. 后续迭代方向

### 22.1 第二版

- `mpt scan <dir>` 扫描 Git 仓库；
- 支持删除项目；
- 支持修改项目 path；
- 支持 favorite 过滤；
- 支持 group 过滤；
- 支持直接打开 group 下所有项目；
- 支持 fuzzy score；
- 支持读取 JetBrains recent projects；
- 支持最近打开项目视图。

### 22.2 第三版

- 插件式模块；
- snippets 小工具；
- port killer 小工具；
- env switcher 小工具；
- git branch cleaner 小工具；
- workspace launcher；
- agent launcher；
- 支持 OpenCode / Codex / Claude Code 的参数模板。

---

## 23. 第一版开发建议

第一版开发顺序建议：

1. 初始化 Go 项目；
2. 实现 config 读取和默认配置创建；
3. 实现 project 数据模型；
4. 实现 `mpt add .`；
5. 实现 opener 数据模型；
6. 实现项目打开逻辑；
7. 实现 Bubble Tea 项目列表界面；
8. 实现搜索过滤；
9. 实现上下键选择；
10. 实现 Enter 打开；
11. 实现无默认 opener 时的 opener 选择界面；
12. 实现打开后更新统计字段；
13. 补充 README；
14. 在 Windows/macOS/Linux 分别测试。

---

## 24. 核心原则

1. 优先轻量；
2. 不使用 Electron；
3. 不使用 WebView；
4. 不依赖 GUI；
5. 不引入数据库；
6. 配置文件可读、可手工修改；
7. 第一版优先可用，不追求功能完整；
8. 后续按模块逐步扩展为个人工具箱。

---

## 25. 一句话总结

`my-power-toys` 第一阶段要做的是一个跨平台 TUI 项目打开器：记录代码目录，支持模糊搜索、上下选择、默认 opener，并能用 IDEA、OpenCode、Codex CLI、Claude Code 等工具快速打开项目。
