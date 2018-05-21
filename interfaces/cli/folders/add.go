package folders

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerAddCommand(rootCmd *cobra.Command, logger logs.Logger) {
	var clientID string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add folder",
		Long:  "Add folder",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("Folder '%v'", args[0])
		},
	}
	cmd.Flags().StringVarP(&clientID, "session", "s", "", "Client ID")
	cmd.MarkFlagRequired("session")
	rootCmd.AddCommand(cmd)
}
