package formula

import (
	"gopkg.in/yaml.v2"
	"math/cmplx"
	"wallpaper/entities/utility"
)

// CoefficientRelationship relates how a pair of coordinates should be applied.
type CoefficientRelationship string

// CoefficientRelationship s determine the order and sign of powers n and m.
//   Plus means *1, Minus means *-1
//   If N appears first the powers then power N is applied to the number and power M to the complex conjugate.
//   If M appears first the powers then power M is applied to the number and power N to the complex conjugate.
//	 MaybeFlipScale will multiply the scale by -1 if N + M is odd.
const (
	PlusNPlusM CoefficientRelationship = "+N+M"
	PlusMPlusN                         = "+M+N"
	MinusNMinusM                       = "-N-M"
	MinusMMinusN                       = "-M-N"
	PlusMPlusNMaybeFlipScale           = "+M+NF"
	MinusMMinusNMaybeFlipScale         = "-M-NF"
)

// SetCoefficientsBasedOnRelationship will rearrange powerN and powerM according to their relationship.
func SetCoefficientsBasedOnRelationship(powerN, powerM int, scale complex128, relationship CoefficientRelationship) (int, int, complex128) {
	var power1, power2 int
	switch relationship {
	case PlusNPlusM:
		power1 = powerN
		power2 = powerM
	case PlusMPlusN, PlusMPlusNMaybeFlipScale:
		power1 = powerM
		power2 = powerN
	case MinusMMinusN, MinusMMinusNMaybeFlipScale:
		power1 = -1 * powerM
		power2 = -1 * powerN
	case MinusNMinusM:
		power1 = -1 * powerN
		power2 = -1 * powerM
	}

	sumOfPowersIsOdd := (powerN + powerM) % 2 == 1
	relationshipMayFlipScale := relationship == PlusMPlusNMaybeFlipScale || relationship == MinusMMinusNMaybeFlipScale
	if sumOfPowersIsOdd && relationshipMayFlipScale {
		scale *= -1
	}

	return power1, power2, scale
}

// ZExponentialFormulaTerm describes a formula of the form Scale * z^PowerN * zConjugate^PowerM.
type ZExponentialFormulaTerm struct {
	Scale                  complex128
	PowerN                 int
	PowerM                 int
	// IgnoreComplexConjugate will make sure zConjugate is not used in this calculation
	//    (effectively setting it to 1 + 0i)
	IgnoreComplexConjugate bool
	// CoefficientPairs will create similar terms to add to this one when calculating.
	//    This is useful when trying to force symmetry by adding another term with swapped
	//    PowerN & PowerM, or multiplying by -1.
	CoefficientPairs LockedCoefficientPair
}

// Calculate returns the result of using the formula on the given complex number.
func (term ZExponentialFormulaTerm) Calculate(z complex128) complex128 {
	sum := CalculateExponentTerm(z, term.PowerN, term.PowerM, term.Scale, term.IgnoreComplexConjugate)

	for _, relationship := range term.CoefficientPairs.OtherCoefficientRelationships {
		power1, power2, scale := SetCoefficientsBasedOnRelationship(term.PowerN, term.PowerM, term.Scale, relationship)
		relationshipScale := scale * complex(term.CoefficientPairs.Multiplier, 0)

		sum += CalculateExponentTerm(z, power1, power2, relationshipScale, term.IgnoreComplexConjugate)
	}
	return sum
}

// LockedCoefficientPair describes how to create a new Term based on the current one.
type LockedCoefficientPair struct {
	Multiplier                    float64					`json:"multiplier" yaml:"multiplier"`
	OtherCoefficientRelationships []CoefficientRelationship	`json:"relationships" yaml:"relationships"`
}

// NewLockedCoefficientPairFromYAML reads the data and returns a LockedCoefficientPair.
func NewLockedCoefficientPairFromYAML(data []byte) (*LockedCoefficientPair, error) {
	return newLockedCoefficientPairFromDatastream(data, yaml.Unmarshal)
}

// newLockedCoefficientPairFromDatastream consumes a given bytestream and tries to create a new object from it.
func newLockedCoefficientPairFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*LockedCoefficientPair, error) {
	var unmarshalError error
	var lockedCoefficientPairToCreate LockedCoefficientPair
	unmarshalError = unmarshal(data, &lockedCoefficientPairToCreate)

	if unmarshalError != nil {
		return nil, unmarshalError
	}
	return &lockedCoefficientPairToCreate, nil
}

