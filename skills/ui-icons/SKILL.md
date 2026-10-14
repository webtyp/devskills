---
name: ui-icons
description: Shipping icons from a Go package — webtyp/icons per-glyph packages, IconSvg() *sprite.Sprite in svg.go behind !wasm, reference-only in WASM, currentColor. Use when a component, layout or module needs an icon.
---

# UI icons

The reference is WASM-safe, the payload backend-only, the filename (`svg.go`) the contract. Placement in a component: **components**.

An icon has two halves that must reach different places:

- the **reference** — the symbol id, a plain string: `svg.Icon("trash")`. This
  is all that may reach the browser.
- the **geometry** — the `<path>` data and viewBox. Backend-only: `webtyp/ssr`
  extracts it at build time, injects it once into the page, and the markup
  points at it with `<use href="#trash">`. Shipping geometry to the browser
  would drag `webtyp/json` + `webtyp/model` into the WASM bundle for nothing.

## Prefer a shared glyph from `webtyp/icons`

`webtyp.com/icons` is the shared set, **one package per glyph** (`trash`,
`pencil`, `plus`, `undo`, `selectall`). Per-glyph packaging is what keeps each
glyph's geometry behind its own `//go:build !wasm`, so importing one glyph for
its `Ref` can never leak another's path data into a WASM build.

```go
// main file — WASM-safe: only the id string crosses
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

func (c *RangeSlider) IconSvg() *sprite.Sprite {
    return sprite.NewSprite(trash.Def(), pencil.Def())
}
```

Adding a glyph to the shared set is a new folder in `webtyp/icons` — it never
touches the ones already there. Copy `trash/`, rename, swap the `Ref` id, and
paste the viewBox + path into `icons.Solid(Ref, viewBox, d)`. Keep it a solid
single-path glyph; anything else is not that family and calls `sprite.Define`
directly.

## A glyph only this package owns

Reference in the main file, geometry in `svg.go`:

```go
const IconPhone = svg.Icon("rangeslider-phone")   // prefix ids with the package
```

```go
//go:build !wasm

func (c *RangeSlider) IconSvg() *sprite.Sprite {
    return sprite.NewSprite(
        sprite.Define(IconPhone, "0 0 24 24",
            sprite.Path("M6.62 10.79a15.053 15.053 0 006.59 6.59l…"),
        ),
    )
}
```

**`IconSvg()` returns `*sprite.Sprite` — never `map[string]string`.** Builders:
`sprite.NewSprite(defs…)`, `sprite.Define(icon, viewBox, body…)`,
`sprite.Path(d)`, `sprite.Raw(s)`, `(*Sprite).AddFile(id, svgFile)`,
`sprite.MergeAll(sprites…)`.

## Hard rules

- `IconSvg()` MUST live in `svg.go` behind `//go:build !wasm`.
- Every path uses `fill="currentColor"` (or `stroke`). **Never hardcode a
  colour in SVG** — the box around the glyph owns the colour, so a white-text
  button yields a white glyph and a red box a red glyph.
- The class passed to `Ref.Render(class)` is where size and colour come from.
- The sprite is injected inline into the `<body>`; there is no
  `/assets/icons.svg` URL, so `href="#id"` always resolves with no network
  request. Never reference the sprite by URL from CSS.
