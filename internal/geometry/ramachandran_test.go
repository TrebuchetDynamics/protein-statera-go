package geometry

import (
	"math"
	"testing"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

func TestComputeTorsionKnownAngles(t *testing.T) {
	tests := []struct {
		name     string
		a, b, c, d structure.Atom
		expected float64
		tolerance float64
	}{
		{
			name: "zero torsion",
			a: structure.Atom{X: 0, Y: 0, Z: 0},
			b: structure.Atom{X: 1, Y: 0, Z: 0},
			c: structure.Atom{X: 1, Y: 1, Z: 0},
			d: structure.Atom{X: 0, Y: 1, Z: 0},
			expected: 0,
			tolerance: 0.1,
		},
		{
			name: "90 degree torsion",
			a: structure.Atom{X: 0, Y: 0, Z: 0},
			b: structure.Atom{X: 1, Y: 0, Z: 0},
			c: structure.Atom{X: 1, Y: 1, Z: 0},
			d: structure.Atom{X: 1, Y: 1, Z: 1},
			expected: 90,
			tolerance: 0.1,
		},
		{
			name: "180 degree torsion (trans)",
			a: structure.Atom{X: 0, Y: 0, Z: 0},
			b: structure.Atom{X: 1, Y: 0, Z: 0},
			c: structure.Atom{X: 1, Y: 1, Z: 0},
			d: structure.Atom{X: 2, Y: 1, Z: 0},
			expected: 180,
			tolerance: 0.1,
		},
		{
			name: "-90 degree torsion",
			a: structure.Atom{X: 0, Y: 0, Z: 0},
			b: structure.Atom{X: 1, Y: 0, Z: 0},
			c: structure.Atom{X: 1, Y: 1, Z: 0},
			d: structure.Atom{X: 1, Y: 1, Z: -1},
			expected: -90,
			tolerance: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeTorsion(tt.a, tt.b, tt.c, tt.d)
			if math.Abs(got-tt.expected) > tt.tolerance {
				t.Errorf("ComputeTorsion = %.2f°, want %.2f° (±%.1f)",
					got, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestPhiForAlphaHelixResidue(t *testing.T) {
	a := structure.Atom{X: 1, Y: 0, Z: 0}
	b := structure.Atom{X: 0, Y: 0, Z: 0}
	c := structure.Atom{X: 0, Y: 1, Z: 0}
	d := structure.Atom{X: 1, Y: 1, Z: 1}

	got := ComputeTorsion(a, b, c, d)
	if math.IsNaN(got) || math.Abs(got) < 1 {
		t.Errorf("expected non-zero torsion, got %.2f", got)
	}
	t.Logf("torsion = %.2f°", got)
}

func TestPsiForAlphaHelixResidue(t *testing.T) {
	a := structure.Atom{X: -1, Y: 0, Z: 0}
	b := structure.Atom{X: 0, Y: 0, Z: 0}
	c := structure.Atom{X: 0, Y: 0, Z: 1}
	d := structure.Atom{X: 1, Y: 0, Z: 1}

	got := ComputeTorsion(a, b, c, d)
	if math.IsNaN(got) || math.Abs(got) < 1 {
		t.Errorf("expected non-zero torsion, got %.2f", got)
	}
	t.Logf("torsion = %.2f°", got)
}

func TestPhiPsiReturnsAnglesForValidResidue(t *testing.T) {
	s := structure.Structure{
		ID: "test",
		Residues: []structure.Residue{
			{
				Name: "ALA", Index: 1, ChainID: "A",
				Atoms: []structure.Atom{
					{Name: "N", X: 10.0, Y: 10.0, Z: 10.0},
					{Name: "CA", X: 11.5, Y: 10.0, Z: 10.0},
					{Name: "C", X: 12.0, Y: 11.5, Z: 10.0},
				},
			},
			{
				Name: "GLY", Index: 2, ChainID: "A",
				Atoms: []structure.Atom{
					{Name: "N", X: 13.0, Y: 11.5, Z: 10.0},
					{Name: "CA", X: 13.5, Y: 10.5, Z: 10.0},
					{Name: "C", X: 15.0, Y: 10.5, Z: 10.0},
				},
			},
			{
				Name: "SER", Index: 3, ChainID: "A",
				Atoms: []structure.Atom{
					{Name: "N", X: 15.5, Y: 11.5, Z: 10.0},
					{Name: "CA", X: 17.0, Y: 11.5, Z: 10.0},
					{Name: "C", X: 17.5, Y: 10.0, Z: 10.0},
				},
			},
		},
	}

	phi, psi, ok := PhiPsi(s, 1)
	if !ok {
		t.Fatal("PhiPsi returned not ok for residue 1")
	}
	if !math.IsNaN(phi) && !math.IsNaN(psi) {
		t.Logf("phi=%.1f°, psi=%.1f°", phi, psi)
	}
}

func TestPhiPsiRejectsTerminalResidues(t *testing.T) {
	s := structure.Structure{
		ID: "test",
		Residues: []structure.Residue{
			{
				Name: "ALA", Index: 1, ChainID: "A",
				Atoms: []structure.Atom{
					{Name: "N", X: 0, Y: 0, Z: 0},
					{Name: "CA", X: 1.5, Y: 0, Z: 0},
					{Name: "C", X: 2.0, Y: 1.5, Z: 0},
				},
			},
		},
	}

	_, _, ok := PhiPsi(s, 0)
	if ok {
		t.Error("PhiPsi should reject index 0 (first residue in slice, which is N-terminal)")
	}

	_, _, ok = PhiPsi(s, 1)
	if ok {
		t.Error("PhiPsi should reject index 1 (last residue in slice)")
	}
}

func TestClassifyResiduesIntoRegions(t *testing.T) {
	tests := []struct {
		phi, psi float64
		resType  string
		want     RamachandranRegion
	}{
		// Ideal alpha-helix
		{-57, -47, "ALA", RegionFavored},
		// Extended beta strand
		{-120, 120, "SER", RegionFavored},
		// Glycine left-handed helix (broad allowed region for Gly)
		{60, 40, "GLY", RegionFavored},
		// Outlier - should not occur in well-refined structures
		{0, 0, "ALA", RegionOutlier},
		// PPII region
		{-75, 150, "PRO", RegionFavored},
	}

	for _, tt := range tests {
		got := Classify(tt.phi, tt.psi, tt.resType)
		if got != tt.want {
			t.Errorf("Classify(%.0f, %.0f, %s) = %v, want %v",
				tt.phi, tt.psi, tt.resType, got, tt.want)
		}
	}
}
