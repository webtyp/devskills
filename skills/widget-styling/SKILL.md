---
name: widget-styling
description: Writing a stylesheet with the webtyp/widget/style DSL — style.For builder, parts and states, composed recipes over hand-composition, allowed tokens, and the rule that a missing recipe is a defect in widget/style. Use when writing or reviewing any css.go.
---

# Stylesheets — the `style` DSL

A stylesheet in this ecosystem is **Go**, not CSS. There are no `.css` files
and nothing is `//go:embed`-ed.

## Where it lives

`css.go`, with `//go:build !wasm`. The filename is the contract — any other
name is invisible to `webtyp/ssr`. (`ssr.go` is an eliminated convention.)

```go
//go:build !wasm

package rangeslider

import (
    "webtyp.com/css"
    "webtyp.com/widget"
    "webtyp.com/widget/style"
)

func (c *RangeSlider) RenderCSS() *css.Stylesheet {
    return style.For(c).
        Root(
            style.Row(style.Space4),
            style.As(style.Inset),
            style.Pad(style.Space2),
        ).
        Part(PartItem,
            style.Row(style.Space2),
            style.CenterContent(),
        ).
        When(widget.Current, PartItem,
            style.As(style.AccentInverse),
        ).
        Cue(widget.Hover, PartItem,
            style.As(style.AccentHover),
        ).
        Stylesheet()
}
```

`RenderCSS()` returns `*css.Stylesheet` — never a `string`.

## Builder methods on `*style.Sheet`

`Root`, `Part`, `Within`, `When`, `WhenWithin`, `Cue`, `CueWithin`,
`CueWithinHover`, `CueAcross`, `StateAcross`, `On`, `OnlyOn`, then
`Stylesheet()`.

- `Part(p, …)` styles one `widget.Part`.
- `When(state, part, …)` reacts to a `widget.State` bound in the markup
  (skill **dom-elements**).
- `Cue(cue, part, …)` reacts to an interaction cue (`widget.Hover`, …).
- `On(device, …)` / `OnlyOn(device, …)` scope to a device.

**`Stylesheet()` runs `Validate()` and panics on an invalid composition.** A
test that merely calls `RenderCSS()` therefore also asserts the stylesheet is
legal — always have one.

## Use the composed recipe, not its parts

When `style` ships a recipe for the thing, use it. The parts are frequently
**not** safely composable by hand, and the recipe's doc says so.

`style.Button(Surface)` is *the* recipe for anything the user presses. Its doc
records why: `KeepSize()` looks like the way to stop a button filling its panel
and is not — it emits flex-shrink/grow, which govern the **main** axis, while a
button inside a `Stack` is stretched across the **cross** axis by
`align-items: stretch`. Every component that composed its own reached a
different answer, and one shipped an 800px-wide "Add row" bar.

Same for `SlideDeck`, `MasterDetail`, `Sidebar`, `Flyout`, `Drawer`,
`ControlBox`, `IconBox`, `Scroll`, `Cover`.

A recipe that already derives a treatment does not want a second one:
`Button` derives hover, focus and press from its own family, so adding a
`Cue(widget.Hover, …)` on top is redundant composition and `Validate()` may
reject it.

## A missing recipe is a defect in `widget/style`

If the rule you need has no recipe, **stop**. That is a defect in
`webtyp/widget/style`, and it gets its own `docs/PLAN.md` in that repository.
It is never hand-composed at the call site. Fix at the root, never at the leaf
(skill **api-design**).

## Never hardcode a value

No hex, no `px`, no `rem`, no magic number. Everything comes from a recipe or an
allowed token. `components/conformance_test.go` holds the list — ~109 tokens:

```
--color-*      (primary/secondary/success/error/surface/outline/muted/hover/
                selection/disabled, each with on-*, -light, -dark, -hover,
                -focus, -press variants)
--text-*  --leading-*  --font-weight-*  --tracking-*
--space-0|1|2|3|4|6|8|12      --radius-sm|md|lg|full
--shadow-*  --duration-*  --ease-*  --z-*  --bp-*  --max-w-*
```

Anything else fails the conformance test.

Scales are typed, not raw: `style.Space2`, `style.TextSm`, `style.IconSm`,
`style.RadiusMd`, `style.MotionBase`, and the `style.Surface` family
(`Page`, `Panel`, `Inset`, `Primary`, `Secondary`, `Highlight`, `Accent`,
`AccentWash`, `AccentInverse`, `AccentHover`, `Success`, `Danger`,
`DangerWash`, `Subtle`, `Bare`, `Inactive`).

## Never declare `:root` or `RootCSS()`

Theme tokens are global state owned by the application (or `webtyp/dom`'s
default theme). A component or widget only **consumes** tokens. A third-party
`RootCSS()` is ignored with a warning — the single-override rule.

## Related skills

- **dom-elements** — binding the states this stylesheet reacts to
- **components** — where `css.go` sits in a component
- **api-design** — the gate for adding a `style` Option or a `widget.Part`
