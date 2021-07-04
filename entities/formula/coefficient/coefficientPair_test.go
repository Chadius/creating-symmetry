package coefficient_test

import (
	. "gopkg.in/check.v1"
	"testing"
	"wallpaper/entities/formula/coefficient"
)

func Test(t *testing.T) { TestingT(t) }

type CoefficientPairFeatures struct {
	evenSumPair *coefficient.Pairing
	oddSumPair *coefficient.Pairing
}

var _ = Suite(&CoefficientPairFeatures{})

func (suite *CoefficientPairFeatures) SetUpTest(checker *C) {
	suite.evenSumPair = &coefficient.Pairing{
		PowerN: 1,
		PowerM: 3,
	}

	suite.oddSumPair = &coefficient.Pairing{
		PowerN: 1,
		PowerM: 2,
	}
}

func (suite *CoefficientPairFeatures) TestPlusNPlusM(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusNPlusM,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, 1)
	checker.Assert(newSets[0].PowerM, Equals, 3)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestPlusMPlusN(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusMPlusN,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, 3)
	checker.Assert(newSets[0].PowerM, Equals, 1)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestReturnsMultiples(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusNPlusM,
		coefficient.PlusMPlusN,
	})

	checker.Assert(newSets, HasLen, 2)
	checker.Assert(newSets[0].PowerN, Equals, 1)
	checker.Assert(newSets[0].PowerM, Equals, 3)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)

	checker.Assert(newSets[1].PowerN, Equals, 3)
	checker.Assert(newSets[1].PowerM, Equals, 1)
	checker.Assert(newSets[1].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestPlusMPlusNMaybeFlipScale(checker *C) {
	newSets := suite.oddSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusMPlusNMaybeFlipScale,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, 2)
	checker.Assert(newSets[0].PowerM, Equals, 1)
	checker.Assert(newSets[0].NegateMultiplier, Equals, true)
}

func (suite *CoefficientPairFeatures) TestMinusNMinusM(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusNMinusM,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, -1)
	checker.Assert(newSets[0].PowerM, Equals, -3)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestMinusMMinusN(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusMMinusN,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, -3)
	checker.Assert(newSets[0].PowerM, Equals, -1)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestMinusMMinusNMaybeFlipScale(checker *C) {
	newSets := suite.oddSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusMMinusNMaybeFlipScale,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, -2)
	checker.Assert(newSets[0].PowerM, Equals, -1)
	checker.Assert(newSets[0].NegateMultiplier, Equals, true)
}

