package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/andyyoon2/lyrix/internal/state"
)

const hookScript = `#!/bin/sh
# lyrix prepare-commit-msg hook

COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2

# Only process regular commits (not merges, squashes, etc.)
if [ -z "$COMMIT_SOURCE" ] || [ "$COMMIT_SOURCE" = "message" ]; then
    lyrix hook-exec "$COMMIT_MSG_FILE"
fi
`

func GetGitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository")
	}
	return strings.TrimSpace(string(output)), nil
}

func InstallHook() error {
	gitDir, err := GetGitDir()
	if err != nil {
		return err
	}
	
	hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")
	
	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		// Read existing hook to check if it's ours
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return fmt.Errorf("failed to read existing hook: %w", err)
		}
		
		if !strings.Contains(string(content), "lyrix") {
			return fmt.Errorf("prepare-commit-msg hook already exists and is not managed by lyrix")
		}
		
		return fmt.Errorf("lyrix hook is already installed")
	}
	
	// Create hooks directory if it doesn't exist
	hooksDir := filepath.Dir(hookPath)
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}
	
	// Write hook script
	if err := os.WriteFile(hookPath, []byte(hookScript), 0755); err != nil {
		return fmt.Errorf("failed to write hook: %w", err)
	}
	
	return nil
}

func UninstallHook() error {
	gitDir, err := GetGitDir()
	if err != nil {
		return err
	}
	
	hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")
	
	// Check if hook exists
	content, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no hook to uninstall")
		}
		return fmt.Errorf("failed to read hook: %w", err)
	}
	
	// Only remove if it's our hook
	if !strings.Contains(string(content), "lyrix") {
		return fmt.Errorf("existing hook is not managed by lyrix")
	}
	
	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("failed to remove hook: %w", err)
	}
	
	return nil
}

func ExecuteHook(commitMsgFile string) error {
	// Load state
	s, err := state.Load()
	if err != nil {
		return nil // Silently fail to not break commits
	}
	
	if !s.Active || s.CurrentSong == nil {
		return nil // Nothing to do
	}
	
	// Get next line
	nextLine, err := s.GetNextLine()
	if err != nil {
		// If no more lyrics, just use regular commit
		if err.Error() == "no more lyrics available" {
			fmt.Println("No more lyrics available. Queue a new song with 'lyrix queue'")
		}
		return nil
	}
	
	// Save state immediately after getting the line
	if err := s.Save(); err != nil {
		return nil // Silently fail
	}
	
	// Read existing commit message
	content, err := os.ReadFile(commitMsgFile)
	if err != nil {
		return nil
	}
	
	// Parse existing message
	lines := strings.Split(string(content), "\n")
	var userMessage []string
	
	// Skip empty lines and comments at the beginning
	startIdx := 0
	for i, line := range lines {
		if line != "" && !strings.HasPrefix(line, "#") {
			startIdx = i
			break
		}
	}
	
	// Collect non-comment lines as user message
	for i := startIdx; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "#") {
			userMessage = append(userMessage, lines[i])
		}
	}
	
	// Build new commit message
	var newMessage strings.Builder
	newMessage.WriteString(nextLine)
	newMessage.WriteString("\n\n")
	
	// Add user message if it exists and is not empty
	userMsg := strings.TrimSpace(strings.Join(userMessage, "\n"))
	if userMsg != "" {
		newMessage.WriteString(userMsg)
		newMessage.WriteString("\n")
	}
	
	// Add comments back
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			newMessage.WriteString("\n")
			newMessage.WriteString(line)
		}
	}
	
	// Write back to file
	if err := os.WriteFile(commitMsgFile, []byte(newMessage.String()), 0644); err != nil {
		return nil
	}
	
	return nil
}