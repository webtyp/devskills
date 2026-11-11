---
name: form-codegen
description: Model authoring and form generation across webtyp/model, webtyp/definitionc, webtyp/form, webtyp/json, and webtyp/dom. Use when creating or modifying models/DTOs, input widgets, definitionc code generation, field validation, widget assignment, or form API design.
---

# Models + Forms: typed `model.Definition` → definitionc → everything else

> **Designing or changing the form API, a widget contract, or any exported
> symbol across `model`/`orm`/`form`/`json`/`dom`?** Skill **api-design** is the
> gate that must pass first, and its five answers belong in the `docs/PLAN.md`
> before code.

## The ONE pattern (source of truth)

A model is authored as a **typed `model.Definition` literal**. The var name
MUST end in `Model` — that is how `definitionc` discovers it:

```go
import (
    "webtyp.com/input"
    "webtyp.com/model"
)

var ServiceModel = model.Definition{
    Name: "service",
    Fields: model.Fields{
        {Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true},
        {Name: "tenant_id", Type: model.Text(), NotNull: true},
        {Name: "specialty_id", Type: input.Select(), Label: "Specialty", Ref: &SpecialtyModel, DB: &model.FieldDB{RefColumn: "id"}, NotNull: true},
        {Name: "name", Type: input.Text(), Label: "Name", NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 255}},
        {Name: "price", Type: input.Decimal(), Label: "Price"},
    },
}
```

A transport format (a JSON-RPC message, a webhook payload) sets `Wire: true`: its field names are protocol, not screen text — **langc** does not extract them and it is never given to a form. `definitionc` generates `Wire() bool` (`model.Wirer`) and `form.New` refuses it.

- **`Type` is a `model.Kind`** — it fixes storage, semantic name and validation in one value. An `input.Xxx()` Kind
  (`webtyp.com/input`, it embeds `model.Kind`) **is** the form control; a plain `model.Text()`/`model.Int()`… Kind gets
  **no** input (ids, tenant, timestamps). There is no `Widget` slot.
- **`Label`** (English key, translated) and **`Help`** sit on the field; `Ref` is a foreign key (`FieldDB.RefColumn`).
- **DB tables** carry `DB: &model.FieldDB{...}`; **form-only DTOs** (login, filters) carry none.
- **A listed model marks its list item on the fields**: `Item: model.ItemLabel` (one), `model.ItemDescription` (several
  are joined with « · » in field order), `model.ItemBadge` (one; a date goes as seconds). `definitionc` generates a typed
  `func (m *X) Item() model.Item` and rejects an invalid mark set (`model.ValidateItem`); never write `Item()` by hand.
  A select fed by another list without a foreign key says `Options: &OtherModel` (`Ref` when it is a foreign key).

Then generate:

```bash
go generate ./...   # runs definitionc (//go:generate definitionc); install:
go install webtyp.com/definitionc/cmd/definitionc@latest
```

