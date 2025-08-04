package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	content = strings.ReplaceAll(content, "{date}", GetDate(defaultDateFormat))
	content = strings.ReplaceAll(content, "{time}", GetTime(defaultTimeFormat))
	content = strings.ReplaceAll(content, "{timestamp}", GetTimestamp(timestampFormat))
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

	}
	return os.Getwd()
}

func GetDate(format string) string {
	// TODO: return current date formatted
	return ""
}

func GetTime(format string) string {
	// TODO: return current time formatted
	return ""
}

func GetTimestamp(format string) string {
	// TODO: return current timestamp formatted
	return ""
}

func GetLog(stateDir string) {
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
		return ""
	}
	return strings.TrimSpace(string(output))
}

func GetAuthorName() string {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func GetAuthorEmail() string {
	cmd := exec.Command("git", "config", "user.email")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
