package create

import (
	"net/url"
	"os"
	"strings"

	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/internal/domain"
	"github.com/mr-smith-org/mr/internal/handlers"
	"github.com/mr-smith-org/mr/internal/helpers"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/style"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	ProjectPath       string
	VariablesFile     string
	FromFile          string
	TemplateVariables map[string]interface{}
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a scaffold for a project based on Go Templates",
	Run: func(cmd *cobra.Command, args []string) {
		Create()
	},
}

func Create() {
	fs := filesystem.NewFileSystem(afero.NewOsFs())
	if VariablesFile != "" {
		var vars interface{}
		_, err := url.ParseRequestURI(VariablesFile)
		if err != nil {
			vars, err = helpers.UnmarshalFile(VariablesFile, fs)
			if err != nil {
				style.ErrorPrint("parsing file error: " + err.Error())
				os.Exit(1)
			}
		} else {
			style.LogPrint("downloading variables file")
			varsContent, err := shared.ReadFileFromURL(VariablesFile)
			if err != nil {
				style.ErrorPrint("reading file error: " + err.Error())
				os.Exit(1)
			}
			splitURL := strings.Split(VariablesFile, "/")
			vars, err = helpers.UnmarshalByExt(splitURL[len(splitURL)-1], []byte(varsContent))
			if err != nil {
				style.ErrorPrint("parsing file error: " + err.Error())
				os.Exit(1)
			}
		}
		TemplateVariables = vars.(map[string]interface{})
		build()
	}
}

// build initializes the Builder and triggers the build process.
// It reads the Mr. Sith configuration file and applies templates to create the project structure.
func build() {
	fs := filesystem.NewFileSystem(afero.NewOsFs())
	// Initialize a new Builder with the provided configurations.
	builder, err := domain.NewBuilder(fs, domain.NewConfig(ProjectPath, shared.FilesPath))
	builder.SetBuilderDataFromFile(shared.FilesPath+"/"+FromFile, TemplateVariables)
	if err != nil {
		style.ErrorPrint(err.Error())
		os.Exit(1)
	}

	// Execute the build process using the BuilderHandler.
	if err = handlers.NewBuilderHandler(builder).Build(); err != nil {
		style.ErrorPrint(err.Error())
		os.Exit(1)
	}
}

// init sets up flags for the 'create' subcommand and binds them to variables.
func init() {
	// Target file directory
	CreateCmd.Flags().StringVarP(&VariablesFile, "variables", "v", "", "path or URL to the variables file")
	CreateCmd.Flags().StringVarP(&ProjectPath, "project", "p", ".", "Path to the project you want to create")
	CreateCmd.Flags().StringVarP(&FromFile, "from", "f", ".", "Path to the YAML file with the structure and templates")
}
