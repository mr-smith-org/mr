package handlers

import (
	"strconv"

	"github.com/charmbracelet/huh"
	execBuilders "github.com/mr-smith-org/mr/cmd/commands/run/builders"
	"github.com/mr-smith-org/mr/cmd/constants"
	"github.com/mr-smith-org/mr/cmd/shared"
)

type ConfirmHandler struct {
}

func NewConfirmHandler() *ConfirmHandler {
	return &ConfirmHandler{}
}

func (h *ConfirmHandler) Handle(data any, vars map[string]any) (huh.Field, string, any, error) {
	return handleConfirm(data.(map[string]interface{}), vars)
}

func handleConfirm(input map[string]interface{}, vars map[string]interface{}) (huh.Field, string, any, error) {
	var err error
	data := vars["data"].(map[string]interface{})

	label, err := execBuilders.BuildStringValue("label", input, vars, false, constants.ConfirmComponent)
	if err != nil {
		return nil, "", nil, err
	}
	description, err := execBuilders.BuildStringValue("description", input, vars, false, constants.ConfirmComponent)
	if err != nil {
		return nil, "", nil, err
	}
	affirmative, err := execBuilders.BuildStringValue("affirmative", input, vars, false, constants.ConfirmComponent)
	if err != nil {
		return nil, "", nil, err
	}
	negative, err := execBuilders.BuildStringValue("negative", input, vars, false, constants.ConfirmComponent)
	if err != nil {
		return nil, "", nil, err
	}
	out, err := execBuilders.BuildStringValue("out", input, vars, true, constants.ConfirmComponent)
	if err != nil {
		return nil, "", nil, err
	}

	if affirmative == "" {
		affirmative = "Yes"
	}
	if negative == "" {
		negative = "No"
	}

	if shared.Vars[out] != "" {
		outResult, err := strconv.ParseBool(shared.Vars[out])
		if err != nil {
			return nil, "", nil, err
		}
		return nil, out, &outResult, nil
	}

	var outValue bool
	h := huh.NewConfirm().
		Title(label).
		Description(description).
		Affirmative(affirmative).
		Negative(negative).
		Value(&outValue)

	data[out] = out

	return h, out, &outValue, nil
}
