package geometry

import (
	"math"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

// DistanceAngstroms returns Euclidean atom distance in angstroms.
func DistanceAngstroms(a, b structure.Atom) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
