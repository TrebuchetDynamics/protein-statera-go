package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/comparison"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/confidence"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/evidence"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/geometry"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/render"
	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "analyze":
		if len(args) != 2 {
			printUsage(stderr)
			return 2
		}
		return analyze(args[1], stdout, stderr, false)
	case "confidence":
		if len(args) != 2 {
			printUsage(stderr)
			return 2
		}
		return confidenceCommand(args[1], stdout, stderr)
	case "compare":
		if len(args) != 3 {
			printUsage(stderr)
			return 2
		}
		return compare(args[1], args[2], stdout, stderr)
	case "report":
		if len(args) != 2 && len(args) != 3 {
			printUsage(stderr)
			return 2
		}
		html := len(args) == 3 && args[2] == "--html"
		if len(args) == 3 && !html {
			printUsage(stderr)
			return 2
		}
		return analyze(args[1], stdout, stderr, html)
	default:
		printUsage(stderr)
		return 2
	}
}

func analyze(path string, stdout, stderr io.Writer, html bool) int {
	s, ok := loadStructure(path, stderr)
	if !ok {
		return 1
	}
	conf := confidence.Analyze(s)
	clashes := geometry.FindStericClashes(s, geometry.DefaultClashThresholdAngstroms)
	report := evidence.BuildStructureReport(s, conf, clashes)
	if html {
		fmt.Fprint(stdout, render.HTMLReport(report))
		return 0
	}
	fmt.Fprint(stdout, render.StructureReportText(report))
	return 0
}

func confidenceCommand(path string, stdout, stderr io.Writer) int {
	s, ok := loadStructure(path, stderr)
	if !ok {
		return 1
	}
	fmt.Fprint(stdout, render.ConfidenceText(confidence.Analyze(s)))
	return 0
}

func compare(pathA, pathB string, stdout, stderr io.Writer) int {
	a, ok := loadStructure(pathA, stderr)
	if !ok {
		return 1
	}
	b, ok := loadStructure(pathB, stderr)
	if !ok {
		return 1
	}
	rmsd, err := comparison.RMSD(a, b)
	if err != nil {
		fmt.Fprintf(stderr, "compare: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, render.ComparisonText(comparison.Result{RMSDAngstroms: rmsd}))
	return 0
}

func loadStructure(path string, stderr io.Writer) (structure.Structure, bool) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(stderr, "open %s: %v\n", path, err)
		return structure.Structure{}, false
	}
	defer f.Close()

	id := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	s, err := structure.ParsePDB(id, f)
	if err != nil {
		fmt.Fprintf(stderr, "parse %s: %v\n", path, err)
		return structure.Structure{}, false
	}
	s.Source = path
	return s, true
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: protein <command> [args]")
	fmt.Fprintln(w, "commands:")
	fmt.Fprintln(w, "  analyze <structure.pdb>")
	fmt.Fprintln(w, "  confidence <structure.pdb>")
	fmt.Fprintln(w, "  compare <a.pdb> <b.pdb>")
	fmt.Fprintln(w, "  report <structure.pdb> [--html]")
}
