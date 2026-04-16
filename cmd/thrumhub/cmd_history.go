package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cmdHistory() *cobra.Command {
	var (
		repo  string
		limit int
	)
	c := &cobra.Command{
		Use:   "history",
		Short: "Show recent queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			recs, err := a.telem.RecentQueries(repo, limit)
			if err != nil {
				return err
			}
			for _, r := range recs {
				fmt.Printf("%s  %s  [%d results]  %s\n", r.Timestamp.Format("2006-01-02 15:04"), r.ID, r.NumResults, r.QueryText)
			}
			return nil
		},
	}
	c.Flags().StringVar(&repo, "repo", "", "filter by source repo")
	c.Flags().IntVar(&limit, "limit", 20, "max rows")
	return c
}
