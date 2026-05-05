package ui

const (
	DefaultTitle  = "Protein Statera Go"
	DefaultWidth  = 1180
	DefaultHeight = 760
)

func DefaultAppTheme() struct {
	Primary string
	Surface string
	Text    string
	Border  string
} {
	return struct {
		Primary string
		Surface string
		Text    string
		Border  string
	}{
		Primary: "2F5D50",
		Surface: "F6FAF7",
		Text:    "183D34",
		Border:  "D4DED8",
	}
}
