---
name: ui-assets
description: Shipping icons and images from a Go package — webtyp/icons glyph packages, IconSvg() *sprite.Sprite in svg.go, and image.go declaring RenderImages() []image.Asset for the raster pipeline. Use when a component, layout or app module needs an icon or a photo.
---

# Icons and images

Both follow one shape: **the reference is WASM-safe, the payload is
backend-only, and the filename is the contract.**

| Asset | Declared in | Returns |
|---|---|---|
| Icons | `svg.go` (`//go:build !wasm`) | `IconSvg() *sprite.Sprite` |
| Images | `image.go` | `RenderImages() []image.Asset` |

A file with any other name is invisible to the tooling that scans for it.

---

# Icons

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

---

# Images

Rasters are not icons: they go through an optimization pipeline that generates
responsive variants.

## Declaring what the pipeline must process — `image.go`

```go
package site

import "webtyp.com/image"

// RenderImages declares every raster this package serves. Read by the pipeline
// via AST: it must stay a hand-written composite literal — no loops, no
// computed paths, no variables.
func RenderImages() []image.Asset {
    return []image.Asset{
        {Path: "img/hero.jpg", Variants: image.VariantM, Alt: "Hero"},
        {Path: "img/badge.png", Variants: image.AllVariants, Alt: "Badge"},
    }
}
```

- **The filename MUST be `image.go`.** The pipeline detects it by name.
- **It is read by AST, not executed.** A loop, a computed path, or a value
  pulled from a variable extracts as nothing, silently.
- `Path` is relative to the module directory. `Alt` is SEO text; empty is
  derived from the filename.
- `Variants` is a bitmask: `VariantS` (480px), `VariantM` (1024px), `VariantL`
  (1600px), `AllVariants`. **Declare only the variants something actually
  renders** — each is a file written to disk and shipped to the CDN. One site
  generated S, M and L per photo while only one was ever linked.

## Rendering — the `webtyp/image` builders

These compile for **both** WASM and backend (no build tag), so they go in the
main file:

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

Sibling packages, none imported directly by a view: `image/min` (backend WebP
pipeline), `image/browser` (client-side compression before upload,
`//go:build wasm`), `image/favicon` (logo → icon set).

## Related skills

- **components** — where `svg.go` and `image.go` sit in a component
- **dom-elements** — building the tree these assets are rendered into
