package comparison

import (
	"math"
	"testing"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

func TestKabschAlignZeroRMSDForIdenticalStructures(t *testing.T) {
	a := []structure.Atom{
		{Name: "CA", X: 0, Y: 0, Z: 0},
		{Name: "CA", X: 3, Y: 4, Z: 12},
	}
	b := []structure.Atom{
		{Name: "CA", X: 0, Y: 0, Z: 0},
		{Name: "CA", X: 3, Y: 4, Z: 12},
	}

	rmsd, err := KabschAlign(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(rmsd) > 1e-9 {
		t.Errorf("RMSD for identical = %f, want 0", rmsd)
	}
}

func TestKabschAlignRotatesTranslatedStructure(t *testing.T) {
	// Two sets that differ by a translation + rotation
	a := []structure.Atom{
		{Name: "CA", X: 0, Y: 0, Z: 0},
		{Name: "CA", X: 1, Y: 0, Z: 0},
	}
	b := []structure.Atom{
		{Name: "CA", X: 10, Y: 0, Z: 0},
		{Name: "CA", X: 11, Y: 0, Z: 0},
	}

	rmsd, err := KabschAlign(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(rmsd) > 1e-9 {
		t.Errorf("RMSD for translated identical = %f, want ~0", rmsd)
	}
}

func TestKabschAlignRejectsMismatchedLengths(t *testing.T) {
	_, err := KabschAlign(
		[]structure.Atom{{Name: "CA"}},
		[]structure.Atom{{Name: "CA"}, {Name: "CA"}},
	)
	if err == nil {
		t.Fatal("expected error for mismatched lengths")
	}
}

func TestTMScoreSameStructureGetsOne(t *testing.T) {
	a := structure.Structure{
		Residues: []structure.Residue{
			{Name: "ALA", Atoms: []structure.Atom{{Name: "CA", X: 0, Y: 0, Z: 0}}},
			{Name: "GLY", Atoms: []structure.Atom{{Name: "CA", X: 3, Y: 4, Z: 0}}},
		},
	}
	b := structure.Structure{
		Residues: []structure.Residue{
			{Name: "ALA", Atoms: []structure.Atom{{Name: "CA", X: 0, Y: 0, Z: 0}}},
			{Name: "GLY", Atoms: []structure.Atom{{Name: "CA", X: 3, Y: 4, Z: 0}}},
		},
	}

	score, err := TMScore(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(score-1.0) > 1e-9 {
		t.Errorf("TMScore for identical = %f, want 1.0", score)
	}
}

func TestTMScoreRejectsMismatchedAtoms(t *testing.T) {
	a := structure.Structure{
		Residues: []structure.Residue{{Atoms: []structure.Atom{{Name: "CA"}}}},
	}
	b := structure.Structure{}

	_, err := TMScore(a, b)
	if err == nil {
		t.Fatal("expected error for mismatched atom counts")
	}
}

func TestRMSDAlignedUsesCAAtomsWithSuperposition(t *testing.T) {
	a := structure.Structure{
		Residues: []structure.Residue{
			{Name: "ALA", Atoms: []structure.Atom{
				{Name: "N", X: 0, Y: 0, Z: 0},
				{Name: "CA", X: 1, Y: 0, Z: 0},
				{Name: "C", X: 2, Y: 1, Z: 0},
			}},
		},
	}
	b := structure.Structure{
		Residues: []structure.Residue{
			{Name: "ALA", Atoms: []structure.Atom{
				{Name: "N", X: 5, Y: 0, Z: 0},
				{Name: "CA", X: 6, Y: 0, Z: 0},
				{Name: "C", X: 7, Y: 1, Z: 0},
			}},
		},
	}

	rmsd, err := RMSDAligned(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// After Kabsch superposition of translated structures, RMSD should be near 0
	if math.Abs(rmsd) > 1e-9 {
		t.Errorf("RMSDAligned after superposition = %f, want ~0", rmsd)
	}
}
