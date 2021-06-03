package formula_test

import (
	. "gopkg.in/check.v1"
	"testing"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/coefficient"
	"wallpaper/entities/utility"
)

func Test(t *testing.T) { TestingT(t) }

type EisensteinFormulaSuite struct {}

var _ = Suite(&EisensteinFormulaSuite{})

func (suite *EisensteinFormulaSuite) SetUpTest(checker *C) {
}

func (suite *EisensteinFormulaSuite) TestCalculateEisensteinTermForGivenPoint(checker *C) {
	squareLatticePair := formula.LatticeVectorPair{
		XLatticeVector: complex(1, 0),
		YLatticeVector: complex(0, 1),
	}

	latticeCoordinate := squareLatticePair.ConvertToLatticeCoordinates(complex(1.0,1.5))

	eisenstein := formula.EisensteinFormulaTerm{
		PowerN: 1,
		PowerM: 1,
		Multiplier: complex(0.5,0),
	}
	calculatedCoordinate := eisenstein.Calculate(latticeCoordinate)

	checker.Assert(real(calculatedCoordinate), utility.NumericallyCloseEnough{}, -0.5, 1e-6)
	checker.Assert(imag(calculatedCoordinate), utility.NumericallyCloseEnough{}, 0.0, 1e-6)
}

func (suite *EisensteinFormulaSuite) TestCreateFormulaWithJSON(checker *C) {
	jsonByteStream := []byte(`{
				"power_n": 12,
				"power_m": -10,
				"multiplier": {
					"real": -1.0,
					"imaginary": 2e-2
				}
			}`)
	term, err := formula.NewEisensteinFormulaTermFromJSON(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(term.PowerN, Equals, 12)
	checker.Assert(term.PowerM, Equals, -10)
	checker.Assert(real(term.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(term.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
}

func (suite *EisensteinFormulaSuite) TestCreateFormulaWithYAML(checker *C) {
	yamlByteStream := []byte(`
power_n: 12
power_m: -10
multiplier:
  real: -1.0
  imaginary: 2e-2
`)
	term, err := formula.NewEisensteinFormulaTermFromYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(term.PowerN, Equals, 12)
	checker.Assert(term.PowerM, Equals, -10)
	checker.Assert(real(term.Multiplier), utility.NumericallyCloseEnough{}, -1.0, 1e-6)
	checker.Assert(imag(term.Multiplier), utility.NumericallyCloseEnough{}, 2e-2, 1e-6)
}

type EisensteinRelationshipTest struct {
	aPlusNPlusMOddTerm *formula.EisensteinFormulaTerm
	aPlusMMinusNOddTerm *formula.EisensteinFormulaTerm
	aPlusMPlusNOddTerm *formula.EisensteinFormulaTerm
	aMinusNMinusMOddTerm *formula.EisensteinFormulaTerm
	aMinusMMinusNOddTerm *formula.EisensteinFormulaTerm
	aMinusMPlusNOddTerm *formula.EisensteinFormulaTerm
	aPlusMPlusNOddNegatedTerm *formula.EisensteinFormulaTerm
	aPlusNPlusMEvenTerm *formula.EisensteinFormulaTerm
	aPlusMPlusNEvenTerm *formula.EisensteinFormulaTerm
	aPlusMPlusNEvenNegatedTerm *formula.EisensteinFormulaTerm
	aMinusMMinusNOddNegatedTerm *formula.EisensteinFormulaTerm
	aMinusMMinusNEvenTerm *formula.EisensteinFormulaTerm
	aMinusMMinusNEvenNegatedTerm *formula.EisensteinFormulaTerm
	aPlusMMinusSumNAndMOddTerm *formula.EisensteinFormulaTerm
	aMinusSumNAndMPlusN *formula.EisensteinFormulaTerm
}

var _ = Suite(&EisensteinRelationshipTest{})

func (suite *EisensteinRelationshipTest) SetUpTest(checker *C) {
	suite.aPlusNPlusMOddTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -1,
		PowerM:     2,
		Multiplier: complex(2, 1),
	}
	suite.aPlusMPlusNOddTerm = &formula.EisensteinFormulaTerm{
		PowerN:     2,
		PowerM:     -1,
		Multiplier: complex(2, 1),
	}
	suite.aMinusNMinusMOddTerm = &formula.EisensteinFormulaTerm{
		PowerN:     1,
		PowerM:     -2,
		Multiplier: complex(2, 1),
	}
	suite.aMinusMMinusNOddTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -2,
		PowerM:     1,
		Multiplier: complex(2, 1),
	}
	suite.aPlusMMinusNOddTerm = &formula.EisensteinFormulaTerm{
		PowerN:     2,
		PowerM:     1,
		Multiplier: complex(2, 1),
	}
	suite.aMinusMPlusNOddTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -2,
		PowerM:     -1,
		Multiplier: complex(2, 1),
	}
	suite.aPlusMPlusNOddNegatedTerm = &formula.EisensteinFormulaTerm{
		PowerN:     2,
		PowerM:     -1,
		Multiplier: complex(-2, -1),
	}
	suite.aMinusMMinusNOddNegatedTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -2,
		PowerM:     1,
		Multiplier: complex(-2, -1),
	}
	suite.aPlusMMinusSumNAndMOddTerm = &formula.EisensteinFormulaTerm{
		PowerN:     2,
		PowerM:     -1,
		Multiplier: complex(2, 1),
	}
	suite.aMinusSumNAndMPlusN = &formula.EisensteinFormulaTerm{
		PowerN:     -1,
		PowerM:     -1,
		Multiplier: complex(2, 1),
	}

	suite.aPlusNPlusMEvenTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -6,
		PowerM:     2,
		Multiplier: complex(2, 1),
	}
	suite.aPlusMPlusNEvenTerm = &formula.EisensteinFormulaTerm{
		PowerN:     2,
		PowerM:     -6,
		Multiplier: complex(2, 1),
	}
	suite.aPlusMPlusNEvenNegatedTerm = &formula.EisensteinFormulaTerm{
		PowerN:     2,
		PowerM:     -6,
		Multiplier: complex(-2, -1),
	}
	suite.aMinusMMinusNEvenTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -2,
		PowerM:     6,
		Multiplier: complex(2, 1),
	}
	suite.aMinusMMinusNEvenNegatedTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -2,
		PowerM:     6,
		Multiplier: complex(-2, -1),
	}
}

