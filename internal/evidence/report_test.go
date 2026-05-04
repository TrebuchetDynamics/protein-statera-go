package evidence

import (
	"strings"
	"testing"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/confidence"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/geometry"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

func TestBuildStructureReportAddsMetricsAndInterpretation(t *testing.T) {
	s := structure.Structure{
		ID:     "AF-SAMPLE",
		Source: "sample_af.pdb",
		Residues: []structure.Residue{
			{Name: "MET", Index: 1, PLDDT: 95, Atoms: []structure.Atom{{ID: 1}}},
			{Name: "GLY", Index: 2, PLDDT: 45, Atoms: []structure.Atom{{ID: 2}}},
		},
	}

	report := BuildStructureReport(s, confidence.Analyze(s), geometry.ClashResult{Total: 2, Severe: 0})

	if report.StructureID != "AF-SAMPLE" {
		t.Fatalf("StructureID = %q, want AF-SAMPLE", report.StructureID)
	}
	if report.Metrics["residues"] != 2 || report.Metrics["atoms"] != 2 || report.Metrics["steric_clashes"] != 2 {
		t.Fatalf("metrics = %#v, want residues=2 atoms=2 steric_clashes=2", report.Metrics)
	}
	if report.Evidence != EvidencePredictedStructure {
		t.Fatalf("Evidence = %q, want predicted_structure", report.Evidence)
	}
	if !strings.Contains(strings.Join(report.Notes, " "), "Low-confidence residues detected") {
		t.Fatalf("notes = %#v, want low-confidence interpretation", report.Notes)
	}
}
