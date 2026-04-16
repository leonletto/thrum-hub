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
			inboxIDs, _ := a.store.InboxList()
			fmt.Printf("inbox depth:   %d\n", len(inboxIDs))
			n, _ := a.index.Count()
			fmt.Printf("indexed:       %d\n", n)
			return nil
		},
	}
}
