package main

import (
	"fmt"
	"time"

	"github.com/leonletto/thrum-hub/internal/telemetry"
	"github.com/spf13/cobra"
)

func cmdFeedback() *cobra.Command {
	var (
		up, down bool
		onEntry  string
		note     string
	)
	c := &cobra.Command{
		Use:   "feedback <query_id>",
		Short: "Record thumbs up/down on a query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if up == down { // both true or both false
				return fmt.Errorf("exactly one of --up or --down required")
			}
			signal := "up"
			if down {
				signal = "down"
			}
			if signal == "down" && note == "" {
				return fmt.Errorf("--note required for --down (explain what you were looking for)")
			}
			return a.telem.RecordFeedback(telemetry.Feedback{
				QueryID:   args[0],
				EntryID:   onEntry,
				Signal:    signal,
				Note:      note,
				Timestamp: time.Now().UTC(),
			})
		},
	}
	c.Flags().BoolVar(&up, "up", false, "positive feedback")
	c.Flags().BoolVar(&down, "down", false, "negative feedback")
	c.Flags().StringVar(&onEntry, "on", "", "scope feedback to a specific entry ID")
	c.Flags().StringVar(&note, "note", "", "optional note (required for --down)")
	return c
}
