package handlers

import (
	"github.com/charmbracelet/huh"
	execBuilders "github.com/mr-smith-org/mr/cmd/commands/exec/builders"
	"github.com/mr-smith-org/mr/cmd/constants"
	"github.com/mr-smith-org/mr/cmd/shared"
)

type TextHandler struct {
}

func NewTextHandler() *TextHandler {
	return &TextHandler{}
}

func (h *TextHandler) Handle(data any, vars map[string]any) (huh.Field, string, any, error) {
	return handleText(data.(map[string]interface{}), vars)
}

func handleText(input map[string]interface{}, vars map[string]interface{}) (huh.Field, string, any, error) {
	var err error

	label, err := execBuilders.BuildStringValue("label", input, vars, false, constants.TextComponent)
	if err != nil {
		return nil, "", nil, err
	}
	description, err := execBuilders.BuildStringValue("description", input, vars, false, constants.TextComponent)
	if err != nil {
		return nil, "", nil, err
	}
	out, err := execBuilders.BuildStringValue("out", input, vars, true, constants.TextComponent)
	if err != nil {
		return nil, "", nil, err
	}
	placeholder, err := execBuilders.BuildStringValue("placeholder", input, vars, false, constants.TextComponent)
	if err != nil {
		return nil, "", nil, err
	}

	if shared.Vars[out] != "" {
		outResult := shared.Vars[out]
		return nil, out, &outResult, nil
	}

	var outValue string
	h := huh.NewText().
		Title(label).
		Description(description).
		Placeholder(placeholder).
		Value(&outValue)

	return h, out, &outValue, nil
}
