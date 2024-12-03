package handlers

import (
	execBuilders "github.com/mr-smith-org/mr/cmd/commands/run/builders"
	"github.com/mr-smith-org/mr/cmd/constants"
)

type DefineHandler struct {
}

func NewDefineHandler() *DefineHandler {
	return &DefineHandler{}
}

func (h *DefineHandler) Handle(data any, vars map[string]any) error {
	return handleDefine(data.(map[string]interface{}), vars)
}

func handleDefine(params map[string]interface{}, vars map[string]interface{}) error {
	data := vars["data"].(map[string]interface{})
	variable, err := execBuilders.BuildStringValue("variable", params, vars, true, constants.DefineHandler)
	if err != nil {
		return err
	}
	var value any
	value, err = execBuilders.BuildBoolValue("value", params, vars, true, constants.DefineHandler)
	if err != nil {
		value, err = execBuilders.BuildIntValue("value", params, vars, true, constants.DefineHandler)
		if err != nil {
			value, err = execBuilders.BuildStringValue("value", params, vars, true, constants.DefineHandler)
			if err != nil {
				return err
			}
		}
	}
	data[variable] = value
	return nil
}
