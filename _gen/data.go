package main

import "strings"

const repoOrg = "thesmos-ai"

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

func (s Sub) GoSource() string {
	base := s.parent.RepoURL()
	return strings.Join([]string{
		s.ImportPath(),
		"    " + base + "/tree/main/" + s.Name,
		"    " + base + "/tree/main/" + s.Name + "{/dir}",
		"    " + base + "/blob/main/" + s.Name + "{/dir}/{file}#L{line}",
	}, "\n")
}

// modules is the source of truth for the site. Edit this slice and re-run
// `go run ./_gen` from the repo root.
var modules = []Module{
	{
		Name:        "protoc-gen-codec",
		Description: "High-performance protobuf codec for Go. Emits marshal/unmarshal on your hand-written types instead of generating new ones. Zero-alloc, deterministic, 100% mutation-tested.",
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
