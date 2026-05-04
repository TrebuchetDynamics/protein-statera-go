# Implementation Plan

Protein Statera Go — research-backed protein structure validation workbench.
Generated: 2026-05-04.

## Phase 1: Parser Robustness (HIGH)

### 1.1 PDB Parser — Multi-Model Support
**Justification:** Real PDB files contain MODEL/ENDMDL records for NMR ensembles.
**Test first:** `internal/structure/parser_test.go`
- Add MODEL record parsing test with 2+ models
- Test MODEL skipping (only first model parsed by default, all models with flag)
- Add ENDMDL termination test
**Implementation:**
- Add `ParsePDBMulti` function returning `[]Structure`
- Track MODEL record state in scanner
- Skip non-ATOM records (TER, END) during residue grouping

### 1.2 PDB Parser — AltLoc and Insertion Codes
**Justification:** PDB columns 17 (altLoc) and 27 (iCode) are part of the standard.
**Test first:** Test parsing of altLoc ('A', 'B') and iCode in ATOM records.
**Implementation:**
- Add `AltLoc` field to Atom struct (default to 'A' or first encountered)
- Add `ICode` field to Residue struct
- Update `residueKey` to include iCode and altLoc

### 1.3 PDB Parser — Occupancy and Additional Fields
**Justification:** Occupancy (< 1.0) indicates partial occupancy or alternate conformations.
**Implementation:**
- Add `Occupancy` field to Atom struct
- Parse columns 55-60 in `parseAtomLine`

---

## Phase 2: Confidence Evidence (HIGH)

### 2.1 Updated pLDDT Bands
**Justification:** EMBL-EBI documentation standardizes 4 bands (not 3).
**Source:** https://www.ebi.ac.uk/training/online/courses/alphafold/
**Implementation:**
- Update `HighConfidencePLDDT` to 90.0, add `ConfidentPLDDT` = 70.0, keep `LowConfidencePLDDT` = 50.0
- Update `Analyze` to return 4-band counts: VeryHigh(>90), Confident(70-90), Low(50-70), VeryLow(<50)
- Update segment tracking for VeryLow residues only

### 2.2 PAE Evidence Handling
**Justification:** PAE complements pLDDT for domain-level confidence.
**Source:** EMBL-EBI PAE guide, AF3 docs
**Implementation:**
- New package: `internal/confidence/pae.go`
- Types: `PAEMatrix` (2D float64), `PAEDomain` (residue range with low intra-PAE)
- Parse PAE from JSON files (AlphaFold DB output format)
- Add `ParsePAEJSON(io.Reader) (*PAEMatrix, error)`
- Domain extraction: find contiguous blocks where intra-PAE < 5Å
- Add PAE evidence to report

### 2.3 Residue Confidence Segment Improvements
**Justification:** Current segment tracking only for very low. Add per-band segments.
**Implementation:**
- Add `Segments` map[string][]Segment to Analysis for all bands
- Track contiguous high, medium, low segments
- Report segment boundaries with residue indices

---

## Phase 3: Geometry Validation (HIGH)

### 3.1 Bond Length Validation (Engh-Huber)
**Justification:** Bond geometry is a primary validation criterion.
**Source:** Engh & Huber (1991), wwPDB Chemical Component Dictionary
**Test first:** Test with known Ala dipeptide coordinates.
**Implementation:**
- New package: `internal/geometry/bondlength.go`
- Embed Engh-Huber standard bond lengths in a map
- Key: "RESNAME:ATOM1-ATOM2" → {mean, sigma}
- Compute observed bond lengths from atomic coordinates
- Flag bonds outside 4σ (outlier) and 6σ (severe)
- Report per-residue and per-structure bond statistics

### 3.2 Bond Angle Validation
**Justification:** Complement to bond length validation.
**Implementation:**
- New package: `internal/geometry/bondangle.go`
- Compute angles from 3-atom sets based on connectivity
- Compare against Engh-Huber reference values
- Flag angle outliers

### 3.3 Steric Clash Refinements
**Justification:** Current clash detection uses ad-hoc threshold; need MolProbity-standard thresholds.
**Source:** MolProbity papers, PROBE algorithm
**Implementation:**
- Update `FindStericClashes` with element-specific vdW radii:
  - C: 1.70, N: 1.55, O: 1.52, S: 1.80, H: 1.00, P: 1.80
- Clash threshold: sum(vdW_A, vdW_B) - 0.4 (overlap ≥ 0.4 Å)
- Severe threshold: sum(vdW_A, vdW_B) - 0.9 (overlap > 0.9 Å)
- Add clashscore metric: clashes per 1000 atoms
- Replace `approximateBonded` with covalent bond lookup table (atom1-atom2 pairs within residues)

