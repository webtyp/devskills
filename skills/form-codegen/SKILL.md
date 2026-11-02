---
name: form-codegen
description: Model authoring and form generation across webtyp/model, webtyp/orm (ormc), webtyp/form, webtyp/json, and webtyp/dom. Use when creating or modifying models/DTOs, input widgets, ormc code generation, field validation, widget assignment, or form API design.
---

# Models + Forms: typed `model.Definition` → ormc → everything else

> **Designing or changing the form API, a widget contract, or any exported
> symbol across `model`/`orm`/`form`/`json`/`dom`?** Skill **api-design** is the
> gate that must pass first, and its five answers belong in the `docs/PLAN.md`
> before code.

## The ONE pattern (source of truth)

A model is authored as a **typed `model.Definition` literal**. The var name
MUST end in `Model` — that is how `ormc` discovers it:

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

A transport format (a JSON-RPC message, a webhook payload) sets `Wire: true`: its field names are protocol, not screen text — **langc** does not extract them and it is never given to a form. `ormc` generates `Wire() bool` (`model.Wirer`) and `form.New` refuses it. STATUS (remove this line when TUI_SIGNAL_MASTER_PLAN closes): local, unpublished `model`, `ormc`, `form`.

- **`Type` is a `model.Kind`** — it fixes storage, semantic name and validation in one value. An `input.Xxx()` Kind
  (`webtyp.com/input`, it embeds `model.Kind`) **is** the form control; a plain `model.Text()`/`model.Int()`… Kind gets
  **no** input (ids, tenant, timestamps). There is no `Widget` slot.
- **`Label`** (English key, translated) and **`Help`** sit on the field; `Ref` is a foreign key (`FieldDB.RefColumn`).
- **DB tables** carry `DB: &model.FieldDB{...}`; **form-only DTOs** (login, filters) carry none.
- **A listed model marks its list item on the fields**: `Item: model.ItemLabel` (one), `model.ItemDescription` (several
  are joined with « · » in field order), `model.ItemBadge` (one; a date goes as seconds). `ormc` generates a typed
  `func (m *X) Item() model.Item` and rejects an invalid mark set (`model.ValidateItem`); never write `Item()` by hand.
  A select fed by another list without a foreign key says `Options: &OtherModel` (`Ref` when it is a foreign key).

Then generate:

```bash
go generate ./...   # runs ormc (//go:generate ormc); install:
go install webtyp.com/orm/cmd/ormc@latest
```

`ormc` parses the Definition literal (including each field's Kind expression)
and **generates the concrete Go struct** plus `ModelName()`, `Schema()`,
`Pointers()`, `Validate()`, `EncodeFields()`/`DecodeFields()` (typed codec)
and the `*List` type, into the generated `*_orm.go` file. **Never edit
generated files — they are overwritten on every run.**

## FORBIDDEN (these are the bugs this skill exists to prevent)

- ❌ **Struct tags as source of truth** (`input:"email"`, `db:"pk"` on
  hand-written structs). The Definition literal is typed; tags are the old,
  removed pattern.
- ❌ **Hand-writing the struct or any generated method**
  (`Schema`, `Pointers`, `ModelName`, `Validate`, `Encode/DecodeFields`).
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
| **orm / ormc** | DB mapping + THE code generator (Definition → struct + methods) | model |
| **form** | `form.New(parentID, &Generated{})` — schema-driven form; SSR + reactive render | model, dom, webtyp/input |
| **input** | Concrete form Kinds with validation (`Email`, `Password`, `Phone`, `Text`, `Number`, `Checkbox`, `Textarea`, …) | model |
| **json** | Zero-reflection codec over generated `Encode/DecodeFields` | model |
| **dom** | Pure HTML layout elements + signals. **NO form functions** (no `Input()`, `Form()`, …) | — |

## Validation flow (same code front + back)

```
model.ValidateFields(action, m)
  → per field: PK/NotNull → the Kind's Validate (semantic) → Permitted (chars/length)
```

Runs identically in WASM (`form.Validate()`) and on the server (after decode,
before persistence).

## How to instruct an agent in a PLAN.md

1. Write the `model.Definition` literal(s) — var name ending in `Model`,
   every form-visible field with an `input.Xxx()` Kind.
2. Delete any hand-written struct the Definition replaces (ormc generates it;
   names would collide).
3. Run `go generate ./...` (ormc). Verify the generated schema carries the
   input Kinds and codecs.
4. Regression test: `form.New(id, &Generated{})` yields exactly the expected
   number of inputs — catches a future regeneration that loses widgets.

## Custom inputs

A project may define a custom input (implement `input.Input`, which embeds
`model.Kind`) and use it in its Definitions like any `input.Xxx()`. It lives
with the consumer; the shared ones live in `webtyp.com/input`.
