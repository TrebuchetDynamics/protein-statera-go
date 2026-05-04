package evidence

import (
	"github.com/TrebuchetDynamics/protein-statera-go/internal/confidence"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/geometry"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

// Report is the evidence-bound output model for CLI and renderers.
type Report struct {
	StructureID string
	Metrics     map[string]float64
	Confidence  map[string]int
	Notes       []string
	Evidence    EvidenceClass
	Source      string
}

// BuildStructureReport combines parsed structure data with validation metrics.
func BuildStructureReport(s structure.Structure, confidenceAnalysis confidence.Analysis, clashes geometry.ClashResult) Report {
	report := Report{
		StructureID: s.ID,
		Metrics: map[string]float64{
			"residues":          float64(len(s.Residues)),
			"atoms":             float64(s.AtomCount()),
			"steric_clashes":    float64(clashes.Total),
			"severe_clashes":    float64(clashes.Severe),
			"very_low_conf":     float64(confidenceAnalysis.VeryLow),
			"low_conf":          float64(confidenceAnalysis.Low),
			"confident_conf":    float64(confidenceAnalysis.Confident),
			"very_high_conf":    float64(confidenceAnalysis.VeryHigh),
		},
		Confidence: map[string]int{
			"very_high": confidenceAnalysis.VeryHigh,
			"confident": confidenceAnalysis.Confident,
			"low":       confidenceAnalysis.Low,
			"very_low":  confidenceAnalysis.VeryLow,
		},
		Evidence: EvidencePredictedStructure,
		Source:   s.Source,
	}

	if confidenceAnalysis.VeryHigh > confidenceAnalysis.Confident+confidenceAnalysis.Low+confidenceAnalysis.VeryLow {
		report.Notes = append(report.Notes, "Stable high-confidence core detected")
	}
	if confidenceAnalysis.VeryLow > 0 {
		report.Notes = append(report.Notes, "Very-low-confidence residues detected; inspect contiguous flexible or disordered regions")
	}
	if clashes.Total > 0 {
		report.Notes = append(report.Notes, "Steric clashes detected; review local geometry around reported atom pairs")
	}
	if len(report.Notes) == 0 {
		report.Notes = append(report.Notes, "No very-low-confidence residues or steric clashes detected by MVP checks")
	}

	return report
}
