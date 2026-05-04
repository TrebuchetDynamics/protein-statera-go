# Protein Statera Go Design

## Goal

Build a pure-Go starter repo for an evidence-first protein structure validation
engine, matching the supplied PRD and publishing it as
`TrebuchetDynamics/protein-statera-go`.

## Scope

MVP includes PDB parsing, AlphaFold pLDDT extraction from B-factor fields,
confidence band analysis, simple steric clash detection, RMSD comparison for
already aligned structures, text/HTML evidence reports, a CLI, examples, and a
decoupled `protein-ui` placeholder.

Out of scope: protein prediction, ML training, Python runtime, mmCIF,
structural alignment, Ramachandran plots, and graphical gogpu/ui integration.

## Architecture

Core packages live under `internal/` and expose small typed functions. CLI
commands compose those packages without owning scientific logic. Rendering is
separate from evidence construction so future UI and HTML paths can reuse the
same report model.

## Verification

Use `CGO_ENABLED=0 go test ./...`, `CGO_ENABLED=0 go vet ./...`,
`CGO_ENABLED=0 go build ./cmd/protein`, and
`CGO_ENABLED=0 go build ./cmd/protein-ui`.
