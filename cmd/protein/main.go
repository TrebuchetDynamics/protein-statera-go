package main

import (
	"fmt"
	"io"
	"io/fs"
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
		return handleAnalyze(args[1:], stdout, stderr)
	case "confidence":
		if len(args) != 2 {
			printUsage(stderr)
			return 2
		}
		return confidenceCommand(args[1], stdout, stderr)
	case "compare":
		if len(args) < 3 {
			printUsage(stderr)
			return 2
		}
		return compare(args[1], args[2], stdout, stderr)
	case "report":
		return handleReport(args[1:], stdout, stderr)
	case "batch":
		if len(args) != 2 {
			printUsage(stderr)
			return 2
		}
		return batch(args[1], stdout, stderr)
	default:
		printUsage(stderr)
		return 2
	}
}

func handleAnalyze(rest []string, stdout, stderr io.Writer) int {
	jsonFlag := hasFlag(rest, "--json")
	path := rest[0]
	for _, a := range rest {
		if a != "--json" && a != "--html" {
			path = a
			break
		}
	}
	if path == "" {
		printUsage(stderr)
		return 2
	}
	htmlFlag := hasFlag(rest, "--html")
	s, ok := loadStructure(path, stderr)
	if !ok {
		return 1
	}
	conf := confidence.Analyze(s)
	clashes := geometry.FindStericClashesVDW(s)
	report := evidence.BuildStructureReport(s, conf, clashes)
	if htmlFlag {
		fmt.Fprint(stdout, render.HTMLReport(report))
		return 0
	}
	if jsonFlag {
		fmt.Fprint(stdout, render.JSONReport(report))
		return 0
	}
	fmt.Fprint(stdout, render.StructureReportText(report))
	return 0
}

func handleReport(rest []string, stdout, stderr io.Writer) int {
	htmlFlag := hasFlag(rest, "--html")
	jsonFlag := hasFlag(rest, "--json")
	path := ""
	for _, a := range rest {
		if a != "--html" && a != "--json" {
			path = a
			break
		}
	}
	if path == "" {
		printUsage(stderr)
		return 2
	}
	s, ok := loadStructure(path, stderr)
	if !ok {
		return 1
	}
	conf := confidence.Analyze(s)
	clashes := geometry.FindStericClashesVDW(s)
	report := evidence.BuildStructureReport(s, conf, clashes)
	if htmlFlag {
		fmt.Fprint(stdout, render.HTMLReport(report))
		return 0
	}
	if jsonFlag {
		fmt.Fprint(stdout, render.JSONReport(report))
		return 0
	}
	fmt.Fprint(stdout, render.StructureReportText(report))
	return 0
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
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

func batch(dir string, stdout, stderr io.Writer) int {
	entries, err := fs.ReadDir(os.DirFS(dir), ".")
	if err != nil {
		fmt.Fprintf(stderr, "batch: %v\n", err)
		return 1
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".pdb") {
			continue
		}
		path := filepath.Join(dir, name)
		s, ok := loadStructure(path, stderr)
		if !ok {
			continue
		}
		conf := confidence.Analyze(s)
		clashes := geometry.FindStericClashesVDW(s)
		report := evidence.BuildStructureReport(s, conf, clashes)
		fmt.Fprintf(stdout, "%s\t%.0f\t%.0f\t%d\t%d\t%d\t%d\t%d\n",
			path,
			report.Metrics["residues"],
			report.Metrics["atoms"],
			report.Confidence["very_high"],
			report.Confidence["confident"],
			report.Confidence["low"],
			report.Confidence["very_low"],
			clashes.Total,
		)
	}
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
	fmt.Fprintln(w, "  analyze <structure.pdb> [--json] [--html]")
	fmt.Fprintln(w, "  confidence <structure.pdb>")
	fmt.Fprintln(w, "  compare <a.pdb> <b.pdb>")
	fmt.Fprintln(w, "  report <structure.pdb> [--json] [--html]")
	fmt.Fprintln(w, "  batch <directory>")
}
