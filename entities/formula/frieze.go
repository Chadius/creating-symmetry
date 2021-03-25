package formula

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"math/cmplx"
	"wallpaper/entities/utility"
)

// FriezeFormula is used to generate frieze patterns.
type FriezeFormula struct {
	Terms []*EulerFormulaTerm
}

// Calculate applies the Frieze formula to the complex number z.
func (formula FriezeFormula) Calculate(z complex128) *CalculationResultForFormula {
	result := &CalculationResultForFormula{
		Total: complex(0,0),
		ContributionByTerm: []complex128{},
	}

	for _, term := range formula.Terms {
		termResult := term.Calculate(z)
		result.Total += termResult
		result.ContributionByTerm = append(result.ContributionByTerm, termResult)
	}

	return result
}

// FriezeSymmetry notes the kinds of symmetries the formula contains.
type FriezeSymmetry struct {
	P111 bool
	P11m bool
	P211 bool
	P1m1 bool
	P11g bool
	P2mm bool
	P2mg bool
}

//AnalyzeForSymmetry scans the formula and returns a list of symmetries.
func (formula FriezeFormula) AnalyzeForSymmetry() *FriezeSymmetry {
	symmetriesFound := &FriezeSymmetry{
		P111: true,
		P11m: true,
		P211: true,
		P1m1: true,
		P11g: true,
		P2mm: true,
		P2mg: true,
	}
	for _, term := range formula.Terms {
		if term.IgnoreComplexConjugate {
			symmetriesFound.P211 = false
			symmetriesFound.P1m1 = false
			symmetriesFound.P11g = false
			symmetriesFound.P11m = false
			symmetriesFound.P2mm = false
			symmetriesFound.P2mg = false
		}

		powerSumIsEven := (term.PowerN + term.PowerM) % 2 == 0

		containsMinusNMinusM := coefficientPairsIncludes(term.CoefficientPairs.OtherCoefficientRelationships, MinusNMinusM)
		containsMinusMMinusN := coefficientPairsIncludes(term.CoefficientPairs.OtherCoefficientRelationships, MinusMMinusN)
		containsPlusMPlusN := coefficientPairsIncludes(term.CoefficientPairs.OtherCoefficientRelationships, PlusMPlusN)

		containsMinusMMinusNAndPowerSumIsOdd := coefficientPairsIncludes(term.CoefficientPairs.OtherCoefficientRelationships, MinusMMinusNMaybeFlipScale ) && !powerSumIsEven
		containsPlusMPlusNAndPowerSumIsOdd := coefficientPairsIncludes(term.CoefficientPairs.OtherCoefficientRelationships, PlusMPlusNMaybeFlipScale) && !powerSumIsEven

		containsMinusMMinusNAndPowerSumIsEven := coefficientPairsIncludes(term.CoefficientPairs.OtherCoefficientRelationships, MinusMMinusNMaybeFlipScale ) && powerSumIsEven
		containsPlusMPlusNAndPowerSumIsEven := coefficientPairsIncludes(term.CoefficientPairs.OtherCoefficientRelationships, PlusMPlusNMaybeFlipScale) && powerSumIsEven

		if !containsMinusNMinusM {
			symmetriesFound.P211 = false
		}
		if !containsPlusMPlusN {
			symmetriesFound.P1m1 = false
		}
		if !containsMinusMMinusNAndPowerSumIsOdd {
			symmetriesFound.P11g = false
		}
		if !(containsMinusMMinusN || containsMinusMMinusNAndPowerSumIsEven) {
			symmetriesFound.P11m = false
		}
		if !(
			containsMinusNMinusM &&
				(containsPlusMPlusN || containsPlusMPlusNAndPowerSumIsEven) &&
				(containsMinusMMinusN || containsMinusMMinusNAndPowerSumIsEven)) {
			symmetriesFound.P2mm = false
		}
		if !(containsMinusNMinusM && containsPlusMPlusNAndPowerSumIsOdd && containsMinusMMinusNAndPowerSumIsOdd) {
			symmetriesFound.P2mg = false
		}
	}

	return symmetriesFound
}

type eulerFormulaTermMarshalable struct {
	Multiplier				utility.ComplexNumberForMarshal	`json:"multiplier" yaml:"multiplier"`
	PowerN					int								`json:"power_n" yaml:"power_n"`
	PowerM					int								`json:"power_m" yaml:"power_m"`
	IgnoreComplexConjugate	bool							`json:"ignore_complex_conjugate" yaml:"ignore_complex_conjugate"`
	CoefficientPairs		LockedCoefficientPair			`json:"coefficient_pairs" yaml:"coefficient_pairs"`
}

// EulerFormulaTerm calculates e^(i*n*z) * e^(-i*m*zConj)
type EulerFormulaTerm struct {
	Multiplier				complex128
	PowerN					int
	PowerM					int
	// IgnoreComplexConjugate will make sure zConjugate is not used in this calculation
	//    (effectively setting it to 1 + 0i)
	IgnoreComplexConjugate	bool
	// CoefficientPairs will create similar terms to add to this one when calculating.
	//    This is useful when trying to force symmetry by adding another term with swapped
	//    PowerN & PowerM, or multiplying by -1.
	CoefficientPairs		LockedCoefficientPair
}

