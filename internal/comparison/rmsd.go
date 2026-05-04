package comparison

import (
	"errors"
	"math"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

// RMSD returns root-mean-square deviation in angstroms for already aligned atoms.
func RMSD(a, b structure.Structure) (float64, error) {
	atomsA := a.Atoms()
	atomsB := b.Atoms()
	if len(atomsA) == 0 || len(atomsA) != len(atomsB) {
		return 0, errors.New("structures must contain the same nonzero atom count")
	}

	sum := 0.0
	for i := range atomsA {
		dx := atomsA[i].X - atomsB[i].X
		dy := atomsA[i].Y - atomsB[i].Y
		dz := atomsA[i].Z - atomsB[i].Z
		sum += dx*dx + dy*dy + dz*dz
	}

	return math.Sqrt(sum / float64(len(atomsA))), nil
}
