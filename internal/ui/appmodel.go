package ui

import (
	"fmt"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/confidence"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/evidence"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/geometry"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

// Spec holds window parameters.
type Spec struct {
	Title  string
	Width  int
	Height int
}

// ViewSpec identifies a navigable view.
type ViewSpec struct {
	Name        string
	Description string
}

// Summary is the top-level dashboard summary.
type Summary struct {
	StructureName   string
	ResidueCount    int
	AtomCount       int
	VeryHighConf    int
	ConfidentConf   int
	LowConf         int
	VeryLowConf     int
	StericClashes   int
	SevereClashes   int
	BoundaryNotice  string
}

// StructureRecord holds a parsed structure for display.
type StructureRecord struct {
	ID             string
	Source         string
	ResidueCount   int
	AtomCount      int
	VeryHighConf   int
	ConfidentConf  int
	LowConf        int
	VeryLowConf    int
	StericClashes  int
	SevereClashes  int
	PLDDTRange     string
	Notes          []string
}

// ConfidenceRecord holds per-residue confidence data.
type ConfidenceRecord struct {
	ResidueName string
	Index       int
	ChainID     string
	PLDDT       float64
	Band        string
}

// ClashRecord holds a single steric clash.
type ClashRecord struct {
	AtomA      string
	ResidueA   string
	ChainA     string
	AtomB      string
	ResidueB   string
	ChainB     string
	Distance   float64
	Severe     bool
}

// AppModel bundles all data for the gogpu/ui dashboard.
type AppModel struct {
	Spec          Spec
	Views         []ViewSpec
	Summary       Summary
	Structures    []StructureRecord
	Confidence    []ConfidenceRecord
	Clashes       []ClashRecord
	VeryLowSegs   []confidence.Segment
	DemoResidues  []string
}

// DefaultSpec returns the default window spec.
func DefaultSpec() Spec {
	return Spec{
		Title:  "Protein Statera Go",
		Width:  1180,
		Height: 760,
	}
}

// DefaultModel builds the AppModel with default structure data.
func DefaultModel() AppModel {
	spec := DefaultSpec()
	s := defaultStructure()
	conf := confidence.Analyze(s)
	clashes := geometry.FindStericClashesVDW(s)

	summary := Summary{
		StructureName:  s.ID,
		ResidueCount:   len(s.Residues),
		AtomCount:      s.AtomCount(),
		VeryHighConf:   conf.VeryHigh,
		ConfidentConf:  conf.Confident,
		LowConf:        conf.Low,
		VeryLowConf:    conf.VeryLow,
		StericClashes:  clashes.Total,
		SevereClashes:  clashes.Severe,
		BoundaryNotice: "This is a structure validation workbench. Model predictions are not experimental measurements.",
	}

	structRec := StructureRecord{
		ID:            s.ID,
		Source:        s.Source,
		ResidueCount:  len(s.Residues),
		AtomCount:     s.AtomCount(),
		VeryHighConf:  conf.VeryHigh,
		ConfidentConf: conf.Confident,
		LowConf:       conf.Low,
		VeryLowConf:   conf.VeryLow,
		StericClashes: clashes.Total,
		SevereClashes: clashes.Severe,
		PLDDTRange:    fmt.Sprintf("%.1f-%.1f", minPLDDT(s), maxPLDDT(s)),
	}
	report := evidence.BuildStructureReport(s, conf, clashes)
	structRec.Notes = report.Notes

	confidenceRecs := make([]ConfidenceRecord, 0, len(s.Residues))
	for _, r := range s.Residues {
		band := "very-low"
		switch {
		case r.PLDDT >= confidence.VeryHighPLDDT:
			band = "very-high"
		case r.PLDDT >= confidence.ConfidentPLDDT:
			band = "confident"
		case r.PLDDT >= confidence.LowPLDDT:
			band = "low"
		}
		confidenceRecs = append(confidenceRecs, ConfidenceRecord{
			ResidueName: r.Name,
			Index:       r.Index,
			ChainID:     r.ChainID,
			PLDDT:       r.PLDDT,
			Band:        band,
		})
	}

	clashRecs := make([]ClashRecord, 0, len(clashes.Pairs))
	for _, pair := range clashes.Pairs {
		clashRecs = append(clashRecs, ClashRecord{
			AtomA:    pair.AtomA.Name,
			ResidueA: fmt.Sprintf("%s%d", pair.AtomA.ResidueName, pair.AtomA.ResidueIndex),
			ChainA:   pair.AtomA.ChainID,
			AtomB:    pair.AtomB.Name,
			ResidueB: fmt.Sprintf("%s%d", pair.AtomB.ResidueName, pair.AtomB.ResidueIndex),
			ChainB:   pair.AtomB.ChainID,
			Distance: pair.Distance,
			Severe:   pair.Distance < 2.0,
		})
	}

	demoResidues := make([]string, 0)
	for _, r := range s.Residues {
		if len(demoResidues) >= 12 {
			break
		}
		demoResidues = append(demoResidues, fmt.Sprintf("%s%d", r.Name, r.Index))
	}

	return AppModel{
		Spec: spec,
		Views: []ViewSpec{
			{Name: "Dashboard", Description: "Structure overview and quality metrics"},
			{Name: "Confidence", Description: "AlphaFold pLDDT confidence analysis"},
			{Name: "Geometry", Description: "Steric clash detection and geometry validation"},
			{Name: "Evidence", Description: "Evidence report and methodology notes"},
		},
		Summary:      summary,
		Structures:   []StructureRecord{structRec},
		Confidence:   confidenceRecs,
		Clashes:      clashRecs,
		VeryLowSegs:  conf.VeryLowSegments,
		DemoResidues: demoResidues,
	}
}

// defaultStructure returns a synthetic mini protein structure for demo display.
// In production, structures would be loaded via ParsePDB from user-provided files.
func defaultStructure() structure.Structure {
	atoms := func(name string, x, y, z, b float64, elem string) structure.Atom {
		return structure.Atom{Name: name, X: x, Y: y, Z: z, BFactor: b, Element: elem}
	}

	return structure.Structure{
		ID:     "demo_1abc",
		Source: "synthetic-demo (no PDB file loaded)",
		Residues: []structure.Residue{
			{Name: "ALA", Index: 1, ChainID: "A", Atoms: []structure.Atom{
				atoms("N", -1.5, 0.0, 0.0, 95.0, "N"),
				atoms("CA", 0.0, 0.0, 0.0, 95.0, "C"),
				atoms("C", 1.5, 0.0, 0.0, 95.0, "C"),
				atoms("O", 2.0, 1.0, 0.0, 95.0, "O"),
				atoms("CB", -0.5, 1.3, 0.0, 95.0, "C"),
			}, PLDDT: 95.0},
			{Name: "GLY", Index: 2, ChainID: "A", Atoms: []structure.Atom{
				atoms("N", 2.2, -1.0, 0.0, 92.0, "N"),
				atoms("CA", 3.7, -1.0, 0.0, 92.0, "C"),
				atoms("C", 4.5, 0.3, 0.0, 92.0, "C"),
				atoms("O", 4.0, 1.4, 0.0, 92.0, "O"),
			}, PLDDT: 92.0},
			{Name: "SER", Index: 3, ChainID: "A", Atoms: []structure.Atom{
				atoms("N", 5.8, 0.2, 0.0, 88.0, "N"),
				atoms("CA", 6.8, 1.3, 0.0, 88.0, "C"),
				atoms("C", 8.2, 1.0, 0.0, 88.0, "C"),
				atoms("O", 8.8, -0.1, 0.0, 88.0, "O"),
				atoms("CB", 6.4, 2.7, 0.0, 85.0, "C"),
				atoms("OG", 7.0, 3.5, 1.2, 82.0, "O"),
			}, PLDDT: 86.0},
			{Name: "THR", Index: 4, ChainID: "A", Atoms: []structure.Atom{
				atoms("N", 8.9, 2.0, 0.0, 75.0, "N"),
				atoms("CA", 10.3, 1.8, 0.0, 75.0, "C"),
				atoms("C", 11.0, 3.2, 0.0, 75.0, "C"),
				atoms("O", 10.4, 4.3, 0.0, 75.0, "O"),
				atoms("CB", 10.8, 0.8, 1.2, 70.0, "C"),
				atoms("OG1", 10.3, -0.5, 1.0, 65.0, "O"),
			}, PLDDT: 72.0},
			{Name: "VAL", Index: 5, ChainID: "A", Atoms: []structure.Atom{
				atoms("N", 12.3, 3.3, 0.0, 55.0, "N"),
				atoms("CA", 13.1, 4.6, 0.0, 55.0, "C"),
				atoms("C", 14.6, 4.5, 0.0, 55.0, "C"),
				atoms("O", 15.2, 3.4, 0.0, 55.0, "O"),
				atoms("CB", 12.7, 5.5, 1.2, 48.0, "C"),
			}, PLDDT: 52.0},
			{Name: "LYS", Index: 6, ChainID: "A", Atoms: []structure.Atom{
				atoms("N", 15.2, 5.7, 0.0, 42.0, "N"),
				atoms("CA", 16.6, 5.6, 0.0, 42.0, "C"),
				atoms("C", 17.3, 7.0, 0.0, 42.0, "C"),
				atoms("O", 16.7, 8.1, 0.0, 42.0, "O"),
				atoms("CB", 17.1, 4.7, 1.3, 35.0, "C"),
			}, PLDDT: 40.0},
		},
	}
}

func minPLDDT(s structure.Structure) float64 {
	if len(s.Residues) == 0 {
		return 0
	}
	min := s.Residues[0].PLDDT
	for _, r := range s.Residues[1:] {
		if r.PLDDT < min {
			min = r.PLDDT
		}
	}
	return min
}

func maxPLDDT(s structure.Structure) float64 {
	if len(s.Residues) == 0 {
		return 0
	}
	max := s.Residues[0].PLDDT
	for _, r := range s.Residues[1:] {
		if r.PLDDT > max {
			max = r.PLDDT
		}
	}
	return max
}


