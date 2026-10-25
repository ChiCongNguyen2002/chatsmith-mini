# Working Rules

## Code
1. Short, readable, best practice like Anfin (`internal/<domain>`, `pkg/failure`, `pkg/response`).
2. Follow https://go.dev/doc/effective_go.
3. No comments in code.
4. No test files in the repo. Verify with scripts kept outside the repo and SQL invariant checks.
5. No over-engineering. Understand before coding.
6. Errors go through `failure.ErrWithTrace`. Finite variants use enum type + const.

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
