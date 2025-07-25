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
	"github.com/tm-craggs/devnote/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new devnote",
	Long: `The 'new' command creates a new devnote in the specified directory. The devnote will open in your default
text editor and will contain the template text specified in the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// load project root
		projectRoot, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not determine working directory: %w", err)
		}

		// load config file
		configFilePath := filepath.Join(projectRoot, ".devnote.yaml")
		configData, err := os.ReadFile(configFilePath)
		if err != nil {
			return fmt.Errorf("devnote not initalised. run 'devnote init' in your project root")
		}

		var config devnotesConfig
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
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		ext := config.FileExtension
		if ext == "" {
			ext = ".md"
		}

		noteFileName := timeNow + ext
		notePath := filepath.Join(config.NotesPath, noteFileName)

		// create devnote
		file, err := os.Create(notePath)
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}

		// close devnote
		err = file.Close()
		if err != nil {
			return fmt.Errorf("could not close file: %w", err)
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
	rootCmd.AddCommand(newCmd)
}
