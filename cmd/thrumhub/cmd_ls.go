package main

import (
	"fmt"

	"github.com/leonletto/thrum-hub/internal/schema"
	"github.com/leonletto/thrum-hub/internal/store"
	"github.com/spf13/cobra"
)

func cmdLs() *cobra.Command {
	var language, typ, tag string
	c := &cobra.Command{
		Use:   "ls",
		Short: "List curated entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := a.store.EntryList(store.EntryFilter{
				Language: language,
				Type:     schema.EntryType(typ),
				Tag:      tag,
			})
			if err != nil {
				return err
			}
			for _, e := range entries {
				fmt.Printf("[%s] %s (%s/%s) status=%s\n", e.ID, e.Title, e.Language, e.Type, e.Status)
			}
			return nil
		},
	}
	c.Flags().StringVar(&language, "language", "", "filter by language")
	c.Flags().StringVar(&typ, "type", "", "filter by type")
	c.Flags().StringVar(&tag, "tag", "", "filter by tag")
	return c
}
