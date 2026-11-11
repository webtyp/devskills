---
name: langc
description: The webtyp translation dictionary and its CLI — positional lang.json format, where dictionaries live (app config/lang.json is configuration only; texts live in each library lang.json), langc sync and its 8 AST extraction rules. Use when running langc, editing lang.json, or a key is missing or untranslated.
---

# langc — the dictionary

Writing the text itself: **translations**.

## Where dictionaries live

- App: `config/lang.json` holds **configuration only** — `default` and `languages`, no keys. A text is translated once, in the `lang.json` of the library that uses it; an app that needs its own text puts it in a module. `langc sync` moves a translation an app still holds into that library (never overwriting a different one: it reports the conflict).
- Library: `lang.json` at its module root (no `default`). Run on its own (a demo, no `config/`), the bundle's default is the first of its `languages`.
- There is **no** `config/lang.go` and no Go API to register words or force a language: the dictionary's `default` plus the browser's language decide. Never call `lang.OutLang("es")` in an app.
- `sitec` merges them at build (`langc.BundleTranslations`); the library dictionaries filtered to the app's `languages`. Injected as `<script type="application/json" id="webtyp-lang">`, read once by `webtyp/lang` in the browser.

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
2. In a library adds new keys as `["", …]` and keeps unused ones (dynamic keys); a key a dependency already translates without contradiction is dropped (written once, lower in the graph: **translations**). In an app writes no keys: it moves the app's translations to their library and reports `N untranslated in <languages> — <its lang.json>`.
3. Scope: the app's packages and every package they import from the client **or** server build; a module reached only from a file in neither build (`config/deploy.go`, tag `deploy`) is a tool and is not scanned. Keys are read only from the browser build: `GOOS=js GOARCH=wasm -tags browser`, as `app` compiles it — a `!browser` file is server code.
4. Sorts keys; reports each diagnostic once (again only when it changes), under the `LANG` label in the dev TUI. One log call per severity: its label prints once.
5. Checks reuse: flags keys sharing a leading or trailing run of 2+ words, after dropping connectors (`and`, `for`, `of`…) on the edge where they diverge — `and column` or `required for` is never flagged. The warning names the call to write (`fmt.Err("<run>", rest)`) or the key that already exists. Reuse warnings are defects; the usual fix is not a smaller key but **data**: a statement, type, field or code name in the varying part is a variable, never a key (`fmt.Err(errNameMissing, stmtDropColumn)`).
6. Checks identical values: flags non-prose keys whose values end unchanged in every language (`ERROR ... ends with ... unchanged`). Fix by passing codes, prefixes and type names as variables (data), never literals.
7. Checks placeholders: a key whose same value sits in 3+ languages, one of them zh/hi/ar/ru (`fake` filler) is an `Error` — one line, count + value; `grep` the value in the dictionary. Latin-script cognates (es/pt/fr `canal`) are fine; a value equal to its key is exempt.

Then fill the empty values and run `gotest`.

**Don't translate a key into itself to satisfy the sync.** `fmt.Err("int", "invalid")` word-by-word has no good Spanish reading (`"int"` isn't prose — it has no natural translation, so the lazy fix is `"int": ["int"]`, a no-op entry that looks done but says nothing). Reword the *source* call into one phrase instead — `fmt.Err("not a valid integer")` → `"no es un número entero válido"` — so every language gets a real sentence. A dictionary entry whose value equals its key is a sign the extraction is scooping up a non-prose token (a type name, a literal code); fix the call site, not the dictionary.

## The 8 extraction rules

1. `lang.Translate("…")`, `lang.Text("…")`
2. `fmt.Err("…")` and `model.{Invalid,NotFound,Conflict}` constants or conversions
3. `model.Field` `Label` (or `Name` with `_` → space when no Label) — never inside a `model.Definition{Wire: true}` (a transport format)
4. `fmt.KeyValue{Value: "…"}` passed to `webtyp.com/input`
5. `input.SetPlaceholder("…")`, `input.SetTitle("…")`
6. `model.Field` `Help`
7. Composite literals assigning a `lang.Text` field
8. Literals returned by a function whose sole return type is `lang.Text`

`webtyp.com/lang` must never import `langc` (it compiles to WASM).