// NewEulerFormulaTermFromYAML reads the data and returns a formula term from it.
func NewEulerFormulaTermFromYAML(data []byte) (*EulerFormulaTerm, error) {
	return newEulerFormulaTermFromDatastream(data, yaml.Unmarshal)
}

// NewEulerFormulaTermFromJSON reads the data and returns a formula term from it.
func NewEulerFormulaTermFromJSON(data []byte) (*EulerFormulaTerm, error) {
	return newEulerFormulaTermFromDatastream(data, json.Unmarshal)
}

//newEulerFormulaTermFromDatastream consumes a given bytestream and tries to create a new object from it.
func newEulerFormulaTermFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*EulerFormulaTerm, error) {
	var unmarshalError error
	var formulaTermMarshal eulerFormulaTermMarshalable
	unmarshalError = unmarshal(data, &formulaTermMarshal)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	formulaTerm := newEulerFormulaTermFromMarshalObject(formulaTermMarshal)
	return formulaTerm, nil
}

func newEulerFormulaTermFromMarshalObject(marshalObject eulerFormulaTermMarshalable) *EulerFormulaTerm {
	return &EulerFormulaTerm{
		Multiplier:             complex(marshalObject.Multiplier.Real, marshalObject.Multiplier.Imaginary),
		PowerN:                 marshalObject.PowerN,
		PowerM:                 marshalObject.PowerM,
		IgnoreComplexConjugate: marshalObject.IgnoreComplexConjugate,
		CoefficientPairs:       marshalObject.CoefficientPairs,
	}
}

// Calculate returns the result of using the formula on the given complex number.
func (term EulerFormulaTerm) Calculate(z complex128) complex128 {
	sum := CalculateEulerTerm(z, term.PowerN, term.PowerM, term.Multiplier, term.IgnoreComplexConjugate)

	for _, relationship := range term.CoefficientPairs.OtherCoefficientRelationships {
		power1, power2, scale := SetCoefficientsBasedOnRelationship(term.PowerN, term.PowerM, term.Multiplier, relationship)
		relationshipScale := scale * complex(term.CoefficientPairs.Multiplier, 0)
		sum += CalculateEulerTerm(z, power1, power2, relationshipScale, term.IgnoreComplexConjugate)
	}

	return sum
}

// CalculateEulerTerm calculates e^(i*n*z) * e^(-i*m*zConj)
func CalculateEulerTerm(z complex128, power1, power2 int, scale complex128, ignoreComplexConjugate bool) complex128 {
	eRaisedToTheNZi := cmplx.Exp(complex(0,1) * z * complex(float64(power1), 0))
	if ignoreComplexConjugate {
		return eRaisedToTheNZi * scale
	}

	complexConjugate := complex(real(z), -1 * imag(z))
	eRaisedToTheNegativeMZConji := cmplx.Exp(complexConjugate * complex(0, -1 * float64(power2)))
	return eRaisedToTheNZi * eRaisedToTheNegativeMZConji * scale
}

func coefficientPairsIncludes (relationships []CoefficientRelationship, relationshipToFind CoefficientRelationship) bool {
	for _, relationship := range relationships {
		if relationship == relationshipToFind {
			return true
		}
	}
	return false
}

// NewFriezeFormulaFromYAML reads the data and returns a Frieze formula from it.
func NewFriezeFormulaFromYAML(data []byte) (*FriezeFormula, error) {
	return newFriezeFormulaFromDatastream(data, yaml.Unmarshal)
}

// NewFriezeFormulaFromJSON reads the data and returns a Frieze formula from it.
func NewFriezeFormulaFromJSON(data []byte) (*FriezeFormula, error) {
	return newFriezeFormulaFromDatastream(data, yaml.Unmarshal)
}

// FriezeFormulaMarshalable can be marshaled and can be converted into a FriezeFormula.
type FriezeFormulaMarshalable struct {
	Terms []*eulerFormulaTermMarshalable
}

// newFriezeFormulaFromDatastream consumes a given bytestream and tries to create a new object from it.
func newFriezeFormulaFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*FriezeFormula, error) {
	var unmarshalError error
	var friezeFormulaMarshal FriezeFormulaMarshalable
	unmarshalError = unmarshal(data, &friezeFormulaMarshal)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	friezeFormula := NewFriezeFormulaFromMarshalObject(friezeFormulaMarshal)
	return friezeFormula, nil
}

// NewFriezeFormulaFromMarshalObject converts the marshaled object into a FriezeFormula.
func NewFriezeFormulaFromMarshalObject(marshalObject FriezeFormulaMarshalable) *FriezeFormula {
	terms := []*EulerFormulaTerm{}
	for _, termMarshal := range marshalObject.Terms {
		newTerm := newEulerFormulaTermFromMarshalObject(*termMarshal)
		terms = append(terms, newTerm)
	}
	return &FriezeFormula{Terms: terms}
}
