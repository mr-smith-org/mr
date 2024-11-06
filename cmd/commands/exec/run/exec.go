package execRun

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

	handlers "github.com/mr-smith-org/mr/cmd/commands/exec/handlers"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var ExecCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a specific run without a module",
	Run: func(cmd *cobra.Command, args []string) {
		Execute()
	},
}

func Execute() {
	if shared.Run == "" {
		shared.Run = handleTea()
	}
	vars := map[string]interface{}{
		"data": map[string]interface{}{},
	}
	hdl := handlers.NewRunHandler(shared.Run, "")
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
	runService := services.NewRunService(shared.RunsPath, fs)
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
	ExecCmd.Flags().StringVarP(&shared.Run, "run", "r", "", "run to use")
}
