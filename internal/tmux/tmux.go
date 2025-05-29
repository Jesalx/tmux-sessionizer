package tmux

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jesalx/tmux-sessionizer/internal/formatter"
)

func IsRunning() bool {
	cmd := exec.Command("pgrep", "tmux")
	return cmd.Run() == nil
}

func IsInSession() bool {
	return os.Getenv("TMUX") != ""
}

func GetCurrentSession() string {
	cmd := exec.Command("tmux", "display-message", "-p", "#S")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func GetSessions() []string {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	var sessions []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		sessions = append(sessions, scanner.Text())
	}
	return sessions
}

func HasSession(sessionName string) bool {
	sessions := GetSessions()
	for _, session := range sessions {
		if session == sessionName {
			return true
		}
	}
	return false
}

func CreateSession(sessionName, directory string) error {
	cmd := exec.Command("tmux", "new-session", "-ds", sessionName, "-c", directory)
	return cmd.Run()
}

func SwitchToSession(sessionName string) error {
	if !IsInSession() {
		cmd := exec.Command("tmux", "attach-session", "-t", sessionName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else {
		cmd := exec.Command("tmux", "switch-client", "-t", sessionName)
		return cmd.Run()
	}
}

func RenameSession(newName string) error {
	if !IsInSession() {
		return fmt.Errorf("not inside a tmux session")
	}

	newName = formatter.FormatSessionName(newName)

	if HasSession(newName) {
		return fmt.Errorf("session '%s' already exists", newName)
	}

	cmd := exec.Command("tmux", "rename-session", newName)
	return cmd.Run()
}

func KillSession() error {
	if !IsInSession() {
		return fmt.Errorf("not inside a tmux session")
	}

	cmd := exec.Command("tmux", "kill-session")
	return cmd.Run()
}

func DetachSession() error {
	if !IsInSession() {
		return fmt.Errorf("not inside a tmux session")
	}

	cmd := exec.Command("tmux", "detach-client")
	return cmd.Run()
}