func (suite *EisensteinRelationshipTest) TestPlusMPlusN(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddTerm,
		coefficient.PlusMPlusN,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aMinusNMinusMOddTerm,
			suite.aMinusMMinusNOddTerm,
		coefficient.PlusMPlusN,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aMinusMMinusNOddTerm,
		coefficient.PlusMPlusN,
	), Equals, false)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddNegatedTerm,
		coefficient.PlusMPlusN,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestMultipliersMustBeTheSameOrNegated(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			&formula.EisensteinFormulaTerm{
				PowerN: suite.aPlusMPlusNOddTerm.PowerN,
				PowerM: suite.aPlusMPlusNOddTerm.PowerM,
				Multiplier: suite.aPlusMPlusNOddTerm.Multiplier * -2,
			},
		coefficient.PlusMPlusN,
	), Equals, false)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			&formula.EisensteinFormulaTerm{
				PowerN: suite.aPlusMPlusNOddTerm.PowerN,
				PowerM: suite.aPlusMPlusNOddTerm.PowerM,
				Multiplier: suite.aPlusMPlusNOddTerm.Multiplier * 2,
			},
		coefficient.PlusMPlusN,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestMinusNMinusM(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aMinusNMinusMOddTerm,
		coefficient.MinusNMinusM,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddTerm,
		coefficient.MinusNMinusM,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestPlusNPlusM(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusNPlusMOddTerm,
		coefficient.PlusNPlusM,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aMinusMMinusNOddTerm,
		coefficient.PlusNPlusM,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestMinusMMinusN(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aMinusMMinusNOddTerm,
		coefficient.MinusMMinusN,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddTerm,
		coefficient.MinusMMinusN,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestPlusMPlusNMaybeFlipScale(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddNegatedTerm,
		coefficient.PlusMPlusNMaybeFlipScale,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddTerm,
		coefficient.PlusMPlusNMaybeFlipScale,
	), Equals, false)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMEvenTerm,
			suite.aPlusMPlusNEvenTerm,
		coefficient.PlusMPlusNMaybeFlipScale,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMEvenTerm,
			suite.aPlusMPlusNEvenNegatedTerm,
		coefficient.PlusMPlusNMaybeFlipScale,
	), Equals, false)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddTerm,
		coefficient.PlusMPlusNMaybeFlipScale,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestMinusMMinusNMaybeFlipScale(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aMinusMMinusNOddNegatedTerm,
		coefficient.MinusMMinusNMaybeFlipScale,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aMinusMMinusNOddTerm,
		coefficient.MinusMMinusNMaybeFlipScale,
	), Equals, false)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMEvenTerm,
			suite.aMinusMMinusNEvenTerm,
		coefficient.MinusMMinusNMaybeFlipScale,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMEvenTerm,
			suite.aMinusMMinusNEvenNegatedTerm,
		coefficient.MinusMMinusNMaybeFlipScale,
	), Equals, false)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddTerm,
		coefficient.MinusMMinusNMaybeFlipScale,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestPlusMMinusSumNAndM(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMMinusSumNAndMOddTerm,
		coefficient.PlusMMinusSumNAndM,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusNPlusMOddTerm,
		coefficient.PlusMMinusSumNAndM,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestMinusSumNAndMPlusN(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aMinusSumNAndMPlusN,
		coefficient.MinusSumNAndMPlusN,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusNPlusMOddTerm,
		coefficient.MinusSumNAndMPlusN,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestPlusMMinusN(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMMinusNOddTerm,
		coefficient.PlusMMinusN,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddTerm,
		coefficient.PlusMMinusN,
	), Equals, false)
}

