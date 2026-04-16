package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func cmdShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Print an entry by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			path := a.store.EntryPath(id)
			if path == "" {
				return fmt.Errorf("entry not found: %s", id)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
	}
}
