package accounts

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func RegisterCommands(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Manipulate account",
		Long:  "Manipulate account",
	}
	registerCreateCommand(cmd, logger)
	registerRegisterCommand(cmd, logger)
	registerValidateCommand(cmd, logger)
	rootCmd.AddCommand(cmd)
}
