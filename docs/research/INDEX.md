# Research Inventory

Protein Statera Go — evidence-first protein structure validation workbench.
Last updated: 2026-05-04.

---

## 1. AlphaFold Confidence Metrics

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| EMBL-EBI AlphaFold Training | Official docs | https://www.ebi.ac.uk/training/online/courses/alphafold/inputs-and-outputs/evaluating-alphafolds-predicted-structures-using-confidence-scores/ | pLDDT band definitions, PAE interpretation, confidence score integration | pLDDT, PAE, segments |
| AlphaFold 3 Docs (DeepMind) | Official docs | https://google-deepmind-alphafold3.mintlify.app/guides/confidence-metrics | Per-atom pLDDT, pTM/ipTM, ranking_score, has_clash, PAE plot interpretation | pLDDT, confidence |
| UCSF ChimeraX AlphaFold Error | Documentation | https://rbvi.ucsf.edu/chimerax/data/pae-apr2022/pae.html | pLDDT coloring (blue>90, yellow 50-70, red<50), PAE heatmap visualization | Reports, coloring |
| Jumper et al. 2021 (AlphaFold2) | Paper | 10.1038/s41586-021-03819-2 | Original pLDDT and PAE methods | Confidence evidence |
| Mariani et al. 2013 (lDDT) | Paper | 10.1093/bioinformatics/btt473 | lDDT-Cα metric that pLDDT estimates | pLDDT definition |
| Guo et al. 2022 | Paper | 10.1002/pro.4390 | Confidence score integration: pLDDT+PAE interpretation | PAE evidence |

**Standard pLDDT bands (EMBL-EBI / AlphaFold DB):**
- Very high: pLDDT > 90
- Confident: 90 > pLDDT > 70
- Low: 70 > pLDDT > 50
- Very low: pLDDT < 50

**PAE interpretation (Å):**
- < 5 Å: High confidence in relative position
- 5-15 Å: Moderate confidence
- > 15 Å: Low confidence / flexible orientation

---

## 2. Steric Clash Validation

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| MolProbity (Chen et al. 2010) | Paper | 10.1107/S0907444909042073 | PROBE all-atom contact, clashscore, clash thresholds | Steric clashes |
| Word et al. 1999 (PROBE) | Paper | 10.1006/jmbi.1998.2400 | van der Waals radii, contact dot algorithm, 0.5Å probe | Clash detection |
| Williams et al. 2018 (MolProbity update) | Paper | 10.1002/pro.3330 | Updated reference data, clashscore tracking, CaBLAM | Clash thresholds |
| MolProbity Validation Options | Official docs | http://molprobity.biochem.duke.edu/help/validation_options/ | Clash severity levels: mild ≥0.4Å, severe >0.9Å | Clash classification |
| Davis et al. 2007 | Paper | PMID: 17452350 | All-atom contact analysis, clashscore = clashes/1000 atoms | Clashscore metric |

**Thresholds:**
- Steric clash (non-bonded): ≥ 0.4 Å overlap of vdW radii
- Severe clash: > 0.9 Å overlap
- Clashscore: clashes per 1000 atoms
- Standard vdW radii: C=1.7, N=1.55, O=1.52, S=1.8, H=1.0-1.2 Å

---

## 3. Ramachandran Validation

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| Ramachandran et al. 1963 | Paper | 10.1016/S0022-2836(63)80023-6 | Original hard-sphere model, φ/ψ plot | Phi/psi basics |
| Lovell et al. 2003 | Paper | 10.1002/prot.10286 | Updated distributions, Cβ deviation, MolProbity criteria | Ramachandran regions |
| Hintze et al. 2016 | Paper | PMID: 27543030 | Top8000 rotamer library, 0.3% outlier threshold | Rotamer validation |
| Hollingsworth & Karplus 2010 | Paper | 10.1016/j.str.2010.08.012 | Fresh look at Ramachandran distributions, τ-angle dependence | Region definitions |
| Zhou, O'Hern & Regan 2011 | Paper | 10.1002/pro.683 | Bridge region, τ-angle influence on allowed space | Bridge region |
| EMBL-EBI Structure Validation | Course | https://www.ebi.ac.uk/pdbe/modval4 | φ/ψ/ω definitions, validation criteria | Phi/psi calculation |
| EMBL-EBI Foundations | Course | https://www.ebi.ac.uk/training/online/courses/foundations-protein-structure/ | Gly/Pro/pre-Pro distributions | Residue categories |

**Torsion angle definitions:**
- φ (phi): C(i-1)-N(i)-CA(i)-C(i)
- ψ (psi): N(i)-CA(i)-C(i)-N(i+1)
- ω (omega): CA(i)-C(i)-N(i+1)-CA(i+1) — normally ~180° (trans) or ~0° (cis)

**Six residue categories** (MolProbity):
1. General (18 amino acids)
2. Glycine
3. Proline (trans)
4. Proline (cis)
5. Pre-proline (residue before Pro)
6. Isoleucine/Valine (special Cβ branching)

