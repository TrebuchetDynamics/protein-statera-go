package structure

// Residue groups atoms by residue name, index, and chain.
type Residue struct {
	Name    string
	Index   int
	ChainID string
	ICode   string
	Atoms   []Atom
	PLDDT   float64
}

// AtomByName returns the first atom matching the given name, or nil.
func (r Residue) AtomByName(name string) *Atom {
	for i := range r.Atoms {
		if r.Atoms[i].Name == name {
			return &r.Atoms[i]
		}
	}
	return nil
}

// Structure is the parsed protein model used by validation modules.
type Structure struct {
	ID       string
	Residues []Residue
	Models   [][]Residue
	Source   string
}

// AtomCount returns the total number of atoms across all residues.
func (s Structure) AtomCount() int {
	count := 0
	for _, residue := range s.Residues {
		count += len(residue.Atoms)
	}
	return count
}

// Atoms returns all atoms in residue order.
func (s Structure) Atoms() []Atom {
	atoms := make([]Atom, 0, s.AtomCount())
	for _, residue := range s.Residues {
		atoms = append(atoms, residue.Atoms...)
	}
	return atoms
}
