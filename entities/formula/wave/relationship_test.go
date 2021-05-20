package wave_test

import (
	. "gopkg.in/check.v1"
	"math"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/wave"
)

type WaveSymmetryFormulaTests struct {
	baseOddWave *wave.Formula
}

var _ = Suite(&WaveSymmetryFormulaTests{})

func (suite *WaveSymmetryFormulaTests) SetUpTest(checker *C) {
	suite.baseOddWave = &wave.Formula{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				XLatticeVector: complex(1,0),
				YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
				PowerN:         1,
				PowerM:         -2,
			},
		},
		Multiplier: complex(1, 0),
	}
}

func (suite *WaveSymmetryFormulaTests) TestNeedTwoFormulasToDetectSymmetry(checker *C) {
	relationship := wave.FindWaveRelationships([]*wave.Formula{
		suite.baseOddWave,
	})
	checker.Assert(relationship.NoRelation, Equals, true)
}

func (suite *WaveSymmetryFormulaTests) TestPlusMPlusN(checker *C) {
	relationship := wave.FindWaveRelationships([]*wave.Formula{
		suite.baseOddWave,
		{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					XLatticeVector: complex(1,0),
					YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
					PowerN:         suite.baseOddWave.Terms[0].PowerM,
					PowerM:         suite.baseOddWave.Terms[0].PowerN,
				},
			},
			Multiplier: complex(1, 0),
		},
	})
	checker.Assert(relationship.PlusMPlusN, Equals, true)
}

func (suite *WaveSymmetryFormulaTests) TestMinusNMinusM(checker *C) {
	relationship := wave.FindWaveRelationships([]*wave.Formula{
		suite.baseOddWave,
		{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					XLatticeVector: complex(1,0),
					YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
					PowerN:         suite.baseOddWave.Terms[0].PowerN * -1,
					PowerM:         suite.baseOddWave.Terms[0].PowerM * -1,
				},
			},
			Multiplier: complex(1, 0),
		},
	})
	checker.Assert(relationship.MinusNMinusM, Equals, true)
}

func (suite *WaveSymmetryFormulaTests) TestMinusMMinusN(checker *C) {
	relationship := wave.FindWaveRelationships([]*wave.Formula{
		suite.baseOddWave,
		{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					XLatticeVector: complex(1,0),
					YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
					PowerN:         suite.baseOddWave.Terms[0].PowerM * -1,
					PowerM:         suite.baseOddWave.Terms[0].PowerN * -1,
				},
			},
			Multiplier: complex(1, 0),
		},
	})
	checker.Assert(relationship.MinusMMinusN, Equals, true)
}

func (suite *WaveSymmetryFormulaTests) TestMinusNMinusMPlusMPlusNMinusMMinusN(checker *C) {
	relationship := wave.FindWaveRelationships([]*wave.Formula{
		suite.baseOddWave,
		{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					XLatticeVector: complex(1,0),
					YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
					PowerN:         suite.baseOddWave.Terms[0].PowerN * -1,
					PowerM:         suite.baseOddWave.Terms[0].PowerM * -1,
				},
			},
			Multiplier: complex(1, 0),
		},
		{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					XLatticeVector: complex(1,0),
					YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
					PowerN:         suite.baseOddWave.Terms[0].PowerM * -1,
					PowerM:         suite.baseOddWave.Terms[0].PowerN * -1,
				},
			},
			Multiplier: complex(1, 0),
		},
		{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					XLatticeVector: complex(1,0),
					YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
					PowerN:         suite.baseOddWave.Terms[0].PowerM,
					PowerM:         suite.baseOddWave.Terms[0].PowerN,
				},
			},
			Multiplier: complex(1, 0),
		},
	})
	checker.Assert(relationship.MinusNMinusMPlusMPlusNMinusMMinusN, Equals, true)
}

func (suite *WaveSymmetryFormulaTests) TestMultipleSymmetries(checker *C) {
	relationship := wave.FindWaveRelationships([]*wave.Formula{
		{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					XLatticeVector: complex(1,0),
					YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
					PowerN:         1,
					PowerM:         -1,
				},
			},
			Multiplier: complex(1, 0),
		},
		{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					XLatticeVector: complex(1,0),
					YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
					PowerN:         -1,
					PowerM:         1,
				},
			},
			Multiplier: complex(1, 0),
		},
	})
	checker.Assert(relationship.PlusMPlusN, Equals, true)
	checker.Assert(relationship.MinusNMinusM, Equals, true)
}