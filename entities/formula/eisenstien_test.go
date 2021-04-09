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
	}

	latticeCoordinate := squareLatticeFormula.ConvertToLatticeCoordinates(complex(1.0,2.0))
	checker.Assert(real(latticeCoordinate), utility.NumericallyCloseEnough{}, 1.0, 1e-6)
	checker.Assert(imag(latticeCoordinate), utility.NumericallyCloseEnough{}, 2.0, 1e-6)
}

func (suite *EisensteinFormulaSuite) TestConvertToLatticeVectorEvenIfOneVectorIsZero(checker *C) {
	squareLatticeFormula := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		XLatticeVector: complex(0, 1),
		YLatticeVector: complex(1, 0),
	}

	latticeCoordinate := squareLatticeFormula.ConvertToLatticeCoordinates(complex(1.0,2.0))
	checker.Assert(real(latticeCoordinate), utility.NumericallyCloseEnough{}, 2.0, 1e-6)
	checker.Assert(imag(latticeCoordinate), utility.NumericallyCloseEnough{}, 1.0, 1e-6)
}
