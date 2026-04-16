package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cmdStatus() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show hub status",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("hub root:      %s\n", a.cfg.Root)
			fmt.Printf("embedding:     %s (%s)\n", a.cfg.Embedding.Backend, a.cfg.Embedding.Model)
			inboxIDs, err := a.store.InboxList()
			if err != nil {
				return fmt.Errorf("inbox list: %w", err)
			}
			fmt.Printf("inbox depth:   %d\n", len(inboxIDs))
			n, err := a.index.Count()
			if err != nil {
				return fmt.Errorf("index count: %w", err)
			}
			fmt.Printf("indexed:       %d\n", n)
			return nil
		},
	}
}
