package geometry

import (
	"math"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/structure"
)

type BondOutlier struct {
	AtomA    structure.Atom
	AtomB    structure.Atom
	Observed float64
	Ideal    float64
	Sigma    float64
	ZScore   float64
	Severity string
}

type BondResult struct {
	Outliers []BondOutlier
	RMSE     float64
}

// Bond deviation at >4σ is an outlier, >6σ is severe.
const (
	bondOutlierSigma = 4.0
	bondSevereSigma  = 6.0
)

func ValidateBondLengths(s structure.Structure) BondResult {
	result := BondResult{}
	var sumSq float64
	var count int

	for _, r := range s.Residues {
		pairs := expectedBonds(r)
		for _, pair := range pairs {
			ideal, sigma, ok := enghHuberLookup(r.Name, pair.a, pair.b)
			if !ok {
				continue
			}
			atomA := r.AtomByName(pair.a)
			atomB := r.AtomByName(pair.b)
			if atomA == nil || atomB == nil {
				continue
			}
			observed := DistanceAngstroms(*atomA, *atomB)
			delta := math.Abs(observed - ideal)
			z := delta / sigma
			sumSq += delta * delta
			count++
			if z > bondOutlierSigma {
				sev := ""
				if z > bondSevereSigma {
					sev = "severe"
				}
				result.Outliers = append(result.Outliers, BondOutlier{
					AtomA:    *atomA,
					AtomB:    *atomB,
					Observed: observed,
					Ideal:    ideal,
					Sigma:    sigma,
					ZScore:   z,
					Severity: sev,
				})
			}
		}
	}
	if count > 0 {
		result.RMSE = math.Sqrt(sumSq / float64(count))
	}
	return result
}

type bondPair struct{ a, b string }

func expectedBonds(r structure.Residue) []bondPair {
	known := map[string][]bondPair{
		// Standard amino acid backbone connectivity
		"__backbone__": {
			{"N", "CA"}, {"CA", "C"}, {"C", "O"},
		},
	}
	pairs := known["__backbone__"]
	switch r.Name {
	case "ALA":
		pairs = append(pairs, bondPair{"CA", "CB"})
	case "ARG":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD"}, bondPair{"CD", "NE"},
			bondPair{"NE", "CZ"}, bondPair{"CZ", "NH1"},
			bondPair{"CZ", "NH2"},
		)
	case "ASN":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "OD1"}, bondPair{"CG", "ND2"},
		)
	case "ASP":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "OD1"}, bondPair{"CG", "OD2"},
		)
	case "CYS":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "SG"},
		)
	case "GLN":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD"}, bondPair{"CD", "OE1"},
			bondPair{"CD", "NE2"},
		)
	case "GLU":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD"}, bondPair{"CD", "OE1"},
			bondPair{"CD", "OE2"},
		)
	case "GLY":
		// Glycine has no CB
	case "HIS":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "ND1"}, bondPair{"ND1", "CE1"},
			bondPair{"CE1", "NE2"}, bondPair{"NE2", "CD2"},
			bondPair{"CD2", "CG"},
		)
	case "ILE":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG1"},
			bondPair{"CB", "CG2"}, bondPair{"CG1", "CD1"},
		)
	case "LEU":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD1"}, bondPair{"CG", "CD2"},
		)
	case "LYS":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD"}, bondPair{"CD", "CE"},
			bondPair{"CE", "NZ"},
		)
	case "MET":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "SD"}, bondPair{"SD", "CE"},
		)
	case "PHE":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD1"}, bondPair{"CD1", "CE1"},
			bondPair{"CE1", "CZ"}, bondPair{"CZ", "CE2"},
			bondPair{"CE2", "CD2"}, bondPair{"CD2", "CG"},
		)
	case "PRO":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD"}, bondPair{"CD", "N"},
		)
	case "SER":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "OG"},
		)
	case "THR":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "OG1"},
			bondPair{"CB", "CG2"},
		)
	case "TRP":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD1"}, bondPair{"CG", "CD2"},
			bondPair{"CD1", "NE1"}, bondPair{"NE1", "CE2"},
			bondPair{"CE2", "CD2"}, bondPair{"CE2", "CZ2"},
			bondPair{"CZ2", "CH2"}, bondPair{"CH2", "CZ3"},
			bondPair{"CZ3", "CE3"}, bondPair{"CE3", "CD2"},
		)
	case "TYR":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG"},
			bondPair{"CG", "CD1"}, bondPair{"CD1", "CE1"},
			bondPair{"CE1", "CZ"}, bondPair{"CZ", "OH"},
			bondPair{"CZ", "CE2"}, bondPair{"CE2", "CD2"},
			bondPair{"CD2", "CG"},
		)
	case "VAL":
		pairs = append(pairs,
			bondPair{"CA", "CB"}, bondPair{"CB", "CG1"},
			bondPair{"CB", "CG2"},
		)
	}
	return pairs
}

