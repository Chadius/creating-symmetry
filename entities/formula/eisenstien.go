package formula

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"math"
	"math/cmplx"
	"wallpaper/entities/formula/coefficient"
	"wallpaper/entities/utility"
)

// EisensteinFormulaTermMarshal can be marshaled and converted to a EisensteinFormulaTerm
type EisensteinFormulaTermMarshal struct {
	PowerN					int								`json:"power_n" yaml:"power_n"`
	PowerM					int								`json:"power_m" yaml:"power_m"`
	Multiplier utility.ComplexNumberForMarshal	`json:"multiplier" yaml:"multiplier"`
}

// EisensteinFormulaTerm defines the shape of a lattice, a 2D structure that remains consistent
//    in wallpaper symmetry.
type EisensteinFormulaTerm struct {
	PowerN					int
	PowerM					int
	Multiplier 		complex128
}

// PowerSumIsEven returns true if the sum of the term powers is divisible by 2.
func (term EisensteinFormulaTerm)PowerSumIsEven() bool {
	return (term.PowerM + term.PowerN) % 2 == 0
}

// Calculate uses the Eisenstein formula on the complex number z.
// Calculate(z) = e ^ (2 PI i * (nX + mY))
//  where n amd m are PowerN and PowerM,
//  and X and Y are the real and imag parts of (zInLatticeCoordinates)
func(term EisensteinFormulaTerm)Calculate(zInLatticeCoordinates complex128) complex128 {
	powerMultiplier := (float64(term.PowerN) * real(zInLatticeCoordinates)) +
		(float64(term.PowerM) * imag(zInLatticeCoordinates))
	expo := cmplx.Exp(complex(0, 2.0 * math.Pi * powerMultiplier))
	return expo * term.Multiplier
}

// NewEisensteinFormulaTermFromJSON reads the data and returns a formula term from it.
func NewEisensteinFormulaTermFromJSON(data []byte) (*EisensteinFormulaTerm, error) {
	return newEisensteinFormulaTermFromDatastream(data, json.Unmarshal)
}

// NewEisensteinFormulaTermFromYAML reads the data and returns a formula term from it.
func NewEisensteinFormulaTermFromYAML(data []byte) (*EisensteinFormulaTerm, error) {
	return newEisensteinFormulaTermFromDatastream(data, yaml.Unmarshal)
}

//newEisensteinFormulaTermFromDatastream consumes a given bytestream and tries to create a new object from it.
func newEisensteinFormulaTermFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*EisensteinFormulaTerm, error) {
	var unmarshalError error
	var formulaTermMarshal EisensteinFormulaTermMarshal
	unmarshalError = unmarshal(data, &formulaTermMarshal)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	formulaTerm := NewEisensteinFormulaTermFromMarshalObject(formulaTermMarshal)
	return formulaTerm, nil
}

// NewEisensteinFormulaTermFromMarshalObject converts the marshaled intermediary object into a usable object.
func NewEisensteinFormulaTermFromMarshalObject(marshalObject EisensteinFormulaTermMarshal) *EisensteinFormulaTerm {
	return &EisensteinFormulaTerm{
		PowerN:                 marshalObject.PowerN,
		PowerM:                 marshalObject.PowerM,
		Multiplier: complex(marshalObject.Multiplier.Real, marshalObject.Multiplier.Imaginary),
	}
}

// GetAllPossibleTermRelationships returns a list of relationships that all of the terms conform to.
func GetAllPossibleTermRelationships(term1, term2 *EisensteinFormulaTerm) []coefficient.Relationship {
	if term1 == nil || term2 == nil {
		return []coefficient.Relationship{}
	}

	foundRelationships := []coefficient.Relationship{}

	relationshipsToTest := []coefficient.Relationship{
		coefficient.PlusNPlusM,
		coefficient.PlusMPlusN,
		coefficient.MinusNMinusM,
		coefficient.MinusMMinusN,
		coefficient.PlusMPlusNMaybeFlipScale,
		coefficient.MinusMMinusNMaybeFlipScale,
		coefficient.PlusMMinusSumNAndM,
		coefficient.MinusSumNAndMPlusN,
		coefficient.PlusMMinusN,
		coefficient.MinusMPlusN,
		coefficient.PlusNMinusM,
		coefficient.MinusNPlusM,
		coefficient.PlusNMinusMNegateMultiplierIfOddPowerN,
		coefficient.MinusNPlusMNegateMultiplierIfOddPowerN,
		coefficient.PlusNMinusMNegateMultiplierIfOddPowerSum,
		coefficient.MinusNPlusMNegateMultiplierIfOddPowerSum,
	}

	for _, relationshipToTest := range relationshipsToTest {
		if SatisfiesRelationship(term1, term2, relationshipToTest) {
			foundRelationships = append(foundRelationships, relationshipToTest)
		}
	}

	return foundRelationships
}

