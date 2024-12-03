package services

import (
	"fmt"

	"github.com/mr-smith-org/mr/internal/domain"
	"github.com/mr-smith-org/mr/internal/helpers"
	"github.com/mr-smith-org/mr/pkg/filesystem"
)

type PipelineService struct {
	path string
	fs   filesystem.FileSystemInterface
}

func NewPipelineService(path string, fs filesystem.FileSystemInterface) *PipelineService {
	return &PipelineService{
		path: path,
		fs:   fs,
	}
}

func (s *PipelineService) GetAll(onlyVisible bool) (map[string]domain.Pipeline, error) {
	deprecateRunMsg := "\nif your a using runs folder, please rename it to pipelines and try again"
	files, err := s.fs.ReadDir(s.path)
	if err != nil {
		return nil, fmt.Errorf("error reading pipelines directory: %w%s", err, deprecateRunMsg)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no pipeline found in %s%s", s.path, deprecateRunMsg)
	}
	pipelines := make(map[string]domain.Pipeline)
	for _, fileName := range files {
		data, err := helpers.UnmarshalFile(s.path+"/"+fileName, s.fs)
		if err != nil {
			return nil, err
		}
		for key, pipeline := range data {
			if _, ok := pipelines[key]; ok {
				return nil, fmt.Errorf("conflict between pipelines found for the pipeline %s\n rename one of them and try again", key)
			}
			steps, ok := pipeline.(map[string]interface{})["steps"].([]interface{})
			if !ok {
				steps = []interface{}{}
			}
			visible, ok := pipeline.(map[string]interface{})["visible"].(bool)
			if !ok {
				visible = true
			}
			description, ok := pipeline.(map[string]interface{})["description"].(string)
			if !ok {
				description = ""
			}
			if onlyVisible && !visible {
				continue
			}
			pipelines[key] = domain.NewPipeline(
				key,
				description,
				steps,
				fileName,
				visible,
			)
		}
	}

	return pipelines, nil
}

func (s *PipelineService) Get(name string) (*domain.Pipeline, error) {
	pipelines, err := s.GetAll(false)
	if err != nil {
		return nil, err
	}
	pipeline, ok := pipelines[name]

	if !ok {
		return nil, fmt.Errorf("pipeline not found: %s", name)
	}
	return &pipeline, nil
}
