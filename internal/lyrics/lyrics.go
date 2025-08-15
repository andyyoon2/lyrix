package lyrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type LRCLibResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

func FetchLyrics(artist, title string) ([]string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Use LRCLIB API - completely free, no registration required
	baseURL := "https://lrclib.net/api/search"
	
	// Build query parameters
	params := url.Values{}
	params.Add("artist_name", artist)
	params.Add("track_name", title)
	
	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	// Retry logic for transient errors
	var resp *http.Response
	var err error
	maxRetries := 3
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Create request with context for additional timeout control
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		
		// Add user agent as recommended by LRCLIB
		req.Header.Set("User-Agent", "lyrix/1.0 (https://github.com/andyyoon2/lyrix)")
		
		resp, err = client.Do(req)
		if err == nil {
			break // Success!
		}
		
		// Check if we should retry
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("request timed out after 10 seconds")
		}
		
		// For other errors, retry after a short delay
		if attempt < maxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		
		return nil, fmt.Errorf("failed to fetch lyrics after %d attempts: %w", maxRetries, err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}
	
	// Read response with size limit
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // 1MB limit
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Parse response as array
	var results []LRCLibResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("no lyrics found for '%s' by %s", title, artist)
	}
	
	// Use the first result
	result := results[0]
	
	if result.Instrumental {
		return nil, fmt.Errorf("'%s' by %s is an instrumental track", title, artist)
	}
	
	if result.PlainLyrics == "" {
		return nil, fmt.Errorf("no lyrics available for '%s' by %s", title, artist)
	}
	
	lines := parseLyrics(result.PlainLyrics)
	
	return lines, nil
}

func parseLyrics(lyrics string) []string {
	lines := strings.Split(lyrics, "\n")
	
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and metadata lines
		if line != "" && !strings.HasPrefix(line, "Paroles de la chanson") {
			cleanLines = append(cleanLines, line)
		}
	}
	
	return cleanLines
}