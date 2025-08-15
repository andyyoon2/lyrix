package cmd

import (
	"fmt"

	"github.com/andyyoon2/lyrix/internal/hook"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall git hook",
	Long:  `Remove the prepare-commit-msg git hook installed by lyrix.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := hook.UninstallHook(); err != nil {
			return err
		}
		
		fmt.Println("✓ Git hook uninstalled successfully")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}