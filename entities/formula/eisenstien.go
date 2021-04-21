package formula

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"math"
	"math/cmplx"
	"wallpaper/entities/utility"
)

// EisensteinFormulaTermMarshalable can be marshaled and converted to a EisensteinFormulaTerm
type EisensteinFormulaTermMarshalable struct {
	XLatticeVector			utility.ComplexNumberForMarshal	`json:"x_lattice_vector" yaml:"x_lattice_vector"`
	YLatticeVector			utility.ComplexNumberForMarshal	`json:"y_lattice_vector" yaml:"y_lattice_vector"`
	PowerN					int								`json:"power_n" yaml:"power_n"`
	PowerM					int								`json:"power_m" yaml:"power_m"`
}

// EisensteinFormulaTerm defines the shape of a lattice, a 2D structure that remains consistent
//    in wallpaper symmetry.
type EisensteinFormulaTerm struct {
	XLatticeVector			complex128
	YLatticeVector			complex128
	PowerN					int
	PowerM					int
	Multiplier				complex128
}

func vectorIsZero(vector complex128) bool {
	return real(vector) == 0 && imag(vector) == 0
}

// vectorsAreCollinear returns true if both vectors are perfectly lined up
func vectorsAreCollinear(vector1 complex128, vector2 complex128) bool {
	absoluteValueDotProduct := math.Abs((real(vector1) * real(vector2)) + (imag(vector1) * imag(vector2)))
	lengthOfVector1 := math.Sqrt((real(vector1) * real(vector1)) + (imag(vector1) * imag(vector1)))
	lengthOfVector2 := math.Sqrt((real(vector2) * real(vector2)) + (imag(vector2) * imag(vector2)))

	tolerance := 1e-8
	return math.Abs(absoluteValueDotProduct - lengthOfVector1 * lengthOfVector2) < tolerance
}

// Validate returns an error if this is an invalid formula.
func(term EisensteinFormulaTerm)Validate() error {
	if vectorIsZero(term.XLatticeVector) || vectorIsZero(term.YLatticeVector) {
		return errors.New(`lattice vectors cannot be (0,0)`)
	}
	if vectorsAreCollinear(term.XLatticeVector, term.YLatticeVector) {
		return fmt.Errorf(
			`vectors cannot be collinear: (%f,%f) and \(%f,%f)`,
			real(term.XLatticeVector),
			imag(term.XLatticeVector),
			real(term.YLatticeVector),
			imag(term.YLatticeVector),
		)
	}
	return nil
}

// Calculate uses the Eisenstein formula on the complex number z.
// Calculate(z) = e ^ (2 PI i * (nX + mY))
//  where n amd m are PowerN and PowerM,
//  and X and Y are the real and imag parts of (z converted into LatticeCoordinates)
func(term EisensteinFormulaTerm)Calculate(z complex128) complex128 {
	zInLatticeCoordinates := term.ConvertToLatticeCoordinates(z)
	powerMultiplier := (float64(term.PowerN) * real(zInLatticeCoordinates)) +
		(float64(term.PowerM) * imag(zInLatticeCoordinates))
	return cmplx.Exp(complex(0, 2.0 * math.Pi * powerMultiplier))
}

// ConvertToLatticeCoordinates converts a point from cartesian coordinates to the lattice coordinates
func (term EisensteinFormulaTerm) ConvertToLatticeCoordinates(cartesianPoint complex128) complex128 {

	vector1 := term.XLatticeVector
	vector2 := term.YLatticeVector
	swapVectorsDuringCalculation := real(vector1) < 1e-6

	if swapVectorsDuringCalculation == true {
		vector1 = term.YLatticeVector
		vector2 = term.XLatticeVector
	}

	scalarForVector2 := (imag(cartesianPoint) - (real(cartesianPoint) * imag(vector1))) /
		((real(vector1) * imag(vector2)) - (imag(vector1) * real(vector2)))

	scalarForVector1 := (real(cartesianPoint) - (scalarForVector2 * real(vector2)))/ real(vector1)

	if swapVectorsDuringCalculation {
		return complex(scalarForVector2, scalarForVector1)
	}

	return complex(scalarForVector1, scalarForVector2)
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
	var formulaTermMarshal EisensteinFormulaTermMarshalable
	unmarshalError = unmarshal(data, &formulaTermMarshal)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	formulaTerm := newEisensteinFormulaTermFromMarshalObject(formulaTermMarshal)
	return formulaTerm, nil
}

func newEisensteinFormulaTermFromMarshalObject(marshalObject EisensteinFormulaTermMarshalable) *EisensteinFormulaTerm {
	return &EisensteinFormulaTerm{
		PowerN:                 marshalObject.PowerN,
		PowerM:                 marshalObject.PowerM,
		XLatticeVector:			complex(marshalObject.XLatticeVector.Real, marshalObject.XLatticeVector.Imaginary),
		YLatticeVector:			complex(marshalObject.YLatticeVector.Real, marshalObject.YLatticeVector.Imaginary),
	}
}
