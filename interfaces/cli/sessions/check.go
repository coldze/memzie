package sessions

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerCheckCommand(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check sessopm",
		Long:  "Check session",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			panic("NOT IMPLEMENTED")
		},
	}
	rootCmd.AddCommand(cmd)
}
