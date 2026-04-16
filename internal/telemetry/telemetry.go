package telemetry

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaSQL string

type Telemetry struct {
	db *sql.DB
}

type Result struct {
	EntryID string
	Rank    int
	Score   float64
}

type QueryRecord struct {
	ID         string
	Timestamp  time.Time
	QueryText  string
	Agent      string
	SourceRepo string
	NumResults int
	TopScore   float64
	Results    []Result
}

type Feedback struct {
	QueryID   string
	EntryID   string // optional
	Signal    string // "up" or "down"
	Note      string
	Timestamp time.Time
}

func Open(path string) (*Telemetry, error) {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("telemetry: mkdir %s: %w", dir, err)
		}
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, err
	}
	return &Telemetry{db: db}, nil
}

func (t *Telemetry) Close() error { return t.db.Close() }

func (t *Telemetry) RecordQuery(q QueryRecord) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		INSERT INTO queries (id, ts, query_text, agent, source_repo, num_results, top_score)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, q.ID, q.Timestamp, q.QueryText, q.Agent, q.SourceRepo, q.NumResults, q.TopScore); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO query_results (query_id, entry_id, rank, score) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range q.Results {
		if _, err := stmt.Exec(q.ID, r.EntryID, r.Rank, r.Score); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (t *Telemetry) RecordFeedback(f Feedback) error {
	_, err := t.db.Exec(`
		INSERT INTO feedback (query_id, entry_id, signal, note, ts)
		VALUES (?, ?, ?, ?, ?)
	`, f.QueryID, f.EntryID, f.Signal, f.Note, f.Timestamp)
	return err
}

func (t *Telemetry) RecentQueries(sourceRepo string, limit int) ([]QueryRecord, error) {
	q := `SELECT id, ts, query_text, agent, source_repo, num_results, top_score FROM queries`
	args := []any{}
	if sourceRepo != "" {
		q += ` WHERE source_repo = ?`
		args = append(args, sourceRepo)
	}
	q += ` ORDER BY ts DESC LIMIT ?`
	args = append(args, limit)

	rows, err := t.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []QueryRecord
	for rows.Next() {
		var r QueryRecord
		var agent, repo sql.NullString
		if err := rows.Scan(&r.ID, &r.Timestamp, &r.QueryText, &agent, &repo, &r.NumResults, &r.TopScore); err != nil {
			return nil, err
		}
		r.Agent = agent.String
		r.SourceRepo = repo.String
		out = append(out, r)
	}
	return out, rows.Err()
}
