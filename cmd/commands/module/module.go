package module

import "github.com/spf13/cobra"

var ModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Manage modules",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ModuleCmd.AddCommand(ModuleAddCmd)
	ModuleCmd.AddCommand(ModuleRmCmd)
}
