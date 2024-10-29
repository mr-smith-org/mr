package exec

import (
	execModule "github.com/kuma-framework/kuma/v2/cmd/commands/exec/module"
	execRun "github.com/kuma-framework/kuma/v2/cmd/commands/exec/run"
	"github.com/spf13/cobra"
)

var ExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Manage Kuma execs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ExecCmd.AddCommand(execRun.ExecCmd)
	ExecCmd.AddCommand(execModule.ExecModuleCmd)
}
