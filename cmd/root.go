package cmd

import (
	"fmt"
	"os"

	"github.com/mr-smith-org/mr/cmd/commands/create"
	"github.com/mr-smith-org/mr/cmd/commands/modify"
	"github.com/mr-smith-org/mr/cmd/commands/module"
	"github.com/mr-smith-org/mr/cmd/commands/run"
	"github.com/mr-smith-org/mr/cmd/shared"
	"github.com/mr-smith-org/mr/internal/verbose"
	"github.com/spf13/cobra"
)

const version = "v1.3.1"

var (
	showVersion bool
)

var rootCmd = &cobra.Command{
	Use:  "mr",
	Long: fmt.Sprintf("Welcome to Mr Smith! \n A powerful CLI for generating project scaffolds based on Go templates. \n version: %s", version),
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println(version)
			os.Exit(0)
		}

		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(0)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&verbose.Verbose, "verbose", "", false, "Enable verbose logs")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print the version number of Mr. Smith")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		shared.RunCommand("go", "install", "github.com/mr-smith-org/mr@latest")
	}
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(module.ModuleCmd)
	rootCmd.AddCommand(run.RunCmd)
	rootCmd.AddCommand(modify.ModifyCmd)
}
