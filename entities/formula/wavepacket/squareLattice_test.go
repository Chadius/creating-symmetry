package wavepacket_test

import (
	. "gopkg.in/check.v1"
	"math"
	"math/cmplx"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/wavepacket"
	"wallpaper/entities/utility"
)

type SquareLatticeWallpaper struct {
	squareWavePacket *wavepacket.SquareWallpaperFormula
}

var _ = Suite(&SquareLatticeWallpaper{})

func (suite *SquareLatticeWallpaper) SetUpTest(checker *C) {
	suite.squareWavePacket = &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         1,
							PowerM:         -2,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
}

func (suite *SquareLatticeWallpaper) TestSquareLatticeWallpaperImpliesAveragedLockedTerms(checker *C) {
	suite.squareWavePacket.SetUp()
	calculation := suite.squareWavePacket.Calculate(complex(2, 0.5))
	total := calculation.Total

	expectedAnswer :=
		(
			cmplx.Exp(complex(0, 2 * math.Pi)) +
			cmplx.Exp(complex(0, 2 * math.Pi * -3.5)) +
			cmplx.Exp(complex(0, 2 * math.Pi * -1)) +
			cmplx.Exp(complex(0, 2 * math.Pi * 4.5)))/4

	checker.Assert(real(total), utility.NumericallyCloseEnough{}, real(expectedAnswer), 1e-6)
	checker.Assert(imag(total), utility.NumericallyCloseEnough{}, imag(expectedAnswer), 1e-6)
}

