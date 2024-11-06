package handlers

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-smith/mr/cmd/constants"
	"github.com/mr-smith/mr/cmd/shared"
	"github.com/mr-smith/mr/internal/domain"
	"github.com/mr-smith/mr/internal/services"
	"github.com/mr-smith/mr/pkg/filesystem"
	"github.com/spf13/afero"
)

type RunHandler struct {
	name   string
	module string
}

func NewRunHandler(name, module string) *RunHandler {
	return &RunHandler{name: name, module: module}
}

func (h *RunHandler) Handle(data any, vars map[string]any) error {
	return handleRun(h.name, h.module, vars)
}

func handleRun(name, moduleName string, vars map[string]interface{}) error {
	var err error
	var run = &domain.Run{}
	fs := filesystem.NewFileSystem(afero.NewOsFs())
	if moduleName != "" {
		moduleService := services.NewModuleService(shared.FilesPath, fs)
		modules, err := moduleService.GetAll()
		if err != nil {
			return err
		}
		module := modules[moduleName]
		run, err = moduleService.GetRun(&module, name, shared.FilesPath+"/"+moduleName+"/"+shared.RunsPath)

		if err != nil {
			return err
		}
	} else {
		runService := services.NewRunService(shared.RunsPath, fs)
		run, err = runService.Get(name)
		if err != nil {
			return err
		}
	}

	mapHandlers := map[string]shared.Handler{
		constants.CmdHandler:    NewCmdHandler(),
		constants.LogHandler:    NewLogHandler(),
		constants.RunHandler:    NewRunHandler(name, moduleName),
		constants.CreateHandler: NewCreateHandler(moduleName),
		constants.LoadHandler:   NewLoadHandler(),
		constants.WhenHandler:   NewWhenHandler(moduleName),
		constants.ModifyHandler: NewModifyHandler(moduleName),
		constants.FormHandler:   NewFormHandler(),
		constants.DefineHandler: NewDefineHandler(),
	}

	for _, step := range run.Steps {
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
