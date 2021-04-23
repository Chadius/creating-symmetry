package wave

import "wallpaper/entities/formula"

// Formula for Waves mathematically creates repeating, cyclical mathematical patterns
//   in 2D space, similar to waves on the ocean.
type Formula struct {
	Terms 			[]*formula.EisensteinFormulaTerm
	Multiplier 		complex128
}

// Calculate takes the complex number z and processes it using the mathematical terms.
func (waveFormula Formula) Calculate(z complex128) *formula.CalculationResultForFormula {
	result := &formula.CalculationResultForFormula{
		Total: complex(0,0),
		ContributionByTerm: []complex128{},
	}

	for _, term := range waveFormula.Terms {
		termContribution := term.Calculate(z)
		result.Total += termContribution
		result.ContributionByTerm = append(result.ContributionByTerm, termContribution)
	}
	result.Total *= waveFormula.Multiplier

	return result
}
