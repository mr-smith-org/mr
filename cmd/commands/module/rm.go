// module.go
//
// Package get defines the 'module' subcommand for the Kuma CLI.
// It handles generating project scaffolds based on Go templates.
package module

import (
	"os"

	"github.com/arthurbcp/kuma/cmd/shared"
	"github.com/arthurbcp/kuma/internal/services"
	"github.com/arthurbcp/kuma/pkg/filesystem"
	"github.com/arthurbcp/kuma/pkg/style"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	Module             string
	RemoveGitSubmodule bool
)

// Add a Kuma module from a GitHub repository
var ModuleRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a Kuma module",
	Run: func(cmd *cobra.Command, args []string) {
		if Module != "" {
			err := RemoveModule(Module)
			if err != nil {
				style.ErrorPrint("error removing module: " + err.Error())
				os.Exit(1)
			}
		}
	},
}

func RemoveModule(module string) error {
	if RemoveGitSubmodule {
		if err := removeGitSubmodule(module); err != nil {
			return err
		}
	}
	// Remove the module from the kuma-modules.yaml file
	moduleService := services.NewModuleService(shared.KumaFilesPath, filesystem.NewFileSystem(afero.NewOsFs()))
	err := moduleService.Remove(module)
	if err != nil {
		return err
	}
	return nil
}

// remove git submodule removes a submodule from Kuma
func removeGitSubmodule(module string) error {
	// Full path to the submodule inside _kuma
	fullSubmodulePath := shared.KumaFilesPath + "/" + module

	// 1. Remove submodule config from .git/config
	if err := shared.RunCommand("git", "config", "--remove-section", "submodule."+fullSubmodulePath); err != nil {
		return err
	}

	// 2. Remove submodule entry from .gitmodules if it exists
	if err := shared.RunCommand("git", "config", "-f", ".gitmodules", "--remove-section", "submodule."+fullSubmodulePath); err != nil {
		return err
	}

	// 3. Remove the submodule from git cache
	if err := shared.RunCommand("git", "rm", "--cached", fullSubmodulePath); err != nil {
		return err
	}

	return nil
}

// init sets up flags for the 'rm' subcommand and binds them to variables.
func init() {
	// Module name
	ModuleRmCmd.Flags().StringVarP(&Module, "module", "m", "", "module to remove")
	ModuleRmCmd.Flags().BoolVarP(&RemoveGitSubmodule, "rm-git-submodule", "", false, "remove git submodule")
}