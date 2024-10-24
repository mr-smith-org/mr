package handlers

import (
	"fmt"

	execBuilders "github.com/arthurbcp/kuma/v2/cmd/commands/exec/builders"
	"github.com/arthurbcp/kuma/v2/cmd/commands/modify"
	"github.com/arthurbcp/kuma/v2/cmd/constants"
	"github.com/arthurbcp/kuma/v2/cmd/shared"
	"github.com/arthurbcp/kuma/v2/internal/functions"
	"github.com/arthurbcp/kuma/v2/internal/helpers"
	"github.com/arthurbcp/kuma/v2/pkg/filesystem"
	"github.com/arthurbcp/kuma/v2/pkg/style"
	"github.com/spf13/afero"
)

type ModifyHandler struct {
	module string
}

func NewModifyHandler(module string) *ModifyHandler {
	return &ModifyHandler{module: module}
}

func (h *ModifyHandler) Handle(data any, vars map[string]any) error {
	return handleModify(h.module, data.(map[string]interface{}), vars)
}

func handleModify(module string, data map[string]interface{}, vars map[string]interface{}) error {
	path := shared.KumaFilesPath
	fs := filesystem.NewFileSystem(afero.NewOsFs())
	if module != "" {
		path = shared.KumaFilesPath + "/" + module + "/" + shared.KumaFilesPath
	}
	file, err := execBuilders.BuildStringValue("file", data, vars, true, constants.ModifyHandler)
	if err != nil {
		return err
	}
	fileContent, err := fs.ReadFile(file)
	if err != nil {
		_, err = fs.CreateFile(file)
		if err != nil {
			return fmt.Errorf("creating file error: %s", err.Error())
		}
		fileContent = ""
	}
	template, err := execBuilders.BuildStringValue("template", data, vars, true, constants.ModifyHandler)
	if err != nil {
		return err
	}
	codeMark, err := execBuilders.BuildStringValue("mark", data, vars, false, constants.ModifyHandler)
	if err != nil {
		return err
	}
	action, err := execBuilders.BuildStringValue("action", data, vars, false, constants.ModifyHandler)
	if err != nil {
		return err
	}
	templateContent, err := fs.ReadFile(path + "/" + template)
	if err != nil {
		return fmt.Errorf("reading template file error: %s", err.Error())
	}
	templateContent, err = helpers.ReplaceVars(templateContent, vars, functions.GetFuncMap())
	if err != nil {
		return fmt.Errorf("parsing template file error: %s", err.Error())
	}
	fileContent = modify.HandleAction(action, fileContent, templateContent, codeMark)
	err = fs.WriteFile(file, fileContent)
	if err != nil {
		return fmt.Errorf("writing file error: %s", err.Error())
	}
	style.CheckMarkPrint(fmt.Sprintf("file %s modified successfully!", file))
	return nil
}
