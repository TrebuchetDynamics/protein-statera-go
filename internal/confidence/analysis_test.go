package confidence

import (
	"testing"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

func TestAnalyzeCountsPLDDTBandsAndSegments(t *testing.T) {
	s := structure.Structure{
		ID: "AF-SAMPLE",
		Residues: []structure.Residue{
			{Name: "MET", Index: 1, PLDDT: 95},
			{Name: "GLY", Index: 2, PLDDT: 88},
			{Name: "SER", Index: 3, PLDDT: 52},
			{Name: "ALA", Index: 4, PLDDT: 47},
			{Name: "THR", Index: 5, PLDDT: 42},
			{Name: "LYS", Index: 6, PLDDT: 91},
		},
	}

	result := Analyze(s)

	if result.High != 2 || result.Medium != 2 || result.Low != 2 {
		t.Fatalf("bands = high:%d medium:%d low:%d, want 2/2/2", result.High, result.Medium, result.Low)
	}
	if len(result.LowSegments) != 1 {
		t.Fatalf("low segment count = %d, want 1", len(result.LowSegments))
	}
	segment := result.LowSegments[0]
	if segment.Start != 4 || segment.End != 5 || segment.Count != 2 {
		t.Fatalf("low segment = %+v, want 4-5 count 2", segment)
	}
}
