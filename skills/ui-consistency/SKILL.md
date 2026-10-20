---
name: ui-consistency
description: The visual rules every module screen obeys — form field layout (label above, 1 column phone / 2 desktop, validation on blur), one heading level, and the colour roles (primary = brand and actions; amber accent = selected AND focused, same in light and dark; outline = passing state, fill = confirmed; data tags neutral). Use when styling or reviewing any form, list row, tag, focus or selected state.
---

# ui-consistency — one product, many modules

Background, prior art and rejected alternatives: `webtyp/docs/UI_CONSISTENCY_DESIGN.md`.

STATUS (remove this line when UI_CONSISTENCY_MASTER_PLAN closes): not every rule has shipped yet; what is still rolling out is in `webtyp/docs/UI_CONSISTENCY_MASTER_PLAN.md`.

## Who draws what

- **Fields → `form`.** A module declares a model; `form.New(model)` draws it. Never a hand-written `Input(...)` or a field stylesheet in a module.
- **Module title → `platformd`** (`<h1>` in the header, large, beside the app's mark; the app's name and version live in the panel the mark opens, never beside the title). Content has only `<h2>` via `style.Heading()`; no component fixes a title's size or weight.
- A module's own pieces nobody reuses (floor plan, schedule editor) are styled in the module.

## A module explains itself (user rule, 2026-10-20)

A screen must be usable by someone who was never told how it works: no external explanation, no manual.

- **Words of the person, not of the code.** A screen, field or list is named for what it means to its user ("Quién recibe", "Área"), never for the internal structure it maps (`channel`, `rule`, `scope`, `assignment`). Every label and help text is in `lang.json` — an English label on screen is a defect.
- **No free-text codes.** A value the system defines (an alert type, an area, a unit) is chosen from a list with a human name, never typed.
- **One screen per task, not one per table.** Tables that belong to one task (a channel and its settings) are edited together.
- **A field whose meaning is not obvious carries `Field.Help`** (`model/field.go`), shown under it; obvious fields ("Name", "Active") do not.
- **Show where to start.** A module with prerequisites states its first steps; an empty list says what goes there and how to add the first one; a configuration problem is shown on the screen, never only through an operation.
- **A read-only list still offers the actions that make sense** (open the detail, resolve, go to the origin); a list that only displays is a defect.

The shared mechanism (module description, "?" button, welcome card, empty state of a view) is being decided in `webtyp/docs/MODULE_GUIDANCE_MASTER_PLAN.md`; until it ships, follow these rules with what exists.

## Form

- Label **above** the control: translated, never truncated, never on the border. Help, then error, below; the error replaces the help.
- **1 column on a phone, at most 2 on a desktop**, by the form's own width (`style.Grid(2, style.ColumnMedium, …)`); textarea and submit `SpanFull`. No width cap.
- Validate **on blur and on submit**, never per keystroke.
- Every field with a widget declares `Label` (English key). A lone checkbox carries its text once, inside the option (no label above, full row). A select with no value shows a prompt option. A radio group names itself with `role=radiogroup` + `aria-labelledby`; never render a title equal to the field name.
- `Radio` = chips. The renderer writes **no colour** (`--color-*` in markup is a defect; `form/tests` guards it).

## Colour roles

- **Primary** (gradient by default, solid if the theme clears it) = brand and **actions**.
- **Amber accent** = **selected and focused**, one value in light and dark. Never a second colour in the same indicator.
- **Focus** = the global 2 px amber inset ring. The label does not change colour on focus (amber text on the light card fails 1.4.3).
- **Outline = passing state, fill = confirmed:** selected and pressed are an amber **fill**; hover is a **1 px** amber outline, focus **2 px**. A hover never paints the content.
- **Data** (codes, counts, tags) is **neutral** (outlined or tonal), never primary-filled. In a list row it rides the bottom-right edge out of the flow (`OnEdge`), so the label keeps the full width and the tag never wraps.
- Recipes: `style.AccentOutline` for a hover/focus outline that repaints nothing; `style.Picker()` (on a part carrying `widget.PartSelect`) for a `<select>`'s native list, which only `appearance: base-select` can recolour — native elsewhere.
- A colour the theme should own is never written in markup or `css.go`: use a recipe, and a missing recipe is a `widget/style` defect (**widget-styling**).

## Messages and the phone strip

- A toast is `As(Inverse)` on every device and severity; severity is `role=status|alert`, never colour. On a phone it docks below the header (`DockedClear`), never over it. `platformd.Auto()` floors at 4 s.
- An error toast stays until dismissed (`Notify` makes `Auto()` persistent for `Error`); a message names the action, then the cause ("Could not save: …"); an identical one is not stacked.
- On a phone the detail panel of the list/detail strip takes the full width (`MasterDetail(Full)`): no sliver of the list beside the form.
