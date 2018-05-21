package folders

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerRemoveCommand(rootCmd *cobra.Command, logger logs.Logger) {
	var clientID string
	var folderID string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove folder",
		Long:  "Remove folder",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("Remove folder '%v'", folderID)
		},
	}
	cmd.Flags().StringVarP(&clientID, "session", "s", "", "Client ID")
	cmd.Flags().StringVarP(&folderID, "folder", "f", "", "Folder ID")
	cmd.MarkFlagRequired("session")
	cmd.MarkFlagRequired("folder")
	rootCmd.AddCommand(cmd)
}
