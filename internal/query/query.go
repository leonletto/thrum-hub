package query

import (
	"context"
	"time"

	"github.com/leonletto/thrum-hub/internal/embed"
	"github.com/leonletto/thrum-hub/internal/ids"
	"github.com/leonletto/thrum-hub/internal/index"
	"github.com/leonletto/thrum-hub/internal/store"
	"github.com/leonletto/thrum-hub/internal/telemetry"
)

type Deps struct {
	Store     *store.Store
	Index     *index.Index
	Telemetry *telemetry.Telemetry
	Embedder  embed.Embedder
}

type Request struct {
	Query      string
	K          int
	Language   string
	Type       string
	Agent      string
	SourceRepo string
}

type Result struct {
	QueryID string
	Hits    []index.Hit
}

func Do(ctx context.Context, d Deps, r Request) (*Result, error) {
	k := r.K
	if k <= 0 {
		k = 5
	}
	hits, err := d.Index.Search(ctx, r.Query, k, d.Embedder, index.SearchFilter{
		Language: r.Language,
		Type:     r.Type,
	})
	if err != nil {
		return nil, err
	}

	qid := ids.NewQueryID()
	topScore := 0.0
	if len(hits) > 0 {
		topScore = hits[0].Score
	}
	rec := telemetry.QueryRecord{
		ID:         qid,
		Timestamp:  time.Now().UTC(),
		QueryText:  r.Query,
		Agent:      r.Agent,
		SourceRepo: r.SourceRepo,
		NumResults: len(hits),
		TopScore:   topScore,
	}
	for i, h := range hits {
		rec.Results = append(rec.Results, telemetry.Result{
			EntryID: h.ID, Rank: i, Score: h.Score,
		})
	}
	if err := d.Telemetry.RecordQuery(rec); err != nil {
		return nil, err
	}
	return &Result{QueryID: qid, Hits: hits}, nil
}
