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
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type initFlags struct {
	path string
}

type devnotesConfig struct {
	NotesPath     string `yaml:"notesPath"`
	FileExtension string `yaml:"fileExtension"`
}

// helper function to parse flags with error handling
func getInitFlags(cmd *cobra.Command) (initFlags, error) {
	var flags initFlags
	var err error

	flags.path, err = cmd.Flags().GetString("path")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --path flag: %w", err)
	}

	return flags, nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise devnotes in the current directory",
	Long: `The 'init' command initialises a devnotes config in the current directory. You can specify the notePath for
the config file using --path`,
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getInitFlags(cmd)
		if err != nil {
			return err
		}

		projectRoot, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not get current working directory: %w", err)
		}

		notePath := flags.path
		if notePath == "" {
			notePath = "devnotes" // default
		}

		// Expand ~
		if len(notePath) >= 2 && notePath[:2] == "~/" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("could not resolve home directory: %w", err)
			}
			notePath = filepath.Join(home, notePath[2:])
		}

		// Make path absolute
		notesAbsPath := notePath
		if !filepath.IsAbs(notesAbsPath) {
			notesAbsPath = filepath.Join(projectRoot, notesAbsPath)
		}

		// Check existence
		info, err := os.Stat(notesAbsPath)
		if os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			if err := os.MkdirAll(notesAbsPath, 0755); err != nil {
				return fmt.Errorf("failed to create notes directory: %w", err)
			}
			fmt.Printf("Created notes directory at %s\n", notesAbsPath)
		} else if err != nil {
			return fmt.Errorf("error checking notes directory: %w", err)
		} else if !info.IsDir() {
			return fmt.Errorf("path exists and is not a directory: %s", notesAbsPath)
		}

		// Write config
		config := devnotesConfig{
			NotesPath:     notesAbsPath,
			FileExtension: ".md", // default to markdown
		}

		data, err := yaml.Marshal(&config)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		configFilePath := filepath.Join(projectRoot, ".devnote.yaml")
		if err := os.WriteFile(configFilePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		fmt.Printf("Initialised devnotes config at %s\n", configFilePath)
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("path", "p", "", "Define a custom path for your devnotes folder.")
	rootCmd.AddCommand(initCmd)
}
