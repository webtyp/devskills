---
name: components
description: WebTyp component creation standards for webtyp/components. Use when creating or reviewing UI components: file structure, embedding rules, identity and element IDs, the style DSL, SSR split, icons (webtyp/icons + sprite), images (image.go + webtyp/image), and PLAN.md workflow.
---

# WebTyp Component Creation

> **Adding or changing a `widget.Part`, a `style` Option, an exported field, or
> any symbol outside this package?** Skill **api-design** is the gate that must
> pass first, and its five answers belong in the `docs/PLAN.md` before code.
> This document covers how a component is assembled; that one covers what may be
> added to the vocabulary it assembles from.
>
> When a component needs a recipe that `widget/style` does not have, it has
> found a defect **in `widget/style`** — not a licence to hand-compose one
> locally. See api-design's "fix at the root, never at the leaf".

## Module

`webtyp.com/components` — located at `webtyp/components/`.

Each component lives in its own sub-package: `webtyp/components/<name>/`.

## File Structure (per component)

```
webtyp/components/
└── mycomponent/
    ├── mycomponent.go        # Struct, Render(), interactivity
    ├── css.go                # //go:build !wasm — RenderCSS() *css.Stylesheet
    ├── svg.go                # //go:build !wasm — IconSvg() *sprite.Sprite (only if it ships glyphs)
    ├── image.go              # RenderImages() []image.Asset (only if it ships rasters)
    ├── mycomponent_test.go   # //go:build !wasm
    └── README.md
```

**Extension-named files, never `ssr.go`.** CSS, SVG, JS, heavy HTML and image
declarations live in `css.go`, `svg.go`, `js.go`, `html.go` and `image.go`. The
`ssr.go` convention is eliminated — **the filename is the contract**, and a file
with any other name is invisible to the tooling that scans for it.

**There are no `.css` files.** A stylesheet is Go, built with the `style` DSL
(below). Nothing is `//go:embed`-ed.

No `front.go`. WASM interactivity lives in the main file — TinyGo eliminates it
as dead code on SSR builds.

## Struct Definition — Embedding Rule

Always embed `dom.Element` as a **value**, never as a pointer.

```go
// ✅ CORRECT
type MyComponent struct {
    Element
    Title string
}

// ❌ WRONG — double heap allocation, nil-guard boilerplate, GC pressure in TinyGo
type MyComponent struct {
    *Element
}
```

Why: TinyGo (final WASM compiler) has a simple GC. Fewer heap objects = fewer
pauses. Value embedding = 1 allocation instead of 2, better cache locality,
zero nil panic risk.

## Identity: `widget.Name` and `widget.Part` — never a class literal

A component declares its identity once; every class is derived from it. **No
component writes a class string by hand.**

```go
const NameMyComponent = widget.Name("mycomponent")

const (
    PartItem = widget.Part("item")
    PartIcon = widget.Part("icon")
)

var (
    clsRoot = NameMyComponent.Root()
    clsItem = NameMyComponent.Class(PartItem)
    clsIcon = NameMyComponent.Class(PartIcon)
)

func (c *MyComponent) WidgetName() widget.Name { return NameMyComponent }
func (c *MyComponent) WidgetKind() widget.Kind { return widget.Region }
```

`WidgetKind` is not decoration: `Kind.Allows(State)` decides which states are
meaningful for it, and the repository's `TestKindAllowsEveryState` enforces
that every state a component uses is admitted by its Kind. Pick the Kind that
describes the thing (`widget.Tabs`, `widget.Listbox`, `widget.Dialog`,
`widget.Form`, `widget.Region`, …), not the one that is convenient.

Apply a class with `.Set(cls.AsAttr())`. **`.Class("mycomponent__item")` is
forbidden** — a literal drifts from the stylesheet silently, and the
conformance test cross-checks the classes a stylesheet emits against the
classes the markup renders.

## Element IDs are MINTED by `dom` — never composed by a component

This is the rule most often got wrong, and it fails silently or crashes.

`dom` assigns ids itself. `Element.GetID()` returns the element's id, **minting
a unique one on demand** if it has none. `Element.For(other)` does the same for
a `for=` attribute, and its doc names the wider use: *label/input pairing and
`aria-*` references*.

