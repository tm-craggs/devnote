/*
   Copyright 2025 Tom Craggs

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tm-craggs/devnote/internal/utils"
)

const (
	defaultNotesPath = "devnotes"
	devnoteDirName   = ".devnote"
	configFileName   = "devnote.yaml"
)

type initFlags struct {
	path string
}

// getInitFlags parses command flags with error handling
func getInitFlags(cmd *cobra.Command) (initFlags, error) {
	var flags initFlags
	var err error

	flags.path, err = cmd.Flags().GetString("path")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --path flag: %w", err)
	}

	return flags, nil
}

// expandPath handles ~ expansion and converts relative paths to absolute
func expandPath(notePath, projectRoot string) (string, error) {
	// Use default if empty
	if notePath == "" {
		notePath = defaultNotesPath
	}

	// Expand ~ to home directory
	if len(notePath) >= 2 && notePath[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not resolve home directory: %w", err)
		}
		notePath = filepath.Join(home, notePath[2:])
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(notePath) {
		notePath = filepath.Join(projectRoot, notePath)
	}

	return notePath, nil
}

// ensureNotesDirectory creates the notes directory if it doesn't exist
func ensureNotesDirectory(notesPath string) error {
	info, err := os.Stat(notesPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(notesPath, 0755); err != nil {
			return fmt.Errorf("failed to create notes directory: %w", err)
		}
		fmt.Printf("Created notes directory at %s\n", notesPath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("error checking notes directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path exists and is not a directory: %s", notesPath)
	}
	return nil
}

// createDevnoteStructure creates the .devnote directory structure and files
func createDevnoteStructure(projectRoot string) error {
	devnoteDir := filepath.Join(projectRoot, devnoteDirName)
	stateDir := filepath.Join(devnoteDir, "state")
	templateDir := filepath.Join(devnoteDir, "templates")

	// Create .devnote directory
	if err := os.MkdirAll(devnoteDir, 0755); err != nil {
		return fmt.Errorf("failed to create devnote directory: %w", err)
	}

	// Create devnote.yaml config file
	// TODO: Import global config if flag set here
	configFilePath := filepath.Join(devnoteDir, configFileName)
	if err := utils.CreateDefaultConfig(configFilePath); err != nil {
		return fmt.Errorf("failed to create devnote config: %w", err)
	}

	// Create state and templates directories
	for _, dir := range []string{stateDir, templateDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create state files
	stateFiles := map[string]string{
		"current-template.txt": filepath.Join(stateDir, "current-template.txt"),
		"last_commit.txt":      filepath.Join(stateDir, "last_commit.txt"),
	}

	for name, path := range stateFiles {
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", name, err)
		}
	}

	fmt.Printf("Initialised devnotes config at %s\n", configFilePath)
	return nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise devnote in the current directory",
	Long: `The 'init' command initialises a devnote config in the current directory. 
You can specify the notePath for the config file using --path`,
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getInitFlags(cmd)
		if err != nil {
			return err
		}

		projectRoot, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not get current working directory: %w", err)
		}

		notesPath, err := expandPath(flags.path, projectRoot)
		if err != nil {
			return err
		}

		if err := ensureNotesDirectory(notesPath); err != nil {
			return err
		}

		return createDevnoteStructure(projectRoot)
	},
}

func init() {
	initCmd.Flags().StringP("path", "p", "", "Define a custom path for your notes folder.")
	rootCmd.AddCommand(initCmd)
}
