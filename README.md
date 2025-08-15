# lyrix

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
go build -o lyrix .
sudo mv lyrix /usr/local/bin/
```

## Quick Start

1. Install the git hook in your repository:
```bash
cd your-repo
lyrix install
```

2. Start with a song:
```bash
lyrix start "The Beatles" "Hey Jude"
```

3. Make commits normally:
```bash
git add .
git commit -m "Fixed bug in authentication"
```

Your commit will have the next lyric line as the title, with your message as the description.

## Commands

### `lyrix start <artist> <title>`
Start using lyrics from a song. Fetches lyrics and sets as active.

### `lyrix queue <artist> <title>`
Add a song to the queue. Will be used after current song ends.

### `lyrix status`
Show current song, position, and queue.

### `lyrix next`
Skip to the next song in queue.

### `lyrix stop`
Deactivate lyrix and clear all songs.

### `lyrix install`
Install git hook in current repository.

### `lyrix uninstall`
Remove git hook from current repository.

## How It Works

1. lyrix fetches lyrics from LRCLIB API
2. State is stored in `~/.lyrix/state.json`
3. A git `prepare-commit-msg` hook intercepts commits
4. Each commit uses the next line, preserving your message

## Example Workflow

```bash
# Install hook in your repo
lyrix install

# Start with a classic
lyrix start "Queen" "Bohemian Rhapsody"

# Queue another song
lyrix queue "David Bowie" "Space Oddity"

# Check status
lyrix status

# Make commits - each will use the next line
git add feature.go
git commit -m "Add new feature"
# Commit title: "Is this the real life?"

git add bugfix.go
git commit -m "Fix critical bug"
# Commit title: "Is this just fantasy?"

# When done
lyrix stop
```

## Notes

- Requires internet connection to fetch lyrics
- Empty lines in lyrics are skipped
- When a song ends, automatically moves to next in queue
- Hook only affects regular commits (not merges/rebases)