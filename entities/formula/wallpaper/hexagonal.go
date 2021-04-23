package wallpaper

import (
	"math"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/wave"
)

type PowerTerm struct {
	PowerN			int
	PowerM			int
	Multiplier		complex128
}

type HexagonalFormula struct {
	PowerTerms  []*PowerTerm
	WavePackets []*wave.Formula
}

// Calculate uses the complex number z and applies it to the WavePackets in this pattern.
func (hexagonalFormula *HexagonalFormula) Calculate(z complex128) *formula.CalculationResultForFormula {
	result := &formula.CalculationResultForFormula{
		Total: complex(0,0),
		ContributionByTerm: []complex128{},
	}

	for _, wavePacket := range hexagonalFormula.WavePackets {
		packetResult := wavePacket.Calculate(z)
		result.Total += packetResult.Total
		result.ContributionByTerm = append(result.ContributionByTerm, packetResult.Total)
	}

	return result
}

// Setup uses the defined PowerTerm objects to create the needed helper objects.
func (hexagonalFormula *HexagonalFormula) Setup() {
	hexagonalFormula.WavePackets = []*wave.Formula{}

	for _, powerTerm := range hexagonalFormula.PowerTerms{
		waveFormula := hexagonalFormula.createWavePacket(powerTerm)
		hexagonalFormula.WavePackets = append(hexagonalFormula.WavePackets, waveFormula)
	}
}

// createWavePacket uses a term to create the mathematical formulas needed to create
//   3 fold symmetrical patterns.
func (hexagonalFormula *HexagonalFormula) createWavePacket(term *PowerTerm) *wave.Formula {
	xLatticeVector := complex(1,0)
	yLatticeVector := complex(-0.5, math.Sqrt(3.0)/2.0)

	return &wave.Formula{
		Terms: []*formula.EisensteinFormulaTerm{
			{
				XLatticeVector: xLatticeVector,
				YLatticeVector: yLatticeVector,
				PowerN:         term.PowerN,
				PowerM:         term.PowerM,
			},
			{
				XLatticeVector: xLatticeVector,
				YLatticeVector: yLatticeVector,
				PowerN:         term.PowerM,
				PowerM:         -1 * (term.PowerM + term.PowerN),
			},
			{
				XLatticeVector: xLatticeVector,
				YLatticeVector: yLatticeVector,
				PowerN:         -1 * (term.PowerN + term.PowerM),
				PowerM:         term.PowerN,
			},
		},
		Multiplier: complex(1/3.0, 0) * term.Multiplier,
	}
}