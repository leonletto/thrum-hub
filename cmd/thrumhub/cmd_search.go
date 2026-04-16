package main

import (
	"fmt"
	"strings"

	"github.com/leonletto/thrum-hub/internal/query"
	"github.com/spf13/cobra"
)

func cmdSearch() *cobra.Command {
	var (
		k        int
		language string
		typ      string
	)
	c := &cobra.Command{
		Use:   "search <query>",
		Short: "Semantic search over curated entries",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := strings.Join(args, " ")
			res, err := query.Do(cmd.Context(), query.Deps{
				Store: a.store, Index: a.index, Telemetry: a.telem, Embedder: a.embed,
			}, query.Request{
				Query: q, K: k, Language: language, Type: typ,
				Agent: detectAuthor(), SourceRepo: detectSourceRepo(),
			})
			if err != nil {
				return err
			}
			fmt.Printf("query_id: %s\n", res.QueryID)
			fmt.Printf("results: %d\n\n", len(res.Hits))
			for i, h := range res.Hits {
				fmt.Printf("%d. [%s] %s (%s/%s, score=%.3f)\n", i+1, h.ID, h.Title, h.Language, h.Type, h.Score)
				if h.SourceRepo != "" {
					fmt.Printf("   repo: %s\n", h.SourceRepo)
				}
				fmt.Printf("   path: %s\n\n", h.Path)
			}
			return nil
		},
	}
	c.Flags().IntVar(&k, "k", 5, "number of results")
	c.Flags().StringVar(&language, "language", "", "filter by language")
	c.Flags().StringVar(&typ, "type", "", "filter by type")
	return c
}

func detectSourceRepo() string {
	// Stub for Plan 1. Plan 3 wires git-remote detection.
	return ""
}
