# Working Rules

## Structure (Anfin)
`cmd/app/main.go` (<=15 lines) | `configs/{config,default,load}.go` + yaml | `constants/` | `internal/app/{app,infra,module}.go` (manual DI) | `internal/<domain>/{module,repo,service,model,dto,constants}.go` | `internal/controller/<domain>.go` | `internal/router/<domain>.go` | `internal/middleware` | `internal/server` | `pkg/{failure,response,logger,databases,uid}`

## Domain
1. `NewModule(dto *CreateModuleDTO) *Module`; `newRepo`/`newService` unexported; return structs.
2. Interface only with 2+ implementations, defined at the consumer, 1-3 methods; one method is named `<Method>er`.
3. Repo uses GORM and `WithTx` for transactions. Raw SQL only where GORM cannot express it (LATERAL, FILTER), as a named const.

## Clean
1. No hard string or number. Header, context key, event, provider, status, error code, table name are const or enum type. Numbers live in config. Adding a provider, model, error code, route or plan is one registry line.
2. No dead code, no comments. File <=200 lines, func <=40 lines. Same logic twice goes to one place. Write nothing unused yet.
3. Names say what the thing does; one function does one job.

## Errors and logs
1. Error codes are enum values in `pkg/failure`. Wrap with `failure.ErrWithTrace` only where the error is born or crosses an outside boundary (repo/DB, provider/HTTP, config); upper layers return it unchanged. Handlers answer through one helper in `pkg/response`.
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
