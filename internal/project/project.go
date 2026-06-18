package project

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/faustine8/my-power-toys/internal/config"
	"golang.org/x/term"
)

func Add(file config.File, dir string, name string) (config.File, config.Project, error) {
	path, err := filepath.Abs(strings.TrimSpace(dir))
	if err != nil {
		return config.File{}, config.Project{}, err
	}
	path = filepath.Clean(path)

	projectName := strings.TrimSpace(name)
	if projectName == "" {
		projectName = filepath.Base(path)
	}
	if projectName == "" || projectName == "." || projectName == string(filepath.Separator) {
		return config.File{}, config.Project{}, fmt.Errorf("project name cannot be empty")
	}

	for _, existing := range file.Projects {
		if existing.Name == projectName {
			return config.File{}, config.Project{}, fmt.Errorf("project name already exists: %s", projectName)
		}
		if samePath(existing.Path, path) {
			return config.File{}, config.Project{}, fmt.Errorf("project path already exists: %s", path)
		}
	}

	added := config.Project{Name: projectName, Path: path}
	file.Projects = append(file.Projects, added)
	return file, added, nil
}

func Remove(file config.File, name string) (config.File, config.Project, error) {
	name = strings.TrimSpace(name)
	for i, existing := range file.Projects {
		if existing.Name == name {
			file.Projects = append(file.Projects[:i], file.Projects[i+1:]...)
			return file, existing, nil
		}
	}
	return config.File{}, config.Project{}, fmt.Errorf("project not found: %s", name)
}

func List(file config.File) []config.Project {
	projects := append([]config.Project(nil), file.Projects...)
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
	return projects
}

func Search(projects []config.Project, query string) []config.Project {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return append([]config.Project(nil), projects...)
	}

	nameMatches := make([]config.Project, 0, len(projects))
	pathMatches := make([]config.Project, 0, len(projects))
	for _, project := range projects {
		if strings.Contains(strings.ToLower(project.Name), query) {
			nameMatches = append(nameMatches, project)
			continue
		}
		if strings.Contains(strings.ToLower(project.Path), query) {
			pathMatches = append(pathMatches, project)
		}
	}
	return append(nameMatches, pathMatches...)
}

func Select(projects []config.Project, input io.Reader, prompt io.Writer) (config.Project, bool, error) {
	if len(projects) == 0 {
		return config.Project{}, false, nil
	}
	if prompt == nil {
		prompt = io.Discard
	}

	restore, err := makeInputRaw(input)
	if err != nil {
		return config.Project{}, false, err
	}
	defer restore()

	fmt.Fprint(prompt, "\x1b[?25l")
	defer fmt.Fprint(prompt, "\x1b[?25h")

	reader := bufio.NewReader(input)
	query := ""
	filtered := Search(projects, query)
	selected := 0
	renderedLines := renderSelection(prompt, query, filtered, selected, 0)

	for {
		key, err := readSelectionKey(reader)
		if err != nil {
			if err == io.EOF {
				return config.Project{}, false, nil
			}
			return config.Project{}, false, err
		}

		switch key {
		case selectionKeyUp:
			if n := len(filtered); n > 0 {
				selected = (selected - 1 + n) % n
				renderedLines = renderSelection(prompt, query, filtered, selected, renderedLines)
			}
		case selectionKeyDown:
			if n := len(filtered); n > 0 {
				selected = (selected + 1) % n
				renderedLines = renderSelection(prompt, query, filtered, selected, renderedLines)
			}
		case selectionKeyEnter:
			if len(filtered) > 0 {
				return filtered[selected], true, nil
			}
		case selectionKeyCancel:
			return config.Project{}, false, nil
		case selectionKeyBackspace:
			if query != "" {
				runes := []rune(query)
				query = string(runes[:len(runes)-1])
				filtered = Search(projects, query)
				selected = clampSelection(selected, len(filtered))
				renderedLines = renderSelection(prompt, query, filtered, selected, renderedLines)
			}
		case selectionKeyClearQuery:
			if query != "" {
				query = ""
				filtered = Search(projects, query)
				selected = 0
				renderedLines = renderSelection(prompt, query, filtered, selected, renderedLines)
			}
		default:
			if key.isPrintable() {
				query += string(rune(key))
				filtered = Search(projects, query)
				selected = clampSelection(selected, len(filtered))
				renderedLines = renderSelection(prompt, query, filtered, selected, renderedLines)
			}
		}
	}
}

type selectionKey int

const (
	selectionKeyUnknown selectionKey = iota
	selectionKeyUp
	selectionKeyDown
	selectionKeyEnter
	selectionKeyCancel
	selectionKeyBackspace
	selectionKeyClearQuery
)

func (key selectionKey) isPrintable() bool {
	return key >= 0x20 && key != 0x7f
}

func clampSelection(selected int, projectCount int) int {
	if projectCount == 0 || selected < 0 {
		return 0
	}
	if selected >= projectCount {
		return projectCount - 1
	}
	return selected
}

