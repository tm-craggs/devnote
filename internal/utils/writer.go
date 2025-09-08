package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/storer"

	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	templateFileName  = "templates.txt"
	stateDirName      = "state"
	lastCommitFile    = "last_commit.txt"
	defaultDateFormat = "2006-01-02"
	defaultTimeFormat = "15:04:05"
	timestampFormat   = "2006-01-02 15:04:05"
)

// CreateNoteContent builds a note using a templates and inserts dynamic placeholders.
func CreateNoteContent(devnoteDir string) (string, error) {
	templatePath := filepath.Join(devnoteDir, templateFileName)
	stateDir := filepath.Join(devnoteDir, stateDirName)

	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return "", fmt.Errorf("could not create state directory: %w", err)
	}

	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("could not read templates: %w", err)
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
Placeholder functions to collect data for the following variables used in the note templates:

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

	r, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Sprintf("Error opening git repo: %v", err)
	}

	ref, err := r.Head()
	if err != nil {
		return fmt.Sprintf("Error getting HEAD: %v", err)
	}

	var sinceHash string
	if data, err := os.ReadFile(lastCommitPath); err == nil {
		sinceHash = strings.TrimSpace(string(data))
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return fmt.Sprintf("Error getting git log: %v", err)
	}

	var result strings.Builder
	stopAt := sinceHash

	// Iterate commits until we reach stopAt
	err = cIter.ForEach(func(c *object.Commit) error {
		if stopAt != "" && c.Hash.String() == stopAt {
			return storer.ErrStop // stop iteration
		}

		shortHash := c.Hash.String()[:7]
		shortDate := c.Author.When.Format("2006-01-02")
		message := strings.TrimSpace(c.Message)

		result.WriteString(fmt.Sprintf("- %s %s (%s, %s)\n",
			shortHash,
			message,
			c.Author.Name,
			shortDate,
		))
		return nil
	})

	if err != nil && err != storer.ErrStop {
		return fmt.Sprintf("Error iterating commits: %v", err)
	}

	// Save current HEAD for next run
	_ = os.WriteFile(lastCommitPath, []byte(ref.Hash().String()), 0644)

	return result.String()
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
