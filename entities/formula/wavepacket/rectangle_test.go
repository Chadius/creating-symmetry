package wavepacket_test

import (
	. "gopkg.in/check.v1"
	"math"
	"math/cmplx"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/wavepacket"
	"wallpaper/entities/utility"
)

type RectangularWallpaper struct {
	RectangularFormula *wavepacket.RectangularWallpaperFormula
}

var _ = Suite(&RectangularWallpaper{})

func (suite *RectangularWallpaper) SetUpTest (checker *C) {
	suite.RectangularFormula = &wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
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
		LatticeHeight: 0.5,
	}
	suite.RectangularFormula.SetUp()
}

func (suite *RectangularWallpaper) TestSetupCreatesLatticeVectors (checker *C) {
	checker.Assert(real(suite.RectangularFormula.Formula.Lattice.XLatticeVector), utility.NumericallyCloseEnough{}, 1, 1e-6)
	checker.Assert(imag(suite.RectangularFormula.Formula.Lattice.XLatticeVector), utility.NumericallyCloseEnough{}, 0, 1e-6)

	checker.Assert(real(suite.RectangularFormula.Formula.Lattice.YLatticeVector), utility.NumericallyCloseEnough{}, 0, 1e-6)
	checker.Assert(imag(suite.RectangularFormula.Formula.Lattice.YLatticeVector), utility.NumericallyCloseEnough{}, 0.5, 1e-6)
}

func (suite *RectangularWallpaper) TestRaiseErrorIfHeightIsZero (checker *C) {
	RectangularFormulaWithNoHeight := &wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
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
		LatticeHeight: 0,
	}
	err := RectangularFormulaWithNoHeight.SetUp()

	checker.Assert(err, ErrorMatches, "lattice vectors cannot be \\(0,0\\)")
}

func (suite *RectangularWallpaper) TestSetupLocksPairs (checker *C) {
	checker.Assert(suite.RectangularFormula.Formula.WavePackets[0].Terms, HasLen, 1)
}

func (suite *RectangularWallpaper) TestCalculationOfPoints (checker *C) {
	calculation := suite.RectangularFormula.Calculate(complex(0.75, -0.25))
	total := calculation.Total

	expectedAnswer := cmplx.Exp(complex(0, math.Pi * 7 / 2))

	checker.Assert(real(total), utility.NumericallyCloseEnough{}, real(expectedAnswer), 1e-6)
	checker.Assert(imag(total), utility.NumericallyCloseEnough{}, imag(expectedAnswer), 1e-6)
}

