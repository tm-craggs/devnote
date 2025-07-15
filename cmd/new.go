package cmd

import "github.com/spf13/cobra"

import (
	"fmt"
	"os"
	"time"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new devnote",
	Long: `The 'new' command creates a new devnote in the specified directory. The devnote will open in your default
text editor and will contain the template text specified in the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		timeNow := time.Now().Format("2006-01-02 15:04:05")

		// create devnote directory if none exists

		// check if the directory exists, if not create.
		if _, err := os.Stat("devnotes"); os.IsNotExist(err) {
			// does not exist
			err := os.Mkdir("devnotes", 0755)
			if err != nil {
				return fmt.Errorf("could not create devnotes subdirectory: %w", err)
			}
		}

		// create devnote
		file, err := os.Create("devnotes/" + timeNow + ".md")

		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}

		fmt.Println("devnote created")
		return file.Close()

	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
