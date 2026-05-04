package evidence

// EvidenceClass describes the source and strength of the report evidence.
type EvidenceClass string

const (
	EvidencePredictedStructure EvidenceClass = "predicted_structure"
	EvidenceExperimental       EvidenceClass = "experimental_structure"
	EvidenceComparison         EvidenceClass = "comparative_analysis"
)
