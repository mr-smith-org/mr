package execModule

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	handlers "github.com/mr-smith-org/mr/cmd/commands/exec/handlers"
	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/cmd/ui/selectInput"
	"github.com/mr-smith-org/mr/cmd/ui/utils/program"
	"github.com/mr-smith-org/mr/cmd/ui/utils/steps"
	"github.com/mr-smith-org/mr/internal/services"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/style"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var ExecModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Execute a specific run from a module",
	Run: func(cmd *cobra.Command, args []string) {
		Execute()
	},
}

func Execute() {
	if shared.Run == "" || shared.Module == "" {
		shared.Run = handleTea()
	}
	vars := map[string]interface{}{
		"data": map[string]interface{}{},
	}
	hdl := handlers.NewRunHandler(shared.Run, shared.Module)
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

	runService := services.NewRunService(shared.FilesPath+"/"+shared.Module+"/"+shared.RunsPath, fs)
	runs, err := runService.GetAll(true)
	if err != nil {
		style.ErrorPrint("getting runs error: " + err.Error())
		os.Exit(1)
	}
	var options = make([]steps.Item, 0)
	for key, run := range runs {
		options = append(options, steps.NewItem(
			key,
			key,
			run.Description,
			[]string{},
		))
	}

	output := &selectInput.Selection{}
	p := tea.NewProgram(selectInput.InitialSelectInputModel(options, output, "Select a run", false, program))
	_, err = p.Run()

	program.ExitCLI(p)

	if err != nil {
		style.ErrorPrint("error running program: " + err.Error())
		os.Exit(1)
	}
	return output.Choice
}

func init() {
	ExecModuleCmd.Flags().StringVarP(&shared.Run, "run", "r", "", "run to use")
	ExecModuleCmd.Flags().StringVarP(&shared.Module, "module", "m", "", "runs module")
}
