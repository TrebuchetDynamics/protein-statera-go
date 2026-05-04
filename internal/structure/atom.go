package structure

// Atom stores one Cartesian atom coordinate from a PDB ATOM record.
type Atom struct {
	ID           int
	Name         string
	Element      string
	X, Y, Z      float64
	BFactor      float64
	Occupancy    float64
	ResidueName  string
	ResidueIndex int
	ChainID      string
	AltLoc       string
	ICode        string
}
