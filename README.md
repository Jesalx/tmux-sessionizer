# tmux-sessionizer

A tmux session manager inspired by [ThePrimeagen's script](https://github.com/ThePrimeagen/tmux-sessionizer). Create and switch between tmux sessions based on your project directories.

## Prerequisites

Make sure you have these tools installed:

- `tmux`
- `fzf`
- `git`

## Installation

```bash
go install github.com/jesalx/tmux-sessionizer/cmd/tms@latest
```

Make sure your Go binary path is in your `$PATH`. You can add this to your shell profile:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

| Command                      | Description                                     |
| ---------------------------- | ----------------------------------------------- |
| `tms`                        | Launch interactive session selector             |
| `tms [directory]`            | Create/switch to session for specific directory |
| `tms new [session-name]`     | Create a new session in current directory       |
| `tms rename [new-name]`      | Rename current session                          |
| `tms clone <repository-url>` | Clone repo and create session                   |
| `tms kill`                   | Kill current session                            |
| `tms exit`                   | Exit (detach from) current session              |

## Configuration

Create a config file at `~/.config/tmux-sessionizer/config.yaml`:

```yaml
# Maximum depth to search in directories
max_depth: 2

# Search paths for projects
search_paths:
  - path: "~/Developer"
    depth: 1
  - path: "~/.config"
    depth: 1
  - path: "~/work"
    depth: 2
```

### Configuration Options

- `max_depth`: Default depth for directory searching (default: 1)
- `search_paths`: Array of paths to search for projects
  - `path`: Directory path (supports `~/` expansion)
  - `depth`: How deep to search in this path (optional, uses max_depth if not set)
