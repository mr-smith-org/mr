package handlers

import (
	execBuilders "github.com/kuma-framework/kuma/v2/cmd/commands/exec/builders"
	"github.com/kuma-framework/kuma/v2/cmd/constants"
)

type WhenHandler struct {
	module string
}

func NewWhenHandler(module string) *WhenHandler {
	return &WhenHandler{module: module}
}

func (h *WhenHandler) Handle(data any, vars map[string]any) error {
	return handleWhen(h.module, data.(map[string]interface{}), vars)
}

func handleWhen(module string, params map[string]interface{}, vars map[string]interface{}) error {
	isTrue, err := execBuilders.BuildBoolValue("condition", params, vars, true, constants.WhenHandler)
	if err != nil {
		return err
	}
	run, err := execBuilders.BuildStringValue("run", params, vars, true, constants.WhenHandler)
	if err != nil {
		return err
	}
	if isTrue {
		hdl := NewRunHandler(run, module)
		err := hdl.Handle(nil, vars)
		if err != nil {
			return err
		}
	}

	return nil
}
