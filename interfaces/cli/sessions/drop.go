package sessions

import (
	"github.com/coldze/primitives/logs"
	"github.com/spf13/cobra"
)

func registerRegisterCommand(rootCmd *cobra.Command, logger logs.Logger) {
	cmd := &cobra.Command{
		Use:   "drop",
		Short: "Drop session",
		Long:  "Drop session",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			panic("NOT IMPLEMENTED")
		},
	}
	rootCmd.AddCommand(cmd)
}
