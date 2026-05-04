package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAnalyzePrintsStructureReport(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"analyze", filepath.Join("..", "..", "data", "examples", "sample_af.pdb")}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "=== Protein Structure Report ===") {
		t.Fatalf("stdout missing report header:\n%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "evidence_class=predicted_structure") {
		t.Fatalf("stdout missing evidence class:\n%s", stdout.String())
	}
}

func TestRunComparePrintsRMSD(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"compare",
		filepath.Join("..", "..", "data", "examples", "sample_af.pdb"),
		filepath.Join("..", "..", "data", "examples", "sample_exp.pdb"),
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "=== Structure Comparison ===") || !strings.Contains(stdout.String(), "RMSD =") {
		t.Fatalf("stdout missing comparison report:\n%s", stdout.String())
	}
}

func TestRunReturnsUsageErrorForUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"unknown"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "usage: protein") {
		t.Fatalf("stderr missing usage:\n%s", stderr.String())
	}
}
