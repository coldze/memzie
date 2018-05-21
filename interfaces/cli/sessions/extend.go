package sessions

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerExtendCommand(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Extend session",
		Long:  "Extend sessions",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			panic("NOT IMPLEMENTED")
		},
	}
	rootCmd.AddCommand(cmd)
}
