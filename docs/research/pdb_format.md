# PDB Format MVP

The parser reads PDB `ATOM` records only.

Fields used:

- atom serial
- atom name
- residue name
- chain ID
- residue index
- x, y, z coordinates in angstroms
- B-factor as pLDDT source
- element symbol when present

`HETATM`, alternate locations, insertion codes, occupancy interpretation, and
mmCIF are future scope.