func enghHuberLookup(resName, atomA, atomB string) (ideal, sigma float64, ok bool) {
	key := enghKey(resName, atomA, atomB)
	v, found := enghHuberBonds[key]
	return v.mean, v.sigma, found
}

func enghKey(res, a, b string) string {
	if a > b {
		a, b = b, a
	}
	if res == "" {
		return a + "-" + b
	}
	return res + ":" + a + "-" + b
}

type enghBond struct{ mean, sigma float64 }

var enghHuberBonds = map[string]enghBond{
	// Backbone (Engh & Huber 1991, Table 5)
	"N-CA":    {1.458, 0.019},
	"CA-C":    {1.525, 0.021},
	"C-O":     {1.231, 0.020},
	"C-N":     {1.329, 0.014}, // peptide bond (inter-residue)

	// Sidechain C-C bonds (generic aliphatic)
	"CA-CB":   {1.530, 0.020},
	"CB-CG":   {1.520, 0.020},
	"CG-CD":   {1.520, 0.020},
	"CD-CE":   {1.520, 0.020},
	"CB-CG1":  {1.530, 0.020},
	"CB-CG2":  {1.530, 0.020},
	"CG1-CD1": {1.520, 0.020},

	// C-O bonds
	"CB-OG":   {1.430, 0.020},
	"CB-OG1":  {1.430, 0.020},
	"CZ-OH":   {1.380, 0.020},
	"CG-OD1":  {1.250, 0.020},
	"CG-OD2":  {1.250, 0.020},
	"CD-OE1":  {1.250, 0.020},
	"CD-OE2":  {1.250, 0.020},

	// C-N bonds (sidechain)
	"CB-SG":   {1.810, 0.020},
	"CG-SD":   {1.810, 0.020},
	"SD-CE":   {1.810, 0.020},
	"CG-ND1":  {1.380, 0.020},
	"CE-NE2":  {1.370, 0.020},
	"CD-NE2":  {1.340, 0.020},
	"CG-ND2":  {1.330, 0.020},
	"CD-NE":   {1.460, 0.020},
	"NE-CZ":   {1.330, 0.020},
	"CZ-NH1":  {1.330, 0.020},
	"CZ-NH2":  {1.330, 0.020},
	"CE-NZ":   {1.490, 0.020},
	"CD-N":    {1.470, 0.020},

	// Ring bonds (aromatic)
	"CG-CD1":  {1.390, 0.020},
	"CD1-CE1": {1.390, 0.020},
	"CE1-CZ":  {1.390, 0.020},
	"CZ-CE2":  {1.390, 0.020},
	"CE2-CD2": {1.390, 0.020},
	"ND1-CE1": {1.370, 0.020},
	"NE2-CD2": {1.370, 0.020},
	"NE2-CE1": {1.370, 0.020},
	"ND1-CG":  {1.370, 0.020},
	"CD2-CG":  {1.390, 0.020},
	"CD1-NE1": {1.370, 0.020},
	"NE1-CE2": {1.370, 0.020},
	"CE2-CZ2": {1.400, 0.020},
	"CZ2-CH2": {1.370, 0.020},
	"CH2-CZ3": {1.370, 0.020},
	"CZ3-CE3": {1.400, 0.020},
	"CE3-CD2": {1.400, 0.020},
	"TRP:CE2-CD2": {1.400, 0.020},
	"TRP:CG-CD2": {1.430, 0.020},
	"CD1-CG":  {1.360, 0.020},
}
