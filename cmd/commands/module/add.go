package module

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	execModule "github.com/mr-smith-org/mr/cmd/commands/run/module"
	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/cmd/ui/selectInput"
	"github.com/mr-smith-org/mr/cmd/ui/utils/program"
	"github.com/mr-smith-org/mr/cmd/ui/utils/steps"
	"github.com/mr-smith-org/mr/internal/services"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/style"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	Repository string
)

var ModuleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an Mr. Smith module from a GitHub repository",
	Run: func(cmd *cobra.Command, args []string) {
		if Repository == "" {
			Repository = handleTea()
		}
		download(cmd)
	},
}

func handleTea() string {
	program := program.NewProgram()
	var options = make([]steps.Item, 0)

	for repository, template := range shared.Modules {
		options = append(options, steps.NewItem(
			template.Name,
			repository,
			template.Description,
			template.Tags,
		))
	}

	output := &selectInput.Selection{}
	p := tea.NewProgram(selectInput.InitialSelectInputModel(options, output, "Select a template or type \"o\" to use a different repository", true, program))
	_, err := p.Run()

	program.ExitCLI(p)

	if err != nil {
		style.ErrorPrint("error running program: " + err.Error())
		os.Exit(1)
	}
	return output.Choice
}

func download(cmd *cobra.Command) {
	if Repository == "" {
		cmd.Help()
		style.LogPrint("\nplease specify a repository")
		os.Exit(1)
	}

	style.LogPrint("getting templates from github repository...")
	fs := filesystem.NewFileSystem(afero.NewOsFs())

	err := fs.CreateDirectoryIfNotExists(shared.FilesPath)
	if err != nil {
		style.ErrorPrint("error creating files directory: " + err.Error())
		os.Exit(1)
	}

	err = addModule(Repository)
	if err != nil {
		style.ErrorPrint("error adding module: " + err.Error())
		os.Exit(1)
	}
	style.CheckMarkPrint("templates downloaded successfully!\n")

	moduleService := services.NewModuleService(shared.FilesPath, fs)
	moduleName := moduleService.GetModuleName(Repository)
	err = moduleService.Add(moduleName)
	if err != nil {
		style.ErrorPrint("error adding module: " + err.Error())
		os.Exit(1)
	}
	shared.Module = moduleName
	execModule.Execute()
	os.Exit(0)
}

func addModule(module string) error {
	fs := filesystem.NewFileSystem(afero.NewOsFs())
	moduleService := services.NewModuleService(shared.FilesPath, fs)
	if err := shared.RunCommand("git", "clone", shared.GitHubURL+"/"+module, shared.FilesPath+"/"+moduleService.GetModuleName(module)); err != nil {
		return err
	}
	if err := shared.RunCommand("rm", "-rf", shared.FilesPath+"/"+moduleService.GetModuleName(module)+"/.git"); err != nil {
		return err
	}
	return nil
}

func init() {
	ModuleAddCmd.Flags().StringVarP(&Repository, "repository", "r", "", "Github repository with a Mr. Smith module")
}