### 3.4 Ramachandran Validation
**Justification:** φ/ψ torsion angles are the most sensitive structural quality indicator.
**Source:** Ramachandran et al. (1963), Lovell et al. (2003), MolProbity
**Test first:** Test with known φ/ψ values for α-helix (-57,-47) and β-sheet (-120,120).
**Implementation:**
- New package: `internal/geometry/ramachandran.go`
- `ComputeTorsion(a, b, c, d Atom) float64` — general torsion angle
- `PhiPsi(residues []Residue, i int) (phi, psi float64, ok bool)`
  - φ: C(i-1)-N(i)-CA(i)-C(i)
  - ψ: N(i)-CA(i)-C(i)-N(i+1)
- `Classify(phi, psi float64, resType string) RamachandranRegion`
  - 6 residue categories: General, Gly, Pro, pre-Pro, Ile/Val
  - 3 quality levels: Favored (98%), Allowed (~2%), Outlier (<0.05%)
- Embed reference region boundaries as polygon arrays
- Report per-residue Ramachandran classification

---

## Phase 4: Structural Comparison (HIGH)

### 4.1 Kabsch Alignment
**Justification:** Structural superposition is prerequisite for RMSD comparison.
**Source:** Kabsch (1976)
**Implementation:**
- New package: `internal/comparison/superpose.go`
- `KabschRotation(mobile, target []Vec3) (rotation [3][3]float64, rmsd float64)`
- Center both coordinate sets at origin
- Compute covariance matrix
- SVD decomposition (implement inline or use lightweight math)
- Apply rotation to mobile coordinates
- Return aligned RMSD

### 4.2 TM-score
**Justification:** Length-independent structural similarity metric (standard in CASP).
**Source:** Zhang & Skolnick (2004, 2005)
**Implementation:**
- `TMScore(alignedA, alignedB []Atom) float64`
- Use Cα atoms only
- d_0 = 1.24 * cbrt(L - 15) - 1.8 where L = target length
- TM = (1/L_target) * Σ 1/(1 + (d_i/d_0)²)
- Value 0-1, > 0.5 = same fold

### 4.3 GDT-TS
**Justification:** CASP gold standard for model quality assessment.
**Implementation:**
- `GDTHA(target, model []Vec3) (gdt_ha float64)`
- GDT at thresholds 1,2,4,8 Å using iterative superposition
- Sliding window of size 7 for initial alignment seeds

---

## Phase 5: Secondary Structure (MEDIUM)

### 5.1 DSSP-compatible SS Assignment
**Justification:** Standard secondary structure classification is essential for structure interpretation.
**Source:** Kabsch & Sander (1983), DSSP C++ source (BSD-2)
**Test first:** Test with mini α-helix (ideal φ=-57, ψ=-47) and β-strand coordinates.
**Implementation:**
- New package: `internal/structure/dssp.go`
- `ComputeHBondEnergy(donor, acceptor struct{N,H,C,O Atom}) float64`
- `AssignSecondaryStructure(structure Structure) []SSState`
- 8-state output: H,G,I,E,B,T,S,C
- Virtual H placement (1.01 Å along bisector of C-N-CA)
- Minimum helix length: 3 residues (α), 2 residues (3₁₀)
- Key limitation: no H-atom reconstruction needed for initial implementation (use geometric H placement)

---

## Phase 6: Solvent Accessibility (MEDIUM)

### 6.1 Shrake-Rupley SASA
**Justification:** Solvent exposure is key for assessing surface vs buried residues.
**Source:** Shrake & Rupley (1973)
**Implementation:**
- New package: `internal/geometry/sasa.go`
- `ShrakeRupley(atoms []Atom, radii map[string]float64, nPoints int) []float64`
- Tessellated sphere (Fibonacci lattice, nPoints=256)
- Probe radius: 1.4 Å (water)
- vdW radii per element: C=1.70, N=1.55, O=1.52, S=1.80, P=1.80
- Per-atom SASA and per-residue SASA
- Relative SASA: normalize by residue-specific maximum SASA (Gly=84, Ala=113, etc. from Miller et al. 1987)

---

## Phase 7: mmCIF Support (MEDIUM)

### 7.1 mmCIF Parser
**Justification:** mmCIF is the wwPDB archive format since 2014. PDB format is legacy.
**Source:** wwPDB PDBx/mmCIF dictionary v5, BurntSushi/cif (Go CIF parser, Unlicense)
**Implementation:**
- Use `BurntSushi/cif` as CIF parsing layer (Unlicense-compatible)
- New package: `internal/structure/parser_mmcif.go`
- Map `_atom_site` loop to Atom struct:
  - `_atom_site.id` → ID
  - `_atom_site.Cartn_x/y/z` → X,Y,Z
  - `_atom_site.label_atom_id` → Name
  - `_atom_site.label_comp_id` → ResidueName
  - `_atom_site.label_seq_id` → ResidueIndex
  - `_atom_site.auth_asym_id` → ChainID
  - `_atom_site.B_iso_or_equiv` → BFactor
  - `_atom_site.type_symbol` → Element
