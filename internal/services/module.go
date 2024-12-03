package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/internal/domain"
	"github.com/mr-smith-org/mr/internal/helpers"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"gopkg.in/yaml.v3"
)

type ModuleService struct {
	path string
	fs   filesystem.FileSystemInterface
}

func NewModuleService(path string, fs filesystem.FileSystemInterface) *ModuleService {
	return &ModuleService{
		path: path,
		fs:   fs,
	}
}

func (s *ModuleService) Add(newModule string) error {
	modulesFile := shared.FilesPath + "/" + shared.ModulesFileName
	_, err := s.fs.CreateFileIfNotExists(modulesFile)
	if err != nil {
		return err
	}
	modules, err := helpers.UnmarshalFile(modulesFile, s.fs)
	if err != nil {
		return err
	}

	module, err := s.Get(newModule)
	if err != nil {
		return err
	}

	mapModule, err := helpers.StructToMap(module)
	if err != nil {
		return err
	}
	modules[newModule] = mapModule

	yamlContent, err := yaml.Marshal(modules)
	if err != nil {
		return err
	}
	s.fs.WriteFile(modulesFile, string(yamlContent))
	return nil
}

func (s *ModuleService) Remove(module string) error {
	modulesFile := shared.FilesPath + "/" + shared.ModulesFileName
	modules, err := helpers.UnmarshalFile(modulesFile, s.fs)
	if err != nil {
		return err
	}

	delete(modules, module)

	if len(modules) == 0 {
		s.fs.WriteFile(modulesFile, "")
	}

	yamlContent, err := yaml.Marshal(modules)
	if err != nil {
		return err
	}
	s.fs.WriteFile(modulesFile, string(yamlContent))
	return nil
}

func (s *ModuleService) Get(module string) (domain.Module, error) {
	configData, err := helpers.UnmarshalFile(s.path+"/"+module+"/"+shared.ConfigFileName, s.fs)
	if err != nil {
		return domain.Module{}, err
	}
	pipelinesService := NewPipelineService(s.path+"/"+module+"/"+shared.PipelinesPath, s.fs)
	pipelines, err := pipelinesService.GetAll(false)
	if err != nil {
		return domain.Module{}, err
	}
	return domain.NewModule(configData, pipelines), nil
}

func (s *ModuleService) GetAll() (map[string]domain.Module, error) {
	modules, err := helpers.UnmarshalFile(shared.FilesPath+"/"+shared.ModulesFileName, s.fs)
	if err != nil {
		return nil, err
	}
	modulesMap := map[string]domain.Module{}
	for key, module := range modules {
		moduleString, err := json.Marshal(module)
		if err != nil {
			return nil, err
		}
		m := &domain.Module{}
		err = json.Unmarshal(moduleString, m)
		if err != nil {
			return nil, err
		}
		modulesMap[key] = *m
	}
	return modulesMap, nil
}

func (s *ModuleService) GetModuleName(repo string) string {
	splitRepo := strings.Split(repo, "/")
	return splitRepo[1]
}

func (s *ModuleService) GetPipeline(module *domain.Module, pipelineKey string, modulePath string) (*domain.Pipeline, error) {
	modulePipeline := module.Pipelines[pipelineKey]
	pipelines, err := helpers.UnmarshalFile(modulePath+"/"+modulePipeline.File, s.fs)
	if err != nil {
		return nil, err
	}
	pipelineContent, ok := pipelines[pipelineKey]
	if !ok {
		return nil, fmt.Errorf("pipeline not found: %s", pipelineKey)
	}
	description, ok := pipelineContent.(map[string]interface{})["description"].(string)
	if !ok {
		description = ""
	}
	steps, ok := pipelineContent.(map[string]interface{})["steps"].([]interface{})
	if !ok {
		steps = []interface{}{}
	}
	visible, ok := pipelineContent.(map[string]interface{})["visible"].(bool)
	if !ok {
		visible = true
	}
	pipeline := domain.NewPipeline(
		pipelineKey,
		description,
		steps,
		modulePipeline.File,
		visible,
	)
	return &pipeline, nil
}
