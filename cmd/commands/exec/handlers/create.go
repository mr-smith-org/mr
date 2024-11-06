package handlers

import (
	execBuilders "github.com/mr-smith/mr/cmd/commands/exec/builders"
	"github.com/mr-smith/mr/cmd/constants"
	"github.com/mr-smith/mr/cmd/shared"
	"github.com/mr-smith/mr/internal/domain"
	"github.com/mr-smith/mr/internal/handlers"
	"github.com/mr-smith/mr/pkg/filesystem"
	"github.com/spf13/afero"
)

type CreateHandler struct {
	module string
}

func NewCreateHandler(module string) *CreateHandler {
	return &CreateHandler{module: module}
}

func (h *CreateHandler) Handle(data any, vars map[string]any) error {
	return handleCreate(h.module, data.(map[string]interface{}), vars)
}

func handleCreate(module string, data map[string]interface{}, vars map[string]interface{}) error {
	path := shared.FilesPath
	fs := filesystem.NewFileSystem(afero.NewOsFs())
	if module != "" {
		path = shared.FilesPath + "/" + module + "/" + shared.FilesPath
	}
	builder, err := domain.NewBuilder(fs, domain.NewConfig(".", path))
	if err != nil {
		return err
	}
	from, err := execBuilders.BuildStringValue("from", data, vars, true, constants.CreateHandler)
	if err != nil {
		return err
	}
	err = builder.SetBuilderDataFromFile(path+"/"+from, vars)
	if err != nil {
		return err
	}

	if err = handlers.NewBuilderHandler(builder).Build(); err != nil {
		return err
	}
	return nil
}
