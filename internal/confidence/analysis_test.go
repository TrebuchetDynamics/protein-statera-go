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

	if result.VeryHigh != 2 || result.Confident != 1 || result.Low != 1 || result.VeryLow != 2 {
		t.Fatalf("bands = VHigh:%d Conf:%d Low:%d VLow:%d, want 2/1/1/2",
			result.VeryHigh, result.Confident, result.Low, result.VeryLow)
	}
	if len(result.VeryLowSegments) != 1 {
		t.Fatalf("very-low segment count = %d, want 1", len(result.VeryLowSegments))
	}
	segment := result.VeryLowSegments[0]
	if segment.Start != 4 || segment.End != 5 || segment.Count != 2 {
		t.Fatalf("very-low segment = %+v, want 4-5 count 2", segment)
	}
}
