---
name: components
description: Creating a component in webtyp/components — the two-word naming rule, file layout, widget identity (Name/Part/Kind), conformance registration, and tests. Use when creating or reviewing a UI component. Delegates the dom tree to dom-elements, the stylesheet to widget-styling, and icons/images to ui-assets.
---

# WebTyp Component Creation

`webtyp.com/components` — one sub-package per component: `webtyp/components/<name>/`.

This skill covers **what a component is and how its package is assembled**.
Three neighbours own the rest, and they are not repeated here:

| Concern | Skill |
|---|---|
| Building the dom tree — embedding, element ids, events, state binding | **dom-elements** |
| The stylesheet — `style` DSL, parts, states, tokens | **widget-styling** |
| Icons and images — `svg.go`, `image.go` | **ui-assets** |
| stdlib replacements in WASM code | **wasm** |
| CSS-first interactivity, SRP, DI | **core-principles** |

> **Adding or changing a `widget.Part`, a `style` Option, an exported field, or
> any symbol outside this package?** Skill **api-design** is the gate that must
> pass first, and its five answers belong in the `docs/PLAN.md` before code.
> A recipe `widget/style` lacks is a defect **there**, not a licence to
> hand-compose one locally.

## Naming — two words, and the second must name the generic class

A component's package/folder name **and** its Go struct name must be **at least
two lowercase words**, and the pair must say *which style of the thing* it is —
never the generic noun alone.

```
✅ modaldialog/ ModalDialog    themetoggle/ ThemeToggle    actionbutton/ ActionButton
✅ decktabs/    DeckTabs       calendarslider/ CalendarSlider
❌ dialog/      Dialog         toggle/ Toggle    button/ Button    tabs/ Tabs
```

The **second word is the generic UI concept** (dialog, toggle, button, tabs,
slider, card, table, grid, bar, editor, nav). The **first word is the
characteristic** that distinguishes this implementation from every other one of
that concept.

**Why:** a bare generic name claims the whole concept for one implementation.
The day a consumer needs a drawer instead of a centred modal, a segmented
switch instead of a click-to-cycle toggle, or tabs that unmount their inactive
panels instead of keeping them, there is no name left — only a breaking rename.
Naming the class up front keeps every style addressable and coexisting:
`ModalDialog` beside `DrawerDialog`, `DeckTabs` beside a future `LazyTabs`.

**Pick the characteristic that survives its own settings.** `DeckTabs` is named
for the deck — every panel stays mounted as a layer — not for the slide, because
`MotionNone` removes the slide and the component is still the same thing. A name
describing a setting is wrong as soon as someone changes it.

Ask: *"what specific variant is this, and what would a different one be
called?"* If you cannot name the sibling that would coexist with it, the name is
still too generic.

> Violated twice: `dialog`/`DialogWidget` → `modaldialog`/`ModalDialog`, and
> `tabs`/`Tabs` → `decktabs`/`DeckTabs`. Both caught after the code was written;
> the rule is cheaper applied first.

## File structure

```
webtyp/components/
└── rangeslider/                 # two words: <characteristic><generic class>
    ├── rangeslider.go           # Struct, Render(), interactivity
    ├── css.go                   # //go:build !wasm — see widget-styling
    ├── svg.go                   # //go:build !wasm — see ui-assets (only if it ships glyphs)
    ├── image.go                 # see ui-assets (only if it ships rasters)
    ├── rangeslider_test.go      # //go:build !wasm
    └── README.md
```

**The filename is the contract.** CSS, SVG, JS, heavy HTML and image
declarations live in `css.go`, `svg.go`, `js.go`, `html.go` and `image.go`; a
file with any other name is invisible to the tooling that scans for it.
`ssr.go` is an eliminated convention, and there are no `.css` files.

No `front.go`: WASM interactivity lives in the main file — TinyGo eliminates it
as dead code on SSR builds.

## Identity — `widget.Name`, `widget.Part`, `widget.Kind`

A component declares its identity once; every class derives from it. **No
component writes a class string by hand.**

```go
const NameRangeSlider = widget.Name("rangeslider")

const (
    PartTrack = widget.Part("track")
    PartThumb = widget.Part("thumb")
)

var (
    clsRoot  = NameRangeSlider.Root()
    clsTrack = NameRangeSlider.Class(PartTrack)
    clsThumb = NameRangeSlider.Class(PartThumb)
)

func (c *RangeSlider) WidgetName() widget.Name { return NameRangeSlider }
func (c *RangeSlider) WidgetKind() widget.Kind { return widget.Region }
```

Apply a class with `.Set(cls.AsAttr())`. **`.Class("rangeslider__track")` is
forbidden** — a literal drifts from the stylesheet silently, and the conformance
test cross-checks the classes a stylesheet emits against those the markup
renders.

`WidgetKind` is not decoration: `Kind.Allows(State)` decides which states are
meaningful for it, and `TestKindAllowsEveryState` enforces that every state a
component uses is admitted by its Kind. Pick the Kind that describes the thing
(`widget.Tabs`, `widget.Listbox`, `widget.Dialog`, `widget.Menu`,
`widget.Combobox`, `widget.Form`, `widget.Toolbar`, `widget.Grid`,
`widget.Alert`, `widget.Disclosure`, `widget.Region`), not the convenient one.

## Render

```go
package rangeslider

import (
    . "webtyp.com/dom"
    . "webtyp.com/html"
)

func (c *RangeSlider) Render() *Element {
    return Div().Set(clsRoot.AsAttr()).
        Child(Span().Set(clsTrack.AsAttr()))
}
```

Dot-import `webtyp.com/dom` and `webtyp.com/html` — the house style across every
existing component. Everything about ids, events and state binding: skill
**dom-elements**.

## Register in the conformance test

`components/conformance_test.go` enumerates every component **twice**, and one
missing from either list is silently unverified. Add all three:

1. the import, alphabetically;
2. `&rangeslider.RangeSlider{},` in the `components := []interface{ RenderCSS() *css.Stylesheet }{…}` slice;
3. `"rangeslider": &rangeslider.RangeSlider{},` in `TestKindAllowsEveryState`'s map.

## Tests

```go
//go:build !wasm

package rangeslider

import (
    "strings"
    "testing"
)

func TestRangeSlider_Render(t *testing.T) {
    html := (&RangeSlider{}).Render().String()

    if !strings.Contains(html, "class='rangeslider'") {
        t.Error("expected rangeslider class")
    }
}
```

A `_test.go` carrying `//go:build !wasm` is backend-only and may use stdlib
(`regexp`, `strings`, `testing`). Attributes serialize **single-quoted**. Always
include a test that calls `RenderCSS()` — it validates the stylesheet and panics
on an invalid composition. Never assert on a composed element id (see
**dom-elements**).

## When to create a `PLAN.md`

Any code change — a new component or a modification — requires `docs/PLAN.md` at
the module root for user review before dispatching to an external agent.
Documentation-only changes (`.md`) can be edited directly. See skill
**agents-workflow**.
