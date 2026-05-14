# go.thesmos.sh

Vanity import host for the thesmos Go ecosystem. Serves the `go-import` meta
tags that let `go get go.thesmos.sh/<module>` resolve to the canonical sources
at [github.com/thesm-os](https://github.com/thesm-os).

Live at <https://go.thesmos.sh/>.

## Structure

```
.
├── index.html          generated — module index
├── 404.html            generated — styled fallback
├── <module>/index.html generated — per-module landing + meta tags
├── assets/
│   ├── site.css        hand-written stylesheet (Tailwind tokens, buildless)
│   └── site.js         click-to-copy for `go get` snippets
└── _gen/               site generator (excluded from Pages by Jekyll's
                        leading-underscore convention)
    ├── main.go
    ├── data.go         source of truth — modules list
    ├── go.mod
    └── templates/
```

The `_gen/` directory is the source of truth for everything generated. The
rendered `*.html` files are committed so GitHub Pages can serve them directly
with zero build pipeline.

## Adding or editing a module

1. Edit `_gen/data.go` — add a `Module{}` to the `modules` slice.
2. From the repo root: `cd _gen && go run .`
3. Commit both `_gen/data.go` and the regenerated HTML.

The generator runs `npx prettier` on its output by default. Pass `-no-fmt` to
skip if npx isn't available.

```go
{
    Name:        "ergon",
    Description: "Task runner for Go projects: ...",
    Public:      true,                       // shown under Public modules
    Langs:       []string{"go", "rust"},     // optional; defaults to ["go"]
},
```

## Submodules (nested go.mod)

When a module has its own nested `go.mod` (e.g. `eidos/cli/go.mod` declaring
`module go.thesmos.sh/eidos/cli`), the vanity host emits an **explicit
`go-import` meta tag** at the submodule path, with the parent repo as
repo-root:

```html
<meta name="go-import"
      content="go.thesmos.sh/eidos/cli git https://github.com/thesm-os/eidos" />
```

Go's tooling then clones the parent repo and locates the nested `go.mod` by
matching the declared module path. Multi-segment submodule names (e.g.
`frontend/golang`) are written verbatim in the `Sub.Name` field.

```go
Subs: []Sub{
    {Name: "cli", Description: "...", Public: true},
    {Name: "frontend/golang", Description: "...", Public: true},
},
```

> **Note on nested-module `go.mod`s.** A `replace ... => ../` paired with a
> `v0.0.0-00010101000000-...` pseudo-version works only for local dev — the
> `replace` is ignored by downstream consumers. For published submodules,
> require a real tagged version of the parent, or move the `replace` into a
> top-level `go.work` (gitignored) so the published `go.mod` is
> consumer-clean.

## Private modules

Set `Public: false`. The landing page surfaces a one-liner `GOPRIVATE` hint
since `go get` against a private repo needs:

```sh
go env -w GOPRIVATE=go.thesmos.sh/*
```

…plus git credentials that can reach the GitHub org.

## Local preview

```sh
python3 -m http.server 8000     # then open http://localhost:8000/
```

To test meta-tag resolution:

```sh
curl -sL 'http://localhost:8000/eidos/cli/?go-get=1' | grep go-import
```

## Why buildless on the served side

GitHub Pages serves the repo as-is. No CI, no Actions, no deploy pipeline.
The generator is a developer-ergonomics layer; the production artifact is
plain HTML + CSS + ~50 lines of JS.
