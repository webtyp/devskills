---
name: dom-elements
description: Building and wiring a webtyp/dom Element — value embedding, element IDs that dom mints (never compose one), Key vs id, typed event handlers, and state binding. Use when writing any Render() that builds a dom tree, in a component, a layout, or an app module.
---

# Building a `dom.Element`

Applies to any package that builds a DOM tree: `webtyp/components`,
`webtyp/layout`, and an application's own views.

## Embed `dom.Element` as a value, never a pointer

```go
// ✅ CORRECT
type MyThing struct {
    Element
    Title string
}

// ❌ WRONG — double heap allocation, nil-guard boilerplate, GC pressure
type MyThing struct {
    *Element
}
```

TinyGo (the final WASM compiler) has a simple GC. Fewer heap objects = fewer
pauses. Value embedding is 1 allocation instead of 2, better cache locality,
and zero nil-panic risk.

## Element IDs are MINTED by `dom` — never composed

The rule most often got wrong. It fails silently, or crashes.

`dom` assigns ids itself. `Element.GetID()` returns the element's id, **minting
a unique one on demand** if it has none. `Element.For(other)` does the same for
a `for=` attribute, and its doc names the wider use: *label/input pairing and
`aria-*` references*.

```go
// ✅ CORRECT — dom mints the ids; the caller only references them
panel := Section().Attr("role", "tabpanel")
tab := Button().Attr("role", "tab")

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
  many instances as a screen needs.

Use `Attr("data-id", domainID)` when a test or a consumer needs the domain
identifier back. It is data, not something the framework resolves.

**`Element.Key(k)` is not an id.** It is the stable identity for keyed
reconciliation in `BindChildren`. Setting it on an element that is not inside a
`BindChildren` is inert noise — leave it off.

## Wire interactivity on the element as it is built

Use the typed handlers `Element` already exposes — `OnClick`, `OnChange`,
`OnInput`, `OnKeyDown`, `OnFocusIn`, `OnFocusOut`, `OnMouseEnter`,
`OnMouseLeave`, `OnSubmit`, `OnToggle`, `OnBlur` — never a second pass that
looks nodes up by id after render:

```go
tab := Button().
    BindStateFunc(widget.Current, func() bool { return c.isActive(id) }).
    OnClick(func(Event) { c.activate(id) })
```

No `front.go`, and no build tag on the file that carries them: TinyGo
eliminates handlers as dead code on SSR builds.

Prefer a CSS-only mechanism where one exists (see skill **core-principles**),
but only where the mechanism supports it: a recipe driven by a `widget.State`
reads a state attribute, which CSS alone cannot move. Inventing a parallel
CSS-only path beside a state-driven one is a leaf fix.

## State reaches the markup as a `widget.State`, never a class

```go
el.BindState(widget.Open, someSignalBool)              // bound to a signal
el.BindStateFunc(widget.Current, func() bool { … })    // computed per render
```

The stylesheet reacts with `When(...)` / `Cue(...)` (skill **widget-styling**).
A state written as a class drifts from the stylesheet silently.

Binding helpers: `Bind`, `BindText`, `BindTextFunc`, `BindAttr`,
`BindAttrFunc`, `BindAttrBool`, `BindAttrBoolFunc`, `BindClass`,
`BindClassFunc`, `BindChildren`, `BindState`, `BindStateFunc`.

## Serializing

`Element.String()` renders the tree. Attributes come out **single-quoted**
(`class='x'`, `role='tab'`) — match that in test assertions.

**Never assert on an element id you composed**: ids are minted, so a test that
spells one out is asserting the very thing this skill forbids. Parse the markup
and follow the references, or assert on `data-id`.

## Related skills

- **components** — creating a component in `webtyp/components`
- **widget-styling** — the stylesheet that reacts to the states bound here
- **wasm** — stdlib replacements (`webtyp/fmt`, `webtyp/time`, `webtyp/json`)
