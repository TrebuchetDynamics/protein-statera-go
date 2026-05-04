package structure

import (
	"strings"
	"testing"
)

const samplePDB = `HEADER    SAMPLE ALPHAFOLD STRUCTURE
ATOM      1  N   MET A   1      11.104  13.207  14.118  1.00 95.10           N
ATOM      2  CA  MET A   1      12.300  13.900  14.500  1.00 94.20           C
ATOM      3  C   MET A   1      13.100  12.800  15.200  1.00 94.80           C
ATOM      4  N   GLY A   2      14.500  12.900  15.800  1.00 48.70           N
ATOM      5  CA  GLY A   2      15.300  11.800  16.500  1.00 49.20           C
HETATM    6  O   HOH A  10      20.000  20.000  20.000  1.00 10.00           O
END
`

func TestParsePDBGroupsAtomRecordsIntoResidues(t *testing.T) {
	structure, err := ParsePDB("AF-SAMPLE", strings.NewReader(samplePDB))
	if err != nil {
		t.Fatalf("ParsePDB returned error: %v", err)
	}

	if structure.ID != "AF-SAMPLE" {
		t.Fatalf("ID = %q, want AF-SAMPLE", structure.ID)
	}
	if structure.AtomCount() != 5 {
		t.Fatalf("AtomCount = %d, want 5", structure.AtomCount())
	}
	if len(structure.Residues) != 2 {
		t.Fatalf("residue count = %d, want 2", len(structure.Residues))
	}
	if structure.Residues[0].Name != "MET" || structure.Residues[0].Index != 1 || structure.Residues[0].ChainID != "A" {
		t.Fatalf("first residue = %+v, want MET A 1", structure.Residues[0])
	}
	if got := structure.Residues[0].Atoms[1].Name; got != "CA" {
		t.Fatalf("atom name = %q, want CA", got)
	}
	if got := structure.Residues[0].Atoms[1].BFactor; got != 94.20 {
		t.Fatalf("BFactor = %.2f, want 94.20", got)
	}
}

func TestParsePDBReportsMalformedCoordinatesWithLineNumber(t *testing.T) {
	_, err := ParsePDB("bad", strings.NewReader("ATOM      1  N   MET A   1      XX.XXX  13.207  14.118  1.00 95.10           N\n"))
	if err == nil {
		t.Fatal("ParsePDB returned nil error, want malformed coordinate error")
	}
	if !strings.Contains(err.Error(), "line 1") {
		t.Fatalf("error = %q, want line number", err.Error())
	}
}
