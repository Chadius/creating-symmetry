package formula

import (
	"errors"
	"fmt"
	"math"
)

// EisensteinFormulaTerm defines the shape of a lattice, a 2D structure that remains consistent
//    in wallpaper symmetry.
type EisensteinFormulaTerm struct {
	XLatticeVector			complex128
	YLatticeVector			complex128
	PowerN					int
	PowerM					int
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
func(term EisensteinFormulaTerm)Calculate(z complex128) complex128 {
	return complex(0,0)
}