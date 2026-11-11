---
name: crud-module-migration
description: The answers an executor needs to migrate a veltylabs module to the crud-by-declaration framework (argument shapes, crud.Tx, read models, composite keys, tests through LoopbackActingAs, definitionc, the traps that stalled 9 of 12 sessions in the first wave) and how to split such a plan so it needs no questions. Use when writing, splitting or answering questions about a module-migration PLAN.md.
---

# Migrating a module to crud-by-declaration — the answers, so the plan carries them

Process around a plan: **plan-authoring**, **agents-workflow**. Reference module (copy its shape): `veltylabs/patient_directory` — `definition.go`, `rule_patient_*.go`, `store.go`, `import.go`, `module/server.go`, `seed/demo.go`, `tests/setup_test.go`, `tests/crud_test.go`, `tests/demo_test.go`.

A plan for a module migration pastes the blocks below **inside the stage that needs each one**, never as an appendix: an executor reads a stage, not a list of notes. Each block was a question that stalled a session on 2026-11-10.

## 1. Argument shapes (tests and callers)

- `crud.CreateArgs{Records []model.Model}` · `crud.UpdateArgs{IDs []string, Fields []string, Record model.Model}` · `crud.DeleteArgs{IDs []string}` · `crud.ReadArgs{ID, Search, After string, Limit int64, Filter model.Model}`. Never a bare entity pointer, a generated `*List` or a map. `crud.Change` has `Record model.Model` and `Fields []string` only.
- The module name in an operation is the constant `<pkg>.ModelName`, never a literal (`"cashmanager"` ≠ `"cash_manager"` and answers `unknown operation`); always assert the callback error before reading the result. The acting user of a test is the first argument after the publisher (`"u1"`): a read model that filters by `a.UserID` returns nothing for rows created for another user.
- A call is `crud.LoopbackActingAs(pub, "u1", mod).Call(router.Operation{Module: ModelName, Entity: "<entity>", Action: model.Create}, args, &into, func(err error){…})`; events reach `pub`, never `Deps.Publisher`. **The `id` of a table (a Text primary key) is ALWAYS required by `ValidateInput`, unless the field is `ServerOnly`.** Choose per table: (a) the server mints the id (a form creates rows without one): declare `{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true, ServerOnly: true}`; `crud` then stamps it and the payload (and the tests) OMIT the id; sending one answers `ErrServerOnlyField`; (b) the client brings it (a key that means something): no `ServerOnly`, and every create record carries it. `OmitEmpty` alone changes nothing: `id required`. A composite key (`room_id/user_id`) is always client-sent. Never send a `ServerOnly` field.
- An action outside `Allow` is not mounted: the error text contains `unknown operation`. A refusal is `loopback: <status> <detail>`; `model.KindOf(err)` gives `ErrorNotFound`/`ErrorInvalid`/`ErrorConflict`; the failed field is `err.(interface{ FailedField() string }).FailedField()`. A duplicate key is a conflict answered by `crud`, not by the rule.

## 2. Inside a rule

- `tx` is `crud.Tx`: `ReadOne(&rec, id)`, `ReadAll(newRow, crud.Query{Conditions, Order, After, Limit}, onRow)`, `Create`, `UpdateFields(&rec, id, cols)`, `Delete(&rec, id)`. The field is `Conditions` (not `Cond`); conditions are `storage.Eq/Neq/Gt/Gte/Lt/Lte/In/Like`, AND only; never a `tenant_id` condition (the store applies it). `time.Nano` values go as `int64(...)`. A read returns at most `crud.MaxPage` rows: loop when you delete in bulk.
- Rule signatures are the ones `definitionc` writes; keep them. In `Update`, `crud` writes `next` AFTER the rule with exactly the columns in `fields` (the ones the caller sent): assigning `next.X = v` changes what is stored ONLY for an `X` that is in `fields`. A column the rule derives and the caller did not send (a state, a computed total, a difference) is written by the rule itself with `tx.UpdateFields(&row, id, []string{"col", …})` before it returns; a rule that only assigns it leaves the row unchanged and the test sees the old value. Composite key id: `model.JoinKey(a, b)`.
- Logic shared by a rule and a kept method lives in one unexported function both call; a method never reimplements a rule.

## 3. Read models

`Allow: model.Read`, `Rules{Read: AddMyBusinessRule}`, `Wire: true` (or no `DB:` marks), listed in `Definition.Entities`; flat kinds only (`model.Pointer`/`List`/`Map` do not exist: one row per result); the values the caller sends are fields of the read model, `OmitEmpty`, and arrive in `filter` (`definitionc ≥ v0.1.42`: `Read(a, filter *X, args, m)`). A caller passes `&crud.ReadArgs{Filter: &X{…}}`. «cannot read the request» = the filter fails its own `Permitted` rules or the generated code is stale: encode the args and decode them back with `json.Decode` to read the failing field. Every model to mount, tables and read models, must be in `Entities`.

## 4. Tooling traps

