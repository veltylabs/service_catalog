---
PLAN: "refactor: the tenant comes from the server (Deps.TenantID), never from the request; permissions named by the tables (model.ResourceOf)"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 11156536684367955011
---

> Read AGENTS.md at the repo root before any change.

# PLAN — `service_catalog`: tenant from the server, permissions from the tables (M2l + M2h step 2)

## Why

Phases M2l and M2h (step 2) of `veltylabs/mjosefa-cms/docs/MVP_DEPLOYMENT_ROADMAP.md`, decisions D32
and D26.

- **D32 — the tenant never comes from the browser.** Today every list/get/delete handler reads
  `args.TenantId` from the request body (`ops.go:204,224,284,304,320,425,457,493`,
  `recipe.go:137,173`), and the saves keep whatever `TenantId` the record carries. A person can send
  another clinic's tenant and read or write its catalog. Worse, an empty `tenant_id` turns the
  filter off (`module.go:69`: `if tenantId != ""`): it lists **every** tenant. `device_manager`,
  `staff_manager` and `room_layout` already take the tenant from `Deps.TenantID`
  (`device_manager/module.go:20,42`); this module gets the same.
- **D26 — a permission is a table, typed.** `Requires("specialty", …)`, `"service"`,
  `"service_agreement"`, `resServiceSupply` and `Permission: "service"`/`"specialty"` in `ui/browser.go`
  are string literals; `model.ResourceOf(&SpecialtyModel)` names the same resource from the table
  definition, so a typo cannot compile. The names do not change (`specialty`, `service`,
  `service_agreement`, `service_supply`), so stored role grants keep working.

The exported Go methods (`GetService(tenantId, id)`, `ServiceCard`, `ServiceExists`, …) keep their
`tenantId` parameter: they are called by server code of other modules and by `seed`, never by the
browser. Only the **handlers** (what a request reaches) change.

## APIs used

| Symbol | Exact | Where |
|---|---|---|
| `model.ResourceOf` | `func ResourceOf(def *Definition) Resource` | `webtyp.com/model` v0.2.12 |
| `rpc.WriteError` | `func WriteError(ctx router.Context, err error)` | `webtyp.com/rpc` (already used here) |
| `model.KindOf` | `func KindOf(err error) ErrorKind` | `webtyp.com/model` |

## Tools

```
go get webtyp.com/model@v0.2.12
```

## Stage 1 — `Deps.TenantID` (`module.go`)

- `Deps` gets `TenantID string // requerido — el tenant de esta instalación; las operaciones nunca lo leen de la petición`.
- `Module` gets `tenantID string`; `New` returns `fmt.Err("service_catalog: Deps.TenantID is required")`
  when it is empty (same place as the `Deps.IDs` check) and stores it.

## Stage 2 — handlers use `m.tenantID` (`ops.go`, `recipe.go`)

- Every `args.TenantId` in a handler becomes `m.tenantID` (the 10 lines listed above).
- Every handler that saves a record received from the request sets its tenant **before** calling the
  domain method, overwriting whatever came: `spec.TenantId = m.tenantID` in `opUpsertSpecialty`
  (each record of `args.Records`), and the same in `opCreateService`, `opUpdateService`,
  `opUpsertService`, `opUpsertServiceAgreement`, `opUpsertServiceSupply`. Read each handler: if it
  decodes a wrapper (`UpdateServiceArgs`, `UpsertServiceArgs`), set it on the record inside.
- `opListCurrencies`: if it reads `args.TenantId`, stop.

## Stage 3 — args without `tenant_id` (`model.go`, then `ormc`)

- Remove the `{Name: "tenant_id", …}` field from every `*ArgsModel`: `ListServiceSuppliesArgsModel`,
  `DeleteServiceSupplyArgsModel`, `ListServicesArgsModel`, `GetServiceArgsModel`,
  `FindServiceBySKUArgsModel`, `DeactivateServiceArgsModel`, `DeleteServiceArgsModel`,
  `ListAgreementsArgsModel`, `DeleteAgreementArgsModel`, `ListCurrenciesArgsModel`,
  `ListSpecialtiesArgsModel`, `GetSpecialtyArgsModel`, `DeleteSpecialtyArgsModel`.
