package comparison

import (
	"math"
	"testing"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

func TestRMSDAssumesAlignedAtoms(t *testing.T) {
	a := structure.Structure{Residues: []structure.Residue{
		{Index: 1, Atoms: []structure.Atom{{Name: "CA", X: 0, Y: 0, Z: 0}}},
		{Index: 2, Atoms: []structure.Atom{{Name: "CA", X: 1, Y: 0, Z: 0}}},
	}}
	b := structure.Structure{Residues: []structure.Residue{
		{Index: 1, Atoms: []structure.Atom{{Name: "CA", X: 0, Y: 1, Z: 0}}},
		{Index: 2, Atoms: []structure.Atom{{Name: "CA", X: 1, Y: 1, Z: 0}}},
	}}

	got, err := RMSD(a, b)
	if err != nil {
		t.Fatalf("RMSD returned error: %v", err)
	}
	if math.Abs(got-1.0) > 1e-9 {
		t.Fatalf("RMSD = %.6f A, want 1.000000 A", got)
	}
}

func TestRMSDRejectsMismatchedAtomCounts(t *testing.T) {
	a := structure.Structure{Residues: []structure.Residue{{Atoms: []structure.Atom{{Name: "CA"}}}}}
	b := structure.Structure{}

	if _, err := RMSD(a, b); err == nil {
		t.Fatal("RMSD returned nil error, want mismatched atom count error")
	}
}