type fileDescriptor interface {
	Fd() uintptr
}

func makeInputRaw(input io.Reader) (func(), error) {
	file, ok := input.(fileDescriptor)
	if !ok {
		return func() {}, nil
	}

	fd := int(file.Fd())
	if !term.IsTerminal(fd) {
		return func() {}, nil
	}

	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	return func() {
		_ = term.Restore(fd, state)
	}, nil
}

func renderSelection(prompt io.Writer, query string, projects []config.Project, selected int, previousLines int) int {
	if previousLines > 0 {
		fmt.Fprint(prompt, "\r")
		fmt.Fprintf(prompt, "\x1b[%dA", previousLines)
	}

	width := terminalWidth(prompt)
	lines := 0
	renderTerminalLine(prompt, "Select project:", width)
	lines++
	renderTerminalLine(prompt, "Filter: "+query+"\u2588", width)
	lines++
	if len(projects) == 0 {
		renderTerminalLine(prompt, "No matched projects", width)
		lines++
	} else {
		for _, line := range renderProjectRows(projects, selected) {
			renderTerminalLine(prompt, line, width)
			lines++
		}
	}
	renderTerminalLine(prompt, "Type to filter, Ctrl+U to clear, Up/Down to move, Enter to select, Esc/Ctrl+C to cancel.", width)
	lines++
	for lines < previousLines {
		renderTerminalLine(prompt, "", width)
		lines++
	}
	return lines
}

func renderProjectRows(projects []config.Project, selected int) []string {
	var buffer bytes.Buffer
	writer := tabwriter.NewWriter(&buffer, 0, 0, 2, ' ', 0)
	for i, project := range projects {
		cursor := " "
		if i == selected {
			cursor = ">"
		}
		fmt.Fprintf(writer, "%s %s\t%s\n", cursor, project.Name, project.Path)
	}
	_ = writer.Flush()

	output := strings.TrimSuffix(buffer.String(), "\n")
	if output == "" {
		return nil
	}
	return strings.Split(output, "\n")
}

func terminalWidth(output io.Writer) int {
	file, ok := output.(fileDescriptor)
	if !ok {
		return 0
	}

	fd := int(file.Fd())
	if !term.IsTerminal(fd) {
		return 0
	}

	width, _, err := term.GetSize(fd)
	if err != nil || width <= 0 {
		return 0
	}
	return width
}

func renderTerminalLine(output io.Writer, line string, width int) {
	if width > 1 {
		line = fitTerminalLine(line, width-1)
	}
	fmt.Fprintf(output, "\r\x1b[2K%s\r\n", line)
}

func fitTerminalLine(line string, width int) string {
	if width <= 0 {
		return line
	}

	runes := []rune(line)
	if len(runes) <= width {
		return line
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
}

func readSelectionKey(reader *bufio.Reader) (selectionKey, error) {
	key, err := reader.ReadByte()
	if err != nil {
		return selectionKeyUnknown, err
	}

	switch key {
	case '\r', '\n':
		return selectionKeyEnter, nil
	case 0x03:
		return selectionKeyCancel, nil
	case 0x15:
		return selectionKeyClearQuery, nil
	case 0x1b:
		return readEscapeKey(reader)
	case 0x00, 0xe0:
		return readWindowsConsoleKey(reader)
	case 0x7f, 0x08:
		return selectionKeyBackspace, nil
	default:
		if key >= 0x20 {
			return selectionKey(key), nil
		}
		return selectionKeyUnknown, nil
	}
}

func readEscapeKey(reader *bufio.Reader) (selectionKey, error) {
	if reader.Buffered() == 0 {
		return selectionKeyCancel, nil
	}

	key, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return selectionKeyCancel, nil
		}
		return selectionKeyUnknown, err
	}
	if key != '[' && key != 'O' {
		return selectionKeyUnknown, nil
	}
	if reader.Buffered() == 0 {
		return selectionKeyUnknown, nil
	}

	code, err := reader.ReadByte()
	if err != nil {
		return selectionKeyUnknown, err
	}
	return selectionKeyFromCode(code), nil
}

func readWindowsConsoleKey(reader *bufio.Reader) (selectionKey, error) {
	if reader.Buffered() == 0 {
		return selectionKeyUnknown, nil
	}

	code, err := reader.ReadByte()
	if err != nil {
		return selectionKeyUnknown, err
	}
	return selectionKeyFromCode(code), nil
}

func selectionKeyFromCode(code byte) selectionKey {
	switch code {
	case 'A', 'H':
		return selectionKeyUp
	case 'B', 'P':
		return selectionKeyDown
	default:
		return selectionKeyUnknown
	}
}

func samePath(a, b string) bool {
	a = filepath.Clean(a)
	b = filepath.Clean(b)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(a, b)
	}
	return a == b
}
