package cmd

import (
	"github.com/andyyoon2/lyrix/internal/hook"
	"github.com/spf13/cobra"
)

var hookExecCmd = &cobra.Command{
	Use:    "hook-exec <commit-msg-file>",
	Short:  "Execute the git hook (internal use)",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return hook.ExecuteHook(args[0])
	},
}

func init() {
	rootCmd.AddCommand(hookExecCmd)
}