func (suite *RectangularWallpaper) TestUnmarshalFromJSON(checker *C) {
	jsonByteStream := []byte(`{
				"lattice_height": 0.3,
				"formula": {
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
				}
			}`)
	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaFromJSON(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(RectangularFormula.LatticeHeight, utility.NumericallyCloseEnough{}, 0.3, 1e-6)
	checker.Assert(real(RectangularFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(RectangularFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 1)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms[0].PowerN, Equals, 12)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms[0].PowerM, Equals, -10)
}

func (suite *RectangularWallpaper) TestUnmarshalFromYAML(checker *C) {
	yamlByteStream := []byte(`
lattice_height: 0.3
formula:
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
	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaFromYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(RectangularFormula.LatticeHeight, utility.NumericallyCloseEnough{}, 0.3, 1e-6)
	checker.Assert(real(RectangularFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(RectangularFormula.Formula.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 1)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms[0].PowerM, Equals, -10)
}

type RectangularWaveSymmetry struct {
	baseWavePacketWithEvenPowerNAndOddPowerSum *wavepacket.WavePacket
	baseWavePacketWithOddPowerNAndEvenPowerSum *wavepacket.WavePacket
}

var _ = Suite(&RectangularWaveSymmetry{})

func (suite *RectangularWaveSymmetry) SetUpTest(checker *C) {
	suite.baseWavePacketWithEvenPowerNAndOddPowerSum = &wavepacket.WavePacket{
		Terms:[]*formula.EisensteinFormulaTerm{
			{
				PowerN:         8,
				PowerM:         -3,
				Multiplier: complex(1, 0),
			},
		},
		Multiplier: complex(1, 0),
	}

	suite.baseWavePacketWithOddPowerNAndEvenPowerSum = &wavepacket.WavePacket{
		Terms:[]*formula.EisensteinFormulaTerm{
			{
				PowerN:         7,
				PowerM:         -3,
				Multiplier: complex(1, 0),
			},
		},
		Multiplier: complex(1, 0),
	}
}

func (suite *RectangularWaveSymmetry) TestNoSymmetryFound(checker *C) {
	rectangularFormula := wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				suite.baseWavePacketWithEvenPowerNAndOddPowerSum,
			},
			Multiplier: complex(1, 0),
		},
	}

	rectangularFormula.SetUp()
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pm), Equals, false)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pg), Equals, false)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularWaveSymmetry) TestPm(checker *C) {
	rectangularFormula := wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				suite.baseWavePacketWithEvenPowerNAndOddPowerSum,
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerN,
							PowerM:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerM * -1,
							Multiplier: suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Multiplier,
						},
					},
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	rectangularFormula.SetUp()
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pm), Equals, true)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularWaveSymmetry) TestPg(checker *C) {
	rectangularFormulaWithEvenPowerN := wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				suite.baseWavePacketWithEvenPowerNAndOddPowerSum,
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerN,
							PowerM:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerM * -1,
							Multiplier: suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Multiplier,
						},
					},
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	rectangularFormulaWithEvenPowerN.SetUp()
	checker.Assert(rectangularFormulaWithEvenPowerN.HasSymmetry(wavepacket.Pm), Equals, true)
	checker.Assert(rectangularFormulaWithEvenPowerN.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(rectangularFormulaWithEvenPowerN.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(rectangularFormulaWithEvenPowerN.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(rectangularFormulaWithEvenPowerN.HasSymmetry(wavepacket.Pgg), Equals, false)

	rectangularFormulaWithOddPowerN := wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				suite.baseWavePacketWithOddPowerNAndEvenPowerSum,
				{
					Terms: []*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerN,
							PowerM:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerM * -1,
							Multiplier: 	suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Multiplier * -1,
						},
					},
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	rectangularFormulaWithOddPowerN.SetUp()
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pm), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularWaveSymmetry) TestPmmAndPmgWithEvenPowerN(checker *C) {
	rectangularFormula := wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				suite.baseWavePacketWithEvenPowerNAndOddPowerSum,
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerM * -1,
							Multiplier: complex(1, 0),
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerN,
							PowerM:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerM * -1,
							Multiplier: complex(1, 0),
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerM,
							Multiplier: complex(1, 0),
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	rectangularFormula.SetUp()
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pm), Equals, true)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, true)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, true)
	checker.Assert(rectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularWaveSymmetry) TestPmgWithOddPowerN(checker *C) {
	rectangularFormulaWithOddPowerN := wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				suite.baseWavePacketWithOddPowerNAndEvenPowerSum,
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerM * -1,
							Multiplier: suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Multiplier,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerN,
							PowerM:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerM * -1,
							Multiplier: suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Multiplier * -1,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerM,
							Multiplier: suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Multiplier * -1,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	rectangularFormulaWithOddPowerN.SetUp()
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pm), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pmg), Equals, true)
	checker.Assert(rectangularFormulaWithOddPowerN.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularWaveSymmetry) TestPgg(checker *C) {
	rectangularFormulaWithOddPowerSum := wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				suite.baseWavePacketWithEvenPowerNAndOddPowerSum,
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerM * -1,
							Multiplier: suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Multiplier,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerN,
							PowerM:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerM * -1,
							Multiplier: -1 * suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Multiplier,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Terms[0].PowerM,
							Multiplier: -1 * suite.baseWavePacketWithEvenPowerNAndOddPowerSum.Multiplier,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	rectangularFormulaWithOddPowerSum.SetUp()
	checker.Assert(rectangularFormulaWithOddPowerSum.HasSymmetry(wavepacket.Pm), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerSum.HasSymmetry(wavepacket.Pg), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerSum.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerSum.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(rectangularFormulaWithOddPowerSum.HasSymmetry(wavepacket.Pgg), Equals, true)

	rectangularFormulaWithEvenPowerSum := wavepacket.RectangularWallpaperFormula{
		Formula: &wavepacket.WallpaperFormula{
			WavePackets: []*wavepacket.WavePacket{
				suite.baseWavePacketWithOddPowerNAndEvenPowerSum,
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerM * -1,
							Multiplier: suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Multiplier,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerN,
							PowerM:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerM * -1,
							Multiplier: suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Multiplier,
						},
					},
					Multiplier: complex(1, 0),
				},
				{
					Terms:[]*formula.EisensteinFormulaTerm{
						{
							PowerN:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerN * -1,
							PowerM:         suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Terms[0].PowerM,
							Multiplier: suite.baseWavePacketWithOddPowerNAndEvenPowerSum.Multiplier,
						},
					},
					Multiplier: complex(1, 0),
				},
			},
			Multiplier: complex(1, 0),
		},
	}

	rectangularFormulaWithEvenPowerSum.SetUp()
	checker.Assert(rectangularFormulaWithEvenPowerSum.HasSymmetry(wavepacket.Pm), Equals, true)
	checker.Assert(rectangularFormulaWithEvenPowerSum.HasSymmetry(wavepacket.Pg), Equals, false)
	checker.Assert(rectangularFormulaWithEvenPowerSum.HasSymmetry(wavepacket.Pmm), Equals, true)
	checker.Assert(rectangularFormulaWithEvenPowerSum.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(rectangularFormulaWithEvenPowerSum.HasSymmetry(wavepacket.Pgg), Equals, true)
}

type RectangularCreatedWithDesiredSymmetry struct {
	eisensteinTermWithEvenPowerNAndOddPowerSum []*formula.EisensteinFormulaTerm
	eisensteinTermWithOddPowerNAndEvenPowerSum []*formula.EisensteinFormulaTerm
	wallpaperMultiplier complex128
	LatticeHeight float64
}

var _ = Suite(&RectangularCreatedWithDesiredSymmetry{})

func (suite *RectangularCreatedWithDesiredSymmetry) SetUpTest (checker *C) {
	suite.eisensteinTermWithEvenPowerNAndOddPowerSum = []*formula.EisensteinFormulaTerm{
		{
			PowerN:         8,
			PowerM:         -3,
			Multiplier: complex(1, 0),
		},
	}
	suite.eisensteinTermWithOddPowerNAndEvenPowerSum = []*formula.EisensteinFormulaTerm{
			{
				PowerN:         7,
				PowerM:         -3,
				Multiplier: complex(1, 0),
			},
		}

	suite.wallpaperMultiplier = complex(1, 0)
	suite.LatticeHeight = 2.0
}

func (suite *RectangularCreatedWithDesiredSymmetry) TestCreateWallpaperWithPm(checker *C) {
	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaWithSymmetry(
		suite.eisensteinTermWithOddPowerNAndEvenPowerSum,
		suite.wallpaperMultiplier,
		suite.LatticeHeight,
		wavepacket.Pm,
	)

	checker.Assert(err, IsNil)
	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 2)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms, HasLen, 1)

	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM * -1)

	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pm), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pg), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularCreatedWithDesiredSymmetry) TestCreateWallpaperWithPg(checker *C) {
	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaWithSymmetry(
		suite.eisensteinTermWithOddPowerNAndEvenPowerSum,
		suite.wallpaperMultiplier,
		suite.LatticeHeight,
		wavepacket.Pg,
	)

	checker.Assert(err, IsNil)
	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 2)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms, HasLen, 1)

	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM * -1)

	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pm), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularCreatedWithDesiredSymmetry) TestCreateWallpaperWithPmm(checker *C) {
	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaWithSymmetry(
		suite.eisensteinTermWithOddPowerNAndEvenPowerSum,
		suite.wallpaperMultiplier,
		suite.LatticeHeight,
		wavepacket.Pmm,
	)

	checker.Assert(err, IsNil)
	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 4)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms, HasLen, 1)

	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM * -1)

	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier)
	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM)

	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier)
	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN)
	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM * -1)

	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pm), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pg), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, true)
}

func (suite *RectangularCreatedWithDesiredSymmetry) TestCreateWallpaperWithPmg(checker *C) {
	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaWithSymmetry(
		suite.eisensteinTermWithOddPowerNAndEvenPowerSum,
		suite.wallpaperMultiplier,
		suite.LatticeHeight,
		wavepacket.Pmg,
	)

	checker.Assert(err, IsNil)
	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 4)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms, HasLen, 1)

	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM * -1)

	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM)

	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN)
	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM * -1)

	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pm), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularCreatedWithDesiredSymmetry) TestCreateWallpaperWithPgg(checker *C) {
	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaWithSymmetry(
		suite.eisensteinTermWithOddPowerNAndEvenPowerSum,
		suite.wallpaperMultiplier,
		suite.LatticeHeight,
		wavepacket.Pgg,
	)

	checker.Assert(err, IsNil)
	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 4)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms, HasLen, 1)

	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM * -1)

	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier)
	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN * -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM)

	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].Multiplier, Equals, suite.wallpaperMultiplier)
	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].PowerN, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerN)
	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].PowerM, Equals, suite.eisensteinTermWithOddPowerNAndEvenPowerSum[0].PowerM * -1)

	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pm), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pg), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, true)
}

func (suite *RectangularCreatedWithDesiredSymmetry) TestCreateDesiredSymmetryWithYAML(checker *C) {
	yamlByteStream := []byte(`
lattice_height: 0.3
formula:
 desired_symmetry: pmm
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
         power_m: -9
`)

	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaFromYAML(yamlByteStream)
	checker.Assert(err, IsNil)

	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 4)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms[0].PowerM, Equals, -9)

	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, -12)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, 9)

	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pm), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, false)
}

func (suite *RectangularCreatedWithDesiredSymmetry) TestCreateDesiredSymmetryWithJSON(checker *C) {
	jsonByteStream := []byte(`{
				"lattice_height": 0.3,
				"formula": {
					"desired_symmetry": "pmg",
					"multiplier": {
						"real": 1.0,
						"imaginary": 0
					},
					"wave_packets": [
						{
							"multiplier": {
								"real": 1.0,
								"imaginary": 0
							},
							"terms": [
								{
									"power_n": 1,
									"power_m": -2
								}
							]
						}
					]
				}
			}`)
	RectangularFormula, err := wavepacket.NewRectangularWallpaperFormulaFromJSON(jsonByteStream)
	checker.Assert(err, IsNil)

	checker.Assert(RectangularFormula.Formula.WavePackets, HasLen, 4)
	checker.Assert(RectangularFormula.Formula.WavePackets[0].Terms[0].PowerM, Equals, -2)

	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].Multiplier, Equals, complex(1, 0))
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerN, Equals, -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[1].Terms[0].PowerM, Equals, 2)

	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].Multiplier, Equals, complex(-1, 0))
	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].PowerN, Equals, -1)
	checker.Assert(RectangularFormula.Formula.WavePackets[2].Terms[0].PowerM, Equals, -2)

	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].Multiplier, Equals, complex(-1, 0))
	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].PowerN, Equals, 1)
	checker.Assert(RectangularFormula.Formula.WavePackets[3].Terms[0].PowerM, Equals, 2)

	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pm), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pg), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmm), Equals, false)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pmg), Equals, true)
	checker.Assert(RectangularFormula.HasSymmetry(wavepacket.Pgg), Equals, true)
}