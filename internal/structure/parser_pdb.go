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
	result := Structure{ID: id, Models: [][]Residue{}}
	lineNumber := 0
	currentKey := residueKey{}
	current := Residue{}
	haveCurrent := false
	bFactorSum := 0.0
	inModel := false

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

		if strings.HasPrefix(line, "MODEL ") {
			flush()
			haveCurrent = false
			inModel = true
			result.Residues = nil
			continue
		}
		if strings.HasPrefix(line, "ENDMDL") {
			flush()
			haveCurrent = false
			inModel = false
			result.Models = append(result.Models, result.Residues)
			result.Residues = nil
			continue
		}
		if !strings.HasPrefix(line, "ATOM  ") {
			continue
		}

		atom, err := parseAtomLine(line)
		if err != nil {
			return Structure{}, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		key := residueKey{
			name:    atom.ResidueName,
			index:   atom.ResidueIndex,
			chainID: atom.ChainID,
			iCode:   atom.ICode,
		}
		if !haveCurrent || key != currentKey {
			flush()
			currentKey = key
			current = Residue{
				Name:    atom.ResidueName,
				Index:   atom.ResidueIndex,
				ChainID: atom.ChainID,
				ICode:   atom.ICode,
			}
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

	if inModel {
		result.Models = append(result.Models, result.Residues)
	}

	return result, nil
}

type residueKey struct {
	name    string
	index   int
	chainID string
	iCode   string
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

	bfactor := 0.0
	if len(line) >= 66 {
		bfactor, err = parseFloatField(line, 60, 66, "B-factor")
		if err != nil {
			return Atom{}, err
		}
	}

	occupancy := 1.0
	if len(line) >= 60 {
		occ, err := parseFloatField(line, 54, 60, "occupancy")
		if err == nil {
			occupancy = occ
		}
	}

	altLoc := ""
	if len(line) >= 17 {
		altLoc = strings.TrimSpace(line[16:17])
	}

	iCode := ""
	if len(line) >= 27 {
		iCode = strings.TrimSpace(line[26:27])
	}

	elem := ""
	if len(line) >= 78 {
		elem = strings.TrimSpace(slice(line, 76, 78))
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
		BFactor:      bfactor,
		Occupancy:    occupancy,
		AltLoc:       altLoc,
		ICode:        iCode,
		Element:      elem,
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
