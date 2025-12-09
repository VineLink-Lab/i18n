# i18n

Lightweight i18n helper for Go — a small library and example for loading JSON message bundles, selecting languages, parameterized messages, and bundling messages into directories or embedded files.

Module: `github.com/VineLink-Lab/i18n`

Status: Example/demo-quality library intended to show a simple approach to organizing translations in JSON files and loading them from disk or embed.FS.

## Table of Contents

- [Install](#install)
- [Quick Start](#quick-start)
- [Examples](#examples)
- [Message file layout](#message-file-layout)
- [Behavior & fallbacks](#behavior--fallbacks)
- [CLI (generate)](#cli-generate)
- [Contributing](#contributing)

## Install

Requires Go 1.25+.

To use as a module in your project:

```bash
go get github.com/VineLink-Lab/i18n
```

Or add `github.com/VineLink-Lab/i18n` to your `go.mod` and import the packages you need.

## Quick Start

This repository includes a runnable example at `example/i18n.go`. The core usage patterns are:

- load message bundles from a directory on disk
- load message bundles from an embed.FS
- translate a message key with optional parameters and language override
- organize messages into bundles (subdirectories)

Example imports used in the code snippets below:

```go
import (
    "embed"
    "github.com/VineLink-Lab/i18n/pkg/i18n"
    "golang.org/x/text/language"
)
```

### Load from a directory

```go
// Load translations from a directory (example/message)
translator, err := i18n.NewTranslator("example/message", language.English)
if err != nil {
    panic(err)
}

// Translate using the default language
msg := translator.Translate("hello")
```

### Load from embed.FS

```go
//go:embed message
var messageFS embed.FS

translator, err := i18n.NewI18nFromFS(messageFS, language.English)
if err != nil {
    panic(err)
}
```

### Translate with options

You can pass options like target language and parameters. The library provides option helpers such as `WithLanguage` and `WithParams`.

```go
// Simple translate with explicit language
out := translator.Translate("hello", i18n.WithLanguage(language.French))

// Translate with parameters (map-like type i18n.M)
out = translator.Translate("welcome", i18n.WithParams(i18n.M{"name": "Alice"}))

// Combine options
out = translator.Translate("greeting", i18n.WithLanguage(language.French), i18n.WithParams(i18n.M{"name": "Alice"}))
```

### Bundles

Messages can be organized into bundles (subdirectories). Switch the default bundle with:

```go
translator.SetDefaultBundle("with_parameters")
```

Then translations will be looked up inside that bundle (e.g. `with_parameters/en.json`).

## Message file layout

This project uses JSON message files in `example/message/`:

- `example/message/en.json`
- `example/message/fr.json`
- `example/message/with_parameters/en.json`
- `example/message/with_parameters/fr.json`

Each file maps message keys to strings. Strings may contain placeholders that the translator replaces using the parameters map.

Example `en.json`:

```json
{
  "hello": "Hello",
  "welcome": "Welcome, {name}!"
}
```

## Behavior & fallbacks

- Default language: set when creating the I18n instance (e.g. `language.English`).
- Language fallback: if a requested language doesn't have a translation, the library falls back to the default language.
- Parent language fallback: if a specific sublanguage is requested (e.g. `en-US`) and no message exists for it, the library will attempt to fall back to a parent language (e.g. `en`).
- Bundle fallback: if a message is not found in the selected bundle, behavior depends on the library's lookup order (see example code); the default bundle can be changed with `SetDefaultBundle`.

## CLI (generate)

This repo contains a small CLI under `cmd/i18n` — see `cmd/i18n/main.go` and `cmd/i18n/generate.go` for details. It can be used to run generation tasks provided in `internal/generate`.

A quick way to run the CLI locally:

```bash
go run ./cmd/i18n --help
```

## Contributing

Contributions, issues and feature requests are welcome. If you want to extend the project (support more formats such as YAML, add richer templating, or provide stricter validation), please open an issue first to discuss the design.

When contributing:

- Keep changes small and focused
- Add tests for new behavior
- Update or add example usage in `example/`

## License

This repository does not include a LICENSE file by default — add one if you plan to open-source the project publicly.

---

For more detailed usage, check `example/i18n.go` and the `pkg/i18n` package.
