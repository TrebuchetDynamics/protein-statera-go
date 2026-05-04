# Protein Statera Go

Evidence-first protein structure validation and interpretation engine in Go.

Protein Statera Go ingests PDB structures, extracts AlphaFold pLDDT confidence
from B-factor fields, runs basic physical plausibility checks, compares aligned
structures with RMSD, and emits interpretable reports.

## Boundaries

- No protein prediction models
- No ML training
- No Python runtime
- Pure Go core, intended for `CGO_ENABLED=0`
- PDB-first MVP; mmCIF is future scope
- Optional UI is decoupled from the core engine

## Commands

```bash
protein analyze data/examples/sample_af.pdb
protein confidence data/examples/sample_af.pdb
protein compare data/examples/sample_af.pdb data/examples/sample_exp.pdb
protein report data/examples/sample_af.pdb --html
protein-ui
```

## Development

```bash
make test
make build
```

The MVP comparison assumes atom-order aligned structures. It does not perform
structural superposition or sequence alignment yet.
