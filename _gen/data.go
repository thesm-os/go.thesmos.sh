package main

import "strings"

const repoOrg = "thesm-os"

// Module describes a top-level vanity module exposed at go.thesmos.sh/<Name>.
// The repo is assumed to be github.com/<repoOrg>/<Name>.
type Module struct {
	Name        string
	Description string
	Public      bool
	Langs       []string // optional; defaults to []string{"go"}
	Subs        []Sub
}

// Sub is a nested module with its own go.mod, exposed at
// go.thesmos.sh/<Parent>/<Name>. Always inherits the parent's repo.
type Sub struct {
	Name        string
	Description string
	Public      bool
	Langs       []string

	parent *Module // wired in main()
}

func (m Module) ImportPath() string { return "go.thesmos.sh/" + m.Name }
func (m Module) GoGet() string      { return "go get " + m.ImportPath() }
func (m Module) RepoURL() string    { return "https://github.com/" + repoOrg + "/" + m.Name }

func (m Module) Languages() []string {
	if len(m.Langs) == 0 {
		return []string{"go"}
	}
	return m.Langs
}

func (m Module) GoSource() string {
	return strings.Join([]string{
		m.ImportPath(),
		"    " + m.RepoURL(),
		"    " + m.RepoURL() + "/tree/main{/dir}",
		"    " + m.RepoURL() + "/blob/main{/dir}/{file}#L{line}",
	}, "\n")
}

func (s Sub) Path() string       { return s.parent.Name + "/" + s.Name }
func (s Sub) ImportPath() string { return "go.thesmos.sh/" + s.Path() }
func (s Sub) GoGet() string      { return "go get " + s.ImportPath() }
func (s Sub) ParentName() string { return s.parent.Name }
func (s Sub) ParentImport() string {
	return s.parent.ImportPath()
}

// RepoURL points at the subfolder on the parent repo — Go finds the nested
// go.mod there; the submodule's go-import meta tag still uses the parent repo
// as repo-root.
func (s Sub) RepoURL() string {
	return s.parent.RepoURL() + "/tree/main/" + s.Name
}

func (s Sub) Languages() []string {
	if len(s.Langs) == 0 {
		return []string{"go"}
	}
	return s.Langs
}

// modules is the source of truth for the site. Edit this slice and re-run
// `go run ./_gen` from the repo root.
var modules = []Module{
	{
		Name:        "eidos",
		Description: "A composable, plugin-driven code-generation framework. Typed metadata, queryable IR, slot injection, byte-deterministic output.",
		Public:      true,
		Subs: []Sub{
			{
				Name:        "backend/golang",
				Description: "Go emitter for eidos — renders the IR back to deterministic Go source.",
				Public:      true,
			},
			{
				Name:        "bridge/protogo",
				Description: "Protobuf↔Go bridge for eidos — interop between the two frontends.",
				Public:      true,
			},
			{
				Name:        "cli",
				Description: "Command-line driver for eidos pipelines.",
				Public:      true,
			},
			{
				Name:        "cmd/eidos-reference",
				Description: "Reference-runner binary for eidos — exercises the reference plugin pipeline end-to-end.",
				Public:      true,
			},
			{
				Name:        "eidostest",
				Description: "Plugin test harness for eidos — golden-output and byte-determinism checks.",
				Public:      true,
			},
			{
				Name:        "frontend/golang",
				Description: "Go source frontend for eidos — parses Go into the typed IR.",
				Public:      true,
			},
			{
				Name:        "frontend/protobuf",
				Description: "Protobuf frontend for eidos — parses .proto descriptors into the typed IR.",
				Public:      true,
			},
			{
				Name:        "plugins",
				Description: "Bundled plugin set for eidos pipelines.",
				Public:      true,
			},
			{
				Name:        "reference",
				Description: "Reference plugin and worked-example pipeline for eidos.",
				Public:      true,
			},
		},
	},
	{
		Name:        "ergon",
		Description: "Task runner for Go projects: format, lint, test, benchmark, release, with a multi-stage check umbrella.",
		Public:      true,
	},
	{
		Name:        "protoc-gen-codec",
		Description: "High-performance protobuf codec for Go. Emits marshal/unmarshal on your hand-written types instead of generating new ones. Zero-alloc, deterministic, 100% mutation-tested.",
		Public:      true,
	},
	{
		Name:        "techne",
		Description: "Atomic, build-gated developer tools for AI coding agents: type-checked refactors, semantic search, and verify→fix loops over MCP, CLI, and TUI.",
		Public:      true,
	},
	{
		Name:        "testkit",
		Description: "High-integrity testing for Go. You write domain logic. Testkit generates the plumbing, the tests, and the proof that it all works.",
		Public:      true,
		Subs: []Sub{
			{
				Name:        "model",
				Description: "Submodule of testkit — domain model primitives for high-integrity tests.",
				Public:      true,
			},
		},
	},
	{
		Name:        "thesmos",
		Description: "The operating system for autonomous AI agents. Deterministic execution kernel with cryptographic audit trails, built for 100k+ concurrent agents.",
	},
	{
		Name:        "core",
		Description: "Shared primitives and contracts for the thesmos runtime.",
	},
	{
		Name:        "kernel",
		Description: "Deterministic execution kernel — agent scheduling and replay.",
	},
	{
		Name:        "ledger",
		Description: "Append-only audit trail with cryptographic verification.",
	},
	{
		Name:        "runtime",
		Description: "Runtime services and orchestration for the thesmos kernel.",
	},
	{
		Name:        "thesmos-tools",
		Description: "MCP tools for agentic usage — Go, Rust, TypeScript, and filesystem helpers.",
		Langs:       []string{"go", "rust", "ts"},
	},
}
