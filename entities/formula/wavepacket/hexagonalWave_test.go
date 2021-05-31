package wavepacket_test

import (
	. "gopkg.in/check.v1"
	"math"
	"math/cmplx"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/wavepacket"
	"wallpaper/entities/utility"
)

type HexagonalWaveFormula struct {
	hexagonalWavePacket *wavepacket.HexagonalWallpaperFormula
}

var _ = Suite(&HexagonalWaveFormula{})

func (suite *HexagonalWaveFormula) SetUpTest(checker *C) {
	suite.hexagonalWavePacket = &wavepacket.HexagonalWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.Formula{
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
}

func (suite *HexagonalWaveFormula) TestHexagonalWallpaperImpliesAveragedLockedTerms(checker *C) {
	suite.hexagonalWavePacket.SetUp()
	calculation := suite.hexagonalWavePacket.Calculate(complex(math.Sqrt(3), -1 * math.Sqrt(3)))
	total := calculation.Total

	expectedAnswer := (cmplx.Exp(complex(0, 2 * math.Pi * (3 + math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-2 * math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-3 + math.Sqrt(3))))) / 3

	checker.Assert(real(total), utility.NumericallyCloseEnough{}, real(expectedAnswer), 1e-6)
	checker.Assert(imag(total), utility.NumericallyCloseEnough{}, imag(expectedAnswer), 1e-6)
}

func (suite *HexagonalWaveFormula) TestUnmarshalFromJSON(checker *C) {
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
	hexFormula, err := wavepacket.NewHexagonalWallpaperFormulaFromJSON(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(real(hexFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(hexFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(hexFormula.Formula.WavePackets, HasLen, 1)
	checker.Assert(hexFormula.Formula.WavePackets[0].Terms[0].PowerM, Equals, -10)
}

func (suite *HexagonalWaveFormula) TestUnmarshalFromYAML(checker *C) {
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
	hexFormula, err := wavepacket.NewHexagonalWallpaperFormulaFromYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(real(hexFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(hexFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(hexFormula.Formula.WavePackets, HasLen, 1)
	checker.Assert(hexFormula.Formula.WavePackets[0].Terms[0].PowerM, Equals, -10)
}

type HexagonalWaveSymmetry struct {
	baseWavePacket *wavepacket.Formula
}

var _ = Suite(&HexagonalWaveSymmetry{})

func (suite *HexagonalWaveSymmetry) SetUpTest(checker *C) {
	suite.baseWavePacket = &wavepacket.Formula{
	Terms:[]*formula.EisensteinFormulaTerm{
			{
				PowerN:         8,
				PowerM:         -3,
			},
		},
		Multiplier: complex(1, 0),
	}
}

func (suite *HexagonalWaveSymmetry) TestNoSymmetryFound(checker *C) {
	arbitraryHex := wavepacket.HexagonalWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.Formula{
				suite.baseWavePacket,
			},
			Multiplier: complex(1, 0),
		},
	}

	arbitraryHex.SetUp()

	symmetriesFound := arbitraryHex.FindSymmetries()
	checker.Assert(symmetriesFound.P31m, Equals, false)
	checker.Assert(symmetriesFound.P3m1, Equals, false)
	checker.Assert(symmetriesFound.P6, Equals, false)
	checker.Assert(symmetriesFound.P6m, Equals, false)
}

func (suite *HexagonalWaveSymmetry) TestP31m(checker *C) {
	arbitraryHex := wavepacket.HexagonalWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.Formula{
				suite.baseWavePacket,
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacket.Terms[0].PowerM,
							PowerM:         suite.baseWavePacket.Terms[0].PowerN,
						},
					},
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	arbitraryHex.SetUp()

	symmetriesFound := arbitraryHex.FindSymmetries()
	checker.Assert(symmetriesFound.P31m, Equals, true)
	checker.Assert(symmetriesFound.P3m1, Equals, false)
	checker.Assert(symmetriesFound.P6, Equals, false)
	checker.Assert(symmetriesFound.P6m, Equals, false)
}

func (suite *HexagonalWaveSymmetry) TestP3m1(checker *C) {
	arbitraryHex := wavepacket.HexagonalWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.Formula{
				suite.baseWavePacket,
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacket.Terms[0].PowerM * -1,
							PowerM:         suite.baseWavePacket.Terms[0].PowerN * -1,
						},
					},
					Multiplier: suite.baseWavePacket.Multiplier,
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	arbitraryHex.SetUp()

	symmetriesFound := arbitraryHex.FindSymmetries()
	checker.Assert(symmetriesFound.P31m, Equals, false)
	checker.Assert(symmetriesFound.P3m1, Equals, true)
	checker.Assert(symmetriesFound.P6, Equals, false)
	checker.Assert(symmetriesFound.P6m, Equals, false)
}

func (suite *HexagonalWaveSymmetry) TestP6(checker *C) {
	arbitraryHex := wavepacket.HexagonalWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.Formula{
				suite.baseWavePacket,
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacket.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacket.Terms[0].PowerM * -1,
						},
					},
					Multiplier: suite.baseWavePacket.Multiplier,
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	arbitraryHex.SetUp()

	symmetriesFound := arbitraryHex.FindSymmetries()
	checker.Assert(symmetriesFound.P31m, Equals, false)
	checker.Assert(symmetriesFound.P3m1, Equals, false)
	checker.Assert(symmetriesFound.P6, Equals, true)
	checker.Assert(symmetriesFound.P6m, Equals, false)
}

func (suite *HexagonalWaveSymmetry) TestP6m(checker *C) {
	arbitraryHex := wavepacket.HexagonalWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.Formula{
				suite.baseWavePacket,
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacket.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacket.Terms[0].PowerM * -1,
						},
					},
					Multiplier: suite.baseWavePacket.Multiplier,
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacket.Terms[0].PowerM,
							PowerM:         suite.baseWavePacket.Terms[0].PowerN,
						},
					},
					Multiplier: suite.baseWavePacket.Multiplier,
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacket.Terms[0].PowerM * -1,
							PowerM:         suite.baseWavePacket.Terms[0].PowerN * -1,
						},
					},
					Multiplier: suite.baseWavePacket.Multiplier,
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	arbitraryHex.SetUp()

	symmetriesFound := arbitraryHex.FindSymmetries()
	checker.Assert(symmetriesFound.P31m, Equals, false)
	checker.Assert(symmetriesFound.P3m1, Equals, false)
	checker.Assert(symmetriesFound.P6, Equals, false)
	checker.Assert(symmetriesFound.P6m, Equals, true)
}

