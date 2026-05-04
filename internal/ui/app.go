package ui

import (
	"fmt"
	"io"
)

// Run prints the decoupled MVP dashboard until a gogpu/ui frontend is added.
func Run(w io.Writer) {
	fmt.Fprint(w, DashboardText())
}
