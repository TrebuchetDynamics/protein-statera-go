package structure

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParsePDB parses ATOM records from a PDB stream.
func ParsePDB(id string, r io.Reader) (Structure, error) {
	scanner := bufio.NewScanner(r)
	result := Structure{ID: id}
	lineNumber := 0
	currentKey := residueKey{}
	current := Residue{}
	haveCurrent := false
	bFactorSum := 0.0

	flush := func() {
		if !haveCurrent {
			return
		}
		if len(current.Atoms) > 0 {
			current.PLDDT = bFactorSum / float64(len(current.Atoms))
		}
		result.Residues = append(result.Residues, current)
	}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if !strings.HasPrefix(line, "ATOM  ") {
			continue
		}

		atom, err := parseAtomLine(line)
		if err != nil {
			return Structure{}, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		key := residueKey{name: atom.ResidueName, index: atom.ResidueIndex, chainID: atom.ChainID}
		if !haveCurrent || key != currentKey {
			flush()
			currentKey = key
			current = Residue{Name: atom.ResidueName, Index: atom.ResidueIndex, ChainID: atom.ChainID}
			bFactorSum = 0
			haveCurrent = true
		}
		current.Atoms = append(current.Atoms, atom)
		bFactorSum += atom.BFactor
	}
	if err := scanner.Err(); err != nil {
		return Structure{}, err
	}
	flush()

	return result, nil
}

type residueKey struct {
	name    string
	index   int
	chainID string
}

func parseAtomLine(line string) (Atom, error) {
	if len(line) < 54 {
		return Atom{}, fmt.Errorf("ATOM record too short: length=%d", len(line))
	}

	id, err := parseIntField(line, 6, 11, "atom serial")
	if err != nil {
		return Atom{}, err
	}
	residueIndex, err := parseIntField(line, 22, 26, "residue index")
	if err != nil {
		return Atom{}, err
	}
	x, err := parseFloatField(line, 30, 38, "x")
	if err != nil {
		return Atom{}, err
	}
	y, err := parseFloatField(line, 38, 46, "y")
	if err != nil {
		return Atom{}, err
	}
	z, err := parseFloatField(line, 46, 54, "z")
	if err != nil {
		return Atom{}, err
	}

	bFactor := 0.0
	if len(line) >= 66 {
		bFactor, err = parseFloatField(line, 60, 66, "B-factor")
		if err != nil {
			return Atom{}, err
		}
	}

	return Atom{
		ID:           id,
		Name:         strings.TrimSpace(slice(line, 12, 16)),
		ResidueName:  strings.TrimSpace(slice(line, 17, 20)),
		ChainID:      strings.TrimSpace(slice(line, 21, 22)),
		ResidueIndex: residueIndex,
		X:            x,
		Y:            y,
		Z:            z,
		BFactor:      bFactor,
		Element:      strings.TrimSpace(slice(line, 76, 78)),
	}, nil
}

func parseIntField(line string, start, end int, label string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(slice(line, start, end)))
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", label, err)
	}
	return value, nil
}

func parseFloatField(line string, start, end int, label string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(slice(line, start, end)), 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", label, err)
	}
	return value, nil
}

func slice(line string, start, end int) string {
	if start >= len(line) {
		return ""
	}
	if end > len(line) {
		end = len(line)
	}
	return line[start:end]
}