func (suite *HexagonalWaveSymmetry) TestMultipleSymmetries(checker *C) {
	arbitraryHex := wavepacket.HexagonalWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.Formula{
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         1,
							PowerM:         -1,
						},
					},
					Multiplier: suite.baseWavePacket.Multiplier,
				},
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         -1,
							PowerM:         1,
						},
					},
					Multiplier: suite.baseWavePacket.Multiplier,
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	arbitraryHex.SetUp()

	symmetriesFound := arbitraryHex.FindSymmetries()
	checker.Assert(symmetriesFound.P31m, Equals, true)
	checker.Assert(symmetriesFound.P3m1, Equals, false)
	checker.Assert(symmetriesFound.P6, Equals, true)
	checker.Assert(symmetriesFound.P6m, Equals, false)
}

type HexagonalCreatedWithDesiredSymmetry struct {
	singleEisensteinFormulaTerm []*formula.EisensteinFormulaTerm
	wallpaperMultiplier complex128
}

var _ = Suite(&HexagonalCreatedWithDesiredSymmetry{})

func (suite *HexagonalCreatedWithDesiredSymmetry) SetUpTest (checker *C) {
	suite.singleEisensteinFormulaTerm = []*formula.EisensteinFormulaTerm{
		{
			PowerM: -2,
			PowerN: 1,
			Multiplier: complex(1, 0),
		},
	}
	suite.wallpaperMultiplier = complex(1, 0)
}

