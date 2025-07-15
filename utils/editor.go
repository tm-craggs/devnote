package utils

import (
	"os"
	"os/exec"
)

func OpenEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// fallback editor
		editor = "nano"
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