- `definitionc` rewrites `go.mod` to older versions: after EVERY run repeat the plan's `go get` line and `go mod tidy`. Install tools and `export PATH="$PATH:$(go env GOPATH)/bin"`.
- Never edit `definition_gen.go`; never write `go.sum` lines; never script edits over the source (sed/regex/python broke imports in four repos): one file at a time, `gotest` after each.
- `gotest` must end green on EVERY destination (`wasm`, `browser`); `docs/TESTS.md` rows are `| [TestName](../tests/file_test.go) | rule |` under `| Test | Rule |` (no File column).
- Keep `module/` (an app imports it): only `appkit.ModuleBrowser` → `view.ModuleBrowser`. A table without `tenant_id` (a link table) uses `orm.GlobalStore`, not `TenantStore`. `webtyp.com/auth` is not needed for RUT: `input.NormalizeRUT`. Seeds load through an `Import<Entity>` method (`import.go`: set the id, `model.ValidateFields`, then `m.store.Atomic` create-or-update), never the router.
- `docs/PLAN.md` rule for the executor: open the PR after the third stage (Jules publishes only by submitting; a `git push` does not reach GitHub) and keep committing to it; never ask a yes/no about progress.

## 5. Splitting the plan (the 3-hour rule, applied)

A module with more than ~6 rules or ~1000 lines is split into stage files `docs/PLAN_STAGE_<N>_<TOPIC>.md`, each with exact code: 1 tools + `definition.go`; 2 `Module.store` + the two read helpers (`firstRow`, `allRows`) + a translation table «old `m.db` call → crud call»; 3 every rule body as Go; 4 readers and writers with a worked example; 5 browser + demo; 6 tests, each row as setup → call → expected; 7 docs. The index `PLAN.md` holds the frontmatter, the operation mapping, the stages table and the delivery check. The `chat_room` plan of 2026-11-10 is the worked example.

## 6. A module too big for one session, or a session that failed twice — split it in TWO plans, dispatch one, then the other

Applies when the module has more than ~8 tables, or more than ~1500 lines of logic, or a session already ended with an empty or gutted PR (deleted logic, deleted tests, stray scripts, tree that does not compile). The 7 stages do not shrink; the work is cut in **two plans, each staged and each ending with a tree that compiles and a green `gotest`**.

1. Close the bad PR with a comment (`gh pr close <n> --comment "<why>"`), `git checkout main`, remove `STATUS`/`SESSION` from the plan header. Save the old plan in the scratchpad.
2. **Plan 1 (dispatched now)** declares EVERY table and read model in `Definition` (the generated mount is complete from the start), does the whole easy half (its rules, methods on the store, UI, tests, docs) and cuts the hard writers out: the entity of a hard writer is `Allow: model.Read` only, and its writer methods stay **exactly as they are today** (same bodies, on `m.db`, untouched, no operation mounted). Read models whose logic lives in a kept method are rules that call that method. Say in the plan that those writes leave the wire until plan 2.
3. **Plan 2 (written now, parked as `docs/PLAN_<TOPIC>.md`, never touched while plan 1 runs)** turns the hard writers on: `Allow` with the rules, the method bodies moved onto the store with the shared-function pattern, their screens, their tests, the docs. It starts with «Plan 1 is published: `git pull`».
4. When plan 1 is closed, `mv docs/PLAN_<TOPIC>.md docs/PLAN.md`, adjust what plan 1 changed, `codejob dispatch`. Record both in the master (one row, two states: «plan 1 despachado · plan 2 escrito»).
5. Never cut by file or by layer (a half that does not compile is a failed delivery); cut by **writer**: what stays unmounted is a gap in the wire, never in the build.

## 7. Rehacer las reglas tras D58 (`crud.ApplyAndPublish()` función)

- En `go.mod`: actualizar `webtyp.com/crud` a `≥ v0.0.13`.
- En `rule_<entity>_<verb>.go`:
  - `return crud.ApplyAndPublish, nil` → `return crud.ApplyAndPublish(), nil` (o con llamadas si las pide).
  - `return crud.ApplyWithoutPublishing, nil` → `return crud.ApplyWithoutPublishing(), nil`.
  - En rechazos por error: `return 0, err` → `return crud.RuleAnswer{}, err`.
- En pruebas (`tests/`): reemplazar `return 0, nil` (prueba de rechazo por `ErrNoAnswer`) por `return crud.RuleAnswer{}, nil` y llamadas directas con `crud.ApplyAndPublish()`.
- Correr `gotest`.

## 8. External calls (`Definition.ExternalCalls`) — what the API really is (read from `crud` v0.0.14)

- A rule asks for a call with `return crud.ApplyAndPublish(<Entity><Call>{}), nil`; `definitionc` writes the type `<Entity><Call>` and the skeleton `call_<entity>_<call>.go`.
- The call is a method on the **stored record**: `func (p *X) Call(a crud.CurrentUser, m *Module, key crud.IdempotencyKey) (crud.FieldsToSave, error)`. It **sets fields on `p`** and returns their names: `p.Status = StatusDone; return crud.Fields("status"), nil`. `FieldsToSave` is `[]string`, never a map. On failure return `nil, err`: crud marks the attempt failed in the queue; the record keeps its old values (a `failed` status written by a call does not exist; the reason lives in the queue row).
- `WaitForReply` runs inside the request (loopback tests see the result at once). `RetryInBackground` runs later through `RunPendingCalls` with an **empty `CurrentUser`** (never use `a.UserID` in it); its tests call the method on the row directly.
- A module that declares calls needs a queue: tests use `crud.LoopbackActingAsWithQueue(nil, "u1", q, mod)` with a small `crud.CallQueue` double (`Add`, `MarkDone`, `MarkFailedAttempt`, `ConfirmCall`, `RejectCall`); the app passes `appkit.NewCallQueue(q)`.
