package main

import (
	"context"
	"fmt"

	"github.com/leonletto/thrum-hub/internal/embed"
	"github.com/leonletto/thrum-hub/internal/index"
	"github.com/leonletto/thrum-hub/internal/store"
	"github.com/spf13/cobra"
)

func cmdDebug() *cobra.Command {
	c := &cobra.Command{
		Use:   "debug",
		Short: "Debug and test helpers",
	}
	c.AddCommand(cmdDebugProcessInbox())
	return c
}

func cmdDebugProcessInbox() *cobra.Command {
	var acceptAll bool
	c := &cobra.Command{
		Use:   "process-inbox",
		Short: "Move all pending inbox entries into the curated corpus (test helper)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !acceptAll {
				return fmt.Errorf("--accept-all required (this helper has no curation logic)")
			}
			n, err := DebugProcessInbox(cmd.Context(), a.store, a.index, a.embed)
			if err != nil {
				return err
			}
			fmt.Printf("processed %d entries\n", n)
			return nil
		},
	}
	c.Flags().BoolVar(&acceptAll, "accept-all", false, "accept every pending entry without judgment")
	return c
}

// DebugProcessInbox runs the naive inbox→entries promotion and indexing loop.
// Plan 2 replaces the --accept-all semantics with real LLM curation.
func DebugProcessInbox(ctx context.Context, s *store.Store, idx *index.Index, emb embed.Embedder) (int, error) {
	ids, err := s.InboxList()
	if err != nil {
		return 0, err
	}
	processed := 0
	for _, id := range ids {
		e, err := s.InboxRead(id)
		if err != nil {
			return processed, fmt.Errorf("read %s: %w", id, err)
		}
		if err := s.EntryWrite(e); err != nil {
			return processed, fmt.Errorf("promote %s: %w", id, err)
		}
		path := s.EntryPath(id)
		if err := idx.Upsert(ctx, e, path, emb); err != nil {
			return processed, fmt.Errorf("index %s: %w", id, err)
		}
		if err := s.InboxDelete(id); err != nil {
			return processed, fmt.Errorf("delete inbox %s: %w", id, err)
		}
		processed++
	}
	return processed, nil
}
