package cmd

import (
	"fmt"
	"os"

	"github.com/kuma-framework/kuma/v2/cmd/commands/create"
	execRun "github.com/kuma-framework/kuma/v2/cmd/commands/exec"
	"github.com/kuma-framework/kuma/v2/cmd/commands/modify"
	"github.com/kuma-framework/kuma/v2/cmd/commands/module"
	"github.com/kuma-framework/kuma/v2/internal/debug"
	"github.com/spf13/cobra"
)

const (
	UnicodeLogo = `
	
	`
)

var rootCmd = &cobra.Command{
	Use:  "kuma",
	Long: fmt.Sprintf("%s \n\n Welcome to Kuma! \n A powerful CLI for generating project scaffolds based on Go templates.", UnicodeLogo),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&debug.Debug, "debug", "", false, "Enable debug mode")
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(module.ModuleCmd)
	rootCmd.AddCommand(execRun.ExecCmd)
	rootCmd.AddCommand(modify.ModifyCmd)
}
