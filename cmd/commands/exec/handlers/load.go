package handlers

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	execBuilders "github.com/mr-smith-org/mr/cmd/commands/exec/builders"
	"github.com/mr-smith-org/mr/cmd/constants"
	"github.com/mr-smith-org/mr/internal/helpers"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/functions"
	"github.com/mr-smith-org/mr/pkg/style"
	"github.com/spf13/afero"
)

type LoadHandler struct {
}

func NewLoadHandler() *LoadHandler {
	return &LoadHandler{}
}

func (h *LoadHandler) Handle(data any, vars map[string]any) error {
	return handleLoad(data.(map[string]interface{}), vars)
}

func handleLoad(load map[string]interface{}, vars map[string]interface{}) error {
	var err error
	data := vars["data"].(map[string]interface{})
	fs := filesystem.NewFileSystem(afero.NewOsFs())

	from, err := execBuilders.BuildStringValue("from", load, vars, true, constants.LoadHandler)
	if err != nil {
		return err
	}

	out, err := execBuilders.BuildStringValue("out", load, vars, true, constants.LoadHandler)
	if err != nil {
		return err
	}

	var fileVars map[string]interface{}
	parsedURI, err := url.ParseRequestURI(from)
	if err != nil {
		fileVars, err = helpers.UnmarshalFileAndReplaceVars(from, vars, fs)
		if err != nil {
			return fmt.Errorf("[handler:load] - parsing file error: %s", err.Error())
		}
	} else {
		err = spinner.New().
			Title("Downloading variables file").
			Action(func() {
				varsContent, err := fs.ReadFileFromURL(from)
				if err != nil {
					style.ErrorPrint("[handler:load] - reading file error: " + err.Error())
					os.Exit(1)
				}
				varsContent, err = helpers.ReplaceVars(varsContent, vars, functions.GetFuncMap())
				if err != nil {
					style.ErrorPrint("[handler:load] - parsing file error: " + err.Error())
					os.Exit(1)
				}
				splitURIPath := strings.Split(parsedURI.Path, "/")
				fileVars, err = helpers.UnmarshalByExt(splitURIPath[len(splitURIPath)-1], []byte(varsContent))
				if err != nil {
					style.ErrorPrint("[handler:load] - parsing file error: " + err.Error())
					os.Exit(1)
				}
			}).
			Run()

		if err != nil {
			return fmt.Errorf("[handler:load] - downloading variables file error: %s", err.Error())
		}
	}

	data[out] = fileVars
	return nil
}
