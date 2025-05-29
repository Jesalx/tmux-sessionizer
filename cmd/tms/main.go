package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jesalx/tmux-sessionizer/internal/session"
	"github.com/jesalx/tmux-sessionizer/internal/tmux"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tms [directory]",
	Short: "Tmux sessionizer - create and switch between tmux sessions",
	Long:  "A tmux session manager that helps you quickly create and switch between tmux sessions based on directories.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var target string
		if len(args) == 1 {
			target = args[0]
		}
		return session.Run(target)
	},
}

var newCmd = &cobra.Command{
	Use:   "new [session-name]",
	Short: "Create a new tmux session",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var sessionName string
		if len(args) == 1 {
			sessionName = args[0]
		} else {
			fmt.Print("Enter session name: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				sessionName = scanner.Text()
			}
		}

		if sessionName == "" {
			return fmt.Errorf("session name cannot be empty")
		}

		return session.CreateNew(sessionName)
	},
}

var renameCmd = &cobra.Command{
	Use:   "rename [new-name]",
	Short: "Rename the current tmux session",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if we're in a tmux session first
		if !tmux.IsInSession() {
			return fmt.Errorf("not inside a tmux session. Cannot rename session")
		}

		var newName string
		if len(args) == 1 {
			newName = args[0]
		} else {
			fmt.Print("Enter new session name: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				newName = scanner.Text()
			}
		}

		if newName == "" {
			return fmt.Errorf("new session name cannot be empty")
		}

		if err := tmux.RenameSession(newName); err != nil {
			return fmt.Errorf("failed to rename session: %w", err)
		}

		fmt.Printf("Session renamed to '%s'\n", newName)
		return nil
	},
}

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill the current tmux session",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !tmux.IsInSession() {
			return fmt.Errorf("not inside a tmux session")
		}

		return tmux.KillSession()
	},
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exit current tmux session without killing it",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !tmux.IsInSession() {
			return fmt.Errorf("not inside a tmux session")
		}

		return tmux.DetachSession()
	},
}

var cloneCmd = &cobra.Command{
	Use:   "clone [repository-url]",
	Short: "Clone a git repository and create a tmux session in it",
	Long:  "Clone a git repository to a specified directory and automatically create a tmux session in the cloned repository.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var repoURL string
		if len(args) == 1 {
			repoURL = args[0]
		} else {
			fmt.Print("Enter repository: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				repoURL = scanner.Text()
			}
		}

		if repoURL == "" {
			return fmt.Errorf("repository cannot be empty")
		}

		return session.CloneAndRun(repoURL)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(exitCmd)
	rootCmd.AddCommand(killCmd)
	rootCmd.AddCommand(cloneCmd)
}

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}
