package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type SearchPath struct {
	Path  string `yaml:"path"`
	Depth int    `yaml:"depth,omitempty"`
}

type Config struct {
	SearchPaths []SearchPath `yaml:"search_paths,omitempty"`
	MaxDepth    int          `yaml:"max_depth,omitempty"`
}

var (
	globalConfig Config
	configDir    string
	configFile   string
)

func init() {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		homeDir, _ := os.UserHomeDir()
		xdgConfigHome = filepath.Join(homeDir, ".config")
	}
	configDir = filepath.Join(xdgConfigHome, "tmux-sessionizer")
	configFile = filepath.Join(configDir, "config.yaml")

	loadConfig()
}

func Get() *Config {
	return &globalConfig
}

func loadConfig() {
	globalConfig.SearchPaths = []SearchPath{}
	globalConfig.MaxDepth = 1

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Warning: could not read config file: %v\n", err)
		return
	}

	err = yaml.Unmarshal(data, &globalConfig)
	if err != nil {
		fmt.Printf("Warning: could not parse config file: %v\n", err)
		return
	}

	homeDir, _ := os.UserHomeDir()
	expandPaths(globalConfig.SearchPaths, homeDir)
}

func expandPaths(paths []SearchPath, homeDir string) {
	for i := range paths {
		if strings.HasPrefix(paths[i].Path, "~/") {
			paths[i].Path = filepath.Join(homeDir, paths[i].Path[2:])
		}
		if paths[i].Depth == 0 {
			paths[i].Depth = globalConfig.MaxDepth
		}
	}
}
