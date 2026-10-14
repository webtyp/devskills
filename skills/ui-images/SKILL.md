---
name: ui-images
description: Shipping raster images — image.go declaring RenderImages() []image.Asset (read by AST), variants S/M/L, and the webtyp/image Responsive/Img builders with sizes and lazy. Use when a component, layout or module needs a photo or raster.
---

# UI images

The filename (`image.go`) is the contract. Placement in a component: **components**.

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
