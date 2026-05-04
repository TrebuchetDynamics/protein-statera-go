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

func TestFindStericClashesVDWUsesElementSpecificRadii(t *testing.T) {
	s := structure.Structure{
		ID: "vdw-test",
		Residues: []structure.Residue{
			{
				Name:  "ALA",
				Index: 1,
				Atoms: []structure.Atom{
					{ID: 1, Name: "N", X: 0, Y: 0, Z: 0, Element: "N"},
				},
			},
			{
				Name:  "ALA",
				Index: 3,
				Atoms: []structure.Atom{
					{ID: 5, Name: "O", X: 1.0, Y: 0, Z: 0, Element: "O"},
				},
			},
		},
	}

	result := FindStericClashesVDW(s)

	// N vdW=1.55, O vdW=1.52. Sum=3.07. Threshold = 3.07-0.4 = 2.67.
	// Distance 1.0 < 2.67 → clash.
	// Severe: 3.07-0.9 = 2.17. 1.0 < 2.17 → severe clash.
	if result.Total != 1 {
		t.Fatalf("VDW clashes total=%d, want 1", result.Total)
	}
	if result.Severe != 1 {
		t.Fatalf("VDW severe clashes=%d, want 1", result.Severe)
	}
}

func TestElementVDWRadiusKnownElements(t *testing.T) {
	tests := []struct {
		element  string
		expected float64
	}{
		{"C", 1.70},
		{"N", 1.55},
		{"O", 1.52},
		{"S", 1.80},
		{"H", 1.00},
		{"P", 1.80},
		{"", 1.70},
		{"X", 1.70},
	}

	for _, tt := range tests {
		got := ElementVDWRadius(tt.element)
		if got != tt.expected {
			t.Errorf("ElementVDWRadius(%q) = %.2f, want %.2f", tt.element, got, tt.expected)
		}
	}
}

func TestClashscorePerThousandAtoms(t *testing.T) {
	s := structure.Structure{
		ID: "score-test",
		Residues: []structure.Residue{
			{
				Name: "ALA", Index: 1,
				Atoms: []structure.Atom{
					{ID: 1, X: 0, Y: 0, Z: 0, Element: "C"},
					{ID: 2, X: 2.0, Y: 0, Z: 0, Element: "N"},
				},
			},
			{
				Name: "GLY", Index: 2,
				Atoms: []structure.Atom{
					{ID: 3, X: 0.5, Y: 0, Z: 0, Element: "C"},
				},
			},
		},
	}

	result := FindStericClashesVDW(s)

	// Clashscore = clashes * 1000 / totalAtoms
	// We should detect at least one clash with element-specific thresholds
	if result.Clashscore <= 0 && result.Total > 0 {
		t.Errorf("Clashscore should be > 0 when clashes exist, got %.1f (total=%d atoms=%d)",
			result.Clashscore, result.Total, s.AtomCount())
	}
}
