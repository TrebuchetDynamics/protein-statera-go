package confidence

import "github.com/TrebuchetDynamics/protein-statera-go/internal/structure"

// Segment identifies a contiguous residue-index interval.
type Segment struct {
	Start int
	End   int
	Count int
}

// Analysis summarizes AlphaFold pLDDT confidence bands.
type Analysis struct {
	High        int
	Medium      int
	Low         int
	LowSegments []Segment
}

// Analyze counts high, medium, and low pLDDT residues.
func Analyze(s structure.Structure) Analysis {
	result := Analysis{}
	activeLow := Segment{}

	for _, residue := range s.Residues {
		switch {
		case residue.PLDDT >= HighConfidencePLDDT:
			result.High++
			activeLow = closeLowSegment(activeLow, &result)
		case residue.PLDDT < LowConfidencePLDDT:
			result.Low++
			if activeLow.Count == 0 {
				activeLow = Segment{Start: residue.Index, End: residue.Index, Count: 1}
			} else {
				activeLow.End = residue.Index
				activeLow.Count++
			}
		default:
			result.Medium++
			activeLow = closeLowSegment(activeLow, &result)
		}
	}
	closeLowSegment(activeLow, &result)

	return result
}

func closeLowSegment(segment Segment, result *Analysis) Segment {
	if segment.Count > 0 {
		result.LowSegments = append(result.LowSegments, segment)
	}
	return Segment{}
}
