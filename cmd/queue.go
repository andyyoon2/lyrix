package cmd

import (
	"fmt"
	"strings"

	"github.com/andyyoon/git-lyrics/internal/lyrics"
	"github.com/andyyoon/git-lyrics/internal/state"
	"github.com/spf13/cobra"
)

var queueCmd = &cobra.Command{
	Use:   "queue <artist> <title>",
	Short: "Add a song to the queue",
	Long:  `Add a song to the queue. It will be used after the current song's lyrics are exhausted.`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		artist := args[0]
		title := strings.Join(args[1:], " ")
		
		fmt.Printf("Fetching lyrics for '%s' by %s...\n", title, artist)
		
		songLyrics, err := lyrics.FetchLyrics(artist, title)
		if err != nil {
			return fmt.Errorf("failed to fetch lyrics: %w", err)
		}
		
		s, err := state.Load()
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}
		
		newSong := state.Song{
			Title:    title,
			Artist:   artist,
			Lyrics:   songLyrics,
			Position: 0,
		}
		
		s.Queue = append(s.Queue, newSong)
		
		if err := s.Save(); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		
		fmt.Printf("✓ Added '%s' by %s to queue\n", title, artist)
		fmt.Printf("  Queue position: %d\n", len(s.Queue))
		fmt.Printf("  Song has %d lines\n", len(songLyrics))
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(queueCmd)
}