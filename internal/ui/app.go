package ui

import "io"

func Run(w io.Writer) {
	_, _ = w.Write([]byte("Protein Statera Go — run `go run ./cmd/protein-ui` for the gogpu/ui dashboard.\n"))
}