**Quality thresholds:**
- Favored: 98% expected
- Allowed: ~2% expected
- Outlier: < 0.05% of quality-filtered reference data

---

## 4. Structural Comparison (RMSD, TM-score, GDT-TS)

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| Zhang & Skolnick 2004 | Paper | 10.1002/prot.20264 | TM-score definition, length-independent metric | TM-score |
| Zhang & Skolnick 2005 | Paper | 10.1093/nar/gki524 | TM-align algorithm, heuristic DP alignment | TM-score, alignment |
| Zhang et al. 2022 (US-align) | Paper | 10.1038/s41592-022-01585-1 | Universal TM-score for proteins/RNA/DNA | TM-score formula |
| Kabsch 1976 | Paper | 10.1107/S0567739476001873 | Kabsch algorithm for optimal rotation matrix | RMSD superposition |
| OpenStructure docs | Official | https://openstructure.org/docs/2.11/mol/alg/gdt/ | GDT implementation, sliding window, distance thresholds 1,2,4,8 | GDT-TS |
| Zemla 2003 (LGA/GDT) | Paper | 10.1093/nar/gkg500 | GDT-TS definition, CASP evaluation | GDT-TS |
| EMBL-EBI Training | Course | https://www.ebi.ac.uk/training/online/courses/alphafold/validation-and-impact/how-accurate-are-alphafold-structure-predictions/ | RMSD vs experimental: median 0.6Å for high-confidence, ≥2Å for low | RMSD benchmarks |

**Key formulas:**
- RMSD: sqrt(Σ(d_i²)/N) where d_i is distance between corresponding Cα atoms after alignment
- TM-score: (1/L) Σ 1/(1+(d_i/d_0)²) where d_0 = 1.24·∛(L-15) - 1.8
- TM-score > 0.5: same fold; < 0.17: random
- GDT-TS: average of GDT at thresholds 1,2,4,8 Å

---

## 5. Secondary Structure (DSSP)

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| Kabsch & Sander 1983 | Paper | 10.1002/bip.360221211 | Original DSSP algorithm, H-bond energy, 8-state assignment | DSSP |
| DSSP Documentation (CMBI) | Official docs | https://swift.cmbi.umcn.nl/gv/dssp/HTML/descrip.html | H-bond energy formula, output format | DSSP |
| DSSP Wikipedia | Reference | https://en.wikipedia.org/wiki/DSSP_(protein) | Algorithm summary, H-bond formula | DSSP reference |
| PDB-REDO DSSP (Hekkelman) | Source code | https://github.com/PDB-REDO/dssp | Modern C++ rewrite, BSD-2 license | Implementation reference |
| Oxford Protein Informatics | Blog | https://www.blopig.com/blog/2014/08/dssp/ | Algorithm explanation, 8-state types | DSSP explanation |
| ChimeraX DSSP docs | Official | https://rbvi.ucsf.edu/chimerax/docs/user/commands/dssp.html | Virtual H placement, energyCutoff=-0.5 | H-bond params |

**H-bond energy formula (Coulomb approximation):**
```
E = 0.084 · (1/r_ON + 1/r_CH - 1/r_OH - 1/r_CN) · 332 kcal/mol
```
Partial charges: C=+0.42, O=-0.42, N=-0.20, H=+0.20

**8 DSSP states:**
- H: α-helix (4-turn, i→i+4)
- G: 3₁₀-helix (3-turn, i→i+3)
- I: π-helix (5-turn, i→i+5)
- E: β-strand (extended, in ladder)
- B: β-bridge (single bridge pair)
- T: hydrogen-bonded turn
- S: bend (high curvature, ≥70°)
- C: coil/loop (none of above)

---

## 6. Solvent Accessible Surface Area (SASA)

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| Shrake & Rupley 1973 | Paper | 10.1016/0022-2836(73)90015-7 | Original SASA algorithm, water probe sphere | SASA |
| Lee & Richards 1971 | Paper | 10.1016/0022-2836(71)90324-X | Solvent accessibility definition | Solvent exposure |

**Algorithm:** Roll a 1.4 Å probe sphere around each atom's vdW surface. Count accessible points on a tessellated sphere (typically 100-960 points per atom). SASA = fraction of accessible points × atom surface area.

---

## 7. Bond Geometry (Engh-Huber)

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| Engh & Huber 1991 | Paper | 10.1107/S0108767391001640 | Standard bond lengths and angles for amino acids | Bond geometry |
| wwPDB Chemical Component Dictionary | Database | http://mmcif.wwpdb.org/dictionaries/ | CHEM_COMP_BOND, CHEM_COMP_ANGLE tables | Bond/angle ref |
| wwPDB bond_distance_limits | Dictionary | http://mmcif.pdb.org/dictionaries/mmcif_pdbx_v50.dic/Categories/pdbx_bond_distance_limits.html | Element-pair bond distance limits | Bond detection |

