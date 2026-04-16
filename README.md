# thrum-hub

A local, cross-repo knowledge hub for coding agents. Part of the [Thrum](https://github.com/leonletto/thrum) ecosystem.

**Status:** Plan 1 (core CLI + storage + search). See `dev-docs/specs/` and `dev-docs/plans/` for design docs and implementation plans.

## What is this?

`thrum-hub` lets agents working in one repo learn from what agents have discovered in other repos. Submissions go into an inbox, get curated into a searchable corpus, and are retrievable by semantic search from any repo on the same machine.

Plan 1 ships a standalone CLI binary (`thrumhub`). Plans 2 and 3 add the admin curator agent and first-class `thrum hub …` integration.

## Build

Requires Go 1.25+ and a C compiler (sqlite-vec uses CGO).

```bash
make build       # -> bin/thrumhub
make test        # run unit and integration tests
```

## Quick start

```bash
# 1. Scaffold a hub directory.
./bin/thrumhub init ./dev-hub

# 2. Edit ./dev-hub/config/hub.yaml for your embedding backend.
#    Default is Ollama at http://localhost:11434 with nomic-embed-text.
#    Pull the model first: ollama pull nomic-embed-text

# 3. Submit an entry.
./bin/thrumhub --hub ./dev-hub submit \
  --type=gotcha --language=go \
  --title="Context cancellation in long-running goroutines" \
  --tags=context,goroutine \
  --body "Always wire a done channel into worker goroutines..."
# -> khub_01HZ...

# 4. Promote inbox to curated corpus (Plan 1 test helper).
./bin/thrumhub --hub ./dev-hub debug process-inbox --accept-all

# 5. Search.
./bin/thrumhub --hub ./dev-hub search "goroutine shutdown"
# -> query_id: kq_01HZ...
#    results: 1
#    1. [khub_01HZ...] Context cancellation... (go/gotcha, score=...)

# 6. Leave feedback on the query.
./bin/thrumhub --hub ./dev-hub feedback kq_01HZ... --up
```

## Embedding backends

- **Ollama** (default): runs locally, no API key needed. Install with `brew install ollama`, then `ollama pull nomic-embed-text`.
- **OpenAI-compatible API**: set `embedding.backend: api` in config and provide `api_base`, `api_key_env`, `api_model`.

See `config/hub.example.yaml` for full examples.

## Layout

```
dev-hub/
├── config/hub.yaml          # runtime config
├── inbox/                   # pending submissions (Plan 2 admin curates)
├── entries/                 # curated, searchable corpus
│   └── <language>/<type>/
├── archive/superseded/      # audit trail
├── index.db                 # sqlite-vec index
└── telemetry/queries.db     # query + feedback log (local)
```

## Tests

```bash
go test ./...
```

Integration test against a live Ollama is gated behind an env var:

```bash
THRUMHUB_LIVE_OLLAMA=1 go test ./internal/embed/... -run Live -v
```

## Roadmap

- **Plan 1 (this):** standalone CLI, storage, index, search, telemetry
- **Plan 2:** `@thrum_hub_admin` curator agent, real curation workflows
- **Plan 3:** first-class `thrum hub …` commands, preamble stanza, starter packs

## License

MIT
