package execPipeline

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/cmd/ui/selectInput"
	"github.com/mr-smith-org/mr/cmd/ui/utils/program"
	"github.com/mr-smith-org/mr/cmd/ui/utils/steps"
	"github.com/mr-smith-org/mr/internal/services"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/style"
	"github.com/spf13/afero"

	handlers "github.com/mr-smith-org/mr/cmd/commands/run/handlers"
)

func Execute() {
	if shared.Pipeline == "" {
		shared.Pipeline = handleTea()
	}
	vars := map[string]interface{}{
		"data": map[string]interface{}{},
	}
	hdl := handlers.NewPipelineHandler(shared.Module)
	err := hdl.Handle(shared.Pipeline, vars)
	if err != nil {
		style.ErrorPrint(err.Error())
		os.Exit(1)
	}
}

func handleTea() string {
	var err error
	program := program.NewProgram()

	fs := filesystem.NewFileSystem(afero.NewOsFs())
	pipelineService := services.NewPipelineService(shared.PipelinesPath, fs)
	pipelines, err := pipelineService.GetAll(true)
	if err != nil {
		style.ErrorPrint("getting pipelines error: " + err.Error())
		os.Exit(1)
	}
	pipelinesSlice := pipelineService.ToSliceAndSort(pipelines)

	var options = make([]steps.Item, 0)
	for _, pipeline := range pipelinesSlice {
		options = append(options, steps.NewItem(
			pipeline.Key,
			pipeline.Key,
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
