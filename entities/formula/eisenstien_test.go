package formula_test

import (
	. "gopkg.in/check.v1"
	"wallpaper/entities/formula"
	"wallpaper/entities/utility"
)

type EisensteinFormulaSuite struct {
}

var _ = Suite(&EisensteinFormulaSuite{})

func (suite *EisensteinFormulaSuite) SetUpTest(checker *C) {
}

func (suite *EisensteinFormulaSuite) TestVectorCannotBeZero(checker *C) {
	badLatticeFormula := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		XLatticeVector: complex(0, 0),
		YLatticeVector: complex(0, 1),
		Multiplier: complex(1, 0),
	}
	err := badLatticeFormula.Validate()
	checker.Assert(err, ErrorMatches, "lattice vectors cannot be \\(0,0\\)")
}

func (suite *EisensteinFormulaSuite) TestVectorsCannotBeCollinear(checker *C) {
	badLatticeFormula := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		XLatticeVector: complex(1, 1),
		YLatticeVector: complex(-2, -2),
		Multiplier: complex(1, 0),
	}
	err := badLatticeFormula.Validate()
	checker.Assert(err, ErrorMatches, "vectors cannot be collinear: (.*,.*) and (.*,.*)")
}

func (suite *EisensteinFormulaSuite) TestGoodLatticeVectorsAreValid(checker *C) {
	squareLatticeFormula := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		XLatticeVector: complex(1, 0),
		YLatticeVector: complex(0, 1),
		Multiplier: complex(1, 0),
	}
	err := squareLatticeFormula.Validate()
	checker.Assert(err, IsNil)
}

func (suite *EisensteinFormulaSuite) TestConvertToLatticeVector(checker *C) {
	squareLatticeFormula := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		XLatticeVector: complex(1, 0),
		YLatticeVector: complex(0, 1),
		Multiplier: complex(1, 0),
	}

	latticeCoordinate := squareLatticeFormula.ConvertToLatticeCoordinates(complex(1.0,2.0))
	checker.Assert(real(latticeCoordinate), utility.NumericallyCloseEnough{}, 1.0, 1e-6)
	checker.Assert(imag(latticeCoordinate), utility.NumericallyCloseEnough{}, 2.0, 1e-6)
}

func (suite *EisensteinFormulaSuite) TestConvertToLatticeVectorEvenIfFirstVectorHasZeroRealComponent(checker *C) {
	squareLatticeFormulaWithFlippedVectors := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		XLatticeVector: complex(0, 1),
		YLatticeVector: complex(1, 0),
		Multiplier: complex(1, 0),
	}

	latticeCoordinate := squareLatticeFormulaWithFlippedVectors.ConvertToLatticeCoordinates(complex(1.0,2.0))
	checker.Assert(real(latticeCoordinate), utility.NumericallyCloseEnough{}, 2.0, 1e-6)
	checker.Assert(imag(latticeCoordinate), utility.NumericallyCloseEnough{}, 1.0, 1e-6)
}

func (suite *EisensteinFormulaSuite) TestCalculateEisensteinTermForGivenPoint(checker *C) {
	squareLatticeFormula := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		XLatticeVector: complex(1, 0),
		YLatticeVector: complex(0, 1),
		Multiplier: complex(1, 0),
	}

	latticeCoordinate := squareLatticeFormula.Calculate(complex(1.0,1.5))
	checker.Assert(real(latticeCoordinate), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(latticeCoordinate), utility.NumericallyCloseEnough{}, 0.0, 1e-6)
}

func (suite *EisensteinFormulaSuite) TestCreateFormulaWithJSON(checker *C) {
	jsonByteStream := []byte(`{
				"x_lattice_vector": {
					"real": -1.0,
					"imaginary": 2e-2
				},
				"y_lattice_vector": {
					"real": 100,
					"imaginary": -9000
				},
				"power_n": 12,
				"power_m": -10
			}`)
	term, err := formula.NewEisensteinFormulaTermFromJSON(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(real(term.XLatticeVector), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(term.XLatticeVector), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(real(term.YLatticeVector), utility.NumericallyCloseEnough{}, 100, 1e-6)
	checker.Assert(imag(term.YLatticeVector), utility.NumericallyCloseEnough{}, -9000, 1e-6)
	checker.Assert(term.PowerN, Equals, 12)
	checker.Assert(term.PowerM, Equals, -10)
}

func (suite *EisensteinFormulaSuite) TestCreateFormulaWithYAML(checker *C) {
	yamlByteStream := []byte(`
x_lattice_vector:
  real: -1.0
  imaginary: 2e-2
y_lattice_vector:
  real: 100
  imaginary: -9000
power_n: 12
power_m: -10
`)
	term, err := formula.NewEisensteinFormulaTermFromYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(real(term.XLatticeVector), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(term.XLatticeVector), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(real(term.YLatticeVector), utility.NumericallyCloseEnough{}, 100, 1e-6)
	checker.Assert(imag(term.YLatticeVector), utility.NumericallyCloseEnough{}, -9000, 1e-6)
	checker.Assert(term.PowerN, Equals, 12)
	checker.Assert(term.PowerM, Equals, -10)
}

func (suite *EisensteinFormulaSuite) TestCalculateEisensteinWithMultiplier(checker *C) {
	squareLatticeFormula := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		XLatticeVector: complex(1, 0),
		YLatticeVector: complex(0, 1),
		Multiplier: complex(2, 0),
	}

	latticeCoordinate := squareLatticeFormula.Calculate(complex(1.0,1.5))
	checker.Assert(real(latticeCoordinate), utility.NumericallyCloseEnough{}, -2.0, 1e-6)
	checker.Assert(imag(latticeCoordinate), utility.NumericallyCloseEnough{}, 0.0, 1e-6)
}