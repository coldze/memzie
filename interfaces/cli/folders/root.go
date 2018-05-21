package folders

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func RegisterCommands(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "words",
		Short: "Manipulate words",
		Long:  "Manipulate words",
	}
	registerAddCommand(cmd, logger)
	registerListCommand(cmd, logger)
	registerRemoveCommand(cmd, logger)
	rootCmd.AddCommand(cmd)
}
