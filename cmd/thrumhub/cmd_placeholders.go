package main

import "github.com/spf13/cobra"

func stub(name, short string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func cmdInit() *cobra.Command     { return stub("init", "Scaffold a hub at the given path") }
func cmdStatus() *cobra.Command   { return stub("status", "Show hub status") }
func cmdFeedback() *cobra.Command { return stub("feedback", "Record thumbs up/down on a query") }
func cmdHistory() *cobra.Command  { return stub("history", "Show recent queries") }
func cmdDebug() *cobra.Command {
	c := stub("debug", "Debug and test helpers")
	return c
}
