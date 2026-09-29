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
				Name:        "bridge/protogo",
				Description: "Protobuf→Go bridge annotator for eidos — stamps Go-namespaced translation metadata so the Go backend stays proto-agnostic.",
				Public:      true,
			},
			{
				Name:        "cli",
				Description: "Command kernels for eidos — run, plan, explain, check, prune, version.",
				Public:      true,
			},
			{
				Name:        "cmd/eidos-reference",
				Description: "Reference-runner binary for eidos — exercises the in-tree plugin ensemble end-to-end.",
				Public:      true,
			},
			{
				Name:        "eidostest",
				Description: "Test harnesses for eidos plugin, frontend and backend authors — golden-output diffing and conformance suites.",
				Public:      true,
			},
			{
				Name:        "lang/golang",
				Description: "Go language module for eidos — shared Go conventions plus the frontend, the renderer, and the SDK a Go-generating plugin embeds.",
				Public:      true,
			},
			{
				Name:        "lang/protobuf",
				Description: "Protobuf language module for eidos — a proto3 frontend turning descriptors into the typed IR. Proto is read, never written.",
				Public:      true,
			},
			{
				Name:        "lang/typescript",
				Description: "TypeScript language module for eidos — shared conventions plus a tree-sitter frontend over .ts and .tsx.",
				Public:      true,
			},
			{
				Name:        "plugins",
				Description: "Production plugin set for eidos — signature-shape annotation plus the builder, enum and sentinel generators.",
				Public:      true,
			},
			{
				Name:        "reference",
				Description: "Reference plugins for eidos, kept runnable as worked examples.",
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
		Name:        "kanon",
		Description: "kanon generates binary codecs for Go struct types. Equal values encode to identical bytes. Decoded strings can share memory with the input buffer.",
		Public:      true,
	},
	{
		Name:        "protoc-gen-codec",
		Description: "High-performance protobuf codec for Go. Emits marshal/unmarshal on your hand-written types instead of generating new ones. Zero-alloc, deterministic, 100% mutation-tested.",
		Public:      true,
	},
	{
		Name:        "service",
		Description: "Service generators for Golang gRPC services.",
		Public:      false,
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
		Public:      true,
	},
	{
		Name:        "ledger",
		Description: "Append-only audit trail with cryptographic verification.",
		Public:      true,
		Subs: []Sub{
			{
				Name:        "cmd",
				Description: "ledger-verify verifies a ledger offline and monitors its published objects. ledger-witness runs a C2SP tlog-witness that cosigns the ledger's checkpoints.",
				Public:      true,
			},
			{
				Name:        "pkg/evidence",
				Description: "The evidence module contains the ledger's persisted formats, witnessing, timestamping and verification.",
				Public:      true,
			},
			{
				Name:        "pkg/pkcs11",
				Description: "The pkcs11 module implements core's sign.Signer and crypto.Keeper with keys in a PKCS#11 device. Private signing keys and key-encryption keys never leave the device.",
				Public:      true,
			},
			{
				Name:        "pkg/qualified",
				Description: "The qualified module contains the eIDAS trust services: AdES seals, trusted lists, and OCSP and CRL clients.",
				Public:      true,
			},
			{
				Name:        "pkg/storage/s3",
				Description: "The s3 module implements core's blob.Store for S3-compatible object storage, with ranged reads and Object Lock retention.",
				Public:      true,
			},
		},
	},
	{
		Name:        "runtime",
		Description: "Runtime services and orchestration for the thesmos kernel.",
	},
	{
		Name:        "sharder",
		Description: "sharder assigns the shards of a service to its processes and moves ownership between them through fenced handoffs. Each namespace selects its partitioning and placement strategies and its replica state model.",
		Public:      true,
	},
	{
		Name:        "thesmos-tools",
		Description: "MCP tools for agentic usage — Go, Rust, TypeScript, and filesystem helpers.",
		Langs:       []string{"go", "rust", "ts"},
	},
}