// CalculateExponentTerm calculates (z^power * zConj^conjugatePower)
//   where z is a complex number, zConj is the complex conjugate
//   and power and conjugatePower are integers.
func CalculateExponentTerm(z complex128, power1, power2 int, scale complex128, ignoreComplexConjugate bool) complex128 {
	zRaisedToN := cmplx.Pow(z, complex(float64(power1), 0))
	if ignoreComplexConjugate {
		return zRaisedToN * scale
	}

	complexConjugate := complex(real(z), -1 * imag(z))
	complexConjugateRaisedToM := cmplx.Pow(complexConjugate, complex(float64(power2), 0))
	return zRaisedToN * complexConjugateRaisedToM * scale
}

// RosetteFormula uses a collection of z^m terms to calculate results.
//    This transforms the input into a circular pattern rotating around the
//    origin.
type RosetteFormula struct {
	Terms []*ZExponentialFormulaTerm
}

// CalculationResultForFormula shows the results of a calculation
type CalculationResultForFormula struct {
	Total                 complex128
	ContributionByTerm []complex128
}

// Calculate applies the Rosette formula to the complex number z.
func (r RosetteFormula) Calculate(z complex128) *CalculationResultForFormula {
	result := &CalculationResultForFormula{
		Total: complex(0,0),
		ContributionByTerm: []complex128{},
	}

	for _, term := range r.Terms {
		termResult := term.Calculate(z)
		result.Total += termResult
		result.ContributionByTerm = append(result.ContributionByTerm, termResult)
	}

	return result
}

// RosetteSymmetry notes the kinds of symmetries the rosette formula contains.
type RosetteSymmetry struct {
	Multifold int
}

// AnalyzeForSymmetry analyzes the formula for symmetries.
func (r RosetteFormula) AnalyzeForSymmetry() *RosetteSymmetry {
	symmetriesFound := &RosetteSymmetry{
		Multifold: 1,
	}

	r.calculateMultifoldSymmetry(symmetriesFound)
	return symmetriesFound
}

func (r RosetteFormula) calculateMultifoldSymmetry(symmetriesFound *RosetteSymmetry) {
	termPowerDifferences := []int{}

	for _, term := range r.Terms {
		powerDifference := term.PowerN - term.PowerM
		if powerDifference < 0 {
			powerDifference *= -1
		}
		termPowerDifferences = append(termPowerDifferences, powerDifference)
	}

	if len(termPowerDifferences) == 1 {
		symmetriesFound.Multifold = termPowerDifferences[0]
	} else if len(termPowerDifferences) > 1 {
		var currentGreatestCommonDenominator int
		for index := range termPowerDifferences {
			if index >= len(termPowerDifferences) - 1 {
				break
			}
			currentGreatestCommonDenominator = getGreatestCommonDenominator(
				termPowerDifferences[index],
				termPowerDifferences[index + 1])
		}
		symmetriesFound.Multifold = currentGreatestCommonDenominator
	}
}

// getGreatestCommonDenominator finds the largest integer that divides into
//   integers a and b, leaving 0 behind.
func getGreatestCommonDenominator(a, b int) int {
	if a == b {
		return a
	}

	larger := a
	smaller := b

	remainder := larger % smaller
	if remainder == 0 {
		return smaller
	}
	return getGreatestCommonDenominator(smaller, remainder)
}

// EulerFormulaTerm calculates e^(i*n*z) * e^(-i*m*zConj)
type EulerFormulaTerm struct {
	Scale                  complex128
	PowerN                 int
	PowerM                 int
	// IgnoreComplexConjugate will make sure zConjugate is not used in this calculation
	//    (effectively setting it to 1 + 0i)
	IgnoreComplexConjugate bool
	// CoefficientPairs will create similar terms to add to this one when calculating.
	//    This is useful when trying to force symmetry by adding another term with swapped
	//    PowerN & PowerM, or multiplying by -1.
	CoefficientPairs LockedCoefficientPair
}

// Calculate returns the result of using the formula on the given complex number.
func (term EulerFormulaTerm) Calculate(z complex128) complex128 {
	sum := CalculateEulerTerm(z, term.PowerN, term.PowerM, term.Scale, term.IgnoreComplexConjugate)

	for _, relationship := range term.CoefficientPairs.OtherCoefficientRelationships {
		power1, power2, scale := SetCoefficientsBasedOnRelationship(term.PowerN, term.PowerM, term.Scale, relationship)
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

func coefficientPairsIncludes (relationships []CoefficientRelationship, relationshipToFind CoefficientRelationship) bool {
	for _, relationship := range relationships {
		if relationship == relationshipToFind {
			return true
		}
	}
	return false
}