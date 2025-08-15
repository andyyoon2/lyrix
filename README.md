# git-lyrics

A CLI tool that uses song lyrics as your git commit message titles, line by line.

## Features

- Uses song lyrics progressively as commit titles
- Preserves your commit messages as descriptions
- Queue multiple songs
- Simple, performant Go binary
- Seamless git integration via hooks

## Installation

### From Source

```bash
go build -o git-lyrics .
sudo mv git-lyrics /usr/local/bin/
```

## Quick Start

1. Install the git hook in your repository:
```bash
cd your-repo
git-lyrics install
```

2. Start with a song:
```bash
git-lyrics start "The Beatles" "Hey Jude"
```

3. Make commits normally:
```bash
git add .
git commit -m "Fixed bug in authentication"
```

Your commit will have the next lyric line as the title, with your message as the description.

## Commands

### `git-lyrics start <artist> <title>`
Start using lyrics from a song. Fetches lyrics and sets as active.

### `git-lyrics queue <artist> <title>`
Add a song to the queue. Will be used after current song ends.

### `git-lyrics status`
Show current song, position, and queue.

### `git-lyrics next`
Skip to the next song in queue.

### `git-lyrics stop`
Deactivate git-lyrics and clear all songs.

### `git-lyrics install`
Install git hook in current repository.

### `git-lyrics uninstall`
Remove git hook from current repository.

## How It Works

1. git-lyrics fetches lyrics from lyrics.ovh API
2. State is stored in `~/.git-lyrics/state.json`
3. A git `prepare-commit-msg` hook intercepts commits
4. Each commit uses the next line, preserving your message

## Example Workflow

```bash
# Install hook in your repo
git-lyrics install

# Start with a classic
git-lyrics start "Queen" "Bohemian Rhapsody"

# Queue another song
git-lyrics queue "David Bowie" "Space Oddity"

# Check status
git-lyrics status

# Make commits - each will use the next line
git add feature.go
git commit -m "Add new feature"
# Commit title: "Is this the real life?"

git add bugfix.go
git commit -m "Fix critical bug"
# Commit title: "Is this just fantasy?"

# When done
git-lyrics stop
```

## Notes

- Requires internet connection to fetch lyrics
- Empty lines in lyrics are skipped
- When a song ends, automatically moves to next in queue
- Hook only affects regular commits (not merges/rebases)