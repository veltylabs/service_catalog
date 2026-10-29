---
PLAN: "refactor: errors through rpc.WriteError — the errors are model refusals; errors.go keeps only its constants"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 522937650972762527
---

> Read AGENTS.md at the repo root before any change.

# PLAN — `service_catalog` answers its errors with `rpc.WriteError`

## Why

Decision D22 of `veltylabs/mjosefa-cms/docs/MVP_DEPLOYMENT_ROADMAP.md`: one piece answers every
operation's error. It exists now: `webtyp.com/rpc` v0.2.7 `rpc.WriteError(ctx, err)` (status from
`model.KindOf`, the failed field from `model.FieldOf`, a 500 never sends `err.Error()`), on the
vocabulary of `webtyp.com/model` v0.2.10, and `webtyp.com/orm`'s `ErrNotFound` is already a
not-found refusal. Today this module writes `err.Error()` for every status, 500 included
(`errors.go:37`), so SQL text can reach the browser. This plan deletes the module's own
`domainError`, `ValidationError`, `decodeError`, `statusOf` and `writeError`. **Statuses do not
change** (same mapping as today's `statusOf`).

## APIs used

| Symbol | Exact | Tag |
|---|---|---|
| `rpc.WriteError` | `func WriteError(ctx router.Context, err error)` | rpc v0.2.7 |
| `model.Invalid`, `model.NotFound`, `model.Conflict` | `type X string` with `Error()` and `Kind()`; declare `const ErrX model.NotFound = "…"` | model v0.2.10 |
| `model.FieldError` | `struct { Field string; Err error }`, `*FieldError` is an error of kind invalid | model v0.2.10 |
| `model.KindOf`, `model.FieldOf` | `func KindOf(err error) ErrorKind`, `func FieldOf(err error) string` | model v0.2.10 |

## Tools

```
go get webtyp.com/rpc@v0.2.7 webtyp.com/model@v0.2.10
```

## Stage 1 — `errors.go`

Keep the same names and texts, with the kind that gives today's status:

| Constant | Type | Status (unchanged) |
|---|---|---|
| `ErrNotFound`, `ErrSpecialtyNotFound`, `ErrSupplyLineNotFound` | `model.NotFound` | 404 |
| `ErrAlreadyExists`, `ErrSpecialtyPrefixExists`, `ErrSupplyLineExists` | `model.Conflict` | 409 |
| `ErrSpecialtyInUse`, `ErrNoRecords`, `ErrNoIds`, `ErrQuantity`, `ErrUnsupportedCurrency` | `model.Invalid` | 400 |

Add `ErrCannotRead model.Invalid = "service_catalog: cannot read the request"` (replaces
`msgDecode`). Delete `domainError`, `decodeError`, `msgDecode`, `writeError`, `statusOf`,
`ValidationError`. The file's comments stay in Spanish.

## Stage 2 — callers

- Every `writeError(ctx, decodeError{err})` → `rpc.WriteError(ctx, ErrCannotRead)`.
- Every other `writeError(ctx, err)` → `rpc.WriteError(ctx, err)` (57 calls in total, `ops.go`,
  `recipe.go`).
- `ValidationError{Err: err}` where `err` comes from a model `Validate` → return `err` as is (it is
  already a `*model.FieldError`).
- `ValidationError{Err: fmt.Err("invalid tax treatment")}` (`module.go:292`, `:332`) →
  `&model.FieldError{Field: "tax_treatment", Err: fmt.Err("invalid tax treatment")}`.
- `err == ErrNotFound || orm.IsNotFound(err)` (`ops.go:410`) stays as it is (no `==` between
  interfaces is added anywhere).

## Stage 3 — tests

- `tests/error_test.go`: keep the text table; add a second table asserting `model.KindOf` of each
  constant (Stage 1 table).
- `tests/recipe_test.go:74`: the zero-quantity case asserts `model.FieldOf(err) != ""` and
  `model.KindOf(err) == model.ErrorInvalid` instead of the type `ValidationError`.
- `tests/orm_test.go:120-123`: delete the `ValidationError.Error` coverage block.
- `tests/ops_test.go` and any test asserting a status: statuses are unchanged; where a test reads
  the body text of an error, it now reads the RFC 9457 JSON: assert on `"detail"` containing the
  same text. Add `TestOps_InternalErrorHidesDetail`: an operation whose store returns a plain
  `fmt.Err("pq: boom")` answers 500 and the body does not contain `pq:`.

## Stage 4 — docs (Spanish)

`README.md` and `docs/` where errors are described: «los errores salen por `rpc.WriteError`:
400/404/409 con el texto y el campo que falló; un 500 no muestra el detalle».

## Delivery check

```
GOPRIVATE=webtyp.com/*,github.com/veltylabs/* GOFLAGS=-mod=mod gotest
```

- A PR that does not compile, or whose tests could not run, is a failed delivery.
- Never pause to ask permission: when the check passes, commit, push and open the PR. Stuck 30
  minutes: open the PR listing what fails.
- `git diff --stat main` lists only: `errors.go`, `ops.go`, `recipe.go`, `module.go`, files under
  `tests/`, `README.md`, docs, `go.mod`, `go.sum`, this plan.

## Stages

| # | Files | Done when |
|---|---|---|
| 1 | `errors.go`, `go.mod`, `go.sum` | compiles (callers broken is expected until 2) |
| 2 | `ops.go`, `recipe.go`, `module.go` | compiles |
| 3 | `tests/*` | suite green |
| 4 | `README.md`, docs | text present |
