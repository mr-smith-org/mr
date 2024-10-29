package module

import (
	"os"

	"github.com/kuma-framework/kuma/v2/cmd/shared"
	"github.com/kuma-framework/kuma/v2/internal/services"
	"github.com/kuma-framework/kuma/v2/pkg/filesystem"
	"github.com/kuma-framework/kuma/v2/pkg/style"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	Module string
)

// Add a Kuma module from a GitHub repository
var ModuleRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a Kuma module",
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
	moduleService := services.NewModuleService(shared.KumaFilesPath, filesystem.NewFileSystem(afero.NewOsFs()))
	if err := shared.RunCommand("rm", "-rf", shared.KumaFilesPath+"/"+module); err != nil {
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
