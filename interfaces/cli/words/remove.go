package words

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerRemoveCommand(rootCmd *cobra.Command, logger logs.Logger) {
	var clientID string
	var folderID string
	var wordID string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove word",
		Long:  "Remove word",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("Remove Word '%v'", wordID)
		},
	}
	cmd.Flags().StringVarP(&clientID, "session", "s", "", "Client ID")
	cmd.Flags().StringVarP(&wordID, "id", "i", "", "Word ID")
	cmd.Flags().StringVarP(&folderID, "folder", "f", "", "Folder ID")
	cmd.MarkFlagRequired("session")
	cmd.MarkFlagRequired("id")
	rootCmd.AddCommand(cmd)
}
