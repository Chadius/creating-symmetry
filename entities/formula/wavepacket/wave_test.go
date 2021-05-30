package wavepacket_test

import (
	. "gopkg.in/check.v1"
	"math"
	"math/cmplx"
	"testing"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/wavepacket"
	"wallpaper/entities/utility"
)

func Test(t *testing.T) { TestingT(t) }

type WaveFormulaTests struct {
	hexLatticeVectors *formula.LatticeVectorPair
	hexagonalWavePacket *wavepacket.Formula
}

var _ = Suite(&WaveFormulaTests{})

func (suite *WaveFormulaTests) SetUpTest(checker *C) {
	suite.hexLatticeVectors = &formula.LatticeVectorPair{
		XLatticeVector: complex(1,0),
		YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
	}
	suite.hexagonalWavePacket = &wavepacket.Formula{

		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:         1,
				PowerM:         -2,
				Multiplier: complex(1, 0),
			},
			{
				PowerN:         -2,
				PowerM:         1,
				Multiplier: complex(1, 0),
			},
			{
				PowerN:         1,
				PowerM:         1,
				Multiplier: complex(1, 0),
			},
		},
		Multiplier: complex(1, 0),
	}
}

func (suite *WaveFormulaTests) TestWaveFormulaCombinesEisensteinTerms(checker *C) {
	zInLatticeCoordinates := suite.hexLatticeVectors.ConvertToLatticeCoordinates(complex(math.Sqrt(3), -1 * math.Sqrt(3)))
	calculation := suite.hexagonalWavePacket.Calculate(zInLatticeCoordinates)
	total := calculation.Total

	expectedAnswer := cmplx.Exp(complex(0, 2 * math.Pi * (3 + math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-2 * math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-3 + math.Sqrt(3))))

	checker.Assert(real(total), utility.NumericallyCloseEnough{}, real(expectedAnswer), 1e-6)
	checker.Assert(imag(total), utility.NumericallyCloseEnough{}, imag(expectedAnswer), 1e-6)
}

func (suite *WaveFormulaTests) TestWaveFormulaShowsContributionsPerTerm(checker *C) {
	zInLatticeCoordinates := suite.hexLatticeVectors.ConvertToLatticeCoordinates(complex(math.Sqrt(3), -1 * math.Sqrt(3)))
	calculation := suite.hexagonalWavePacket.Calculate(zInLatticeCoordinates)

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

func (suite *WaveFormulaTests) TestWaveFormulaUsesMultiplier(checker *C) {
	suite.hexagonalWavePacket.Multiplier = complex(1/3.0, 0)
	zInLatticeCoordinates := suite.hexLatticeVectors.ConvertToLatticeCoordinates(complex(math.Sqrt(3), -1 * math.Sqrt(3)))
	calculation := suite.hexagonalWavePacket.Calculate(zInLatticeCoordinates)
	total := calculation.Total

	expectedAnswer := (cmplx.Exp(complex(0, 2 * math.Pi * (3 + math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-2 * math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-3 + math.Sqrt(3))))) / 3

	checker.Assert(real(total), utility.NumericallyCloseEnough{}, real(expectedAnswer), 1e-6)
	checker.Assert(imag(total), utility.NumericallyCloseEnough{}, imag(expectedAnswer), 1e-6)
}

func (suite *WaveFormulaTests) TestWaveFormulaMarshalFromJson(checker *C) {
	jsonByteStream := []byte(`{
				"multiplier": {
					"real": -1.0,
					"imaginary": 2e-2
				},
				"terms": [
					{
						"power_n": 12,
						"power_m": -10,
						"multiplier": {
							"real": -1.0,
							"imaginary": 2e-2
						}
					}
				]
	}`)
	wave, err := wavepacket.NewWaveFormulaFromJSON(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(real(wave.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(wave.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(wave.Terms, HasLen, 1)
	checker.Assert(wave.Terms[0].PowerN, Equals, 12)
}

func (suite *WaveFormulaTests) TestWaveFormulaMarshalFromYAML(checker *C) {
	yamlByteStream := []byte(`
multiplier:
  real: -1.0
  imaginary: 2e-2
terms:
  -
    power_n: 12
    power_m: -10
    multiplier:
      real: -1.0
      imaginary: 2e-2
`)
	wave, err := wavepacket.NewWaveFormulaFromYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(real(wave.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(wave.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(wave.Terms, HasLen, 1)
	checker.Assert(wave.Terms[0].PowerN, Equals, 12)
}