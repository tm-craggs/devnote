package utils

import (
	"fmt"
	"os"

	"github.com/tm-craggs/devnote/internal/config"
	"gopkg.in/yaml.v3"
)

// CreateDefaultConfig creates a new devnote config with default values, and writes it to the specified filepath
func CreateDefaultConfig(configPath string) error {

	//TODO: Add comments and explanations to the default config

	defaultConfig := config.DevnoteConfig{
		NotesPath:     "./devnotes",
		FileExtension: ".md",

		Editor:     "auto",
		EditorArgs: []string{},

		OpenAfterNew: true,

		AutoCommit:   false,
		AutoAdd:      false,
		GitCommitMsg: nil,
	}

	// marshal to YAML
	yamlData, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("cannot marshal devnote config to yaml: %w", err)
	}

	// write to file
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		return fmt.Errorf("cannot write devnote config to file: %w", err)
	}

	return nil

}