```go
// ✅ CORRECT — dom mints the ids, the component only references them
panel := Section().Set(clsPanel.AsAttr()).Attr("role", "tabpanel")
tab := Button().Set(clsTab.AsAttr()).Attr("role", "tab")

tab.Attr("aria-controls", panel.GetID())
panel.Attr("aria-labelledby", tab.GetID())
```

```go
// ❌ WRONG — a composed id, and a live bug
tab := Button().Attr("id", "tab-"+item.ID).
    Attr("aria-controls", "panel-"+item.ID)
```

Why the wrong version is not a style nit:

- **`dom` resolves every handler and every signal patch by id.** An event is
  wired to `#3`; a signal patches `#3`. Two nodes sharing one id is *"a
  component that renders, looks right, and does nothing, because the runtime
  resolved the other one"* (`dom/dom.go`, the "one id, one node" section).
- **`claimID` panics** when one render pass writes the same id twice. Two
  instances of the same component on one page, fed the same domain ids, take
  the whole render down.
- A minted id is unique by construction, so the same domain id may appear in as
  many instances as the screen needs.

Use `Attr("data-id", domainID)` when a test or a consumer needs the domain
identifier back. It is data, not something the framework resolves.

`Element.Key(k)` is **not** an id: it is the stable identity for keyed
reconciliation in `BindChildren`. Setting it on an element that is not inside a
`BindChildren` is inert noise — leave it off.

## Render Pattern

```go
package mycomponent

import (
    . "webtyp.com/dom"
    . "webtyp.com/html"
    "webtyp.com/widget"
)

func (c *MyComponent) Render() *Element {
    return Div().Set(clsRoot.AsAttr()).
        Child(Span().Set(clsItem.AsAttr()).Text(c.Title))
}
```

Dot-import `webtyp.com/dom` and `webtyp.com/html` — the house style across
every existing component.

**Interactivity binds on the element as it is built**, with the typed handlers
`Element` already exposes — `OnClick`, `OnChange`, `OnInput`, `OnKeyDown`,
`OnFocusIn`, `OnSubmit`, `OnToggle`, … — not a second pass that looks nodes up
by id:

```go
tab := Button().
    BindStateFunc(widget.Current, func() bool { return c.isActive(id) }).
    OnClick(func(Event) { c.activate(id) })
```

`BindState(state, signal)` / `BindStateFunc(state, fn)` are how a `widget.State`
reaches the markup; the stylesheet reacts to it with `When(...)`. State is never
a class.

## Stylesheets — the `style` DSL, in `css.go`

```go
//go:build !wasm

package mycomponent

import (
    "webtyp.com/css"
    "webtyp.com/widget"
    "webtyp.com/widget/style"
)

func (c *MyComponent) RenderCSS() *css.Stylesheet {
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

Builder methods on `*style.Sheet`: `Root`, `Part`, `Within`, `When`,
`WhenWithin`, `Cue`, `CueWithin`, `CueWithinHover`, `CueAcross`, `StateAcross`,
`On`, `OnlyOn`, then `Stylesheet()`.

**`Stylesheet()` runs `Validate()` and panics on an invalid composition.** A
test that merely calls `RenderCSS()` therefore also asserts the stylesheet is
legal — always have one.

Rules the DSL and the conformance test enforce:

- **Never hardcode a value.** No hex, no `px`, no `rem`. Everything comes from a
  recipe or an allowed token; `conformance_test.go` holds the list (~109 tokens:
  `--color-*`, `--text-*`, `--space-*`, `--radius-*`, `--shadow-*`,
  `--duration-*`, `--ease-*`, `--z-*`, `--bp-*`, `--max-w-*`) and fails on
  anything else.
- **A component MUST NOT define `:root`** and MUST NOT declare `RootCSS()`.
  Theme tokens are global state owned by the app or `webtyp/dom`'s default
  theme; a third-party `RootCSS()` is ignored with a warning.
- **Use the composed recipe, not its parts.** `style.Button(Surface)` is *the*
  recipe for anything the user presses — its own doc records that the parts are
  not safely composable by hand, and that a component which composed its own
  once shipped an 800px-wide button. Same for `SlideDeck`, `MasterDetail`,
  `Sidebar`, `Flyout`.
- **A missing recipe is a defect in `widget/style`**, and gets its own plan in
  that repository. Never hand-compose it locally.

## Register a new component in the conformance test

`components/conformance_test.go` enumerates every component **twice**, and one
missing from either list is silently unverified. Add all three:

1. the import, alphabetically;
2. `&mycomponent.MyComponent{},` in the `components := []interface{ RenderCSS() *css.Stylesheet }{…}` slice;
3. `"mycomponent": &mycomponent.MyComponent{},` in `TestKindAllowsEveryState`'s map.

## Icons — `webtyp/icons` glyphs + `svg.go`

An icon has two halves that must reach different places, and the split is the
whole design:

- the **reference** — the symbol id, a plain string. This is all that may reach
  the browser: `svg.Icon("trash")`.
- the **geometry** — the `<path>` data and viewBox. Backend-only: `webtyp/ssr`
  pulls it out at build time, injects it once into the page, and the markup
  points at it with `<use href="#trash">`. Shipping geometry to the browser
  would drag `webtyp/json` + `webtyp/model` into the WASM bundle for nothing.

### Prefer a shared glyph from `webtyp/icons`

`webtyp.com/icons` is the shared set, **one package per glyph** (`trash`,
`pencil`, `plus`, `undo`, `selectall`). Per-glyph packaging is what keeps each
glyph's geometry behind its own `//go:build !wasm`, so importing one glyph for
its `Ref` can never leak another's path data into a WASM build.