func (suite *EisensteinRelationshipTest) TestMinusMPlusN(checker *C) {
	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aMinusMPlusNOddTerm,
		coefficient.MinusMPlusN,
	), Equals, true)

	checker.Assert(formula.SatisfiesRelationship(
			suite.aPlusNPlusMOddTerm,
			suite.aPlusMPlusNOddTerm,
		coefficient.MinusMPlusN,
	), Equals, false)
}


type EisensteinTermSymmetryTest struct {
	aPlusNPlusMOddEisensteinTerm *formula.EisensteinFormulaTerm
	aPlusMPlusNOddEisensteinTerm *formula.EisensteinFormulaTerm

	aPlusNPlusMEvenEisensteinTerm *formula.EisensteinFormulaTerm
	aMinusMMinusNEvenEisensteinTerm *formula.EisensteinFormulaTerm
}

var _ = Suite(&EisensteinTermSymmetryTest{})

func (suite *EisensteinTermSymmetryTest) SetUpTest(checker *C) {
	suite.aPlusNPlusMOddEisensteinTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -1,
		PowerM:     4,
		Multiplier: complex(2, 1),
	}
	suite.aPlusMPlusNOddEisensteinTerm = &formula.EisensteinFormulaTerm{
		PowerN:     4,
		PowerM:     -1,
		Multiplier: complex(2, 1),
	}


	suite.aPlusNPlusMEvenEisensteinTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -1,
		PowerM:     3,
		Multiplier: complex(2, 1),
	}

	suite.aMinusMMinusNEvenEisensteinTerm = &formula.EisensteinFormulaTerm{
		PowerN:     -3,
		PowerM:     1,
		Multiplier: complex(2, 1),
	}
}

func (suite *EisensteinTermSymmetryTest) TestNoBaseEisensteinTermMeansSymmetryNotFound (checker *C) {
	checker.Assert(formula.SelectTermsToSatisfyRelationships(
		nil,
		[]*formula.EisensteinFormulaTerm{
			suite.aPlusNPlusMOddEisensteinTerm,
		},
		[]coefficient.Relationship{
			coefficient.PlusNPlusM,
		}), HasLen, 0)
}

func (suite *EisensteinTermSymmetryTest) TestNoRelationshipsMeansNoMatchingTermsFound (checker *C) {
	matchingTerms := formula.SelectTermsToSatisfyRelationships(
		suite.aPlusNPlusMOddEisensteinTerm,
		[]*formula.EisensteinFormulaTerm{
			suite.aPlusNPlusMOddEisensteinTerm,
		},
		[]coefficient.Relationship{})
	checker.Assert(matchingTerms, HasLen, 0)
}

