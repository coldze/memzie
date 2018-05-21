package accounts

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerValidateCommand(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate account",
		Long:  "Validate account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			panic("NOT IMPLEMENTED")
		},
	}
	rootCmd.AddCommand(cmd)
}
