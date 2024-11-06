package module

import (
	"os"

	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/internal/services"
	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/style"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	Module string
)

var ModuleRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a Mr. Smith module",
	Run: func(cmd *cobra.Command, args []string) {
		if Module == "" {
			style.ErrorPrint("module is required")
			os.Exit(1)
		}
		err := removeModule(Module)
		if err != nil {
			style.ErrorPrint("error removing module: " + err.Error())
			os.Exit(1)
		}
	},
}

func removeModule(module string) error {
	moduleService := services.NewModuleService(shared.FilesPath, filesystem.NewFileSystem(afero.NewOsFs()))
	if err := shared.RunCommand("rm", "-rf", shared.FilesPath+"/"+module); err != nil {
		return err
	}
	err := moduleService.Remove(module)
	if err != nil {
		return err
	}
	return nil
}

// init sets up flags for the 'rm' subcommand and binds them to variables.
func init() {
	// Module name
	ModuleRmCmd.Flags().StringVarP(&Module, "module", "m", "", "module to remove")
}
