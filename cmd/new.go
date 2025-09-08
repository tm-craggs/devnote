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
	"time"

	"github.com/spf13/cobra"
	"github.com/tm-craggs/devnote/internal/config"
	"github.com/tm-craggs/devnote/internal/utils"
	"gopkg.in/yaml.v3"
)

type newFlags struct {
	name string
	path string
}

// helper function to parse flags with error handling
func getNewFlags(cmd *cobra.Command) (newFlags, error) {
	var flags newFlags
	var err error

	flags.name, err = cmd.Flags().GetString("name")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --name flag: %w", err)
	}

	flags.path, err = cmd.Flags().GetString("path")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --path flag: %w", err)
	}

	return flags, nil
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new devnote",
	Long: `The 'new' command creates a new devnote in the specified directory. The devnote will open in your default
text editor and will contain the templates text specified in the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// get flags
		flags, err := getNewFlags(cmd)
		if err != nil {
			return err
		}

		// get devnote directory
		projectRoot, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not determine working directory: %w", err)
		}

		devnoteDir := filepath.Join(projectRoot, ".devnote")

		// load config file
		configFilePath := filepath.Join(devnoteDir, "devnote.yaml")
		configData, err := os.ReadFile(configFilePath)
		if err != nil {
			return fmt.Errorf("devnote not initalised. run 'devnote init' in your project root")
		}

		var config config.DevnoteConfig
		if err := yaml.Unmarshal(configData, &config); err != nil {
			return fmt.Errorf("could not parse config: %w", err)
		}

		if config.NotesPath == "" {
			return fmt.Errorf("NotesPath not set in config")
		}

		// check notes directory exists
		if _, err := os.Stat(config.NotesPath); os.IsNotExist(err) {
			return fmt.Errorf("notes directory does not exist at: %s", config.NotesPath)
		}

		// create file name
		var noteFileName string
		ext := config.FileExtension
		if ext == "" {
			ext = ".md"
		}
		if flags.name == "" {
			timeNow := time.Now().Format("2006-01-02 15:04:05")
			noteFileName = timeNow + ext
		} else {
			noteFileName = flags.name + ext
		}

		// create file path
		var notePath string
		if flags.path == "" {
			// use path from configuration file
			notePath = filepath.Join(config.NotesPath, noteFileName)
		} else {
			notePath = flags.path
		}

		// create devnote
		content, err := utils.CreateNoteContent(devnoteDir)
		if err != nil {
			return fmt.Errorf("could not generate note content: %w", err)
		}

		err = os.WriteFile(notePath, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("could not write to devnote: %w", err)
		}

		// open user text editor
		err = utils.OpenEditor(notePath)
		if err != nil {
			return fmt.Errorf("could not open editor: %w", err)
		}

		fmt.Println("devnote created at", notePath)
		return nil

	},
}

func init() {

	newCmd.Flags().StringP("name", "n", "", "name of the devnote")
	newCmd.Flags().StringP("path", "p", "", "path of the devnote")

	rootCmd.AddCommand(newCmd)
}
