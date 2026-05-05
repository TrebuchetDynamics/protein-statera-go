package ui

type Panel struct {
	Name        string
	Description string
}

func DefaultPanels() []Panel {
	return []Panel{
		{Name: "Dashboard", Description: "Structure overview and quality metrics"},
		{Name: "Confidence", Description: "AlphaFold pLDDT confidence analysis"},
		{Name: "Geometry", Description: "Steric clash detection and geometry validation"},
		{Name: "Evidence", Description: "Evidence report and methodology notes"},
	}
}
