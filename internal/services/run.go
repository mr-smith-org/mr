package services

import (
	"fmt"

	"github.com/kuma-framework/kuma/v2/internal/domain"
	"github.com/kuma-framework/kuma/v2/internal/helpers"
	"github.com/kuma-framework/kuma/v2/pkg/filesystem"
)

type RunService struct {
	path string
	fs   filesystem.FileSystemInterface
}

func NewRunService(path string, fs filesystem.FileSystemInterface) *RunService {
	return &RunService{
		path: path,
		fs:   fs,
	}
}

func (s *RunService) GetAll(onlyVisible bool) (map[string]domain.Run, error) {
	deprecateRunsFileMsg := "\nif your a using runs.yaml file, please move it to the runs folder"
	files, err := s.fs.ReadDir(s.path)
	if err != nil {
		return nil, fmt.Errorf("error reading runs directory: %w%s", err, deprecateRunsFileMsg)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no runs found in %s%s", s.path, deprecateRunsFileMsg)
	}
	runs := make(map[string]domain.Run)
	for _, fileName := range files {
		data, err := helpers.UnmarshalFile(s.path+"/"+fileName, s.fs)
		if err != nil {
			return nil, err
		}
		for key, run := range data {
			if _, ok := runs[key]; ok {
				return nil, fmt.Errorf("conflict between runs found for the run %s\n rename one of them and try again", key)
			}
			steps, ok := run.(map[string]interface{})["steps"].([]interface{})
			if !ok {
				steps = []interface{}{}
			}
			visible, ok := run.(map[string]interface{})["visible"].(bool)
			if !ok {
				visible = true
			}
			description, ok := run.(map[string]interface{})["description"].(string)
			if !ok {
				description = ""
			}
			if onlyVisible && !visible {
				continue
			}
			runs[key] = domain.NewRun(
				key,
				description,
				steps,
				fileName,
				visible,
			)
		}
	}

	return runs, nil
}

func (s *RunService) Get(name string) (*domain.Run, error) {
	runs, err := s.GetAll(false)
	if err != nil {
		return nil, err
	}
	run, ok := runs[name]
	if !ok {
		return nil, fmt.Errorf("run not found: %s", name)
	}
	return &run, nil
}
