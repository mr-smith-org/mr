package modify

import (
	"net/url"
	"os"
	"strings"

	"github.com/kuma-framework/kuma/v2/cmd/shared"
	"github.com/kuma-framework/kuma/v2/internal/helpers"
	"github.com/kuma-framework/kuma/v2/pkg/filesystem"
	"github.com/kuma-framework/kuma/v2/pkg/functions"
	"github.com/kuma-framework/kuma/v2/pkg/style"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	ReplaceAction        = "replace"
	InsertBeforeAction   = "insert-before"
	InsertAfterAction    = "insert-after"
	ReplaceBetweenAction = "replace-between"
)

var (
	FilePath string

	VariablesFile string

	TemplateFile string

	TemplateVariables map[string]interface{}

	CodeMark string

	Action string
)

var ModifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify a scaffold for a project based on Go Templates",
	Run: func(cmd *cobra.Command, args []string) {
		Modify()
	},
}

func Modify() {
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

func build() {
	fs := filesystem.NewFileSystem(afero.NewOsFs())
	fileContent, err := fs.ReadFile(FilePath)
	if err != nil {
		style.ErrorPrint("reading file error: " + err.Error())
		os.Exit(1)
	}
	templateContent, err := fs.ReadFile(TemplateFile)
	if err != nil {
		style.ErrorPrint("reading template file error: " + err.Error())
		os.Exit(1)
	}
	templateContent, err = helpers.ReplaceVars(templateContent, map[string]interface{}{"data": TemplateVariables}, functions.GetFuncMap())
	if err != nil {
		style.ErrorPrint("parsing template file error: " + err.Error())
		os.Exit(1)
	}
	fileContent = HandleAction(Action, fileContent, templateContent, CodeMark)
	err = fs.WriteFile(FilePath, fileContent)
	if err != nil {
		style.ErrorPrint("writing file error: " + err.Error())
		os.Exit(1)
	}
}

func HandleAction(action string, fileContent string, templateContent string, codeMark string) string {
	if codeMark == "" {
		fileContent = templateContent
		return fileContent
	}
	if action == "" {
		action = ReplaceAction
	}
	switch action {
	case InsertBeforeAction:
		return strings.ReplaceAll(fileContent, codeMark, templateContent+codeMark)
	case InsertAfterAction:
		return strings.ReplaceAll(fileContent, codeMark, codeMark+templateContent)
	case ReplaceAction:
		return strings.ReplaceAll(fileContent, codeMark, templateContent)
	case ReplaceBetweenAction:
		splitCodeMark := strings.Split(codeMark, ",")
		if len(splitCodeMark) != 2 {
			style.ErrorPrint("invalid code mark: " + codeMark)
			os.Exit(1)
		}
		startIdx := strings.Index(fileContent, splitCodeMark[0])
		endIdx := strings.Index(fileContent, splitCodeMark[1])

		if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
			return fileContent
		}

		startIdx += len(splitCodeMark[0])

		return fileContent[:startIdx] + templateContent + fileContent[endIdx:]
	default:
		return fileContent
	}

}

func init() {
	ModifyCmd.Flags().StringVarP(&VariablesFile, "variables", "v", "", "path or URL to the variables file")
	ModifyCmd.Flags().StringVarP(&FilePath, "file", "f", "", "Path to the file you want to modify")
	ModifyCmd.Flags().StringVarP(&TemplateFile, "template", "t", ".", "Path to the template file that be added after the code mark")
	ModifyCmd.Flags().StringVarP(&CodeMark, "mark", "m", "", "Mark inside the file to be identify what part of the code needs to be modified")
	ModifyCmd.Flags().StringVarP(&Action, "action", "r", InsertBeforeAction, "Replace the code mark with the template content")
	ModifyCmd.MarkFlagRequired("file")
}
