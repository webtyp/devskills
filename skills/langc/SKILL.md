---
name: langc
description: The webtyp translation dictionary and its CLI — positional lang.json format, where dictionaries live (app config/lang.json vs library lang.json, project wins), langc sync and its 8 AST extraction rules. Use when running langc, editing lang.json, or a key is missing or untranslated.
---

# langc — the dictionary

Writing the text itself: **translations**.

## Where dictionaries live

- App: `config/lang.json` (has `default` and `languages`).
- Library: `lang.json` at its module root (no `default`). Run on its own (a demo, no `config/`), the bundle's default is the first of its `languages`.
- There is **no** `config/lang.go` and no Go API to register words or force a language: the dictionary's `default` plus the browser's language decide. Never call `lang.OutLang("es")` in an app.
- `sitec` merges them at build (`langc.BundleTranslations`); the app wins. Injected as `<script type="application/json" id="webtyp-lang">`, read once by `webtyp/lang` in the browser.

## Format — positional

```json
{
  "default": "es",
  "languages": ["es", "fr"],
  "keys": {
    "Delete": ["Eliminar", "Supprimer"],
    "Send": ["", ""]
  }
}
```

- `value[i]` ↔ `languages[i]`; `""` = untranslated.
- Never reorder or remove a language without shifting every key.

## `langc sync [dir]` (`-tool` for backend CLIs)

1. Walks the AST (skips `*_test.go`, vendor, hidden dirs, `tests/`, `conformance/`).
2. Adds new keys as `["", …]`; in apps deletes unused keys (libraries keep them: dynamic keys).
3. Sorts keys; reports `lang: N untranslated (es) — see config/lang.json`.

Then fill the empty values and run `gotest`.

## The 8 extraction rules

1. `lang.Translate("…")`, `lang.Text("…")`
2. `fmt.Err("…")`
3. `model.Field` `Label` (or `Name` with `_` → space when no Label)
4. `fmt.KeyValue{Value: "…"}` passed to `webtyp.com/input`
5. `input.SetPlaceholder("…")`, `input.SetTitle("…")`
6. `model.Field` `Help`
7. Composite literals assigning a `lang.Text` field
8. Literals returned by a function whose sole return type is `lang.Text`

`webtyp.com/lang` must never import `langc` (it compiles to WASM).
