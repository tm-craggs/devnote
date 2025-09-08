package cmd

import (
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Create a new devnote",
	Long: `The 'new' command creates a new devnote in the specified directory. The devnote will open in your default
text editor and will contain the template text specified in the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil

	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
