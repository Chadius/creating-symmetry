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

	negateMultiplier := (pairing.PowerN + pairing.PowerM) % 2 != 0

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
			PowerN:     pairing.PowerM,
			PowerM:     pairing.PowerN,
			NegateMultiplier: negateMultiplier,
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
			PowerN:     -1 * pairing.PowerM,
			PowerM:     -1 * pairing.PowerN,
			NegateMultiplier: negateMultiplier,
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
	PlusMPlusN                              = "+M+N"
	MinusNMinusM                            = "-N-M"
	MinusMMinusN                            = "-M-N"
	PlusMPlusNMaybeFlipScale                = "+M+NF"
	MinusMMinusNMaybeFlipScale              = "-M-NF"
	PlusMMinusSumNAndM						= "+M-(N+M)"
	MinusSumNAndMPlusN						= "-(N+M)+N"
)

