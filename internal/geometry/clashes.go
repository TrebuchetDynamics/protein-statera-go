package geometry

import "github.com/TrebuchetDynamics/protein-statera-go/internal/structure"

var vdwRadii = map[string]float64{
	"C": 1.70,
	"N": 1.55,
	"O": 1.52,
	"S": 1.80,
	"H": 1.00,
	"P": 1.80,
}

// ElementVDWRadius returns the van der Waals radius for an element, or a default of 1.70 Å.
func ElementVDWRadius(element string) float64 {
	if r, ok := vdwRadii[element]; ok {
		return r
	}
	return 1.70
}

// ClashPair records two non-bonded atoms closer than the configured threshold.
type ClashPair struct {
	AtomA    structure.Atom
	AtomB    structure.Atom
	Distance float64
}

// ClashResult summarizes steric clash counts.
type ClashResult struct {
	Total     int
	Severe    int
	Pairs     []ClashPair
	Clashscore float64
}

// FindStericClashes counts simple distance-threshold steric clashes.
func FindStericClashes(s structure.Structure, thresholdAngstroms float64) ClashResult {
	type placedAtom struct {
		atom         structure.Atom
		residueIndex int
		chainID      string
	}

	atoms := make([]placedAtom, 0, s.AtomCount())
	for _, residue := range s.Residues {
		for _, atom := range residue.Atoms {
			atoms = append(atoms, placedAtom{atom: atom, residueIndex: residue.Index, chainID: residue.ChainID})
		}
	}

	result := ClashResult{}
	for i := 0; i < len(atoms); i++ {
		for j := i + 1; j < len(atoms); j++ {
			if approximateBonded(atoms[i], atoms[j]) {
				continue
			}
			distance := DistanceAngstroms(atoms[i].atom, atoms[j].atom)
			if distance >= thresholdAngstroms {
				continue
			}
			result.Total++
			if distance <= thresholdAngstroms*0.75 {
				result.Severe++
			}
			result.Pairs = append(result.Pairs, ClashPair{AtomA: atoms[i].atom, AtomB: atoms[j].atom, Distance: distance})
		}
	}

	return result
}

// FindStericClashesVDW uses element-specific van der Waals radii per MolProbity.
// Clash: distance < vdW_A + vdW_B - 0.4 Å (≥ 0.4 Å overlap)
// Severe: distance < vdW_A + vdW_B - 0.9 Å (≥ 0.9 Å overlap)
func FindStericClashesVDW(s structure.Structure) ClashResult {
	type placedAtom struct {
		atom         structure.Atom
		residueIndex int
		chainID      string
		vdW          float64
	}

	atoms := make([]placedAtom, 0, s.AtomCount())
	for _, residue := range s.Residues {
		for _, atom := range residue.Atoms {
			atoms = append(atoms, placedAtom{
				atom:         atom,
				residueIndex: residue.Index,
				chainID:      residue.ChainID,
				vdW:          ElementVDWRadius(atom.Element),
			})
		}
	}

	result := ClashResult{}
	for i := 0; i < len(atoms); i++ {
		for j := i + 1; j < len(atoms); j++ {
			ai := atoms[i]
			aj := atoms[j]
			if approximateBonded(struct {
				atom         structure.Atom
				residueIndex int
				chainID      string
			}{ai.atom, ai.residueIndex, ai.chainID},
				struct {
					atom         structure.Atom
					residueIndex int
					chainID      string
				}{aj.atom, aj.residueIndex, aj.chainID}) {
				continue
			}
			distance := DistanceAngstroms(ai.atom, aj.atom)
			vdWSum := ai.vdW + aj.vdW
			if distance >= vdWSum-0.4 {
				continue
			}
			result.Total++
			if distance < vdWSum-0.9 {
				result.Severe++
			}
			result.Pairs = append(result.Pairs, ClashPair{AtomA: ai.atom, AtomB: aj.atom, Distance: distance})
		}
	}

	if len(atoms) > 0 {
		result.Clashscore = float64(result.Total) * 1000.0 / float64(len(atoms))
	}

	return result
}

func approximateBonded(a, b struct {
	atom         structure.Atom
	residueIndex int
	chainID      string
}) bool {
	if a.chainID != b.chainID {
		return false
	}
	delta := a.atom.ID - b.atom.ID
	if delta < 0 {
		delta = -delta
	}
	return delta == 1
}
