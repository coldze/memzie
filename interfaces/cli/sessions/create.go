package sessions

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerCreateCommand(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create session",
		Long:  "Create session",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			panic("NOT IMPLEMENTED")
		},
	}
	rootCmd.AddCommand(cmd)
}
