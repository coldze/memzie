package folders

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerListCommand(rootCmd *cobra.Command, logger logs.Logger) {
	var clientID string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List folders",
		Long:  "List folders",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("List folders")
		},
	}
	cmd.Flags().StringVarP(&clientID, "session", "s", "", "Client ID")
	cmd.MarkFlagRequired("session")
	rootCmd.AddCommand(cmd)
}
