---
name: chatsmith-review
description: Gate for every chatsmith-mini ticket. Run before sending a plan, after coding new code, and whenever an old file is touched. Checks the 6 steps, the decision table, sources, and runs the grep checks.
---

# chatsmith-review

Orchestrated by the main session (sub-agents cannot call sub-agents).

## When
1. Before the plan: call `pm-analyst` with the ticket key.
2. After coding, run in this order. Each step must pass before the next one starts:
   a. Gate and logic: `qa-verifier` (happy, bad and edge case for every BR).
   b. UI: `ui-tester` if the change touches `web/`. Compare with the approved baseline.
   c. PM: `pm-analyst` checks each Acceptance Criterion in the Jira ticket against the evidence from a and b.
   d. Review: `code-reviewer` on the changed files and every old file touched. Fix, run again until no blocker.
   Steps a and d can run in parallel. Steps b and c wait for a.
3. Then send the report to the user. Unclear points are asked, not guessed. Do not report done if step a or b is unverified.

## 6 steps, each a checkable question
1. Business: can you state who uses it, the flow, 3 edge cases, where the money goes?
2. Correct and complete: every Acceptance Criterion and BR has run evidence; nothing extra, nothing missing?
3. Easy to change: does adding a provider, model, error code, route or plan take one registry line?
4. New dev reads 80-90%: do names say the job, one function one job, file <=200 lines, func <=40?
5. Value: one measurable sentence for the user, one for the company (revenue, cost, retention).
6. Scale: 1M users, 10k rps; where is the bottleneck; what the user sees when DB/Redis/stream dies; is every error traceable by `request_id`?

## Decision table (one row per non-trivial decision)
item | why | trade-off | alternative | source | simpler?

Sources: at least 2 independent ones for a big decision (Anfin file, Effective Go, Go Code Review Comments, Uber or Google Go style, official GORM/Gin/Postgres docs, an engineering post from a large company), checked against the versions in go.mod (2026). No source means write "own choice" and why. Never paste AI output unchecked.

## Grep checks (run from repo root, each must be explained or empty)
- Use `grep -rnE` (rg is not installed here) or the Grep tool.
- Hard strings: `grep -rnE '"[^"]{3,}"' internal pkg cmd --include='*.go' --exclude-dir=constants` and judge each hit (SQL const, struct tag, log key are allowed).
- Comments: `grep -rnE '(^|[[:space:]])//' --include='*.go' . | grep -v 'go:embed'` must be empty.
- Size: files over 200 lines `find . -name '*.go' -not -path './.claude/*' | xargs wc -l | awk '$1>200 && $2!="total"'`; funcs over 40 lines judged by reading.
- Dead code: every new func, field, const has a caller (`grep -rn 'Name'`); `go vet ./...` clean.
- Errors: no `_ =` on an error without a reason; no log-then-return of the same error.
- Leftover tests: `find . -name '*_test.go'` must be empty before commit.
- Layout: each file sits where RULES.md says.

## Report format
BR | evidence | result. Unverified items listed as unverified. Findings from agents grouped as blocker / should fix / note.