func (suite *SquareLatticeWallpaper) TestUnmarshalFromJSON(checker *C) {
	jsonByteStream := []byte(`{
				"multiplier": {
					"real": -1.0,
					"imaginary": 2e-2
				},
				"wave_packets": [
					{
						"multiplier": {
							"real": -1.0,
							"imaginary": 2e-2
						},
						"terms": [
							{
								"power_n": 12,
								"power_m": -10
							}
						]
					}
				]
			}`)
	squareFormula, err := wavepacket.NewSquareWallpaperFormulaFromJSON(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(real(squareFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(squareFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(squareFormula.Formula.WavePackets, HasLen, 1)
	checker.Assert(squareFormula.Formula.WavePackets[0].Terms[0].PowerM, Equals, -10)
}

func (suite *SquareLatticeWallpaper) TestUnmarshalFromYAML(checker *C) {
	yamlByteStream := []byte(`
multiplier:
 real: -1.0
 imaginary: 2e-2
wave_packets:
 -
   multiplier:
     real: -1.0
     imaginary: 2e-2
   terms:
     -
       power_n: 12
       power_m: -10
`)
	squareFormula, err := wavepacket.NewSquareWallpaperFormulaFromYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(real(squareFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(squareFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(squareFormula.Formula.WavePackets, HasLen, 1)
	checker.Assert(squareFormula.Formula.WavePackets[0].Terms[0].PowerM, Equals, -10)
}

type SquareCreatedWithDesiredSymmetry struct {
	singleEisensteinOddSumFormulaTerm []*formula.EisensteinFormulaTerm
	singleEisensteinEvenSumFormulaTerm []*formula.EisensteinFormulaTerm
	wallpaperMultiplier complex128
}

var _ = Suite(&SquareCreatedWithDesiredSymmetry{})

func (suite *SquareCreatedWithDesiredSymmetry) SetUpTest (checker *C) {
	suite.singleEisensteinOddSumFormulaTerm = []*formula.EisensteinFormulaTerm{
		{
			PowerM: -2,
			PowerN: 1,
			Multiplier: complex(1, 0),
		},
	}
	suite.singleEisensteinEvenSumFormulaTerm = []*formula.EisensteinFormulaTerm{
		{
			PowerM: 3,
			PowerN: -1,
			Multiplier: complex(1, 0),
		},
	}
	suite.wallpaperMultiplier = complex(1, 0)
}

func (suite *SquareCreatedWithDesiredSymmetry) TestNoSymmetryDoesNotChangePattern(checker *C) {
	squareFormula, err := wavepacket.NewSquareWallpaperFormulaWithSymmetry(
		suite.singleEisensteinOddSumFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{},
	)

	checker.Assert(err, IsNil)
	checker.Assert(squareFormula.Formula.WavePackets, HasLen, 1)
	checker.Assert(squareFormula.Formula.WavePackets[0].Terms, HasLen, 4)

	checker.Assert(squareFormula.HasSymmetry(wavepacket.P4), Equals, true)
}

func (suite *SquareCreatedWithDesiredSymmetry) TestCreateWallpaperWithP4m(checker *C) {
	squareFormula, err := wavepacket.NewSquareWallpaperFormulaWithSymmetry(
		suite.singleEisensteinOddSumFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P4m: true,
		},
	)

	checker.Assert(err, IsNil)
	checker.Assert(squareFormula.Formula.WavePackets, HasLen, 2)
	checker.Assert(squareFormula.Formula.WavePackets[0].Terms, HasLen, 4)

	checker.Assert(squareFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, 1)
	checker.Assert(squareFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, -2)

	checker.Assert(squareFormula.HasSymmetry(wavepacket.P4), Equals, true)
	checker.Assert(squareFormula.HasSymmetry(wavepacket.P4m), Equals, true)
}

func (suite *SquareCreatedWithDesiredSymmetry) TestCreateWallpaperWithP4gAndOddSumPowers(checker *C) {
	squareFormula, errOddSum := wavepacket.NewSquareWallpaperFormulaWithSymmetry(
		suite.singleEisensteinOddSumFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P4g: true,
		},
	)

	checker.Assert(errOddSum, IsNil)
	checker.Assert(squareFormula.Formula.WavePackets, HasLen, 2)
	checker.Assert(squareFormula.Formula.WavePackets[0].Terms, HasLen, 4)

	checker.Assert(real(squareFormula.Formula.WavePackets[1].Terms[0].Multiplier), utility.NumericallyCloseEnough{}, real(suite.wallpaperMultiplier) * -1, 1e-6)
	checker.Assert(imag(squareFormula.Formula.WavePackets[1].Terms[0].Multiplier), utility.NumericallyCloseEnough{}, imag(suite.wallpaperMultiplier) * -1, 1e-6)
	checker.Assert(squareFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, 1)
	checker.Assert(squareFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, -2)

	checker.Assert(squareFormula.HasSymmetry(wavepacket.P4), Equals, true)
	checker.Assert(squareFormula.HasSymmetry(wavepacket.P4g), Equals, true)
}

func (suite *SquareCreatedWithDesiredSymmetry) TestCreateWallpaperWithP4gAndEvenSumPowers(checker *C) {
	squareFormula, errEvenSum := wavepacket.NewSquareWallpaperFormulaWithSymmetry(
		suite.singleEisensteinEvenSumFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P4g: true,
		},
	)

	checker.Assert(errEvenSum, IsNil)
	checker.Assert(squareFormula.Formula.WavePackets, HasLen, 2)
	checker.Assert(squareFormula.Formula.WavePackets[0].Terms, HasLen, 4)

	checker.Assert(real(squareFormula.Formula.WavePackets[1].Terms[0].Multiplier), utility.NumericallyCloseEnough{}, real(suite.wallpaperMultiplier), 1e-6)
	checker.Assert(imag(squareFormula.Formula.WavePackets[1].Terms[0].Multiplier), utility.NumericallyCloseEnough{}, imag(suite.wallpaperMultiplier), 1e-6)
	checker.Assert(squareFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, -1)
	checker.Assert(squareFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, 3)

	checker.Assert(squareFormula.HasSymmetry(wavepacket.P4), Equals, true)
	checker.Assert(squareFormula.HasSymmetry(wavepacket.P4m), Equals, true)
	checker.Assert(squareFormula.HasSymmetry(wavepacket.P4g), Equals, true)
}

func (suite *SquareCreatedWithDesiredSymmetry) TestP4mP4gOnlyCombinesForEvenSumTerms(checker *C) {
	_, errEven := wavepacket.NewSquareWallpaperFormulaWithSymmetry(
		suite.singleEisensteinEvenSumFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P4: true,
			P4m: true,
			P4g: true,
		},
	)

	checker.Assert(errEven, IsNil)

	_, errIncompatible := wavepacket.NewSquareWallpaperFormulaWithSymmetry(
		suite.singleEisensteinOddSumFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P4: true,
			P4m: true,
			P4g: true,
		},
	)

	checker.Assert(errIncompatible, ErrorMatches, "invalid desired symmetry")
}

type SquareLatticeDetectRelationship struct {}

var _= Suite(&SquareLatticeDetectRelationship{})

func (suite *SquareLatticeDetectRelationship) TestP4mSymmetryDetectedAcrossSinglePairs (checker *C) {
	p4mSquare := &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -2,
							PowerM: 1,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 1,
							PowerM: -2,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
	p4mSquare.SetUp()

	checker.Assert(p4mSquare.HasSymmetry(wavepacket.P4m), Equals, true)
}

func (suite *SquareLatticeDetectRelationship) TestP4mSymmetryDetectedAcrossMultiplePairs (checker *C) {
	p4mSquare := &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -5,
							PowerM: 8,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 2,
							PowerM: -1,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 8,
							PowerM: -5,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -1,
							PowerM: 2,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
	p4mSquare.SetUp()

	checker.Assert(p4mSquare.HasSymmetry(wavepacket.P4m), Equals, true)
}

func (suite *SquareLatticeDetectRelationship) TestP4SymmetryIsAlwaysTrueForSquarePatterns (checker *C) {
	p4Square := &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -2,
							PowerM: 1,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
	p4Square.SetUp()

	checker.Assert(p4Square.HasSymmetry(wavepacket.P4), Equals, true)
}

func (suite *SquareLatticeDetectRelationship) TestFewerThanTwoPacketsMeansNoSymmetry (checker *C) {
	p4Square := &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -2,
							PowerM: 1,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
	p4Square.SetUp()

	checker.Assert(p4Square.HasSymmetry(wavepacket.P4m), Equals, false)
	checker.Assert(p4Square.HasSymmetry(wavepacket.P4g), Equals, false)
}

func (suite *SquareLatticeDetectRelationship) TestOddNumberOfWavePacketsMeansNoSymmetry (checker *C) {
	p4Square := &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -2,
							PowerM: 1,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 5,
							PowerM: -2,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -2,
							PowerM: 5,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 1,
							PowerM: -2,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 1,
							PowerM: -2,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
	p4Square.SetUp()

	checker.Assert(p4Square.HasSymmetry(wavepacket.P4m), Equals, false)
}

func (suite *SquareLatticeDetectRelationship) TestP4m (checker *C) {
	p4mSquare := &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -2,
							PowerM: 1,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 1,
							PowerM: -2,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
	p4mSquare.SetUp()

	checker.Assert(p4mSquare.HasSymmetry(wavepacket.P4m), Equals, true)
}

func (suite *SquareLatticeDetectRelationship) TestP4g (checker *C) {
	p4gOddSum := &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -2,
							PowerM: 1,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 1,
							PowerM: -2,
							Multiplier: complex(-1, 0),
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
	p4gOddSum.SetUp()

	checker.Assert(p4gOddSum.HasSymmetry(wavepacket.P4g), Equals, true)

	p4gEvenSum := &wavepacket.SquareWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: -3,
							PowerM: 1,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN: 1,
							PowerM: -3,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}
	p4gEvenSum.SetUp()

	checker.Assert(p4gEvenSum.HasSymmetry(wavepacket.P4g), Equals, true)
}
