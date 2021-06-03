package wavepacket_test

import (
	. "gopkg.in/check.v1"
	"math"
	"math/cmplx"
	"testing"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/coefficient"
	"wallpaper/entities/formula/wavepacket"
	"wallpaper/entities/utility"
)

func Test(t *testing.T) { TestingT(t) }

type WaveFormulaTests struct {
	hexLatticeVectors *formula.LatticeVectorPair
	hexagonalWavePacket *wavepacket.WavePacket
}

var _ = Suite(&WaveFormulaTests{})

func (suite *WaveFormulaTests) SetUpTest(checker *C) {
	suite.hexLatticeVectors = &formula.LatticeVectorPair{
		XLatticeVector: complex(1,0),
		YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
	}
	suite.hexagonalWavePacket = &wavepacket.WavePacket{

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

type WavePacketRelationshipTest struct {
	aPlusNPlusMOddWavePacket *wavepacket.WavePacket
	aPlusMMinusNOddWavePacket *wavepacket.WavePacket
	aPlusMPlusNOddWavePacket *wavepacket.WavePacket
	aMinusNMinusMOddWavePacket *wavepacket.WavePacket
	aMinusMMinusNOddWavePacket *wavepacket.WavePacket
	aMinusMPlusNOddWavePacket *wavepacket.WavePacket
	aPlusMPlusNOddNegatedWavePacket *wavepacket.WavePacket
	aMinusSumNAndMPlusNOddWavePacket *wavepacket.WavePacket

	aPlusNPlusMEvenWavePacket *wavepacket.WavePacket
	aPlusMPlusNEvenWavePacket *wavepacket.WavePacket
	aPlusMPlusNEvenNegatedWavePacket *wavepacket.WavePacket
	aMinusMMinusNOddNegatedWavePacket *wavepacket.WavePacket
	aMinusMMinusNEvenWavePacket *wavepacket.WavePacket
	aMinusMMinusNEvenNegatedWavePacket *wavepacket.WavePacket
	aPlusMMinusSumNAndMOddWavePacket *wavepacket.WavePacket
	aMinusSumNAndMPlusNWavePacket *wavepacket.WavePacket
}

var _ = Suite(&WavePacketRelationshipTest{})

func (suite *WavePacketRelationshipTest) SetUpTest(checker *C) {
	suite.aPlusNPlusMOddWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -1,
				PowerM:     4,
				Multiplier: complex(2, 1),
			},
		},
	}
	suite.aPlusMPlusNOddWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     4,
				PowerM:     -1,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aMinusNMinusMOddWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     1,
				PowerM:     -4,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aMinusMMinusNOddWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -4,
				PowerM:     1,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aPlusMMinusNOddWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     4,
				PowerM:     1,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aMinusMPlusNOddWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -4,
				PowerM:     -1,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aPlusMPlusNOddNegatedWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     4,
				PowerM:     -1,
				Multiplier: complex(-2, -1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aMinusMMinusNOddNegatedWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -4,
				PowerM:     1,
				Multiplier: complex(-2, -1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aPlusMMinusSumNAndMOddWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     4,
				PowerM:     -3,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aMinusSumNAndMPlusNWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -3,
				PowerM:     -1,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aMinusSumNAndMPlusNOddWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -3,
				PowerM:     -1,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}

	suite.aPlusNPlusMEvenWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -6,
				PowerM:     2,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aPlusMPlusNEvenWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     2,
				PowerM:     -6,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aPlusMPlusNEvenNegatedWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     2,
				PowerM:     -6,
				Multiplier: complex(-2, -1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aMinusMMinusNEvenWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -2,
				PowerM:     6,
				Multiplier: complex(2, 1),
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.aMinusMMinusNEvenNegatedWavePacket = &wavepacket.WavePacket{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				PowerN:     -2,
				PowerM:     6,
				Multiplier: complex(-2, -1),
			},
		},
		Multiplier: complex(1, 0),
	}
}

func (suite *WavePacketRelationshipTest) TestLessThanTwoWavePacketsHasNoRelationship(checker *C) {
	checker.Assert(
		wavepacket.GetWavePacketRelationship(nil, nil),
		HasLen, 0)

	checker.Assert(
		wavepacket.GetWavePacketRelationship(
			suite.aPlusNPlusMOddWavePacket,
			nil,
		),
		HasLen, 0)
}

func (suite *WavePacketRelationshipTest) TestPlusNPlusM(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aPlusNPlusMOddWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(relationshipsFound[0], Equals, coefficient.PlusNPlusM)
}

func (suite *WavePacketRelationshipTest) TestMinusNMinusM(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aMinusNMinusMOddWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(relationshipsFound[0], Equals, coefficient.MinusNMinusM)
}

func (suite *WavePacketRelationshipTest) TestMinusMMinusN(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aMinusMMinusNOddWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(relationshipsFound[0], Equals, coefficient.MinusMMinusN)
}

func (suite *WavePacketRelationshipTest) TestPlusMPlusNMaybeFlipScale(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aPlusMPlusNOddNegatedWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(wavepacket.ContainsRelationship(relationshipsFound, coefficient.PlusMPlusNMaybeFlipScale), Equals, true)

	relationshipsFound = wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMEvenWavePacket,
		suite.aPlusMPlusNEvenNegatedWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 0)

	relationshipsFound = wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMEvenWavePacket,
		suite.aPlusMPlusNEvenWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 2)
	checker.Assert(wavepacket.ContainsRelationship(relationshipsFound, coefficient.PlusMPlusNMaybeFlipScale), Equals, true)
	checker.Assert(wavepacket.ContainsRelationship(relationshipsFound, coefficient.PlusMPlusN), Equals, true)
}

func (suite *WavePacketRelationshipTest) TestMinusMMinusNMaybeFlipScale(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aMinusMMinusNOddNegatedWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(relationshipsFound[0], Equals, coefficient.MinusMMinusNMaybeFlipScale)

	relationshipsFound = wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMEvenWavePacket,
		suite.aMinusMMinusNEvenWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 2)
	checker.Assert(wavepacket.ContainsRelationship(relationshipsFound, coefficient.MinusMMinusNMaybeFlipScale), Equals, true)
	checker.Assert(wavepacket.ContainsRelationship(relationshipsFound, coefficient.MinusMMinusN), Equals, true)

	relationshipsFound = wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMEvenWavePacket,
		suite.aMinusMMinusNEvenNegatedWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 0)
}

func (suite *WavePacketRelationshipTest) TestPlusMMinusSumNAndM(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aPlusMMinusSumNAndMOddWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(relationshipsFound[0], Equals, coefficient.PlusMMinusSumNAndM)
}

func (suite *WavePacketRelationshipTest) TestMinusSumNAndMPlusN(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aMinusSumNAndMPlusNOddWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(relationshipsFound[0], Equals, coefficient.MinusSumNAndMPlusN)
}

func (suite *WavePacketRelationshipTest) TestPlusMMinusN(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aPlusMMinusNOddWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(relationshipsFound[0], Equals, coefficient.PlusMMinusN)
}

func (suite *WavePacketRelationshipTest) TestMinusMPlusN(checker *C) {
	relationshipsFound := wavepacket.GetWavePacketRelationship(
		suite.aPlusNPlusMOddWavePacket,
		suite.aMinusMPlusNOddWavePacket,
	)

	checker.Assert(relationshipsFound, HasLen, 1)
	checker.Assert(relationshipsFound[0], Equals, coefficient.MinusMPlusN)
}
