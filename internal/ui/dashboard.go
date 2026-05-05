package ui

const (
	ProteinSeed       = "2F5D50"
	ProteinAccent     = "1F463C"
	ProteinAccentLite = "6F9C8D"
)

func DashboardText() string {
	return `Protein Statera Go
Structure Viewer: gogpu/ui dashboard
Confidence Map: AlphaFold pLDDT bands
Geometry Report: VDW steric clash detection
Evidence Report: evidence-bound output
Comparison Viewer: RMSD between structures
`
}

type ThemeColors struct {
	Primary      string
	PrimaryDark  string
	PrimaryLight string
	Surface      string
	SurfaceDark  string
	Text         string
	TextLight    string
}

func DefaultThemeColors() ThemeColors {
	return ThemeColors{
		Primary:      "2F5D50",
		PrimaryDark:  "1F463C",
		PrimaryLight: "6F9C8D",
		Surface:      "F6FAF7",
		SurfaceDark:  "E7EFEA",
		Text:         "183D34",
		TextLight:    "44504B",
	}
}
