CREATE TABLE IF NOT EXISTS queries (
  id           TEXT PRIMARY KEY,
  ts           DATETIME NOT NULL,
  query_text   TEXT NOT NULL,
  agent        TEXT,
  source_repo  TEXT,
  num_results  INTEGER NOT NULL,
  top_score    REAL
);

CREATE TABLE IF NOT EXISTS query_results (
  query_id  TEXT NOT NULL,
  entry_id  TEXT NOT NULL,
  rank      INTEGER NOT NULL,
  score     REAL NOT NULL,
  PRIMARY KEY (query_id, entry_id)
);

CREATE TABLE IF NOT EXISTS feedback (
  query_id  TEXT NOT NULL,
  entry_id  TEXT,
  signal    TEXT NOT NULL CHECK (signal IN ('up','down')),
  note      TEXT,
  ts        DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_queries_ts ON queries(ts);
CREATE INDEX IF NOT EXISTS idx_queries_num_results ON queries(num_results);
