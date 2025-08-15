package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lyrix",
	Short: "A tool to use song lyrics as git commit messages",
	Long: `lyrix uses song lyrics line by line as your git commit titles.
Each commit will use the next line of the current song as the title,
while preserving your commit message as the description.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
}