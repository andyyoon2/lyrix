package cmd

import (
	"fmt"
	"strings"

	"github.com/andyyoon2/lyrix/internal/lyrics"
	"github.com/andyyoon2/lyrix/internal/state"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <artist> <title>",
	Short: "Start using lyrics from a song",
	Long:  `Start using lyrics from a specified song for your git commit messages.`,
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
		
		s.Active = true
		s.CurrentSong = &state.Song{
			Title:    title,
			Artist:   artist,
			Lyrics:   songLyrics,
			Position: 0,
		}
		
		if err := s.Save(); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		
		fmt.Printf("✓ Started using lyrics from '%s' by %s\n", title, artist)
		fmt.Printf("  Found %d lines\n", len(songLyrics))
		if len(songLyrics) > 0 {
			fmt.Printf("  Next line: %s\n", songLyrics[0])
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}