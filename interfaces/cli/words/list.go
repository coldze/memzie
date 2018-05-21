package words

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerListCommand(rootCmd *cobra.Command, logger logs.Logger) {
	var clientID string
	var folderID string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list words",
		Long:  "list words",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("List Words for folder: '%v'", folderID)
		},
	}
	cmd.Flags().StringVarP(&clientID, "session", "s", "", "Client ID")
	cmd.Flags().StringVarP(&folderID, "folder", "f", "", "Folder ID")
	cmd.MarkFlagRequired("session")
	rootCmd.AddCommand(cmd)
}
