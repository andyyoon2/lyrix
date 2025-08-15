package cmd

import (
	"fmt"

	"github.com/andyyoon/git-lyrics/internal/state"
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Skip to the next song in queue",
	Long:  `Skip the current song and move to the next song in the queue.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load()
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}
		
		if !s.Active {
			return fmt.Errorf("git-lyrics is not active")
		}
		
		if len(s.Queue) == 0 {
			return fmt.Errorf("no songs in queue")
		}
		
		s.CurrentSong = &s.Queue[0]
		s.Queue = s.Queue[1:]
		s.CurrentSong.Position = 0
		
		if err := s.Save(); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		
		fmt.Printf("✓ Switched to '%s' by %s\n", s.CurrentSong.Title, s.CurrentSong.Artist)
		fmt.Printf("  Next line: %s\n", s.CurrentSong.Lyrics[0])
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}