package project

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/faustine8/my-power-toys/internal/config"
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

	matches := make([]config.Project, 0, len(projects))
	for _, project := range projects {
		if strings.Contains(strings.ToLower(project.Name), query) ||
			strings.Contains(strings.ToLower(project.Path), query) {
			matches = append(matches, project)
		}
	}
	return matches
}

func Select(projects []config.Project, input io.Reader, prompt io.Writer) (config.Project, bool, error) {
	if len(projects) == 0 {
		return config.Project{}, false, nil
	}
	if len(projects) == 1 {
		return projects[0], true, nil
	}
	if prompt == nil {
		prompt = io.Discard
	}

	fmt.Fprintln(prompt, "Select project:")
	for i, project := range projects {
		fmt.Fprintf(prompt, "%d) %s\t%s\n", i+1, project.Name, project.Path)
	}
	fmt.Fprint(prompt, "Enter number: ")

	scanner := bufio.NewScanner(input)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return config.Project{}, false, err
		}
		return config.Project{}, false, nil
	}

	choice := strings.TrimSpace(scanner.Text())
	if choice == "" {
		return config.Project{}, false, nil
	}

	index, err := strconv.Atoi(choice)
	if err != nil || index < 1 || index > len(projects) {
		return config.Project{}, false, fmt.Errorf("invalid selection: %s", choice)
	}
	return projects[index-1], true, nil
}

func samePath(a, b string) bool {
	a = filepath.Clean(a)
	b = filepath.Clean(b)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(a, b)
	}
	return a == b
}
