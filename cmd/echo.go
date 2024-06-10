/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("echo called")

		reader := bufio.NewReader(os.Stdin)
		line := make([]rune, 0, 64)
		for {
			r, _, _ := reader.ReadRune()
			// TODO: handle ANSI code
			fmt.Println("reade a rune", r, "except", rune('\n'))
			if r == '\n' {
				fmt.Fprintln(os.Stdout, strings.TrimRight(string(line), "\n\r"))
				line = line[:0]
				continue
			} else if r == 0x1B {
				second, _ := reader.ReadByte()
				third, _ := reader.ReadByte()
				os.Stdout.Write([]byte{0x1B, second, third})
				continue
			}
			line = append(line, r)
		}
	},
}

func init() {
	rootCmd.AddCommand(echoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// echoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// echoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
