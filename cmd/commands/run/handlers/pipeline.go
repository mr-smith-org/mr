package handlers

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-smith-org/mr/cmd/constants"
	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/internal/domain"
	"github.com/mr-smith-org/mr/internal/services"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/spf13/afero"
)

type PipelineHandler struct {
	name   string
	module string
}

func NewPipelineHandler(name, module string) *PipelineHandler {
	return &PipelineHandler{name: name, module: module}
}

func (h *PipelineHandler) Handle(data any, vars map[string]any) error {
	return handlePipeline(h.name, h.module, vars)
}

func handlePipeline(name, moduleName string, vars map[string]interface{}) error {
	var err error
	var pipeline = &domain.Pipeline{}
	fs := filesystem.NewFileSystem(afero.NewOsFs())
	if moduleName != "" {
		moduleService := services.NewModuleService(shared.FilesPath, fs)
		modules, err := moduleService.GetAll()
		if err != nil {
			return err
		}
		module := modules[moduleName]
		pipeline, err = moduleService.GetPipeline(&module, name, shared.FilesPath+"/"+moduleName+"/"+shared.PipelinesPath)

		if err != nil {
			return err
		}
	} else {
		pipelineService := services.NewPipelineService(shared.PipelinesPath, fs)
		pipeline, err = pipelineService.Get(name)
		if err != nil {
			return err
		}
	}

	mapHandlers := map[string]shared.Handler{
		constants.CmdHandler:      NewCmdHandler(),
		constants.LogHandler:      NewLogHandler(),
		constants.PipelineHandler: NewPipelineHandler(name, moduleName),
		constants.CreateHandler:   NewCreateHandler(moduleName),
		constants.LoadHandler:     NewLoadHandler(),
		constants.WhenHandler:     NewWhenHandler(moduleName),
		constants.ModifyHandler:   NewModifyHandler(moduleName),
		constants.FormHandler:     NewFormHandler(),
		constants.DefineHandler:   NewDefineHandler(),
	}

	for _, step := range pipeline.Steps {
		step := step.(map[string]interface{})
		for key, value := range step {
			hdl := mapHandlers[key]
			if hdl != nil {
				err := hdl.Handle(value, vars)
				if err != nil {
					return fmt.Errorf("[handler: %s] - %s", key, err.Error())
				}
			}
		}
	}
	return nil
}

func ExitCLI(tprogram *tea.Program) {
	if err := tprogram.ReleaseTerminal(); err != nil {
		log.Fatal(err)
	}
	os.Exit(1)
}
