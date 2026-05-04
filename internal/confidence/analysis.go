package confidence

import "github.com/TrebuchetDynamics/protein-statera-go/internal/structure"

// Segment identifies a contiguous residue-index interval.
type Segment struct {
	Start int
	End   int
	Count int
}

// Analysis summarizes AlphaFold pLDDT confidence using EMBL-EBI standard bands.
//
// Bands (matching AlphaFold DB colour scheme):
//
//	VeryHigh   >= 90   (dark blue)
//	Confident  70-90  (light blue)
//	Low        50-70  (yellow)
//	VeryLow    < 50   (orange/red)
type Analysis struct {
	VeryHigh         int
	Confident        int
	Low              int
	VeryLow          int
	VeryLowSegments  []Segment
}

// Analyze counts residues into EMBL-EBI pLDDT bands and tracks very-low segments.
func Analyze(s structure.Structure) Analysis {
	result := Analysis{}
	active := Segment{}

	for _, residue := range s.Residues {
		switch {
		case residue.PLDDT >= VeryHighPLDDT:
			result.VeryHigh++
			active = closeSegment(active, &result)
		case residue.PLDDT >= ConfidentPLDDT:
			result.Confident++
			active = closeSegment(active, &result)
		case residue.PLDDT >= LowPLDDT:
			result.Low++
			active = closeSegment(active, &result)
		default:
			result.VeryLow++
			if active.Count == 0 {
				active = Segment{Start: residue.Index, End: residue.Index, Count: 1}
			} else {
				active.End = residue.Index
				active.Count++
			}
		}
	}
	closeSegment(active, &result)

	return result
}

func closeSegment(segment Segment, result *Analysis) Segment {
	if segment.Count > 0 {
		result.VeryLowSegments = append(result.VeryLowSegments, segment)
	}
	return Segment{}
}
