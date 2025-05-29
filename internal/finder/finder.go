package finder

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jesalx/tmux-sessionizer/internal/config"
	"github.com/jesalx/tmux-sessionizer/internal/tmux"
)

func FindAll() []string {
	var items []string

	if tmux.IsInSession() {
		currentSession := tmux.GetCurrentSession()
		sessions := tmux.GetSessions()
		for _, session := range sessions {
			if session != currentSession {
				items = append(items, "[TMUX] "+session)
			}
		}
	} else {
		sessions := tmux.GetSessions()
		for _, session := range sessions {
			items = append(items, "[TMUX] "+session)
		}
	}

	cfg := config.Get()
	for _, searchPath := range cfg.SearchPaths {
		if _, err := os.Stat(searchPath.Path); err == nil {
			foundDirs := findDirectoriesInPath(searchPath.Path, searchPath.Depth)
			items = append(items, foundDirs...)
		}
	}

	return items
}

func findDirectoriesInPath(rootPath string, maxDepth int) []string {
	var dirs []string

	if maxDepth == 0 {
		return dirs
	}

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if path == rootPath {
			return nil
		}

		if d.Name() == ".git" {
			return filepath.SkipDir
		}

		relPath, _ := filepath.Rel(rootPath, path)
		depth := strings.Count(relPath, string(filepath.Separator)) + 1

		if depth > maxDepth {
			return filepath.SkipDir
		}

		if d.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})

	if err != nil {
		return []string{}
	}

	return dirs
}
