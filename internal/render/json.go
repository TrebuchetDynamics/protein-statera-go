package render

import (
	"encoding/json"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/evidence"
)

type jsonReport struct {
	StructureID string             `json:"structure_id"`
	Source      string             `json:"source,omitempty"`
	Evidence    string             `json:"evidence_class"`
	Metrics     map[string]float64 `json:"metrics"`
	Confidence  map[string]int     `json:"confidence"`
	Notes       []string           `json:"notes"`
}

func JSONReport(report evidence.Report) string {
	jr := jsonReport{
		StructureID: report.StructureID,
		Source:      report.Source,
		Evidence:    string(report.Evidence),
		Metrics:     report.Metrics,
		Confidence:  report.Confidence,
		Notes:       report.Notes,
	}

	b, err := json.MarshalIndent(jr, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
