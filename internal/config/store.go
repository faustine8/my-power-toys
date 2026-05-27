package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const Version = 1

type Project struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type File struct {
	Version  int       `json:"version"`
	Projects []Project `json:"projects"`
}

type Store struct {
	Path string
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".my-power-toys", "projects.json"), nil
}

func (s Store) Load() (File, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultFile(), nil
		}
		return File{}, err
	}

	var file File
	if err := json.Unmarshal(data, &file); err != nil {
		return File{}, err
	}
	if file.Version == 0 {
		file.Version = Version
	}
	if file.Projects == nil {
		file.Projects = []Project{}
	}
	return file, nil
}

func (s Store) Save(file File) error {
	if file.Version == 0 {
		file.Version = Version
	}
	if file.Projects == nil {
		file.Projects = []Project{}
	}

	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(s.Path, data, 0o644)
}

func defaultFile() File {
	return File{
		Version:  Version,
		Projects: []Project{},
	}
}
