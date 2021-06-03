package wavepacket

import "wallpaper/entities/formula"

func addNewWavePacketsBasedOnSymmetry(term *formula.EisensteinFormulaTerm, desiredSymmetry *Symmetry, newWavePackets []*WavePacket) []*WavePacket {
	powerN := term.PowerN
	powerM := term.PowerM

	if desiredSymmetry.P31m || desiredSymmetry.P4m {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerM: powerN,
					PowerN: powerM,
				},
			},
			Multiplier: term.Multiplier,
		})
	}
	if desiredSymmetry.P3m1 {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerM: powerN * -1,
					PowerN: powerM * -1,
				},
			},
			Multiplier: term.Multiplier,
		})
	}
	if desiredSymmetry.P6 {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerM: powerM * -1,
					PowerN: powerN * -1,
				},
			},
			Multiplier: term.Multiplier,
		})
	}
	if desiredSymmetry.P6m {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerM: powerM * -1,
					PowerN: powerN * -1,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerM: powerN,
					PowerN: powerM,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerM: powerN * -1,
					PowerN: powerM * -1,
				},
			},
			Multiplier: term.Multiplier,
		})
	}

	powerSumIsEven := (powerN + powerM) % 2 == 0
	if desiredSymmetry.P4g {
		multiplierMaybeNegatedBasedOnSum := term.Multiplier
		if !powerSumIsEven {
			multiplierMaybeNegatedBasedOnSum *= -1
		}

		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerM: powerN,
					PowerN: powerM,
					Multiplier: multiplierMaybeNegatedBasedOnSum,
				},
			},
			Multiplier: term.Multiplier,
		})
	}

	return newWavePackets
}

// SymmetryType encodes all possible symmetries for wallpaper patterns.
type SymmetryType string

// All possible symmetries for wallpaper patterns, based on crystallography.
const (
	P3 SymmetryType	= "p3"
	P3m1 SymmetryType	= "p3m1"
	P31m SymmetryType	= "p31m"
	P6 SymmetryType	= "p6"
	P6m SymmetryType	= "p6m"
	P4 SymmetryType = "p4"
	P4m SymmetryType = "p4m"
	P4g SymmetryType = "p4g"
)