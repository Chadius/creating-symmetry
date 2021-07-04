package coefficient

// Pairing notes the multiplier and the powers applied to a formula term.
type Pairing struct {
	PowerN				int
	PowerM				int
	NegateMultiplier	bool
}

// GenerateCoefficientSets creates a list of locked coefficient sets (powers and multipliers)
//  based on a given list of relationships.
func (pairing Pairing) GenerateCoefficientSets(relationships []Relationship) []*Pairing {
	pairs := []*Pairing{}

	negateMultiplierIfPowerNIsOdd := pairing.PowerN % 2 != 0
	negateMultiplierIfSumIsOdd := (pairing.PowerN + pairing.PowerM) % 2 != 0

	pairingByRelationship := map[Relationship]*Pairing{
		PlusNPlusM: {
			PowerN: pairing.PowerN,
			PowerM: pairing.PowerM,
			NegateMultiplier: false,
		},
		PlusMPlusN: {
			PowerN: pairing.PowerM,
			PowerM: pairing.PowerN,
			NegateMultiplier: false,
		},
		PlusMPlusNMaybeFlipScale: {
			PowerN:           pairing.PowerM,
			PowerM:           pairing.PowerN,
			NegateMultiplier: negateMultiplierIfSumIsOdd,
		},
		MinusNMinusM: {
			PowerN:     -1 * pairing.PowerN,
			PowerM:     -1 * pairing.PowerM,
			NegateMultiplier: false,
		},
		MinusMMinusN: {
			PowerN:     -1 * pairing.PowerM,
			PowerM:     -1 * pairing.PowerN,
			NegateMultiplier: false,
		},
		MinusMMinusNMaybeFlipScale: {
			PowerN:           -1 * pairing.PowerM,
			PowerM:           -1 * pairing.PowerN,
			NegateMultiplier: negateMultiplierIfSumIsOdd,
		},
		PlusMMinusSumNAndM: {
			PowerN: pairing.PowerM,
			PowerM: -1 * (pairing.PowerN + pairing.PowerM),
			NegateMultiplier: false,
		},
		MinusSumNAndMPlusN: {
			PowerN: -1 * (pairing.PowerN + pairing.PowerM),
			PowerM: pairing.PowerM,
			NegateMultiplier: false,
		},
		PlusNMinusM: {
			PowerN: pairing.PowerN,
			PowerM: -1 * pairing.PowerM,
			NegateMultiplier: false,
		},
		PlusNMinusMNegateMultiplierIfOddPowerN: {
			PowerN: pairing.PowerN,
			PowerM: -1 * pairing.PowerM,
			NegateMultiplier: negateMultiplierIfPowerNIsOdd,
		},
		PlusNMinusMNegateMultiplierIfOddPowerSum: {
			PowerN: pairing.PowerN,
			PowerM: -1 * pairing.PowerM,
			NegateMultiplier: negateMultiplierIfSumIsOdd,
		},
		MinusNPlusMNegateMultiplierIfOddPowerN: {
			PowerN: -1 * pairing.PowerN,
			PowerM: pairing.PowerM,
			NegateMultiplier: negateMultiplierIfPowerNIsOdd,
		},
		MinusNPlusMNegateMultiplierIfOddPowerSum: {
			PowerN: -1 * pairing.PowerN,
			PowerM: pairing.PowerM,
			NegateMultiplier: negateMultiplierIfSumIsOdd,
		},
		PlusMMinusN: {
			PowerN: pairing.PowerM,
			PowerM: -1 * pairing.PowerN,
			NegateMultiplier: false,
		},
		MinusMPlusN: {
			PowerN: -1 * pairing.PowerM,
			PowerM: pairing.PowerN,
			NegateMultiplier: false,
		},
		MinusNPlusM: {
			PowerN:     -1 * pairing.PowerN,
			PowerM:     pairing.PowerM,
			NegateMultiplier: false,
		},
	}

	for _, relationship := range relationships {
		pairWithSameRelationship := pairingByRelationship[relationship]
		newPair := &Pairing{
			PowerN: pairWithSameRelationship.PowerN,
			PowerM: pairWithSameRelationship.PowerM,
			NegateMultiplier: pairWithSameRelationship.NegateMultiplier,
		}
		pairs = append(pairs, newPair)
	}

	return pairs
}

// Relationship relates how a pair of coordinates should be applied.
type Relationship string

// Relationship s determine the order and sign of powers n and m.
//   Plus means *1, Minus means *-1
//   If N appears first the powers then power N is applied to the number and power M to the complex conjugate.
//   If M appears first the powers then power M is applied to the number and power N to the complex conjugate.
//	 MaybeFlipScale will multiply the scale by -1 if N + M is odd.
const (
	PlusNPlusM                 Relationship = "+N+M"
	PlusMPlusN                             Relationship = "+M+N"
	MinusNMinusM                           Relationship = "-N-M"
	MinusMMinusN                           Relationship = "-M-N"
	PlusMPlusNMaybeFlipScale               Relationship = "+M+NF"
	MinusMMinusNMaybeFlipScale             Relationship = "-M-NF"
	PlusMMinusSumNAndM                     Relationship = "+M-(N+M)"
	MinusSumNAndMPlusN                     Relationship = "-(N+M)+N"
	PlusMMinusN                            Relationship = "+M-N"
	MinusMPlusN                            Relationship = "-M+N"
	PlusNMinusM                            Relationship = "+N-M"
	PlusNMinusMNegateMultiplierIfOddPowerN Relationship = "+N-MF(N)"
	MinusNPlusMNegateMultiplierIfOddPowerN Relationship = "-N+MF(N)"
	MinusNPlusM Relationship = "-N+M"
	PlusNMinusMNegateMultiplierIfOddPowerSum Relationship = "+N-MF(N+M)"
	MinusNPlusMNegateMultiplierIfOddPowerSum Relationship = "-N+MF(N+M)"
)

