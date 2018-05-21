package sessions

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func RegisterCommands(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "sessions",
		Short: "Manipulate session",
		Long:  "Manipulate session",
	}
	registerCreateCommand(cmd, logger)
	registerDropCommand(cmd, logger)
	registerExtendCommand(cmd, logger)
	registerCheckCommand(cmd, logger)
	rootCmd.AddCommand(cmd)
}