// SatisfiesRelationship sees if the all of terms match the coefficient relationship.
func SatisfiesRelationship(term1, term2 *EisensteinFormulaTerm, relationship coefficient.Relationship) bool {
	if !termMultipliersAreTheSame(term1, term2) && !termMultipliersAreNegated(term1, term2) {
		return false
	}

	type TwoTermChecker func(term1, term2 *EisensteinFormulaTerm) bool
	relationshipCheckerByRelationship := map[coefficient.Relationship]TwoTermChecker{
		coefficient.PlusMPlusN: satisfiesRelationshipPlusMPlusN,
		coefficient.MinusNMinusM:                           satisfiesRelationshipMinusNMinusM,
		coefficient.PlusNPlusM:                             satisfiesRelationshipPlusNPlusM,
		coefficient.MinusMMinusN:                           satisfiesRelationshipMinusMMinusN,
		coefficient.PlusMPlusNMaybeFlipScale:               satisfiesRelationshipPlusMPlusNMaybeFlipScale,
		coefficient.MinusMMinusNMaybeFlipScale:             satisfiesRelationshipMinusMMinusNMaybeFlipScale,
		coefficient.PlusMMinusSumNAndM:                     satisfiesRelationshipPlusMMinusSumNAndM,
		coefficient.MinusSumNAndMPlusN:                     satisfiesRelationshipMinusSumNAndMPlusN,
		coefficient.PlusMMinusN:                            satisfiesRelationshipPlusMMinusN,
		coefficient.MinusMPlusN:                            satisfiesRelationshipMinusMPlusN,
		coefficient.PlusNMinusM:                            satisfiesRelationshipPlusNMinusM,
		coefficient.MinusNPlusM:                            satisfiesRelationshipMinusNPlusM,
		coefficient.PlusNMinusMNegateMultiplierIfOddPowerN: satisfiesRelationshipPlusNMinusMNegateMultiplierIfOddPowerN,
		coefficient.PlusNMinusMNegateMultiplierIfOddPowerSum: satisfiesRelationshipPlusNMinusMNegateMultiplierIfOddPowerSum,
		coefficient.MinusNPlusMNegateMultiplierIfOddPowerN: satisfiesRelationshipMinusNPlusMNegateMultiplierIfOddPowerN,
		coefficient.MinusNPlusMNegateMultiplierIfOddPowerSum: satisfiesRelationshipMinusNPlusMNegateMultiplierIfOddPowerSum,
	}
	relationshipChecker := relationshipCheckerByRelationship[relationship]
	return relationshipChecker(term1, term2)
}

func termMultipliersAreTheSame(term1, term2 *EisensteinFormulaTerm) bool {
	return real(term1.Multiplier) == real(term2.Multiplier) && imag(term1.Multiplier) == imag(term2.Multiplier)
}

func termMultipliersAreNegated(term1, term2 *EisensteinFormulaTerm) bool {
	return real(term1.Multiplier) == -1 * real(term2.Multiplier) && imag(term1.Multiplier) == -1 * imag(term2.Multiplier)
}

func satisfiesRelationshipPlusMPlusN(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term1.PowerN == term2.PowerM && term1.PowerM == term2.PowerN
}

func satisfiesRelationshipMinusNMinusM(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term1.PowerN == -1 * term2.PowerN && term1.PowerM == -1 * term2.PowerM
}

func satisfiesRelationshipPlusNPlusM(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term1.PowerN == term2.PowerN && term1.PowerM == term2.PowerM
}

func satisfiesRelationshipMinusNPlusM(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term2.PowerN == -1 * term1.PowerN && term2.PowerM == term1.PowerM
}

func satisfiesRelationshipMinusMMinusN(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term1.PowerN == -1 * term2.PowerM && term1.PowerM == -1 * term2.PowerN
}

func satisfiesRelationshipPlusMPlusNMaybeFlipScale(term1, term2 *EisensteinFormulaTerm) bool {
	if term1.PowerSumIsEven() && !termMultipliersAreTheSame(term1, term2) {
		return false
	}

	if !term1.PowerSumIsEven() && !termMultipliersAreNegated(term1, term2) {
		return false
	}

	return term1.PowerN == term2.PowerM && term1.PowerM == term2.PowerN
}

