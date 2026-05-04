package comparison

import (
	"errors"
	"math"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/geometry"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

var errMismatchedAtoms = errors.New("structures must contain the same nonzero atom count")

// KabschAlign computes the optimal rotation matrix aligning mobile onto target
// using the Kabsch algorithm. Returns RMSD after alignment.
func KabschAlign(mobile, target []structure.Atom) (rmsd float64, err error) {
	if len(mobile) == 0 || len(mobile) != len(target) {
		return 0, errMismatchedAtoms
	}

	n := len(mobile)
	cxM, cyM, czM := centroid(mobile)
	cxT, cyT, czT := centroid(target)

	// Center both sets
	m := make([][3]float64, n)
	t := make([][3]float64, n)
	for i := range mobile {
		m[i] = [3]float64{mobile[i].X - cxM, mobile[i].Y - cyM, mobile[i].Z - czM}
		t[i] = [3]float64{target[i].X - cxT, target[i].Y - cyT, target[i].Z - czT}
	}

	// Compute covariance matrix
	var cov [3][3]float64
	for i := 0; i < n; i++ {
		for a := 0; a < 3; a++ {
			for b := 0; b < 3; b++ {
				cov[a][b] += m[i][a] * t[i][b]
			}
		}
	}

	// SVD of covariance matrix via solving 3x3 eigensystem of cov^T * cov
	// For the 3x3 case, we can use a simpler approach with quaternion-based
	// optimal rotation (Horn's method), which avoids full SVD.
	rot := optimalRotationQuaternion(cov)

	// Apply rotation and compute RMSD
	sum := 0.0
	for i := 0; i < n; i++ {
		rx := rot[0][0]*m[i][0] + rot[0][1]*m[i][1] + rot[0][2]*m[i][2]
		ry := rot[1][0]*m[i][0] + rot[1][1]*m[i][1] + rot[1][2]*m[i][2]
		rz := rot[2][0]*m[i][0] + rot[2][1]*m[i][1] + rot[2][2]*m[i][2]
		dx := rx - t[i][0]
		dy := ry - t[i][1]
		dz := rz - t[i][2]
		sum += dx*dx + dy*dy + dz*dz
	}

	return math.Sqrt(sum / float64(n)), nil
}

func centroid(atoms []structure.Atom) (cx, cy, cz float64) {
	for _, a := range atoms {
		cx += a.X
		cy += a.Y
		cz += a.Z
	}
	n := float64(len(atoms))
	return cx / n, cy / n, cz / n
}

func optimalRotationQuaternion(cov [3][3]float64) [3][3]float64 {
	// Build the 4x4 matrix for quaternion-based optimal rotation (Horn 1987)
	qmat := [4][4]float64{
		{cov[0][0] + cov[1][1] + cov[2][2], cov[1][2] - cov[2][1], cov[2][0] - cov[0][2], cov[0][1] - cov[1][0]},
		{cov[1][2] - cov[2][1], cov[0][0] - cov[1][1] - cov[2][2], cov[0][1] + cov[1][0], cov[0][2] + cov[2][0]},
		{cov[2][0] - cov[0][2], cov[0][1] + cov[1][0], -cov[0][0] + cov[1][1] - cov[2][2], cov[1][2] + cov[2][1]},
		{cov[0][1] - cov[1][0], cov[0][2] + cov[2][0], cov[1][2] + cov[2][1], -cov[0][0] - cov[1][1] + cov[2][2]},
	}

	// Power iteration to find the dominant eigenvector of this 4x4 symmetric matrix
	eigen := [4]float64{1, 0, 0, 0}
	for iter := 0; iter < 20; iter++ {
		next := [4]float64{}
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				next[i] += qmat[i][j] * eigen[j]
			}
		}
		norm := math.Sqrt(next[0]*next[0] + next[1]*next[1] + next[2]*next[2] + next[3]*next[3])
		if norm < 1e-12 {
			break
		}
		for k := 0; k < 4; k++ {
			eigen[k] = next[k] / norm
		}
	}

	q0, q1, q2, q3 := eigen[0], eigen[1], eigen[2], eigen[3]

	// Convert quaternion to rotation matrix
	return [3][3]float64{
		{q0*q0 + q1*q1 - q2*q2 - q3*q3, 2*q1*q2 - 2*q0*q3, 2*q1*q3 + 2*q0*q2},
		{2*q1*q2 + 2*q0*q3, q0*q0 - q1*q1 + q2*q2 - q3*q3, 2*q2*q3 - 2*q0*q1},
		{2*q1*q3 - 2*q0*q2, 2*q2*q3 + 2*q0*q1, q0*q0 - q1*q1 - q2*q2 + q3*q3},
	}
}

// TMScore computes the Template Modeling score for two aligned structures.
// Uses Cα atoms only. d0 normalizes by target length.
// Returns a score in (0,1]; <0.17 = random; >0.5 = same fold.
func TMScore(target, model structure.Structure) (float64, error) {
	atomsA := target.Atoms()
	atomsB := model.Atoms()
	if len(atomsA) == 0 || len(atomsA) != len(atomsB) {
		return 0, errMismatchedAtoms
	}
	return tmScoreFromAtoms(atomsA, atomsB), nil
}

func tmScoreFromAtoms(a, b []structure.Atom) float64 {
	L := float64(len(a))
	d0 := 1.24*math.Cbrt(L-15) - 1.8
	if d0 < 0.5 {
		d0 = 0.5
	}

	sum := 0.0
	for i := range a {
		d := geometry.DistanceAngstroms(a[i], b[i])
		sum += 1.0 / (1.0 + (d*d)/(d0*d0))
	}

	return sum / L
}

// RMSDAligned computes RMSD with Kabsch optimal superposition of Cα atoms.
func RMSDAligned(a, b structure.Structure) (float64, error) {
	atomsA := filterCA(a)
	atomsB := filterCA(b)
	if len(atomsA) == 0 || len(atomsA) != len(atomsB) {
		return 0, errMismatchedAtoms
	}
	return KabschAlign(atomsA, atomsB)
}

func filterCA(s structure.Structure) []structure.Atom {
	var result []structure.Atom
	for _, r := range s.Residues {
		if ca := r.AtomByName("CA"); ca != nil {
			result = append(result, *ca)
		}
	}
	return result
}
