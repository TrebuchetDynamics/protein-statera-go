# Protein Statera Go Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a public, pure-Go starter repository for Protein Statera Go.

**Architecture:** Keep parser, confidence, geometry, comparison, evidence, and rendering as separate internal packages. CLI commands compose those packages and the optional UI remains decoupled from the core engine.

**Tech Stack:** Go standard library, `CGO_ENABLED=0`, GitHub CLI for repository creation.

---

### Task 1: Core Model and PDB Parser

**Files:**
- Create: `internal/structure/atom.go`
- Create: `internal/structure/residue.go`
- Create: `internal/structure/parser_pdb.go`
- Test: `internal/structure/parser_test.go`

- [x] Write failing parser tests for ATOM parsing, residue grouping, pLDDT/B-factor capture, and malformed coordinate errors.
- [x] Run `go test ./...` and confirm undefined parser behavior.
- [x] Implement minimal fixed-width PDB ATOM parser.
- [x] Run parser tests until passing.

### Task 2: Validation Packages

**Files:**
- Create: `internal/confidence/*.go`
- Create: `internal/geometry/*.go`
- Create: `internal/comparison/*.go`
- Test: matching package tests

- [x] Write failing tests for pLDDT bands, low segments, distance, clashes, RMSD, and mismatched atom counts.
- [x] Implement minimal behavior to pass tests.
- [x] Run `go test ./...`.

### Task 3: Evidence and Rendering

**Files:**
- Create: `internal/evidence/*.go`
- Create: `internal/render/*.go`
- Test: evidence and render tests

- [x] Write failing tests for report metrics, notes, text rendering, confidence rendering, and HTML escaping.
- [x] Implement report construction and deterministic renderers.
- [x] Run targeted tests.

### Task 4: CLI and Repo Files

**Files:**
- Create: `cmd/protein/main.go`
- Create: `cmd/protein-ui/main.go`
- Create: `data/examples/*.pdb`
- Create: `README.md`, `Makefile`, `LICENSE`, docs

- [x] Write failing CLI tests for analyze, compare, and usage errors.
- [x] Implement CLI command router.
- [x] Add examples and docs.
- [ ] Run full validation gates.

### Task 5: Publish

**Files:**
- Git metadata only

- [ ] Commit verified starter.
- [ ] Create public GitHub repository `TrebuchetDynamics/protein-statera-go`.
- [ ] Push `main` to `origin`.
