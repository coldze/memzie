package accounts

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerRegisterCommand(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register account",
		Long:  "Register account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			panic("NOT IMPLEMENTED")
		},
	}
	rootCmd.AddCommand(cmd)
}
