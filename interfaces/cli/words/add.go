package words

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerAddCommand(rootCmd *cobra.Command, logger logs.Logger) {
	var translation string
	var clientID string
	var folderID string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add word",
		Long:  "Add word",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("Word '%v', translation '%v'", args[0], translation)
		},
	}
	cmd.Flags().StringVarP(&translation, "translation", "t", "", "Word's translation")
	cmd.Flags().StringVarP(&clientID, "session", "s", "", "Client ID")
	cmd.Flags().StringVarP(&folderID, "folder", "f", "", "Folder ID")
	cmd.MarkFlagRequired("translation")
	cmd.MarkFlagRequired("session")
	rootCmd.AddCommand(cmd)
}
