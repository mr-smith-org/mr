package handlers

import (
	execBuilders "github.com/arthurbcp/kuma/v2/cmd/commands/exec/builders"
	"github.com/arthurbcp/kuma/v2/cmd/constants"
	"github.com/charmbracelet/huh"
)

type SelectHandler struct {
}

func NewSelectHandler() *SelectHandler {
	return &SelectHandler{}
}

func (h *SelectHandler) Handle(data any, vars map[string]any) (huh.Field, string, any, error) {
	return handleSelect(data.(map[string]interface{}), vars)
}

func handleSelect(input map[string]interface{}, vars map[string]interface{}) (huh.Field, string, any, error) {
	var err error

	label, err := execBuilders.BuildStringValue("label", input, vars, false, constants.SelectComponent)
	if err != nil {
		return nil, "", nil, err
	}
	description, err := execBuilders.BuildStringValue("description", input, vars, false, constants.SelectComponent)
	if err != nil {
		return nil, "", nil, err
	}
	out, err := execBuilders.BuildStringValue("out", input, vars, true, constants.SelectComponent)
	if err != nil {
		return nil, "", nil, err
	}

	options := []huh.Option[string]{}
	if mapOptions, ok := input["options"].([]interface{}); ok {
		for _, option := range mapOptions {
			optionMap := option.(map[string]interface{})
			label, err := execBuilders.BuildStringValue("label", optionMap, vars, true, constants.SelectOptionComponent)
			if err != nil {
				return nil, "", nil, err
			}
			value, err := execBuilders.BuildStringValue("value", optionMap, vars, false, constants.SelectOptionComponent)
			if err != nil {
				return nil, "", nil, err
			}
			if value == "" {
				value = label
			}
			options = append(options, huh.NewOption[string](label, value))
		}

		var outValue string
		h := huh.NewSelect[string]().
			Title(label).
			Description(description).
			Options(options...).
			Value(&outValue)

		return h, out, outValue, nil
	}
	return nil, out, nil, nil
}
