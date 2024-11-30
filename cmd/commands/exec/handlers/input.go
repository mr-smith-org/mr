package handlers

import (
	"github.com/charmbracelet/huh"
	execBuilders "github.com/mr-smith-org/mr/cmd/commands/exec/builders"
	"github.com/mr-smith-org/mr/cmd/constants"
	"github.com/mr-smith-org/mr/cmd/shared"
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

	if shared.Vars[out] != "" {
		outResult := shared.Vars[out]
		return nil, out, &outResult, nil
	}

	var outValue string
	h := huh.NewInput().
		Title(label).
		Description(description).
		Placeholder(placeholder).
		Value(&outValue)

	return h, out, &outValue, nil
}
