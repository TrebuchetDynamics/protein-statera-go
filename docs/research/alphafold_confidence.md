# AlphaFold Confidence

MVP rule: AlphaFold pLDDT is read from the PDB B-factor field and averaged per
residue when multiple atoms are present.

Current bands:

- High: pLDDT >= 90
- Medium: 50 <= pLDDT < 90
- Low: pLDDT < 50

These bands are report labels, not proof of biological function. Low pLDDT
segments are evidence for uncertainty, flexibility, disorder, or model weakness
that needs domain interpretation.