func satisfiesRelationshipMinusMMinusNMaybeFlipScale(term1, term2 *EisensteinFormulaTerm) bool {
	if term1.PowerSumIsEven() && !termMultipliersAreTheSame(term1, term2) {
		return false
	}

	if !term1.PowerSumIsEven() && !termMultipliersAreNegated(term1, term2) {
		return false
	}

	return term1.PowerN == -1 * term2.PowerM && term1.PowerM == -1 * term2.PowerN
}
func satisfiesRelationshipPlusMMinusSumNAndM(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term2.PowerN == term1.PowerM && term2.PowerM == -1 * (term1.PowerN + term1.PowerM)
}
func satisfiesRelationshipMinusSumNAndMPlusN(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term2.PowerN == -1 * (term1.PowerN + term1.PowerM) && term2.PowerM == term1.PowerN
}
func satisfiesRelationshipPlusMMinusN(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term2.PowerN == term1.PowerM && term2.PowerM == -1 * term1.PowerN
}
func satisfiesRelationshipMinusMPlusN(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term2.PowerN == -1 * term1.PowerM && term2.PowerM == term1.PowerN
}
func satisfiesRelationshipPlusNMinusM(term1, term2 *EisensteinFormulaTerm) bool {
	return termMultipliersAreTheSame(term1, term2) && term2.PowerN == term1.PowerN && term2.PowerM == -1 * term1.PowerM
}
func satisfiesRelationshipPlusNMinusMNegateMultiplierIfOddPowerN(term1, term2 *EisensteinFormulaTerm) bool {
	if term1.PowerN % 2 == 0 && !termMultipliersAreTheSame(term1, term2) {
		return false
	}

	if term1.PowerN % 2 != 0 && termMultipliersAreTheSame(term1, term2) {
		return false
	}

	return term2.PowerN == term1.PowerN && term2.PowerM == -1 * term1.PowerM
}
func satisfiesRelationshipPlusNMinusMNegateMultiplierIfOddPowerSum(term1, term2 *EisensteinFormulaTerm) bool {
	if term1.PowerSumIsEven() && !termMultipliersAreTheSame(term1, term2) {
		return false
	}

	if !term1.PowerSumIsEven() && termMultipliersAreTheSame(term1, term2) {
		return false
	}

	return term2.PowerN == term1.PowerN && term2.PowerM == -1 * term1.PowerM
}

func satisfiesRelationshipMinusNPlusMNegateMultiplierIfOddPowerN(term1, term2 *EisensteinFormulaTerm) bool {
	if term1.PowerN % 2 == 0 && !termMultipliersAreTheSame(term1, term2) {
		return false
	}

	if term1.PowerN % 2 != 0 && termMultipliersAreTheSame(term1, term2) {
		return false
	}

	return term2.PowerN == -1 * term1.PowerN && term2.PowerM == term1.PowerM
}
func satisfiesRelationshipMinusNPlusMNegateMultiplierIfOddPowerSum(term1, term2 *EisensteinFormulaTerm) bool {
	if term1.PowerSumIsEven() && !termMultipliersAreTheSame(term1, term2) {
		return false
	}

	if !term1.PowerSumIsEven() && termMultipliersAreTheSame(term1, term2) {
		return false
	}

	return term2.PowerN == -1 * term1.PowerN && term2.PowerM == term1.PowerM
}
// SelectTermsToSatisfyRelationships tries to find all of the terms that satisfy the given desiredRelationships.
//   The return value will have 1 EisensteinFormulaTerm per desiredRelationship.
//   If there are not enough terms, returns an empty list.
//   desiredRelationships is a list of unique relationships.
func SelectTermsToSatisfyRelationships(baseTerm *EisensteinFormulaTerm,
	otherTerms []*EisensteinFormulaTerm,
	desiredRelationships []coefficient.Relationship) []*EisensteinFormulaTerm {

	if baseTerm == nil {
		return []*EisensteinFormulaTerm{}
	}

	if len(otherTerms) < len(desiredRelationships) {
		return []*EisensteinFormulaTerm{}
	}

	relationshipFound := map[coefficient.Relationship]bool{}
	for _, relationship := range desiredRelationships {
		relationshipFound[relationship] = false
	}

	satisfyingTerms := []*EisensteinFormulaTerm{}

	for _, term := range otherTerms {
		for _, relationship := range desiredRelationships {
			if SatisfiesRelationship(baseTerm, term, relationship) {
				satisfyingTerms = append(satisfyingTerms, term)
				relationshipFound[relationship] = true
				break
			}
		}
	}

	for _, relationship := range desiredRelationships {
		if relationshipFound[relationship] == false {
			return []*EisensteinFormulaTerm{}
		}
	}

	return satisfyingTerms
}