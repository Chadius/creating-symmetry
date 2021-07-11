package wavepacket

import "wallpaper/entities/formula"

func addNewWavePacketsBasedOnSymmetry(term *formula.EisensteinFormulaTerm, desiredSymmetry Symmetry, newWavePackets []*WavePacket) []*WavePacket {
	powerN := term.PowerN
	powerM := term.PowerM
	powerNIsEven := powerN % 2 == 0
	powerSumIsEven := (powerN + powerM) % 2 == 0

	multiplierMaybeNegatedBasedOnSum := term.Multiplier
	if !powerSumIsEven {
		multiplierMaybeNegatedBasedOnSum *= -1
	}

	multiplierMaybeNegatedBasedOnPowerN := term.Multiplier
	if !powerNIsEven {
		multiplierMaybeNegatedBasedOnPowerN *= -1
	}

	if desiredSymmetry == P31m || desiredSymmetry == P4m || desiredSymmetry == Cm {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerM,
					PowerM: powerN,
				},
			},
			Multiplier: term.Multiplier,
		})
	}
	if desiredSymmetry == Pm {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN,
					PowerM: powerM * -1,
				},
			},
			Multiplier: term.Multiplier,
		})
	}
	if desiredSymmetry == Pg {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN,
					PowerM: powerM * -1,
					Multiplier: multiplierMaybeNegatedBasedOnPowerN,
				},
			},
			Multiplier: term.Multiplier,
		})
	}
	if desiredSymmetry == Pmm {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN * -1,
					PowerM: powerM * -1,
					Multiplier: term.Multiplier,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN * -1,
					PowerM: powerM,
					Multiplier: term.Multiplier,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN,
					PowerM: powerM * -1,
					Multiplier: term.Multiplier,
				},
			},
			Multiplier: term.Multiplier,
		})
	}
	if desiredSymmetry == Pmg {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN * -1,
					PowerM: powerM * -1,
					Multiplier: term.Multiplier,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN * -1,
					PowerM: powerM,
					Multiplier: multiplierMaybeNegatedBasedOnPowerN,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN,
					PowerM: powerM * -1,
					Multiplier: multiplierMaybeNegatedBasedOnPowerN,
				},
			},
			Multiplier: term.Multiplier,
		})
	}
	if desiredSymmetry == Pgg {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN * -1,
					PowerM: powerM * -1,
					Multiplier: term.Multiplier,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN * -1,
					PowerM: powerM,
					Multiplier: multiplierMaybeNegatedBasedOnSum,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN,
					PowerM: powerM * -1,
					Multiplier: multiplierMaybeNegatedBasedOnSum,
				},
			},
			Multiplier: term.Multiplier,
		})
	}

	if desiredSymmetry == P3m1 {
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
	if desiredSymmetry == P6 {
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
	if desiredSymmetry == P6m || desiredSymmetry == Cmm {
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerN * -1,
					PowerM: powerM * -1,
					Multiplier: term.Multiplier,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerM,
					PowerM: powerN,
					Multiplier: term.Multiplier,
				},
			},
			Multiplier: term.Multiplier,
		})
		newWavePackets = append(newWavePackets, &WavePacket{
			Terms: []*formula.EisensteinFormulaTerm{
				{
					PowerN: powerM * -1,
					PowerM: powerN * -1,
					Multiplier: term.Multiplier,
				},
			},
			Multiplier: term.Multiplier,
		})
	}

	if desiredSymmetry == P4g {
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

// Symmetry encodes all possible symmetries for wallpaper patterns.
type Symmetry string

// All possible symmetries for wallpaper patterns, based on crystallography.
const (
	P3   Symmetry = "p3"
	P3m1 Symmetry = "p3m1"
	P31m Symmetry = "p31m"
	P6   Symmetry = "p6"
	P6m  Symmetry = "p6m"
	P4   Symmetry = "p4"
	P4m  Symmetry = "p4m"
	P4g  Symmetry = "p4g"
	Cm   Symmetry = "cm"
	Cmm  Symmetry = "cmm"
	Pm   Symmetry = "pm"
	Pg   Symmetry = "pg"
	Pgg   Symmetry = "pgg"
	Pmm   Symmetry = "pmm"
	Pmg   Symmetry = "pmg"
)