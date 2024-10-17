package execHandlers

import (
	"os"

	execBuilders "github.com/arthurbcp/kuma/cmd/commands/exec/builders"
	"github.com/arthurbcp/kuma/cmd/ui/textInput"
	"github.com/arthurbcp/kuma/cmd/ui/utils/program"
	"github.com/arthurbcp/kuma/pkg/style"
	tea "github.com/charmbracelet/bubbletea"
)

func HandleInput(input map[string]interface{}, vars map[string]interface{}) {
	var err error
	program := program.NewProgram()
	data := vars["data"].(map[string]interface{})

	label, err := execBuilders.BuildStringValue("label", input, vars, false)
	if err != nil {
		style.ErrorPrint(err.Error())
		os.Exit(1)
	}

	out, err := execBuilders.BuildStringValue("out", input, vars, true)
	if err != nil {
		style.ErrorPrint(err.Error())
		os.Exit(1)
	}
	placeholder, err := execBuilders.BuildStringValue("placeholder", input, vars, false)
	if err != nil {
		style.ErrorPrint(err.Error())
		os.Exit(1)
	}

	output := &textInput.Output{}
	p := tea.NewProgram(textInput.InitialTextInputModel(output, label, placeholder, program))
	_, err = p.Run()

	program.ExitCLI(p)

	if err != nil {
		style.ErrorPrint("error running program: " + err.Error())
		os.Exit(1)
	}
	data[out] = output.Output

}