func (suite *EisensteinTermSymmetryTest) TestAllOtherTermsSatisfySingleRelationship (checker *C) {
	matchingTerms := formula.SelectTermsToSatisfyRelationships(
		suite.aPlusNPlusMOddEisensteinTerm,
		[]*formula.EisensteinFormulaTerm{
			suite.aPlusNPlusMOddEisensteinTerm,
		},
		[]coefficient.Relationship{
			coefficient.PlusNPlusM,
		})
	checker.Assert(matchingTerms, HasLen, 1)
	checker.Assert(matchingTerms[0], Equals, suite.aPlusNPlusMOddEisensteinTerm)
}

func (suite *EisensteinTermSymmetryTest) TestExtraTermsAreNotNeededToSatisfy (checker *C) {
	matchingTerms := formula.SelectTermsToSatisfyRelationships(
		suite.aPlusNPlusMOddEisensteinTerm,
		[]*formula.EisensteinFormulaTerm{
			suite.aPlusNPlusMOddEisensteinTerm,
			suite.aPlusMPlusNOddEisensteinTerm,
		},
		[]coefficient.Relationship{
			coefficient.PlusMPlusN,
		})
	checker.Assert(matchingTerms, HasLen, 1)
	checker.Assert(matchingTerms[0], Equals, suite.aPlusMPlusNOddEisensteinTerm)
}

func (suite *EisensteinTermSymmetryTest) TestTooManyRelationshipWillNeverSatisfy (checker *C) {
	matchingTerms := formula.SelectTermsToSatisfyRelationships(
		suite.aPlusNPlusMOddEisensteinTerm,
		[]*formula.EisensteinFormulaTerm{
			suite.aPlusNPlusMOddEisensteinTerm,
			suite.aPlusMPlusNOddEisensteinTerm,
		},
		[]coefficient.Relationship{
			coefficient.PlusNPlusM,
			coefficient.PlusNPlusM,
			coefficient.PlusMPlusNMaybeFlipScale,
		})
	checker.Assert(matchingTerms, HasLen, 0)
}

func (suite *EisensteinTermSymmetryTest) TestMultipleRelationshipsCanBeSatisfied (checker *C) {
	matchingTerms := formula.SelectTermsToSatisfyRelationships(
		suite.aPlusNPlusMOddEisensteinTerm,
		[]*formula.EisensteinFormulaTerm{
			suite.aPlusNPlusMOddEisensteinTerm,
			suite.aPlusMPlusNOddEisensteinTerm,
		},
		[]coefficient.Relationship{
			coefficient.PlusNPlusM,
			coefficient.PlusMPlusN,
		})
	checker.Assert(matchingTerms, HasLen, 2)
	checker.Assert(matchingTerms[0], Equals, suite.aPlusNPlusMOddEisensteinTerm)
	checker.Assert(matchingTerms[1], Equals, suite.aPlusMPlusNOddEisensteinTerm)
}

func (suite *EisensteinTermSymmetryTest) TestEachTermCanOnlySatisfyOneRelationship (checker *C) {
	checker.Assert(formula.SelectTermsToSatisfyRelationships(
		suite.aPlusNPlusMEvenEisensteinTerm,
		[]*formula.EisensteinFormulaTerm{
			suite.aMinusMMinusNEvenEisensteinTerm,
		},
		[]coefficient.Relationship{
			coefficient.MinusMMinusN,
		}), HasLen, 1)

	checker.Assert(formula.SelectTermsToSatisfyRelationships(
		suite.aPlusNPlusMEvenEisensteinTerm,
		[]*formula.EisensteinFormulaTerm{
			suite.aMinusMMinusNEvenEisensteinTerm,
		},
		[]coefficient.Relationship{
			coefficient.MinusMMinusNMaybeFlipScale,
		}), HasLen, 1)

	matchingTerms := formula.SelectTermsToSatisfyRelationships(
		suite.aPlusNPlusMEvenEisensteinTerm,
		[]*formula.EisensteinFormulaTerm{
			suite.aMinusMMinusNEvenEisensteinTerm,
		},
		[]coefficient.Relationship{
			coefficient.MinusMMinusN,
			coefficient.MinusMMinusNMaybeFlipScale,
		})
	checker.Assert(matchingTerms, HasLen, 0)
}
