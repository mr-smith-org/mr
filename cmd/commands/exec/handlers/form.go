package handlers

import (
	"fmt"

	execBuilders "github.com/kuma-framework/kuma/v2/cmd/commands/exec/builders"
	"github.com/kuma-framework/kuma/v2/cmd/constants"
	"github.com/kuma-framework/kuma/v2/cmd/shared"
	"github.com/kuma-framework/kuma/v2/pkg/style"
	"github.com/charmbracelet/huh"
)

type FormHandler struct {
}

func NewFormHandler() *FormHandler {
	return &FormHandler{}
}

func (h *FormHandler) Handle(data any, vars map[string]any) error {
	return handleForm(data.(map[string]interface{}), vars)
}

func handleForm(formData map[string]interface{}, vars map[string]interface{}) error {
	data := vars["data"].(map[string]interface{})
	huhFields := []huh.Field{}
	title, err := execBuilders.BuildStringValue("title", formData, vars, false, constants.FormComponent)
	if err != nil {
		return err
	}
	description, err := execBuilders.BuildStringValue("description", formData, vars, false, constants.FormComponent)
	if err != nil {
		return err
	}
	accessibility, err := execBuilders.BuildBoolValue("accessibility", formData, vars, false, constants.FormComponent)
	if err != nil {
		return err
	}

	if _, ok := formData["fields"]; !ok {
		return fmt.Errorf("fields is required")
	}

	hdlMap := map[string]shared.FormFieldHandler{
		constants.SelectComponent:      NewSelectHandler(),
		constants.InputComponent:       NewInputHandler(),
		constants.MultiSelectComponent: NewMultiSelectHandler(),
		constants.TextComponent:        NewTextHandler(),
		constants.ConfirmComponent:     NewConfirmHandler(),
	}

	for _, field := range formData["fields"].([]interface{}) {
		fieldMap, ok := field.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid field map")
		}
		for key, value := range fieldMap {
			if value, ok := value.(map[string]interface{}); ok {
				handler, ok := hdlMap[key]
				if !ok {
					return fmt.Errorf("invalid field type: %s", key)
				}
				huhField, out, outValue, err := handler.Handle(value, vars)
				if err != nil {
					return fmt.Errorf("[field:%s] - %s", key, err.Error())
				}
				huhFields = append(huhFields, huhField)
				data[out] = outValue
			} else {
				return fmt.Errorf("invalid field type: %s", key)
			}
		}
	}
	form := huh.NewForm(
		huh.NewGroup(huhFields...).
			Title(title).
			Description(description),
	)
	form.WithTheme(style.KumaTheme())
	form.WithAccessible(accessibility)
	err = form.Run()
	if err != nil {
		return fmt.Errorf("error running form: %s", err.Error())
	}

	return nil
}
