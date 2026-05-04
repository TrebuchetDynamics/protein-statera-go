package geometry

import (
	"testing"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

func TestDistanceAngstroms(t *testing.T) {
	a := structure.Atom{X: 0, Y: 0, Z: 0}
	b := structure.Atom{X: 3, Y: 4, Z: 12}

	if got := DistanceAngstroms(a, b); got != 13 {
		t.Fatalf("DistanceAngstroms = %.2f A, want 13.00 A", got)
	}
}

func TestFindStericClashesIgnoresBondedNeighbors(t *testing.T) {
	s := structure.Structure{
		ID: "clash",
		Residues: []structure.Residue{
			{
				Name:  "GLY",
				Index: 1,
				Atoms: []structure.Atom{
					{ID: 1, Name: "N", X: 0, Y: 0, Z: 0},
					{ID: 2, Name: "CA", X: 1.45, Y: 0, Z: 0},
				},
			},
			{
				Name:  "SER",
				Index: 2,
				Atoms: []structure.Atom{
					{ID: 3, Name: "OG", X: 0.80, Y: 0, Z: 0},
				},
			},
		},
	}

	result := FindStericClashes(s, 1.10)

	if result.Total != 1 || result.Severe != 1 {
		t.Fatalf("clashes = total:%d severe:%d, want 1/1", result.Total, result.Severe)
	}
	if result.Pairs[0].AtomA.ID != 1 || result.Pairs[0].AtomB.ID != 3 {
		t.Fatalf("clash pair = %d/%d, want 1/3", result.Pairs[0].AtomA.ID, result.Pairs[0].AtomB.ID)
	}
}