```go
// mycomponent.go — WASM-safe: only the id string crosses
import "webtyp.com/icons/trash"

trash.Ref.Render(string(clsIcon))   // <svg class=…><use href="#trash"/></svg>
```

```go
// svg.go — //go:build !wasm — hand the geometry to your sprite
import (
    "webtyp.com/icons/trash"
    "webtyp.com/icons/pencil"
    "webtyp.com/svg/sprite"
)

func (c *MyComponent) IconSvg() *sprite.Sprite {
    return sprite.NewSprite(trash.Def(), pencil.Def())
}
```

Adding a glyph to the shared set is a new folder in `webtyp/icons` — it never
touches the ones already there. Copy `trash/`, rename, swap the `Ref` id, and
paste the viewBox + path into `icons.Solid(Ref, viewBox, d)`. Keep it a solid
single-path glyph; anything else is not that family and calls `sprite.Define`
directly.

### A glyph only this component owns

Define it in the component's own `svg.go`:

```go
//go:build !wasm

package mycomponent

import "webtyp.com/svg/sprite"

func (c *MyComponent) IconSvg() *sprite.Sprite {
    return sprite.NewSprite(
        sprite.Define(IconPhone, "0 0 24 24",
            sprite.Path("M6.62 10.79a15.053 15.053 0 006.59 6.59l…"),
        ),
    )
}
```

with the reference declared in the main file:

```go
const IconPhone = svg.Icon("mycomponent-phone")
```

**`IconSvg()` returns `*sprite.Sprite` — never `map[string]string`.** Builders:
`sprite.NewSprite(defs…)`, `sprite.Define(icon, viewBox, body…)`,
`sprite.Path(d)`, `sprite.Raw(s)`, `(*Sprite).AddFile(id, svgFile)`,
`sprite.MergeAll(sprites…)`.

Hard rules:

- `IconSvg()` MUST live in `svg.go` behind `//go:build !wasm`. Geometry must
  never reach the WASM binary.
- Every path uses `fill="currentColor"` (or `stroke`). **Never hardcode a
  colour in SVG** — the box around the glyph owns the colour, so a white-text
  button yields a white glyph and a red box a red glyph.
- The class passed to `Ref.Render(class)` is where size and colour come from.
- The sprite is injected inline into the `<body>`; there is no
  `/assets/icons.svg` URL, so `href="#id"` always resolves with no network
  request. Never reference the sprite by URL from CSS.

## Images — `image.go` and the `webtyp/image` builders

Rasters are not icons: they go through an optimization pipeline that generates
responsive variants. Two halves again, and again the filename is the contract.

### Declaring what the pipeline must process — `image.go`

```go
package mycomponent

import "webtyp.com/image"

// RenderImages declares every raster this component serves. Read by the
// pipeline via AST: it must stay a hand-written composite literal — no loops,
// no computed paths, no variables.
func RenderImages() []image.Asset {
    return []image.Asset{
        {Path: "img/hero.jpg", Variants: image.AllVariants, Alt: "Hero"},
        {Path: "img/badge.png", Variants: image.VariantM, Alt: "Badge"},
    }
}
```

