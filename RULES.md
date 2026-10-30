# Working Rules

## Structure (Anfin style, versioned routes)
Layout: `backend/` (Go), `frontend/` (static HTML, CSS, JS served at `/`), `deploy/`, `docs/`, `tests/`. Inside `backend/`:
- `cmd/app/main.go`: load config and call `app.Run` (body <= 15 lines).
- `configs/`: `config.go`, `default.go`, `load.go`, `secrets.go` + yaml.
- `constants/`: values shared by several domains (provider, header, context key, log key, route paths).
- `migrations/`: SQL files, applied in order.
- `pkg/`: utilities and external clients. Must not import `internal/`. Packages: `failure`, `response`, `logger`, `databases`, `uid`, `pagination`, `reqctx`, `sse`, `prompt`, `llm`, `oauth`, `sessionstate`.
- `internal/app/`: `app.go` (Run), `oauth.go` (build providers from config). Manual DI only.
- `internal/server/`: start and graceful shutdown.
- `internal/domain/<context>/`: business logic. Contexts: `auth`, `catalog`, `chat`, `conversation`, `user`. Files: `model.go`, `dto.go`, `constants.go`, `module.go`, `repo.go`, `service.go` (a context may omit `dto.go` when it has no API types, and `constants.go` when it has no values).
- `internal/controller/`: one file per context with gin handlers. Handlers answer errors with `response.ErrFrom`. A new version is a new file: `chat_v2.go`.
- `internal/router/`: `engine.go` (mounts `/v1`, and `/v2` when it exists) and one file per context with its routes. A new version is a new file: `chat_v2.go`, with `register<Ctx>V2` mounted on the `/v2` group in `engine.go`. Deprecated route files move to `router/old_router/`.
- `internal/middleware/`: request_id, access_log, recovery, timeout, session.

## Domain (one way to write a domain)
1. Each domain folder holds these files with these names. No suffix files except version files (`dto_v2.go`).
2. `model.go`: structs with gorm tags.
3. `dto.go`: request and response structs with `json` tags, plus mapper functions (`ToResponse`) that build responses from models. This is the contract with clients. A new API version gets `dto_v2.go`.
4. `constants.go`: values used only inside the domain. Query strings (SQL, table, column, where, order) are not constants: write them in `repo.go` where they are used.
5. `module.go`: `NewModule(...)` builds the repo and service and returns `Module{Service}`.
6. `repo.go`: a concrete `Repo` struct using gorm. No interface.
7. `service.go`: business logic. It calls `repo` directly and returns `failure.Code` errors.
8. A domain must not import `gin`, `net/http`, `controller`, `router` or `middleware`.
9. Interfaces only where two or more implementations exist: `pkg/llm.Streamer` and `pkg/oauth.Provider`.
10. A service file over 450 lines is reviewed before it grows further.

## Clean
1. No hard string or number in business code. Shared values go to `constants/`, domain values to `constants.go`. Query strings (SQL, table and column names) live in `repo.go`, the only file that talks to the DB. Numbers that tune behavior live in config.
2. No dead code. Comments only as one doc line per package (`// Package x ...`). Lint uses `.golangci.yml` from Anfin (funlen 100, gocyclo 15, lll 120).
3. Names say what the thing does; one function does one job.

## Errors and logs
1. Error codes are enum values in `pkg/failure`. Wrap with `failure.ErrWithTrace` only where the error is born or crosses an outside boundary (repo/DB, provider/OAuth/HTTP, config); upper layers return it unchanged. Never call `ErrWithTrace(nil)`. Handlers answer through one helper in `pkg/response`.
2. Never log and return the same error. Never ignore an error without a reason.
3. `log/slog` with structured fields; every request-scoped line carries `request_id` from the context.

## Go
1. Packages: short, lower-case, singular, no underscores. Initialisms keep case (`ID`, `URL`). Receivers 1-2 letters. Constants CamelCase. Unused params start with `_`. No `init()` logic.
2. `context.Context` first, never stored in a struct. `defer cancel()` right after creating one. Every goroutine has a stop path.
3. Follow https://go.dev/doc/effective_go.

## Tests
Temporary. Write `*_test.go` (stdlib, table-driven with `t.Run`; DB tests via `TEST_DATABASE_URL`), run, record the result, then DELETE every `*_test.go` before commit. Money and state flows also get SQL invariant checks.

## Git
1. Branches `main`, `dev`. Feature branch `feature/CM-<n>-<slug>` from `dev`, one ticket per branch (Jira: https://congnguyen29102002-1792466035836.atlassian.net/browse/CM-<n>).
2. Commit messages in English, start with the Jira key.
3. Merge to `dev` only after the gate and user approval. `dev` to `main` only on the user's say-so. Never push without "push đi".

## Per-ticket gate
1. Read the ticket: Goal, Value, Scope, Acceptance Criteria.
2. Send the plan and STOP: business flow and value (user, company, scale), BR table (BR-ID | condition | result), 2-3 edge cases with failure thinking (Redis/DB down, stream cut, crash after insert, retry/duplicate: what the user sees, what the data looks like), decision table (item | why | trade-off | alternative | source | simpler?), files to create/change/DELETE, verify steps, Definition of Done.
3. No file is touched before approval. Code only what the plan says; anything else is asked first.
4. Run skill `chatsmith-review`, agents `code-reviewer` and `qa-verifier`; fix; then gate: `gofmt -l .` empty, `go mod tidy`, `go vet ./...`, `go build ./...`, `go test ./...`, run on localhost and hit the endpoints.
5. Report: BR | evidence | result; unverified items listed as unverified. Commit only after approval.

## Definition of Done
Every Acceptance Criterion and BR has evidence. Build and vet are clean. No claim of "100% correct" without evidence.
