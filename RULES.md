# Working Rules

## Code
1. Short, readable, best practice like Anfin (`internal/<domain>`, `pkg/failure`, `pkg/response`).
2. Follow https://go.dev/doc/effective_go.
3. No comments in code.
4. Tests are temporary. Write `*_test.go` (stdlib `testing`, table-driven with `t.Run`; integration tests via `TEST_DATABASE_URL`), run them, record the result in the report, then DELETE every `*_test.go` before commit. No test file is ever committed. Money and state flows also get SQL invariant checks.
5. No over-engineering. Understand before coding.
6. Errors go through `failure.ErrWithTrace`. Finite variants use enum type + const.

## Go conventions
1. Packages: short, lower-case, singular, no underscores; group by domain under `internal/`, entry points in `cmd/`.
2. Interfaces: defined at the consumer, 1-3 methods, one-method interface named `<Method>er`; none when there is a single implementation. Accept interfaces, return structs.
3. Errors: last return value, wrapped with `%w` (here via `failure.ErrWithTrace`), sentinel errors named `ErrXxx`, inspected with `errors.Is/As`. Do not log and return the same error. Never ignore an error without reason.
4. `context.Context` is the first parameter and is never stored in a struct. `defer cancel()` right after creating one. Every goroutine has a stop path.
5. Initialisms keep case (`ID`, `URL`, `HTTP`); receiver names are 1-2 letters; constants are CamelCase; no `init()` logic.
6. Logging: `log/slog`, structured fields, every request-scoped line carries `request_id`.
7. Before each commit: `gofmt -l .` empty, `go mod tidy`, `go vet ./...`, `go build ./...`, `go test ./...` (then delete the test files, see Code #4).

## Git workflow (Gitflow)
1. Branches: `main` (release history), `dev` (integration). Release and hotfix branches only when needed, same Gitflow rules.
2. Feature branch: `feature/CM-<n>-<slug>` where `CM-<n>` is the Jira key (CM-2 = https://congnguyen29102002-1792466035836.atlassian.net/browse/CM-2). Cut from `dev`, one ticket per branch.
3. Commit messages in English, start with the Jira key when the commit is scoped to one ticket.
4. Merge feature into `dev` only after the Per-ticket gate and user approval. Merge `dev` into `main` only on the user's say-so. Never push without "push đi".

## Per-ticket gate
1. Read the ticket: Goal, Scope, Acceptance Criteria.
2. Send the plan and STOP:
   - Goal and business value (revenue, cost, retention).
   - BR table: BR-ID | condition | expected result.
   - Edge cases and failure scenarios (2-3).
   - Trade-offs: A vs B and why.
   - Files to create or change.
   - Verify steps.
   - Definition of Done.
3. No file is touched before the user approves the plan.
4. Code only what the plan says. Anything else is asked first.
5. Self-review against these rules, then verify.
6. Report: BR | evidence | result. Unverified items are listed as unverified.
7. Commit only after approval. Push only on "push đi".

## Every plan must answer
- Why do we do this? Value for the user, value for the company, backend-as-customer view.
- Failure thinking: Redis down, stream cut, row inserted then server crashes, DB down, retry/duplicate. State what the user sees and what the data looks like afterwards.
- Trade-off: at least 2 options, choose one, say why.
- Scope: complete and correct for the business rule, no extra, no missing.
- Scale target: 1M users, 10k rps (read + chat). Every error is traceable by `request_id`.
- Modules have stable boundaries so other teams and partners can use them.

## Definition of Done
- Every Acceptance Criterion is checked with evidence.
- Every BR has proof.
- Build and vet are clean on the user's machine.
- No claim of "100% correct" without evidence.
