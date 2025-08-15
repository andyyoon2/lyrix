package cmd

import (
	"fmt"

	"github.com/andyyoon2/lyrix/internal/state"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current song and position",
	Long:  `Display the current active song, position in lyrics, and queue status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load()
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}
		
		if !s.Active {
			fmt.Println("lyrix is not active")
			fmt.Println("Get started with: lyrix start <artist> <title>")
			return nil
		}
		
		if s.CurrentSong == nil {
			fmt.Println("No current song")
			return nil
		}
		
		fmt.Printf("Active: %t\n", s.Active)
		fmt.Printf("\nCurrent Song:\n")
		fmt.Printf("  Title: %s\n", s.CurrentSong.Title)
		fmt.Printf("  Artist: %s\n", s.CurrentSong.Artist)
		fmt.Printf("  Progress: %d/%d lines\n", s.CurrentSong.Position, len(s.CurrentSong.Lyrics))
		
		if s.CurrentSong.Position < len(s.CurrentSong.Lyrics) {
			fmt.Printf("  Next line: %s\n", s.CurrentSong.Lyrics[s.CurrentSong.Position])
		} else {
			fmt.Printf("  Next line: (end of song)\n")
		}
		
		if len(s.Queue) > 0 {
			fmt.Printf("\nQueue (%d songs):\n", len(s.Queue))
			for i, song := range s.Queue {
				fmt.Printf("  %d. %s by %s (%d lines)\n", i+1, song.Title, song.Artist, len(song.Lyrics))
			}
		} else {
			fmt.Println("\nQueue: Empty")
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}