func (suite *HexagonalCreatedWithDesiredSymmetry) SymmetryCheck(checker *C, hexFormula *wavepacket.HexagonalWallpaperFormula , desiredSymmetry wavepacket.Symmetry) {
	symmetriesFound := hexFormula.FindSymmetries()
	checker.Assert(symmetriesFound.P3, Equals, desiredSymmetry.P3)
	checker.Assert(symmetriesFound.P31m, Equals, desiredSymmetry.P31m)
	checker.Assert(symmetriesFound.P3m1, Equals, desiredSymmetry.P3m1)
	checker.Assert(symmetriesFound.P6, Equals, desiredSymmetry.P6)
	checker.Assert(symmetriesFound.P6m, Equals, desiredSymmetry.P6m)
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestNoSymmetryDoesNotChangePattern(checker *C) {
	hexFormula, err := wavepacket.NewHexagonalWallpaperFormulaWithSymmetry(
		suite.singleEisensteinFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{},
	)

	checker.Assert(err, IsNil)
	checker.Assert(hexFormula.Formula.WavePackets, HasLen, 1)
	checker.Assert(hexFormula.Formula.WavePackets[0].Terms, HasLen, 3)

	suite.SymmetryCheck(checker, hexFormula, wavepacket.Symmetry{
		P3:   true,
	})
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithP31m(checker *C) {
	hexFormula, err := wavepacket.NewHexagonalWallpaperFormulaWithSymmetry(
		suite.singleEisensteinFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P31m: true,
		},
	)

	checker.Assert(err, IsNil)
	checker.Assert(hexFormula.Formula.WavePackets, HasLen, 2)
	checker.Assert(hexFormula.Formula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, 1)
	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, -2)

	suite.SymmetryCheck(checker, hexFormula, wavepacket.Symmetry{
		P3:   true,
		P31m:   true,
	})
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithP3m1(checker *C) {
	hexFormula, err := wavepacket.NewHexagonalWallpaperFormulaWithSymmetry(
		suite.singleEisensteinFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P3m1: true,
		},
	)

	checker.Assert(err, IsNil)
	checker.Assert(hexFormula.Formula.WavePackets, HasLen, 2)
	checker.Assert(hexFormula.Formula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, -1)
	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, 2)

	suite.SymmetryCheck(checker, hexFormula, wavepacket.Symmetry{
		P3:   true,
		P3m1:   true,
	})
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithP6(checker *C) {
	hexFormula, err := wavepacket.NewHexagonalWallpaperFormulaWithSymmetry(
		suite.singleEisensteinFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P6: true,
		},
	)

	checker.Assert(err, IsNil)
	checker.Assert(hexFormula.Formula.WavePackets, HasLen, 2)
	checker.Assert(hexFormula.Formula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, 2)
	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, -1)

	suite.SymmetryCheck(checker, hexFormula, wavepacket.Symmetry{
		P3:   true,
		P6:   true,
	})
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithP6m(checker *C) {
	hexFormula, err := wavepacket.NewHexagonalWallpaperFormulaWithSymmetry(
		suite.singleEisensteinFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P6m: true,
		},
	)

	checker.Assert(err, IsNil)
	checker.Assert(hexFormula.Formula.WavePackets, HasLen, 4)
	checker.Assert(hexFormula.Formula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, 2)
	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, -1)

	checker.Assert(hexFormula.Formula.WavePackets[2].Terms[0].PowerM, Equals, 1)
	checker.Assert(hexFormula.Formula.WavePackets[2].Terms[0].PowerN, Equals, -2)

	checker.Assert(hexFormula.Formula.WavePackets[3].Terms[0].PowerM, Equals, -1)
	checker.Assert(hexFormula.Formula.WavePackets[3].Terms[0].PowerN, Equals, 2)

	suite.SymmetryCheck(checker, hexFormula, wavepacket.Symmetry{
		P3:   true,
		P6m:   true,
	})
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithMultipleTerms(checker *C) {
	hexFormula, err := wavepacket.NewHexagonalWallpaperFormulaWithSymmetry(
		[]*formula.EisensteinFormulaTerm{
			{
				PowerM: -2,
				PowerN: 1,
				Multiplier: complex(1, 0),
			},
			{
				PowerM: -3,
				PowerN: 5,
				Multiplier: complex(0, 1),
			},
		},
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P31m: true,
		},
	)

	checker.Assert(err, IsNil)
	checker.Assert(hexFormula.Formula.WavePackets, HasLen, 4)
	checker.Assert(hexFormula.Formula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, 1)
	checker.Assert(hexFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, -2)

	checker.Assert(hexFormula.Formula.WavePackets[2].Terms[0].PowerM, Equals, -3)
	checker.Assert(hexFormula.Formula.WavePackets[2].Terms[0].PowerN, Equals, 5)

	checker.Assert(hexFormula.Formula.WavePackets[3].Terms[0].PowerM, Equals, 5)
	checker.Assert(hexFormula.Formula.WavePackets[3].Terms[0].PowerN, Equals, -3)
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestOnlyP3CanBeCombined(checker *C) {
	_, err := wavepacket.NewHexagonalWallpaperFormulaWithSymmetry(
		suite.singleEisensteinFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P3: true,
			P31m: true,
		},
	)

	checker.Assert(err, IsNil)

	_, errIncompatible := wavepacket.NewHexagonalWallpaperFormulaWithSymmetry(
		suite.singleEisensteinFormulaTerm,
		suite.wallpaperMultiplier,
		&wavepacket.Symmetry{
			P3: true,
			P31m: true,
			P6: true,
		},
	)

	checker.Assert(errIncompatible, ErrorMatches, "invalid desired symmetry")
}