- Group atoms into Residues → Structure (same as PDB path)
- Parse `_struct_conf` for SS annotations if present

---

## Phase 8: Reports and Batch Processing (HIGH)

### 8.1 JSON Report Format
**Implementation:**
- `internal/render/json.go`
- `JSONReport(report evidence.Report) string`
- Machine-readable JSON with all metrics, confidence, geometry details

### 8.2 Enhanced HTML Report
**Implementation:**
- Update `internal/render/html.go`
- Add Ramachandran summary table
- Add confidence band visualization (CSS bar chart)
- Add clash list table with atom pairs and distances
- Add PAE domain summary if available

### 8.3 Batch Analysis
**Implementation:**
- `cmd/protein batch <dir>` — analyze all PDB/mmCIF files in directory
- CSV output mode for batch results
- Summary statistics across batch
- Parallel processing with worker goroutines

---

## Phase 9: Tests and Fixtures (HIGH)

### 9.1 Test Fixtures
- Add `data/fixtures/` directory with:
  - Mini structures: 5-residue α-helix, 5-residue β-strand
  - Clash test: atoms with 0.3 Å, 0.5 Å, 1.0 Å, 2.0 Å separation
  - Ramachandran test: residues at known φ/ψ (-57,-47), (-120,120), (60,40)
  - SASA test: isolated glycine in vacuum
  - mmCIF test: minimal valid mmCIF file
  - Multi-model PDB: 2-model NMR ensemble
  - PAE mock JSON

### 9.2 Test Coverage Targets
- Every exported function: ≥ 1 test
- Every validation method: ≥ 3 tests (valid, invalid, edge case)
- Parser tests: malformed input, empty input, mixed ATOM/HETATM
- CLI tests: all subcommands, error cases, --html flag

---

## Phase 10: CLI Enhancement (HIGH)

### 10.1 New Subcommands
- `protein rama <file>` — Ramachandran analysis
- `protein geometry <file>` — bond length/angle validation
- `protein dssp <file>` — secondary structure assignment
- `protein sasa <file>` — solvent accessibility report
- `protein batch <dir>` — batch analysis
- `protein validate <file>` — full validation report
- `protein compare <a> <b> --tm` — TM-score instead of RMSD

### 10.2 CLI Flags
- `--json` — JSON output format
- `--csv` — CSV output (batch mode)
- `--chain <id>` — single chain analysis
- `--residues <range>` — residue range filter
- `--model <n>` — specify MODEL number for multi-model files

---

## Implementation Order (TDD)

For each feature, follow strictly:
1. Write test → `CGO_ENABLED=0 go test ./...` → confirm FAIL
2. Implement minimal Go code → `CGO_ENABLED=0 go test ./...` → confirm PASS
3. `CGO_ENABLED=0 go vet ./...` → clean
4. `CGO_ENABLED=0 go build ./cmd/protein` → clean
5. `CGO_ENABLED=0 go build ./cmd/protein-ui` → clean
6. Git commit with descriptive message
7. Push to origin main

**Commit sequence:**
1. `fix: correct approximateBonded to not skip intra-residue non-bonded pairs`
2. `feat: add MODEL/ENDMDL support to PDB parser`
3. `feat: add AltLoc, iCode, and Occupancy fields`
4. `feat: add 4-band pLDDT analysis and PAE evidence types`
5. `feat: add Engh-Huber bond length validation`
6. `feat: add bond angle validation`
7. `feat: add MolProbity-style steric clash with element-specific vdW radii`
8. `feat: add Ramachandran phi/psi calculation and validation`
9. `feat: add Kabsch alignment and TM-score comparison`
10. `feat: add GDT-TS metric`
11. `feat: add DSSP-compatible secondary structure assignment`
12. `feat: add Shrake-Rupley SASA calculation`
13. `feat: add mmCIF parser via BurntSushi/cif`
14. `feat: add JSON report format`
15. `feat: add batch analysis with CSV output`
16. `feat: add enhanced HTML report with visualizations`
17. `test: add comprehensive fixtures and edge-case coverage`
18. `docs: update README with new commands and features`

---

## Deferred Features

| Feature | Reason |
|---------|--------|
| Rotamer validation | Requires full H-atom placement; significant complexity for marginal MVP value |
| CaBLAM validation | Cryo-EM focused; secondary to Ramachandran |
| Full hydrogen placement (REDUCE-like) | Complex optimization; deferred to post-MVP |
| Protein-protein interface analysis | Requires chain-level awareness beyond current scope |
| pTM/ipTM scoring | Only applicable to AF3 multimers; AF2 PAE is sufficient |
| RNA/DNA validation | Out of scope for protein workbench |
| Ligand contact analysis | Requires chemical component dictionary integration |
| Real-space correlation | Requires electron density maps (crystallography) |
