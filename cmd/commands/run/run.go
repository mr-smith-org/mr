package run

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	execModule "github.com/mr-smith-org/mr/cmd/commands/run/module"
	execPipeline "github.com/mr-smith-org/mr/cmd/commands/run/pipeline"
	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/cmd/ui/selectInput"
	"github.com/mr-smith-org/mr/cmd/ui/utils/program"
	"github.com/mr-smith-org/mr/cmd/ui/utils/steps"
	"github.com/mr-smith-org/mr/pkg/style"
	"github.com/spf13/cobra"
)

var varsSlice = make([]string, 0)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a pipeline",
	Run: func(cmd *cobra.Command, args []string) {

		if len(varsSlice) > 0 {
			for _, v := range varsSlice {
				splitVar := strings.Split(v, ":")
				shared.Vars[splitVar[0]] = splitVar[1]
			}
		}
		if shared.Pipeline == "" && shared.Module == "" {
			choice := handleTea()
			if choice == "pipeline" {
				execPipeline.Execute()
				os.Exit(0)
			} else {
				execModule.Execute()
				os.Exit(0)
			}
		}
		if shared.Module != "" {
			execModule.Execute()
			os.Exit(0)
		}

		if shared.Pipeline != "" {
			execPipeline.Execute()
			os.Exit(0)
		}
	},
}

func handleTea() string {
	var options = []steps.Item{
		steps.NewItem("pipeline", "pipeline", "", []string{}),
		steps.NewItem("module", "module", "", []string{}),
	}
	program := program.NewProgram()
	output := &selectInput.Selection{}
	p := tea.NewProgram(selectInput.InitialSelectInputModel(options, output, "Select a option to be executed", false, program))
	_, err := p.Run()
	if err != nil {
		style.ErrorPrint("error running program: " + err.Error())
		os.Exit(1)
	}

	program.ExitCLI(p)

	return output.Choice
}

func init() {
	RunCmd.Flags().StringVarP(&shared.Pipeline, "pipeline", "p", "", "pipeline to use")
	RunCmd.Flags().StringVarP(&shared.Module, "module", "m", "", "module to use")
	RunCmd.Flags().StringSliceVarP(&varsSlice, "vars", "v", []string{}, "variables to use")
}