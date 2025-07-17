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
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new devnote",
	Long: `The 'new' command creates a new devnote in the specified directory. The devnote will open in your default
text editor and will contain the template text specified in the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {

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
		content, err := generateNoteContent(devnoteDir)
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

func generateNoteContent(devnoteDir string) (string, error) {
	stateDir := filepath.Join(devnoteDir, "state")
	lastCommitPath := filepath.Join(stateDir, "last_commit.txt")

	// Ensure state directory exists
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return "", fmt.Errorf("could not create state directory: %w", err)
	}

	// Read last commit
	var sinceHash string
	if data, err := os.ReadFile(lastCommitPath); err == nil {
		sinceHash = strings.TrimSpace(string(data))
	}

	// Run git log
	var gitLogCmd *exec.Cmd
	if sinceHash != "" {
		gitLogCmd = exec.Command("git", "log", sinceHash+"..HEAD", "--pretty=format:- %h %s (%an, %ad)", "--date=short")
	} else {
		gitLogCmd = exec.Command("git", "log", "--pretty=format:- %h %s (%an, %ad)", "--date=short")
	}

	output, err := gitLogCmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not get git log: %w", err)
	}

	// Get latest commit hash
	headCmd := exec.Command("git", "rev-parse", "HEAD")
	headHash, err := headCmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not get current HEAD: %w", err)
	}

	// Save new HEAD as last commit
	if err := os.WriteFile(lastCommitPath, headHash, 0644); err != nil {
		return "", fmt.Errorf("could not write last commit file: %w", err)
	}

	// Return commit log as string
	return string(output), nil
}

func init() {
	rootCmd.AddCommand(newCmd)
}
