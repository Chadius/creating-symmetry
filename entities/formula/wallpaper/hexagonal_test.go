package wallpaper_test

import (
	. "gopkg.in/check.v1"
	"math"
	"math/cmplx"
	"testing"
	"wallpaper/entities/formula/wallpaper"
	"wallpaper/entities/utility"
)

func Test(t *testing.T) { TestingT(t) }

type HexagonalWallpaperSuite struct {
}

var _ = Suite(&HexagonalWallpaperSuite{})

func (suite *HexagonalWallpaperSuite) SetUpTest(checker *C) {
}

func (suite *HexagonalWallpaperSuite) TestHexagonalWallpaperFormulaCalculatesWithAveraging(checker *C) {
	hexagonalFormula := &wallpaper.HexagonalFormula{
		PowerTerms: []*wallpaper.PowerTerm {
			{
				PowerN:         1,
				PowerM:         -2,
				Multiplier:		complex(1,0),
			},
		},
	}
	hexagonalFormula.Setup()

	calculation := hexagonalFormula.Calculate(complex(math.Sqrt(3), -1 * math.Sqrt(3)))
	total := calculation.Total

	expectedAnswer := cmplx.Exp(complex(0, 2 * math.Pi * (3 + math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-2 * math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-3 + math.Sqrt(3))))
	expectedAnswer = expectedAnswer / 3.0

	checker.Assert(real(total), utility.NumericallyCloseEnough{}, real(expectedAnswer), 1e-6)
	checker.Assert(imag(total), utility.NumericallyCloseEnough{}, imag(expectedAnswer), 1e-6)
}