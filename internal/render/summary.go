package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/comparison"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/confidence"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/evidence"
)

// StructureReportText renders a screen-reader friendly evidence report.
func StructureReportText(report evidence.Report) string {
	var b strings.Builder
	fmt.Fprintln(&b, "=== Protein Structure Report ===")
	fmt.Fprintf(&b, "structure=%s\n", report.StructureID)
	if report.Source != "" {
		fmt.Fprintf(&b, "source=%s\n", report.Source)
	}
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "residues=%.0f\n", report.Metrics["residues"])
	fmt.Fprintf(&b, "atoms=%.0f\n", report.Metrics["atoms"])
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "confidence (EMBL-EBI bands):")
	fmt.Fprintf(&b, "  very_high (>90): %d\n", report.Confidence["very_high"])
	fmt.Fprintf(&b, "  confident (70-90): %d\n", report.Confidence["confident"])
	fmt.Fprintf(&b, "  low (50-70): %d\n", report.Confidence["low"])
	fmt.Fprintf(&b, "  very_low (<50): %d\n", report.Confidence["very_low"])
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "geometry:")
	fmt.Fprintf(&b, "  steric_clashes=%.0f\n", report.Metrics["steric_clashes"])
	fmt.Fprintf(&b, "  severe_clashes=%.0f\n", report.Metrics["severe_clashes"])
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "interpretation:")
	for _, note := range report.Notes {
		fmt.Fprintf(&b, "  - %s\n", note)
	}
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "evidence_class=%s\n", report.Evidence)
	return b.String()
}

// ConfidenceText renders pLDDT band counts and very-low-confidence segments.
func ConfidenceText(result confidence.Analysis) string {
	var b strings.Builder
	fmt.Fprintln(&b, "=== Confidence Analysis ===")
	fmt.Fprintf(&b, "very_high (>90): %d\n", result.VeryHigh)
	fmt.Fprintf(&b, "confident (70-90): %d\n", result.Confident)
	fmt.Fprintf(&b, "low (50-70): %d\n", result.Low)
	fmt.Fprintf(&b, "very_low (<50): %d\n", result.VeryLow)
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "very_low_confidence_segments:")
	if len(result.VeryLowSegments) == 0 {
		fmt.Fprintln(&b, "  none")
	} else {
		for _, segment := range result.VeryLowSegments {
			fmt.Fprintf(&b, "  residues %d-%d count=%d\n", segment.Start, segment.End, segment.Count)
		}
	}
	return b.String()
}

// ComparisonText renders quantitative comparison metrics.
func ComparisonText(result comparison.Result) string {
	var b strings.Builder
	fmt.Fprintln(&b, "=== Structure Comparison ===")
	fmt.Fprintf(&b, "RMSD = %.2f A\n", result.RMSDAngstroms)
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "interpretation:")
	fmt.Fprintln(&b, "  - MVP comparison assumes atom-order aligned structures")
	return b.String()
}

func writeAll(w io.Writer, text string) error {
	_, err := io.WriteString(w, text)
	return err
}
