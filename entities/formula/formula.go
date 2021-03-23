package formula

import (
	"gopkg.in/yaml.v2"
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

// CalculationResultForFormula shows the results of a calculation
type CalculationResultForFormula struct {
	Total                 complex128
	ContributionByTerm []complex128
}