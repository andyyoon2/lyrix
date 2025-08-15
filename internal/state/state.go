package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Song struct {
	Title    string   `json:"title"`
	Artist   string   `json:"artist"`
	Lyrics   []string `json:"lyrics"`
	Position int      `json:"position"`
}

type State struct {
	Active      bool   `json:"active"`
	CurrentSong *Song  `json:"current_song,omitempty"`
	Queue       []Song `json:"queue"`
}

func GetStatePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".git-lyrics", "state.json")
}

func Load() (*State, error) {
	statePath := GetStatePath()
	
	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Active: false, Queue: []Song{}}, nil
		}
		return nil, err
	}
	
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	
	return &state, nil
}

func (s *State) Save() error {
	statePath := GetStatePath()
	dir := filepath.Dir(statePath)
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(statePath, data, 0644)
}

func (s *State) GetNextLine() (string, error) {
	if !s.Active || s.CurrentSong == nil {
		return "", fmt.Errorf("no active song")
	}
	
	if s.CurrentSong.Position >= len(s.CurrentSong.Lyrics) {
		if len(s.Queue) > 0 {
			s.CurrentSong = &s.Queue[0]
			s.Queue = s.Queue[1:]
			s.CurrentSong.Position = 0
			return s.GetNextLine()
		}
		s.Active = false
		s.CurrentSong = nil
		return "", fmt.Errorf("no more lyrics available")
	}
	
	line := s.CurrentSong.Lyrics[s.CurrentSong.Position]
	s.CurrentSong.Position++
	return line, nil
}

func (s *State) Reset() {
	s.Active = false
	s.CurrentSong = nil
	s.Queue = []Song{}
}