package cmd

import (
	"fmt"

	"github.com/andyyoon2/lyrix/internal/state"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all song state",
	Long:  `Clear the current song and queue, resetting lyrix to a fresh state.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load()
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}
		
		s.Reset()
		
		if err := s.Save(); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		
		fmt.Println("✓ lyrix state cleared")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}