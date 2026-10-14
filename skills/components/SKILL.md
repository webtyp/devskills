---
name: components
description: Creating a component in webtyp/components — the two-word naming rule, file layout, widget identity (Name/Part/Kind), conformance registration, and tests. Use when creating or reviewing a UI component.
---

# Components

One sub-package per component: `webtyp/components/<name>/`. Neighbours: dom tree **dom-elements** · stylesheet **widget-styling** · icons **ui-icons** · images **ui-images** · WASM stdlib **wasm**. A new `widget.Part`, `style` option or exported field: **api-design** first; a recipe `widget/style` lacks is a defect there, never hand-composed here.

## Naming — two words: `<characteristic><generic class>`

Package and struct: `modaldialog/ModalDialog`, `themetoggle/ThemeToggle`, `actionbutton/ActionButton` — never `dialog/Dialog`, `button/Button`.

- The second word is the generic concept (dialog, toggle, tabs, card, grid…); the first distinguishes this implementation, so siblings can coexist (`ModalDialog` beside `DrawerDialog`).
- Name the trait that survives its settings (`DeckTabs` = panels stay mounted, not "slide": slide can be turned off).
- Test: if you cannot name the sibling that would coexist, the name is too generic.

## File structure — the filename is the contract

```
rangeslider/
├── rangeslider.go       # struct, Render(), interactivity (no front.go)
├── css.go               # //go:build !wasm
├── svg.go               # //go:build !wasm, only if it ships glyphs
├── image.go             # only if it ships rasters
├── rangeslider_test.go  # //go:build !wasm
└── README.md
```

`css.go`/`svg.go`/`js.go`/`html.go`/`image.go` are found by name; any other name is invisible. No `ssr.go`, no `.css` files.

## Identity — no class string by hand

```go
const NameRangeSlider = widget.Name("rangeslider")
const PartTrack = widget.Part("track")
var clsRoot, clsTrack = NameRangeSlider.Root(), NameRangeSlider.Class(PartTrack)

func (c *RangeSlider) WidgetName() widget.Name { return NameRangeSlider }
func (c *RangeSlider) WidgetKind() widget.Kind { return widget.Region }
```

- Apply with `.Set(cls.AsAttr())`; `.Class("rangeslider__track")` is forbidden.
- Pick the `Kind` that describes the thing: `TestKindAllowsEveryState` checks every state used is allowed by it.
- Dot-import `webtyp.com/dom` and `webtyp.com/html` (house style).

## Register in `components/conformance_test.go` — all three

1. the import (alphabetical);
2. `&rangeslider.RangeSlider{},` in the `RenderCSS()` slice;
3. `"rangeslider": &rangeslider.RangeSlider{},` in `TestKindAllowsEveryState`.

## Tests

Only behaviour (keyboard, selection, events, ARIA contract) — never the look (**testing** → «What deserves a test»). The stylesheet is validated by the conformance registration above, not by a per-component test.
