# Usage

## Analyze

```bash
protein analyze structure.pdb
```

Reports residue count, atom count, pLDDT confidence bands, steric clash counts,
and evidence-bound interpretation notes.

## Confidence

```bash
protein confidence structure.pdb
```

Prints high, medium, and low pLDDT residue counts plus contiguous low-confidence
segments.

## Compare

```bash
protein compare af.pdb exp.pdb
```

Computes RMSD in angstroms for atom-order aligned structures. The MVP does not
perform structural alignment.

## HTML Report

```bash
protein report structure.pdb --html
```

Writes a static HTML evidence report to stdout.
