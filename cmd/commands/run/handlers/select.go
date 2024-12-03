package handlers

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	execBuilders "github.com/mr-smith-org/mr/cmd/commands/run/builders"
	"github.com/mr-smith-org/mr/cmd/constants"
	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/internal/helpers"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/functions"
	"github.com/mr-smith-org/mr/pkg/style"
	"github.com/spf13/afero"
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
	optionsFrom, err := execBuilders.BuildStringValue("options-from", input, vars, false, constants.SelectComponent)
	if err != nil {
		return nil, "", nil, err
	}

	out, err := execBuilders.BuildStringValue("out", input, vars, true, constants.SelectComponent)
	if err != nil {
		return nil, "", nil, err
	}

	if shared.Vars[out] != "" {
		outResult := shared.Vars[out]
		return nil, out, &outResult, nil
	}

	options := []huh.Option[string]{}

	if optionsFrom != "" {
		input["options"], err = handleOptionsFrom(optionsFrom, vars)
		if err != nil {
			return nil, "", nil, err
		}
	}
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

		return h, out, &outValue, nil
	}
	return nil, out, nil, nil
}

func handleOptionsFrom(from string, vars map[string]interface{}) ([]interface{}, error) {
	var err error
	fs := filesystem.NewFileSystem(afero.NewOsFs())

	var options map[string]interface{}
	parsedURI, err := url.ParseRequestURI(from)
	if err != nil {
		options, err = helpers.UnmarshalFileAndReplaceVars(from, vars, fs)
		if err != nil {
			return nil, fmt.Errorf("parsing file error: %s", err.Error())
		}
	} else {
		err = spinner.New().
			Title("Downloading variables file").
			Action(func() {
				varsContent, err := fs.ReadFileFromURL(from)
				if err != nil {
					style.ErrorPrint("reading file error: " + err.Error())
					os.Exit(1)
				}
				varsContent, err = helpers.ReplaceVars(varsContent, vars, functions.GetFuncMap())
				if err != nil {
					style.ErrorPrint("parsing file error: " + err.Error())
					os.Exit(1)
				}
				splitURIPath := strings.Split(parsedURI.Path, "/")
				options, err = helpers.UnmarshalByExt(splitURIPath[len(splitURIPath)-1], []byte(varsContent))
				if err != nil {
					style.ErrorPrint("parsing file error: " + err.Error())
					os.Exit(1)
				}
			}).
			Run()

		if err != nil {
			return nil, fmt.Errorf("[handler:load] - downloading variables file error: %s", err.Error())
		}
	}

	return options["options"].([]interface{}), nil
}