func (suite *CoefficientPairFeatures) TestPlusMMinusSumNAndM(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusMMinusSumNAndM,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, 3)
	checker.Assert(newSets[0].PowerM, Equals, -(1 + 3))
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestMinusSumNAndMPlusN(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusSumNAndMPlusN,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, -(1+3))
	checker.Assert(newSets[0].PowerM, Equals, 3)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestPlusMMinusN(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusMMinusN,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, 3)
	checker.Assert(newSets[0].PowerM, Equals, -1)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestMinusMPlusN(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusMPlusN,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, -3)
	checker.Assert(newSets[0].PowerM, Equals, 1)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestPlusNMinusM(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusNMinusM,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, 1)
	checker.Assert(newSets[0].PowerM, Equals, -3)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestPlusNMinusMNegateMultiplierIfOddPowerN(checker *C) {
	newSetsWithOddPowerN := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusNMinusMNegateMultiplierIfOddPowerN,
	})

	checker.Assert(newSetsWithOddPowerN, HasLen, 1)
	checker.Assert(newSetsWithOddPowerN[0].PowerN, Equals, 1)
	checker.Assert(newSetsWithOddPowerN[0].PowerM, Equals, -3)
	checker.Assert(newSetsWithOddPowerN[0].NegateMultiplier, Equals, true)

	coefficientPairWithEvenPowerN := coefficient.Pairing{
		PowerN: 2,
		PowerM: 3,
	}
	newSetsWithEvenPowerN := coefficientPairWithEvenPowerN.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusNMinusMNegateMultiplierIfOddPowerN,
	})

	checker.Assert(newSetsWithEvenPowerN, HasLen, 1)
	checker.Assert(newSetsWithEvenPowerN[0].PowerN, Equals, 2)
	checker.Assert(newSetsWithEvenPowerN[0].PowerM, Equals, -3)
	checker.Assert(newSetsWithEvenPowerN[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestMinusNPlusMNegateMultiplierIfOddPowerN(checker *C) {
	newSetsWithOddPowerN := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusNPlusMNegateMultiplierIfOddPowerN,
	})

	checker.Assert(newSetsWithOddPowerN, HasLen, 1)
	checker.Assert(newSetsWithOddPowerN[0].PowerN, Equals, -1)
	checker.Assert(newSetsWithOddPowerN[0].PowerM, Equals, 3)
	checker.Assert(newSetsWithOddPowerN[0].NegateMultiplier, Equals, true)

	coefficientPairWithEvenPowerN := coefficient.Pairing{
		PowerN: 2,
		PowerM: 3,
	}
	newSetsWithEvenPowerN := coefficientPairWithEvenPowerN.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusNPlusMNegateMultiplierIfOddPowerN,
	})

	checker.Assert(newSetsWithEvenPowerN, HasLen, 1)
	checker.Assert(newSetsWithEvenPowerN[0].PowerN, Equals, -2)
	checker.Assert(newSetsWithEvenPowerN[0].PowerM, Equals, 3)
	checker.Assert(newSetsWithEvenPowerN[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestMinusNPlusM(checker *C) {
	newSets := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusNPlusM,
	})

	checker.Assert(newSets, HasLen, 1)
	checker.Assert(newSets[0].PowerN, Equals, -1)
	checker.Assert(newSets[0].PowerM, Equals, 3)
	checker.Assert(newSets[0].NegateMultiplier, Equals, false)
}

func (suite *CoefficientPairFeatures) TestPlusNMinusMNegateMultiplierIfOddPowerSum(checker *C) {
	newSetsWithEvenSumPower := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusNMinusMNegateMultiplierIfOddPowerSum,
	})

	checker.Assert(newSetsWithEvenSumPower, HasLen, 1)
	checker.Assert(newSetsWithEvenSumPower[0].PowerN, Equals, 1)
	checker.Assert(newSetsWithEvenSumPower[0].PowerM, Equals, -3)
	checker.Assert(newSetsWithEvenSumPower[0].NegateMultiplier, Equals, false)

	newSetsWithOddSumPower := suite.oddSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.PlusNMinusMNegateMultiplierIfOddPowerSum,
	})

	checker.Assert(newSetsWithOddSumPower, HasLen, 1)
	checker.Assert(newSetsWithOddSumPower[0].PowerN, Equals, 1)
	checker.Assert(newSetsWithOddSumPower[0].PowerM, Equals, -2)
	checker.Assert(newSetsWithOddSumPower[0].NegateMultiplier, Equals, true)
}

func (suite *CoefficientPairFeatures) TestMinusNPlusMNegateMultiplierIfOddPowerSum(checker *C) {
	newSetsWithOddPowerN := suite.evenSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusNPlusMNegateMultiplierIfOddPowerSum,
	})

	checker.Assert(newSetsWithOddPowerN, HasLen, 1)
	checker.Assert(newSetsWithOddPowerN[0].PowerN, Equals, -1)
	checker.Assert(newSetsWithOddPowerN[0].PowerM, Equals, 3)
	checker.Assert(newSetsWithOddPowerN[0].NegateMultiplier, Equals, false)

	newSetsWithOddSumPower := suite.oddSumPair.GenerateCoefficientSets([]coefficient.Relationship{
		coefficient.MinusNPlusMNegateMultiplierIfOddPowerSum,
	})

	checker.Assert(newSetsWithOddSumPower, HasLen, 1)
	checker.Assert(newSetsWithOddSumPower[0].PowerN, Equals, -1)
	checker.Assert(newSetsWithOddSumPower[0].PowerM, Equals, 2)
	checker.Assert(newSetsWithOddSumPower[0].NegateMultiplier, Equals, true)
}

