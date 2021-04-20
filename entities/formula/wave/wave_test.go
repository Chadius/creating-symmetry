package wave_test

import (
	. "gopkg.in/check.v1"
	"math"
	"math/cmplx"
	"testing"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/wave"
	"wallpaper/entities/utility"
)

func Test(t *testing.T) { TestingT(t) }

type WaveFormulaTests struct {
}

var _ = Suite(&WaveFormulaTests{})

func (suite *WaveFormulaTests) SetUpTest(checker *C) {
}

func (suite *WaveFormulaTests) TestWaveFormulaCombinesEisensteinTerms(checker *C) {
	waveFormula := &wave.Formula{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				XLatticeVector: complex(1,0),
				YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
				PowerN:         1,
				PowerM:         -2,
			},
			{
				XLatticeVector: complex(1,0),
				YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
				PowerN:         -2,
				PowerM:         1,
			},
			{
				XLatticeVector: complex(1,0),
				YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
				PowerN:         1,
				PowerM:         1,
			},
		},
	}
	calculation := waveFormula.Calculate(complex(math.Sqrt(3), -1 * math.Sqrt(3)))
	total := calculation.Total

	expectedAnswer := cmplx.Exp(complex(0, 2 * math.Pi * (3 + math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-2 * math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-3 + math.Sqrt(3))))

	checker.Assert(real(total), utility.NumericallyCloseEnough{}, real(expectedAnswer), 1e-6)
	checker.Assert(imag(total), utility.NumericallyCloseEnough{}, imag(expectedAnswer), 1e-6)
}

func (suite *WaveFormulaTests) TestWaveFormulaShowsContributionsPerTerm(checker *C) {
	waveFormula := &wave.Formula{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				XLatticeVector: complex(1,0),
				YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
				PowerN:         1,
				PowerM:         -2,
			},
			{
				XLatticeVector: complex(1,0),
				YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
				PowerN:         -2,
				PowerM:         1,
			},
			{
				XLatticeVector: complex(1,0),
				YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
				PowerN:         1,
				PowerM:         1,
			},
		},
	}
	calculation := waveFormula.Calculate(complex(math.Sqrt(3), -1 * math.Sqrt(3)))

	checker.Assert(calculation.ContributionByTerm, HasLen, 3)

	contributionOfTerm1 := cmplx.Exp(complex(0, 2 * math.Pi * (3 + math.Sqrt(3))))
	checker.Assert(real(calculation.ContributionByTerm[0]), utility.NumericallyCloseEnough{}, real(contributionOfTerm1), 1e-6)
	checker.Assert(imag(calculation.ContributionByTerm[0]), utility.NumericallyCloseEnough{}, imag(contributionOfTerm1), 1e-6)

	contributionOfTerm2 := cmplx.Exp(complex(0, 2 * math.Pi * (-2 * math.Sqrt(3))))
	checker.Assert(real(calculation.ContributionByTerm[1]), utility.NumericallyCloseEnough{}, real(contributionOfTerm2), 1e-6)
	checker.Assert(imag(calculation.ContributionByTerm[1]), utility.NumericallyCloseEnough{}, imag(contributionOfTerm2), 1e-6)

	contributionOfTerm3 := cmplx.Exp(complex(0, 2 * math.Pi * (-3 + math.Sqrt(3))))
	checker.Assert(real(calculation.ContributionByTerm[2]), utility.NumericallyCloseEnough{}, real(contributionOfTerm3), 1e-6)
	checker.Assert(imag(calculation.ContributionByTerm[2]), utility.NumericallyCloseEnough{}, imag(contributionOfTerm3), 1e-6)
}

