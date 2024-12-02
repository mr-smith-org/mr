package execModule

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	handlers "github.com/mr-smith-org/mr/cmd/commands/run/handlers"
	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/cmd/ui/selectInput"
	"github.com/mr-smith-org/mr/cmd/ui/utils/program"
	"github.com/mr-smith-org/mr/cmd/ui/utils/steps"
	"github.com/mr-smith-org/mr/internal/services"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/style"

	"github.com/spf13/afero"
)

func Execute() {
	if shared.Pipeline == "" || shared.Module == "" {
		shared.Pipeline = handleTea()
	}
	vars := map[string]interface{}{
		"data": map[string]interface{}{},
	}
	hdl := handlers.NewPipelineHandler(shared.Pipeline, shared.Module)
	err := hdl.Handle(nil, vars)
	if err != nil {
		style.ErrorPrint(err.Error())
		os.Exit(1)
	}
}

func handleTea() string {
	var err error
	program := program.NewProgram()

	fs := filesystem.NewFileSystem(afero.NewOsFs())
	if shared.Module == "" {
		modulesService := services.NewModuleService(shared.FilesPath, fs)
		modules, err := modulesService.GetAll()
		if err != nil {
			style.ErrorPrint("getting modules error: " + err.Error())
			os.Exit(1)
		}
		var options = make([]steps.Item, 0)
		for key, module := range modules {
			options = append(options, steps.NewItem(
				key,
				key,
				module.Description,
				[]string{},
			))
		}

		output := &selectInput.Selection{}
		p := tea.NewProgram(selectInput.InitialSelectInputModel(options, output, "Select a module", false, program))
		_, err = p.Run()
		if err != nil {
			style.ErrorPrint("error running program: " + err.Error())
			os.Exit(1)
		}

		program.ExitCLI(p)

		shared.Module = output.Choice
	}

	pipelineService := services.NewPipelineService(shared.FilesPath+"/"+shared.Module+"/"+shared.PipelinesPath, fs)
	pipelines, err := pipelineService.GetAll(true)
	if err != nil {
		style.ErrorPrint("getting pipelines error: " + err.Error())
		os.Exit(1)
	}
	var options = make([]steps.Item, 0)
	for key, pipeline := range pipelines {
		options = append(options, steps.NewItem(
			key,
			key,
			pipeline.Description,
			[]string{},
		))
	}

	output := &selectInput.Selection{}
	p := tea.NewProgram(selectInput.InitialSelectInputModel(options, output, "Select a pipeline", false, program))
	_, err = p.Run()

	program.ExitCLI(p)

	if err != nil {
		style.ErrorPrint("error running program: " + err.Error())
		os.Exit(1)
	}
	return output.Choice
}