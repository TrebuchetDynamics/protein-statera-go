package geometry

import "github.com/TrebuchetDynamics/protein-statera-go/internal/structure"

// ClashPair records two non-bonded atoms closer than the configured threshold.
type ClashPair struct {
	AtomA    structure.Atom
	AtomB    structure.Atom
	Distance float64
}

// ClashResult summarizes steric clash counts.
type ClashResult struct {
	Total  int
	Severe int
	Pairs  []ClashPair
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
