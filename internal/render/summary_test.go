package render

import (
	"strings"
	"testing"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/confidence"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/evidence"
)

func TestStructureReportTextIncludesEvidenceMetrics(t *testing.T) {
	report := evidence.Report{
		StructureID: "AF-SAMPLE",
		Metrics: map[string]float64{
			"residues":       2,
			"atoms":          5,
			"steric_clashes": 1,
			"severe_clashes": 0,
		},
		Confidence: map[string]int{"very_high": 1, "confident": 0, "low": 0, "very_low": 1},
		Notes:      []string{"Very-low-confidence residues detected"},
		Evidence:   evidence.EvidencePredictedStructure,
		Source:     "sample_af.pdb",
	}

	text := StructureReportText(report)

	for _, want := range []string{
		"=== Protein Structure Report ===",
		"structure=AF-SAMPLE",
		"residues=2",
		"steric_clashes=1",
		"evidence_class=predicted_structure",
		"- Very-low-confidence residues detected",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("report text missing %q:\n%s", want, text)
		}
	}
}

func TestConfidenceTextIncludesLowSegments(t *testing.T) {
	text := ConfidenceText(confidence.Analysis{
		VeryHigh: 2, Confident: 1, Low: 0, VeryLow: 2,
		VeryLowSegments: []confidence.Segment{{Start: 4, End: 5, Count: 2}},
	})

	if !strings.Contains(text, "very_high (>90): 2") || !strings.Contains(text, "residues 4-5 count=2") {
		t.Fatalf("confidence text missing expected counts or segment:\n%s", text)
	}
}

func TestHTMLReportEscapesEvidenceFields(t *testing.T) {
	report := evidence.Report{
		StructureID: "AF-<SAMPLE>",
		Metrics:     map[string]float64{"residues": 1},
		Confidence:  map[string]int{"very_high": 1, "confident": 0, "low": 0, "very_low": 0},
		Notes:       []string{"Check <loop>"},
		Evidence:    evidence.EvidencePredictedStructure,
	}

	html := HTMLReport(report)

	if !strings.Contains(html, "AF-&lt;SAMPLE&gt;") || strings.Contains(html, "Check <loop>") {
		t.Fatalf("HTML escaping failed:\n%s", html)
	}
}
