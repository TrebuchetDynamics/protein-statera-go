# Validation Methods

## Geometry

The MVP includes Euclidean atom distance and simple steric clash detection.
Clashes are counted when non-approximated-bonded atom pairs are closer than the
configured threshold. The default threshold is 1.10 angstroms.

## RMSD

RMSD is computed as:

```text
sqrt(sum((xi - yi)^2) / N)
```

where `N` is the aligned atom count. The MVP assumes atom-order alignment and
rejects mismatched atom counts.

## Report Evidence

Reports explicitly label evidence class:

- `predicted_structure`
- `experimental_structure`
- `comparative_analysis`
