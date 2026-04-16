package main

import (
	"fmt"

	"github.com/leonletto/thrum-hub/internal/index"
	"github.com/spf13/cobra"
)

func cmdRelated() *cobra.Command {
	var k int
	c := &cobra.Command{
		Use:   "related <id>",
		Short: "Find entries related to the given ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			e, err := a.store.EntryRead(id)
			if err != nil {
				return err
			}
			// Use the entry's own title+body as the query text.
			text := e.Title + " " + e.Content
			hits, err := a.index.Search(cmd.Context(), text, k+1, a.embed, index.SearchFilter{})
			if err != nil {
				return err
			}
			printed := 0
			for _, h := range hits {
				if h.ID == id {
					continue
				}
				if printed >= k {
					break
				}
				fmt.Printf("[%s] %s (score=%.3f)\n", h.ID, h.Title, h.Score)
				printed++
			}
			return nil
		},
	}
	c.Flags().IntVar(&k, "k", 5, "number of results")
	return c
}
