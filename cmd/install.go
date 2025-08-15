package cmd

import (
	"fmt"

	"github.com/andyyoon/git-lyrics/internal/hook"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install git hook for automatic lyric insertion",
	Long:  `Install the prepare-commit-msg git hook that automatically inserts song lyrics as commit titles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := hook.InstallHook(); err != nil {
			return err
		}
		
		fmt.Println("✓ Git hook installed successfully")
		fmt.Println("  Your commits will now use song lyrics as titles when git-lyrics is active")
		fmt.Println("  Use 'git-lyrics start <artist> <title>' to begin")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}