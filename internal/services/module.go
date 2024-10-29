package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kuma-framework/kuma/v2/cmd/shared"
	"github.com/kuma-framework/kuma/v2/internal/domain"
	"github.com/kuma-framework/kuma/v2/internal/helpers"
	"github.com/kuma-framework/kuma/v2/pkg/filesystem"
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
	modulesFile := shared.KumaFilesPath + "/kuma-modules.yaml"
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
	modulesFile := shared.KumaFilesPath + "/kuma-modules.yaml"
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
	configData, err := helpers.UnmarshalFile(s.path+"/"+module+"/kuma-config.yaml", s.fs)
	if err != nil {
		return domain.Module{}, err
	}
	runsService := NewRunService(s.path+"/"+module+"/"+shared.KumaRunsPath, s.fs)
	runs, err := runsService.GetAll(false)
	if err != nil {
		return domain.Module{}, err
	}
	return domain.NewModule(configData, runs), nil
}

func (s *ModuleService) GetAll() (map[string]domain.Module, error) {
	modules, err := helpers.UnmarshalFile(shared.KumaFilesPath+"/kuma-modules.yaml", s.fs)
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

func (s *ModuleService) GetRun(module *domain.Module, runKey string, modulePath string) (*domain.Run, error) {
	moduleRun := module.Runs[runKey]
	runs, err := helpers.UnmarshalFile(modulePath+"/"+moduleRun.File, s.fs)
	if err != nil {
		return nil, err
	}
	runContent, ok := runs[runKey]
	if !ok {
		return nil, fmt.Errorf("run not found: %s", runKey)
	}
	description, ok := runContent.(map[string]interface{})["description"].(string)
	if !ok {
		description = ""
	}
	steps, ok := runContent.(map[string]interface{})["steps"].([]interface{})
	if !ok {
		steps = []interface{}{}
	}
	visible, ok := runContent.(map[string]interface{})["visible"].(bool)
	if !ok {
		visible = true
	}
	run := domain.NewRun(
		runKey,
		description,
		steps,
		moduleRun.File,
		visible,
	)
	return &run, nil
}
