package cmd

import (
	"fmt"

	"github.com/andyyoon2/lyrix/internal/state"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop using song lyrics for commits",
	Long:  `Deactivate lyrix and clear the current song and queue.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load()
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}
		
		s.Reset()
		
		if err := s.Save(); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		
		fmt.Println("✓ lyrix stopped and cleared")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}