- **Do not** touch `tenant_id` of the tables (`SpecialtyModel`, `ServiceModel`,
  `ServiceAgreementModel`, `ServiceSupplyModel`): the rows keep their tenant.
- Regenerate `model_orm.go` with the repo's `ormc`; an args type left with no fields stays (an empty
  args model is valid).
- If an args struct becomes empty and its handler then has nothing to decode, keep the decode (the
  request still sends `{}`).

## Stage 4 — permissions from the tables (`ops.go`, `recipe.go`, `ui/browser.go`)

| Today | Becomes |
|---|---|
| `"specialty"` | `model.ResourceOf(&SpecialtyModel)` |
| `"service"` | `model.ResourceOf(&ServiceModel)` |
| `"service_agreement"` | `model.ResourceOf(&ServiceAgreementModel)` |
| `resServiceSupply` (delete the constant, `recipe.go:18`) | `model.ResourceOf(&ServiceSupplyModel)` |
| `Permission: "service"` / `"specialty"` in `ui/browser.go` | `model.ResourceOf(&servicecatalog.ServiceModel)` / `(&servicecatalog.SpecialtyModel)` |

After this stage `grep -rn 'Requires("' --include=*.go .` and `grep -rn 'Permission: "' ui/` return
nothing.

## Stage 5 — tests (`tests/`)

- Every `servicecatalog.New(db, servicecatalog.Deps{…})` gets `TenantID: "tenant-1"` (or the tenant
  the test already uses).
- Tests that send `TenantId` in **args** stop sending it.
- `tests/tenant_test.go`: keep `TestTenantIsolation` (it calls the Go methods with explicit tenants,
  still valid). Add `TestTenantFromServer`: build module A with `TenantID: "tenant-A"` and module B
  with `TenantID: "tenant-B"` on the **same** `db`; save a specialty through A's operations, list
  through B's operations → empty; and an upsert through A's operation whose record says
  `TenantId: "tenant-B"` is stored with `tenant-A`. Call the operations as the existing
  `tests/ops_test.go` does.
- New `TestNew_TenantRequired`: `Deps{IDs: …}` without `TenantID` → error text
  `"service_catalog: Deps.TenantID is required"`.
- Replace `err != error(servicecatalog.ErrX)` comparisons you touch with `model.KindOf(err)` checks.

## Stage 6 — demo and docs

- `web/` demo (if it builds the module): pass `TenantID`.
- `README.md` (Spanish): `Deps.TenantID` es obligatorio; las operaciones nunca leen el tenant de la
  petición; los permisos se nombran con `model.ResourceOf(&Tabla)`.

## Delivery check

```
GOPRIVATE=webtyp.com/*,github.com/veltylabs/* GOFLAGS=-mod=mod gotest
```

- A PR that does not compile, or whose tests could not run, is a failed delivery.
- Never pause to ask permission: when the check passes, commit, push and open the PR. Stuck 30
  minutes: open the PR listing what fails.
- `git diff --stat main` lists only: `module.go`, `ops.go`, `recipe.go`, `model.go`, `model_orm.go`,
  `ui/browser.go`, files under `tests/`, `web/` (only if it builds the module), `README.md`,
  `go.mod`, `go.sum`, this plan. No stray scripts.
- `grep -rn "args.TenantId" --include=*.go . | grep -v tests/` returns nothing.

## Stages

| # | Files | Done when |
|---|---|---|
| 1 | `module.go` | compiles |
| 2 | `ops.go`, `recipe.go` | compiles |
| 3 | `model.go`, `model_orm.go` | compiles |
| 4 | `ops.go`, `recipe.go`, `ui/browser.go` | both greps empty |
| 5 | `tests/*` | suite green |
| 6 | `web/`, `README.md` | text present |
