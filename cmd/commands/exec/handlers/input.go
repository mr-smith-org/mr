package handlers

import (
	execBuilders "github.com/arthurbcp/kuma/v2/cmd/commands/exec/builders"
	"github.com/arthurbcp/kuma/v2/cmd/constants"
	"github.com/charmbracelet/huh"
)

type InputHandler struct {
}

func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

func (h *InputHandler) Handle(data any, vars map[string]any) (huh.Field, string, any, error) {
	return handleInput(data.(map[string]interface{}), vars)
}

func handleInput(input map[string]interface{}, vars map[string]interface{}) (huh.Field, string, any, error) {
	var err error

	label, err := execBuilders.BuildStringValue("label", input, vars, false, constants.InputComponent)
	if err != nil {
		return nil, "", nil, err
	}
	description, err := execBuilders.BuildStringValue("description", input, vars, false, constants.InputComponent)
	if err != nil {
		return nil, "", nil, err
	}
	out, err := execBuilders.BuildStringValue("out", input, vars, true, constants.InputComponent)
	if err != nil {
		return nil, "", nil, err
	}
	placeholder, err := execBuilders.BuildStringValue("placeholder", input, vars, false, constants.InputComponent)
	if err != nil {
		return nil, "", nil, err
	}
	var outValue string
	h := huh.NewInput().
		Title(label).
		Description(description).
		Placeholder(placeholder).
		Value(&outValue)

	return h, out, &outValue, nil
}
