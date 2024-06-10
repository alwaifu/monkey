/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/alwaifu/monkey/pkg/interpreter"
	"github.com/alwaifu/monkey/pkg/vm"

	"github.com/spf13/cobra"
)

var runVersion *int = new(int)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if *runVersion == 1 {
			interpreter.Start(os.Stdin, os.Stdout)
		}
		vm.Start(os.Stdin, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.Flags().IntVar(runVersion, "ver", 2, "run version, version 1 will interprete ast tree directly, version 2 will use virtual machine")
}
