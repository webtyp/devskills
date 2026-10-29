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
- **`en` never gets a column.** The Go string literal *is* the English text (the lookup key), so it never needs a slot — `SyncToolTranslations` confirms this by padding `lang.Supported()` minus `"en"` (`lang/langc/langc.go:200-204`).
- **A library's `languages` array is what it supports TODAY, not a mandate to pre-seed `lang.Supported()`'s full list** (`en, es, zh, hi, ar, pt, fr, de, ru` — `lang/language.go:61`). A brand-new file starts at `["es"]` (`langc/file.go:20`); widen it only when a real consumer app asks for another language — guessing `zh`/`hi`/`ar`/`ru` with no native reviewer is the same defect as the next rule, at the scale of a whole language. Status/rationale: `webtyp/docs/LANG_SUPPORT_MASTER_PLAN.md`.

## `langc sync [dir]` (`-tool` for backend CLIs)

1. Walks the AST (skips `*_test.go`, vendor, hidden dirs, `tests/`, `conformance/`).
2. Adds new keys as `["", …]`; in apps deletes unused keys (libraries keep them: dynamic keys).
3. Sorts keys; reports `lang: N untranslated (es) — see config/lang.json`.
4. Checks reuse: flags keys sharing 2+ words (`share ... — extract it as its own key and reuse it`). Fix by decomposing `fmt.Err` call sites to share vocabulary; reuse warnings are defects.
5. Checks identical values: flags non-prose keys whose values end unchanged in every language (`ERROR ... ends with ... unchanged`). Fix by passing codes, prefixes and type names as variables (data), never literals.

Then fill the empty values and run `gotest`.

**Don't translate a key into itself to satisfy the sync.** `fmt.Err("int", "invalid")` word-by-word has no good Spanish reading (`"int"` isn't prose — it has no natural translation, so the lazy fix is `"int": ["int"]`, a no-op entry that looks done but says nothing). Reword the *source* call into one phrase instead — `fmt.Err("not a valid integer")` → `"no es un número entero válido"` — so every language gets a real sentence. A dictionary entry whose value equals its key is a sign the extraction is scooping up a non-prose token (a type name, a literal code); fix the call site, not the dictionary.

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