- **The filename MUST be `image.go`.** The pipeline detects it by name.
- **It is read by AST, not executed.** A loop, a computed path, or a value
  pulled from a variable extracts as nothing, silently. Keep it a literal.
- `Path` is relative to the module directory. `Alt` is SEO text; empty is
  derived from the filename.
- `Variants` is a bitmask: `VariantS` (480px), `VariantM` (1024px), `VariantL`
  (1600px), `AllVariants`. **Declare only the variants something actually
  renders** — each one is a file written to disk and shipped to the CDN, and a
  variant nothing links to is pure waste.

### Rendering — `webtyp.com/image` builders

These compile for **both** WASM and backend (no build tag), so they belong in
the component's main file:

```go
import . "webtyp.com/image"

// Responsive: emits srcset from the generated variants.
Responsive("/img/hero.jpg", "Hero").
    Sizes("(max-width: 600px) 100vw, 33vw").
    Lazy().
    AsElement()

// Plain image, no srcset:
Img("/img/logo.png", "Logo").Size(200, 50).AsElement()
```

`sizes` matters: without it the browser assumes `100vw` and over-downloads any
image that does not span the full width. Other methods: `Srcset`, `Class`,
`Attr`, `String`.

> `webtyp.com/image` shadows Go's stdlib `image` **on purpose** — the stdlib one
> is too heavy for TinyGo.

Sibling packages, none of which a component imports directly:
`image/min` (backend WebP pipeline), `image/browser` (client-side compression
before upload, `//go:build wasm`), `image/favicon` (logo → icon set).

## CSS-First Interactivity

Prefer CSS-only solutions (`:checked`, `:focus-within`, `:hover`, `~` sibling)
over event handlers — but **only when the mechanism supports it**. A recipe
driven by a `widget.State` (`SlideDeck`, and anything behind `RevealedBy`)
reads a state attribute, which CSS alone cannot move; there you bind the state
and set it from a handler. Inventing a parallel CSS-only reveal beside a
state-driven one is the leaf fix this ecosystem forbids.

## No Standard Library in WASM Packages

Use only webtyp ecosystem modules — never `errors`, `strconv`, `strings`,
`time`, `net/http`, `encoding/json`.

| Stdlib | WebTyp replacement | Docs |
|--------|---------------------|------|
| `errors`, `fmt`, `strconv`, `strings`, `path/filepath` | `webtyp.com/fmt` | [pkg.go.dev](https://pkg.go.dev/webtyp.com/fmt) |
| `encoding/json` | `webtyp.com/json` — zero reflection, requires `fmt.Fielder` (generated by `ormc`) | [pkg.go.dev](https://pkg.go.dev/webtyp.com/json) |
| `time` | `webtyp.com/time` — uses JS `Date` in WASM; `time.Now()` crashes TinyGo | [pkg.go.dev](https://pkg.go.dev/webtyp.com/time) |
| `net/http` | `webtyp.com/fetch` — uses browser `fetch` API in WASM | [pkg.go.dev](https://pkg.go.dev/webtyp.com/fetch) |

A `_test.go` file carrying `//go:build !wasm` is backend-only and may use
stdlib (`regexp`, `strings`, `testing`).

## Tests

```go
//go:build !wasm

package mycomponent

import (
    "strings"
    "testing"
)

func TestMyComponent_Render(t *testing.T) {
    html := (&MyComponent{Title: "Hello"}).Render().String()

    if !strings.Contains(html, "class='mycomponent'") {
        t.Error("expected mycomponent class")
    }
}
```

`Element.String()` serializes the tree; attributes render with **single
quotes** (`class='x'`, `role='tab'`).

**Never assert on an element id you composed** — they are minted, so a test
that spells one out is asserting the very thing a component must not do. Parse
the markup and follow the references, or assert on `data-id`.

## When to Create a PLAN.md

Any code change (new component, modifying an existing one) requires
`docs/PLAN.md` at the module root for user review before dispatching to an
external agent. Documentation-only changes (`.md` files) can be edited directly.

See `agents-workflow` skill for the full PLAN.md and codejob workflow.
