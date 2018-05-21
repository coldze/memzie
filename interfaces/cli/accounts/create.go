package accounts

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerCreateCommand(rootCmd *cobra.Command, logger logs.Logger) {
	var telegramID string
	var name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create account",
		Long:  "Create account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("Creating account with telegram ID '%v', User name: '%v'", telegramID, name)
		},
	}
	cmd.Flags().StringVarP(&telegramID, "telegram", "t", "", "Telegram ID")
	cmd.Flags().StringVarP(&telegramID, "name", "n", "", "Telegram user's name")
	cmd.MarkFlagRequired("telegram")
	cmd.MarkFlagRequired("name")
	rootCmd.AddCommand(cmd)
}