**Engh-Huber ideal values (examples):**
- C-N (peptide): 1.329 Å, σ=0.014
- N-CA: 1.458 Å, σ=0.019
- CA-C: 1.525 Å, σ=0.021
- C=O: 1.231 Å, σ=0.020
- CA-CB: 1.530 Å, σ=0.020

**Validation:** Bonds beyond 4σ from ideal are outliers, beyond 6σ are severe outliers.

---

## 8. mmCIF / PDBx Format

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| wwPDB mmCIF Dictionary v5 | Spec | https://mmcif.rcsb.org/dictionaries/mmcif_pdbx_v50.dic/Index | Complete PDBx/mmCIF specification | mmCIF parser |
| Bourne et al. 1997 (Methods Enzymol) | Paper | https://mmcif.wwpdb.org/docs/pubs/methods-enzymology-paper-1997.html | mmCIF architecture, data categories | mmCIF structure |
| BurntSushi/cif | Go lib | https://github.com/BurntSushi/cif | Pure Go CIF 1.1 parser (Unlicense) | Go CIF parsing |
| PDBx/mmCIF Dictionary Resources | Docs | http://mmcif.rcsb.org/ | Software resources, C++/Python examples | mmCIF reference |
| US-align mmCIF support | Source | https://github.com/pyleb-lab/US-align | Production mmCIF parsing in C++ | Parsing patterns |

**Key mmCIF categories for this project:**
- `_atom_site.*` — atom coordinates (replaces ATOM/HETATM)
- `_entity.*`, `_entity_poly.*`, `_entity_poly_seq.*` — polymer metadata
- `_struct_conf.*` — secondary structure annotations
- `_chem_comp_bond.*`, `_chem_comp_angle.*` — ideal geometry

---

## 9. Go Bioinformatics Projects

| Source | Type | License | Contribution | Feature |
|--------|------|---------|-------------|---------|
| biogo/biogo | Go library | BSD-3 | Comprehensive bioinformatics (sequences, alignment, graphs) | Algorithm reference |
| rmera/gochem | Go library | ? | Chemistry library: PDB/GRO/XYZ, QM interfaces | PDB patterns |
| fogleman/ribbon | Go library | MIT | PDB parser (ATOM, HELIX, SHEET, CONECT), ribbon rendering | PDB parser |
| tikz/bio | Go library | MIT | PDB/DSSP/UniProt/SASA parsers, SMCRA architecture | DSSP, SASA parsing |
| BurntSushi/cif | Go library | Unlicense | CIF 1.1 parser | mmCIF base |
| shenwei356/bio | Go library | MIT | Lightweight bioinformatics, FASTA/FASTQ, kmers | Sequence handling |

**Note:** tikz/bio has DSSP output parsing (wrapper around DSSP binary) and SASA output parsing. No Go-native DSSP or SASA implementation found — gap this project fills.

---

## 10. PDB Format Reference

| Source | Type | DOI/URL | Contribution | Feature |
|--------|------|---------|-------------|---------|
| wwPDB Format Guide v3.3 | Spec | http://www.wwpdb.org/documentation/file-format | Official PDB coordinate entry format | PDB parser |
| PDB ATOM record spec | Spec | http://www.wwpdb.org/documentation/file-format-content/format33/sect9.html | Column-level ATOM field specification | PDB parser |
| PDB HELIX/SHEET records | Spec | http://www.wwpdb.org/documentation/file-format-content/format33/sect5.html | Secondary structure records | PDB parser |

**ATOM record columns (1-indexed):**
| Columns | Field | Type |
|---------|-------|------|
| 1-6 | Record name | "ATOM  " or "HETATM" |
| 7-11 | Serial | Integer |
| 13-16 | Atom name | String |
| 17 | AltLoc | Character |
| 18-20 | Residue name | String |
| 22 | Chain ID | Character |
| 23-26 | Residue seq | Integer |
| 27 | iCode | Character |
| 31-38 | X (Å) | Real(8.3) |
| 39-46 | Y (Å) | Real(8.3) |
| 47-54 | Z (Å) | Real(8.3) |
| 55-60 | Occupancy | Real(6.2) |
| 61-66 | TempFactor | Real(6.2) |
| 77-78 | Element | String(2) |

---

## Summary Statistics

| Category | Count |
|----------|-------|
| Academic papers cited | 20+ |
| Official documentation sources | 12 |
| Open-source projects referenced | 6 Go, 5 Other |
| DOI references | 15+ |
| URL references | 25+ |

---

## Sources Deferred (Future Scope)

- Protein-ligand interaction validation
- Cryo-EM specific validation metrics
- NMR ensemble analysis
- Deep learning-based quality assessment (QMEAN, ProQ3D)
- Full H-bond network analysis
- Crystallographic data fit (R-free, density maps)
