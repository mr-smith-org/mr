package cmd

import (
	"fmt"
	"os"

	"github.com/mr-smith-org/mr/cmd/commands/create"
	execRun "github.com/mr-smith-org/mr/cmd/commands/exec"
	"github.com/mr-smith-org/mr/cmd/commands/modify"
	"github.com/mr-smith-org/mr/cmd/commands/module"
	"github.com/mr-smith-org/mr/internal/debug"
	"github.com/spf13/cobra"
)

const (
	UnicodeLogo = `
	
	`
)

var rootCmd = &cobra.Command{
	Use:  "mr",
	Long: fmt.Sprintf("%s \n\n Welcome to Mr Smith! \n A powerful CLI for generating project scaffolds based on Go templates.", UnicodeLogo),
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
