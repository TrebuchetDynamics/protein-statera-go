package geometry

import (
	"math"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

type RamachandranRegion int

const (
	RegionFavored  RamachandranRegion = iota
	RegionAllowed
	RegionOutlier
)

func (r RamachandranRegion) String() string {
	switch r {
	case RegionFavored:
		return "favored"
	case RegionAllowed:
		return "allowed"
	case RegionOutlier:
		return "outlier"
	default:
		return "unknown"
	}
}

// ComputeTorsion returns the torsion angle (degrees, [-180,180]) defined by four atoms a-b-c-d.
// The angle measures rotation around the b-c bond, with zero meaning cis-planar.
func ComputeTorsion(a, b, c, d structure.Atom) float64 {
	v1 := vec3{b.X - a.X, b.Y - a.Y, b.Z - a.Z}
	v2 := vec3{c.X - b.X, c.Y - b.Y, c.Z - b.Z}
	v3 := vec3{d.X - c.X, d.Y - c.Y, d.Z - c.Z}

	n1 := v1.cross(v2)
	n2 := v2.cross(v3)

	n1Norm := n1.norm()
	n2Norm := n2.norm()
	v2Norm := v2.norm()
	if n1Norm < 1e-12 || n2Norm < 1e-12 || v2Norm < 1e-12 {
		return 0
	}

	n1 = n1.scale(1.0 / n1Norm)
	n2 = n2.scale(1.0 / n2Norm)

	cosPhi := n1.dot(n2)
	if cosPhi < -1.0 {
		cosPhi = -1.0
	}
	if cosPhi > 1.0 {
		cosPhi = 1.0
	}

	sinPhi := n1.cross(n2).dot(v2) / v2Norm

	return math.Atan2(sinPhi, cosPhi) * (180.0 / math.Pi)
}

// PhiPsi computes backbone φ and ψ torsion angles for residue at index in the given Structure.
// Residue 0 is the first in the Residues slice (N-terminal).
// Returns (phi, psi, ok=false) for terminal residues or missing backbone atoms.
func PhiPsi(s structure.Structure, residueIndex int) (phi, psi float64, ok bool) {
	if residueIndex < 0 || residueIndex >= len(s.Residues) {
		return 0, 0, false
	}
	if residueIndex == 0 || residueIndex == len(s.Residues)-1 {
		return 0, 0, false
	}

	prev := s.Residues[residueIndex-1]
	curr := s.Residues[residueIndex]
	next := s.Residues[residueIndex+1]

	if prev.ChainID != curr.ChainID || curr.ChainID != next.ChainID {
		return 0, 0, false
	}

	prevC := prev.AtomByName("C")
	currN := curr.AtomByName("N")
	currCA := curr.AtomByName("CA")
	currC := curr.AtomByName("C")
	nextN := next.AtomByName("N")

	if prevC == nil || currN == nil || currCA == nil || currC == nil || nextN == nil {
		return 0, 0, false
	}

	phi = ComputeTorsion(*prevC, *currN, *currCA, *currC)
	psi = ComputeTorsion(*currN, *currCA, *currC, *nextN)
	return phi, psi, true
}

// Classify assigns a Ramachandran region (favored, allowed, outlier) to the given φ/ψ values
// and residue type. Uses simplified MolProbity-compatible region boundaries.
func Classify(phi, psi float64, resType string) RamachandranRegion {
	switch resType {
	case "GLY":
		return classifyGly(phi, psi)
	case "PRO":
		return classifyPro(phi, psi)
	case "ILE", "VAL":
		return classifyIleVal(phi, psi)
	default:
		return classifyGeneral(phi, psi)
	}
}

func classifyGeneral(phi, psi float64) RamachandranRegion {
	if inRegion(phi, psi, -100, -20, -80, -10) {
		return RegionFavored
	}
	if inRegion(phi, psi, -160, -30, 80, 180) {
		return RegionFavored
	}
	if inRegion(phi, psi, -30, 0, 140, 180) {
		return RegionFavored
	}
	if inRegion(phi, psi, -160, -140, 130, 160) {
		return RegionFavored
	}
	if inRegion(phi, psi, -180, -100, -90, -80) {
		return RegionAllowed
	}
	if inRegion(phi, psi, -30, 0, 120, 140) {
		return RegionAllowed
	}
	if inRegion(phi, psi, 20, 80, -20, 60) {
		return RegionAllowed
	}
	return RegionOutlier
}

func classifyGly(phi, psi float64) RamachandranRegion {
	if inRegion(phi, psi, -100, 20, -80, 20) {
		return RegionFavored
	}
	if inRegion(phi, psi, -180, -40, 60, 180) {
		return RegionFavored
	}
	if inRegion(phi, psi, 20, 100, -60, 60) {
		return RegionFavored
	}
	if inRegion(phi, psi, -180, -80, -90, -60) {
		return RegionAllowed
	}
	if inRegion(phi, psi, 60, 120, 60, 120) {
		return RegionAllowed
	}
	return RegionOutlier
}

func classifyPro(phi, psi float64) RamachandranRegion {
	if inRegion(phi, psi, -80, -50, 120, 160) {
		return RegionFavored
	}
	if inRegion(phi, psi, -100, -50, 60, 120) {
		return RegionFavored
	}
	if inRegion(phi, psi, -80, -50, -60, -20) {
		return RegionFavored
	}
	if inRegion(phi, psi, -120, -50, -60, 60) {
		return RegionAllowed
	}
	return RegionOutlier
}

func classifyIleVal(phi, psi float64) RamachandranRegion {
	if inRegion(phi, psi, -160, -40, 80, 180) {
		return RegionFavored
	}
	if inRegion(phi, psi, -100, -40, -80, -10) {
		return RegionFavored
	}
	if inRegion(phi, psi, -40, -20, 120, 160) {
		return RegionAllowed
	}
	return RegionOutlier
}

func inRegion(phi, psi, phiMin, phiMax, psiMin, psiMax float64) bool {
	return phi >= phiMin && phi <= phiMax && psi >= psiMin && psi <= psiMax
}

type vec3 [3]float64

func (v vec3) cross(w vec3) vec3 {
	return vec3{
		v[1]*w[2] - v[2]*w[1],
		v[2]*w[0] - v[0]*w[2],
		v[0]*w[1] - v[1]*w[0],
	}
}

func (v vec3) dot(w vec3) float64 {
	return v[0]*w[0] + v[1]*w[1] + v[2]*w[2]
}

func (v vec3) norm() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v vec3) scale(s float64) vec3 {
	return vec3{v[0] * s, v[1] * s, v[2] * s}
}
