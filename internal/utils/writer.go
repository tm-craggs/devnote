package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	templateFileName  = "template.txt"
	stateDirName      = "state"
	lastCommitFile    = "last_commit.txt"
	defaultDateFormat = "2006-01-02"
	defaultTimeFormat = "15:04:05"
	timestampFormat   = "2006-01-02 15:04:05"
)

// CreateNoteContent builds a note using a template and inserts dynamic placeholders.
func CreateNoteContent(devnoteDir string) (string, error) {
	templatePath := filepath.Join(devnoteDir, templateFileName)
	stateDir := filepath.Join(devnoteDir, stateDirName)

	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return "", fmt.Errorf("could not create state directory: %w", err)
	}

	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("could not read template: %w", err)
	}
	template := string(templateBytes)

	content := template
	content = strings.ReplaceAll(content, "{project}", GetProject())
	content = strings.ReplaceAll(content, "{date}", GetTime(defaultDateFormat))
	content = strings.ReplaceAll(content, "{time}", GetTime(defaultTimeFormat))
	content = strings.ReplaceAll(content, "{timestamp}", GetTime(timestampFormat))
	content = strings.ReplaceAll(content, "{log}", GetLog(stateDir))
	content = strings.ReplaceAll(content, "{branch}", GetBranch())
	content = strings.ReplaceAll(content, "{author-name}", GetAuthorName())
	content = strings.ReplaceAll(content, "{author-email}", GetAuthorEmail())

	return content, nil
}

/*
Placeholder functions to collect data for the following variables used in the note template:

{project}
{date}
{time}
{timestamp}
{log}
{branch}
{author-name}
{author-email}
*/

func GetProject() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return "project=NULL"
	}
	return filepath.Base(workingDir)
}

func GetTime(format string) string {
	return time.Now().Format(format)
}

func GetLog(stateDir string) string {
	lastCommitPath := filepath.Join(stateDir, lastCommitFile)

	var sinceHash string
	if data, err := os.ReadFile(lastCommitPath); err == nil {
		sinceHash = strings.TrimSpace(string(data))
	}

	var cmd *exec.Cmd
	if sinceHash != "" {
		cmd = exec.Command("git", "log", sinceHash+"..HEAD", "--pretty=format:- %h %s (%an, %ad)", "--date=short")
	} else {
		cmd = exec.Command("git", "log", "--pretty=format:- %h %s (%an, %ad)", "--date=short")
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error getting git log: %v", err)
	}

	headCmd := exec.Command("git", "rev-parse", "HEAD")
	headHash, err := headCmd.Output()
	if err == nil {
		_ = os.WriteFile(lastCommitPath, headHash, 0644)
	}

	return string(output)
}

func GetBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "branch=NULL"
	}
	return strings.TrimSpace(string(output))
}

func GetAuthorName() string {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "author=NULL"
	}
	return strings.TrimSpace(string(output))
}

func GetAuthorEmail() string {
	cmd := exec.Command("git", "config", "user.email")
	output, err := cmd.Output()
	if err != nil {
		return "email=NULL"
	}
	return strings.TrimSpace(string(output))
}
