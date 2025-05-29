package session

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jesalx/tmux-sessionizer/internal/finder"
	"github.com/jesalx/tmux-sessionizer/internal/formatter"
	"github.com/jesalx/tmux-sessionizer/internal/tmux"
)

func SanityCheck() error {
	if _, err := exec.LookPath("tmux"); err != nil {
		return fmt.Errorf("tmux is not installed. Please install it first")
	}
	if _, err := exec.LookPath("fzf"); err != nil {
		return fmt.Errorf("fzf is not installed. Please install it first")
	}
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed. Please install it first")
	}

	return nil
}

func Run(selected string) error {
	if err := SanityCheck(); err != nil {
		return err
	}

	if selected == "" {
		dirs := finder.FindAll()
		if len(dirs) == 0 {
			return fmt.Errorf("no directories or sessions found")
		}

		cmd := exec.Command("fzf")
		cmd.Stdin = strings.NewReader(strings.Join(dirs, "\n"))
		output, err := cmd.Output()
		if err != nil {
			return nil // User cancelled or error
		}
		selected = strings.TrimSpace(string(output))
	}

	if selected == "" {
		return nil
	}

	tmuxSessionRegex := regexp.MustCompile(`^\[TMUX\] (.+)$`)
	if matches := tmuxSessionRegex.FindStringSubmatch(selected); len(matches) > 1 {
		selected = matches[1]
		return tmux.SwitchToSession(selected)
	}

	sessionName := filepath.Base(selected)
	sessionName = formatter.FormatSessionName(sessionName)

	if !tmux.IsInSession() && !tmux.IsRunning() {
		if err := tmux.CreateSession(sessionName, selected); err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
	}

	if !tmux.HasSession(sessionName) {
		if err := tmux.CreateSession(sessionName, selected); err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
	}

	return tmux.SwitchToSession(sessionName)
}

func CreateNew(sessionName string) error {
	if err := SanityCheck(); err != nil {
		return err
	}

	sessionName = formatter.FormatSessionName(sessionName)

	if tmux.HasSession(sessionName) {
		fmt.Printf("Session '%s' already exists\n", sessionName)
		return tmux.SwitchToSession(sessionName)
	}

	currentDir, _ := os.Getwd()
	if err := tmux.CreateSession(sessionName, currentDir); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return tmux.SwitchToSession(sessionName)
}

func CloneAndRun(repoURL string) error {
	if err := SanityCheck(); err != nil {
		return err
	}

	repoName, err := extractRepoName(repoURL)
	if err != nil {
		return fmt.Errorf("failed to parse repository URL: %w", err)
	}

	if _, err := os.Stat(repoName); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", repoURL)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	currentDir, _ := os.Getwd()
	cloneDir := filepath.Join(currentDir, repoName)

	sessionName := formatter.FormatSessionName(repoName)

	if !tmux.HasSession(sessionName) {
		if err := tmux.CreateSession(sessionName, cloneDir); err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
	}

	return tmux.SwitchToSession(sessionName)
}

func extractRepoName(repoURL string) (string, error) {
	repoURL = strings.TrimSpace(repoURL)

	if strings.HasPrefix(repoURL, "git@") {
		parts := strings.Split(repoURL, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid SSH URL format")
		}
		pathPart := parts[1]
		repoName := path.Base(pathPart)
		return strings.TrimSuffix(repoName, ".git"), nil
	} else {
		u, err := url.Parse(repoURL)
		if err != nil {
			return "", err
		}
		repoName := path.Base(u.Path)
		return strings.TrimSuffix(repoName, ".git"), nil
	}
}