`definitionc` parses the Definition literal (including each field's Kind expression)
and **generates the concrete Go struct** plus `ModelName()`, `Schema()`,
`Pointers()`, `Validate()`, `EncodeFields()`/`DecodeFields()` (typed codec)
and the `*List` type, into the generated `definition_gen.go` file. **Never edit
generated files — they are overwritten on every run.**

## FORBIDDEN (these are the bugs this skill exists to prevent)

- ❌ **Struct tags as source of truth** (`input:"email"`, `db:"pk"` on
  hand-written structs). The Definition literal is typed; tags are the old,
  removed pattern.
- ❌ **Hand-writing the struct or any generated method**
  (`Schema`, `Pointers`, `ModelName`, `Validate`, `Encode/DecodeFields`).
- ❌ **Writing a declared name as a string anywhere but its declaration.** An entity, field or effect `Name` is text exactly once, in `definition.go`
  (`definitionc` validates it: snake_case, unique, its files exist); every other use is the generated type or constant (`Payment_.PayURL`, `PaymentCharge{}`,
  `&Payment{}`), so a typo does not compile. An effect is declared inside the entity that owns it (`Effects:` next to `Fields:`), never as a loose variable.
- ❌ **Stdlib `encoding/json`** anywhere. Only `webtyp/json`, which works
  exclusively through the generated typed codec — a hand-written DTO without
  generated `Encode/DecodeFields` cannot travel.
- ❌ **`form.RegisterInput` as a fix for empty forms.** `form.New` binds
  inputs ONLY from fields whose `Type` is an `input.Input` Kind — there is no name matching.
  If a form renders without fields, those fields carry a plain `model.Xxx()` Kind
  → fix the Definition and regenerate; never patch at the consumer.

## Library roles

| Library | Responsibility | Knows about |
|---|---|---|
| **model** | `Definition`, `Field`, `Kind`, `Fielder`, `Permitted`, `ValidateFields`, typed codec contracts | nothing else |
| **orm** / **definitionc** | `orm`: DB at run time · `definitionc`: THE code generator (`definition.go` → struct + methods, `definition_gen.go`) | model |
| **form** | `form.New(parentID, &Generated{})` — schema-driven form; SSR + reactive render | model, dom, webtyp/input |
| **input** | Concrete form Kinds with validation (`Email`, `Password`, `Phone`, `Text`, `Number`, `Checkbox`, `Textarea`, …) | model |
| **json** | Zero-reflection codec over generated `Encode/DecodeFields` | model |
| **dom** | Pure HTML layout elements + signals. **NO form functions** (no `Input()`, `Form()`, …) | — |

STATUS (remove this line when ARCHITECTURE_MASTER_PLAN closes): decided, not yet shipped (D28–D45) — a module's `definition.go` also declares `model.ModuleDefinition{Name, Entities}`, per entity `Rules` (every action in `Allow` is `model.AddMyBusinessRule` = I write its rule, or `model.HaveNoBusinessRule`; undeclared is an error, judged by the `definitioncheck` library, never by `model`), `MFA` (a subset of `Allow`) and, when it calls outside, `Effects` inside the entity (`model.WaitForReply`/`RetryInBackground`/`WaitForConfirmation`); `definitionc` generates the module's `MountOperations`, `Definition()` and creates once `rule_<entity>_<verb>.go` / `effect_<entity>_<effect>.go`, whose body is a `COMPLETA TU REGLA DE NEGOCIO AQUÍ` comment plus output variables declared and unused, so the module does not compile until the agent fills them; a rule answers `(crud.RuleAnswer, error)` with `crud.ApplyAndPublish` or `crud.ApplyWithoutPublishing`; the module keeps one `Store<Entity>() crud.Store` per table entity; `crud` validates first, stores, publishes after commit and runs effects; hand-written `ops.go` mounting disappears.

STATUS (remove this line when ARCHITECTURE_MASTER_PLAN closes): `ormc` was renamed `definitionc` (input `definition.go`, output `definition_gen.go`) and is in local execution — until it is published, `go install webtyp.com/definitionc/cmd/definitionc@latest` does not exist yet and the published `ormc` still reads `model.go`.

STATUS (remove this line when ARCHITECTURE_MASTER_PLAN closes): decided, not yet shipped — `model.Field.ServerOnly` (rejected by `crud` on create/update, hidden by `form`) and `definitionc` generating `model_crud.go` (`//go:build !browser`) with typed `Mount<Entity>*` functions.

## Validation flow (same code front + back)

```
model.ValidateFields(action, m)
  → per field: PK/NotNull → the Kind's Validate (semantic) → Permitted (chars/length)
```

Runs identically in WASM (`form.Validate()`) and on the server (after decode,
before persistence).

**One validator per kind of value, in `webtyp/input`.** Code that checks a value outside a form (an importer, a
login, a lookup) calls the input's exported function — `input.NormalizeRUT`, `input.RutCheckDigit` — never a copy,
and never receives it injected (`Deps.ValidateRUT`): a copy or a fake drifts, and a lookup that normalizes
differently from the save creates duplicates (Pa100T: 354 repeated patients). A validator that fails is fixed
there, once.

## How to instruct an agent in a PLAN.md

1. Write the `model.Definition` literal(s) — var name ending in `Model`,
   every form-visible field with an `input.Xxx()` Kind.
2. Delete any hand-written struct the Definition replaces (definitionc generates it;
   names would collide).
3. Run `go generate ./...` (definitionc). Verify the generated schema carries the
   input Kinds and codecs.
4. Regression test: `form.New(id, &Generated{})` yields exactly the expected
   number of inputs — catches a future regeneration that loses widgets.

## Custom inputs

A project may define a custom input (implement `input.Input`, which embeds
`model.Kind`) and use it in its Definitions like any `input.Xxx()`. It lives
with the consumer; the shared ones live in `webtyp.com/input`.
