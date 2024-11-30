package handlers

import (
	"github.com/charmbracelet/huh"
	execBuilders "github.com/mr-smith-org/mr/cmd/commands/exec/builders"
	"github.com/mr-smith-org/mr/cmd/constants"
)

type MultiSelectHandler struct {
}

func NewMultiSelectHandler() *MultiSelectHandler {
	return &MultiSelectHandler{}
}

func (h *MultiSelectHandler) Handle(data any, vars map[string]any) (huh.Field, string, any, error) {
	return handleMultiSelect(data.(map[string]interface{}), vars)
}

func handleMultiSelect(input map[string]interface{}, vars map[string]interface{}) (huh.Field, string, any, error) {
	var err error

	label, err := execBuilders.BuildStringValue("label", input, vars, false, constants.MultiSelectComponent)
	if err != nil {
		return nil, "", nil, err
	}
	description, err := execBuilders.BuildStringValue("description", input, vars, false, constants.MultiSelectComponent)
	if err != nil {
		return nil, "", nil, err
	}
	out, err := execBuilders.BuildStringValue("out", input, vars, true, constants.MultiSelectComponent)
	if err != nil {
		return nil, "", nil, err
	}
	limit, err := execBuilders.BuildIntValue("limit", input, vars, false, constants.MultiSelectComponent)
	if err != nil {
		return nil, "", nil, err
	}
	optionsFrom, err := execBuilders.BuildStringValue("options-from", input, vars, false, constants.SelectComponent)
	if err != nil {
		return nil, "", nil, err
	}
	if optionsFrom != "" {
		input["options"], err = handleOptionsFrom(optionsFrom, vars)
		if err != nil {
			return nil, "", nil, err
		}
	}
	options := []huh.Option[string]{}
	if mapOptions, ok := input["options"].([]interface{}); ok {
		for _, option := range mapOptions {
			optionMap := option.(map[string]interface{})
			label, err := execBuilders.BuildStringValue("label", optionMap, vars, true, constants.MultiSelectOptionComponent)
			if err != nil {
				return nil, "", nil, err
			}
			value, err := execBuilders.BuildStringValue("value", optionMap, vars, false, constants.MultiSelectOptionComponent)
			if err != nil {
				return nil, "", nil, err
			}
			if value == "" {
				value = label
			}
			options = append(options, huh.NewOption[string](label, value))
		}

		var outValue []string
		h := huh.NewMultiSelect[string]().
			Title(label).
			Description(description).
			Options(options...).
			Value(&outValue)

		if limit > 0 {
			h.Limit(limit)
		}

		return h, out, &outValue, nil
	}
	return nil, out, nil, nil
}
