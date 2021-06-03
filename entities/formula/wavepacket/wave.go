package wavepacket

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/coefficient"
	"wallpaper/entities/utility"
)

// Marshal can be marshaled and converted to a EisensteinFormulaTerm
type Marshal struct {
	Terms []*formula.EisensteinFormulaTermMarshal `json:"terms" yaml:"terms"`
	Multiplier utility.ComplexNumberForMarshal    `json:"multiplier" yaml:"multiplier"`
}

// WavePacket for Waves mathematically creates repeating, cyclical mathematical patterns
//   in 2D space, similar to waves on the ocean.
type WavePacket struct {
	Terms 			[]*formula.EisensteinFormulaTerm
	Multiplier 		complex128
}

// Calculate takes the complex number zInLatticeCoordinates and processes it using the mathematical terms.
func (waveFormula WavePacket) Calculate(zInLatticeCoordinates complex128) *formula.CalculationResultForFormula {
	result := &formula.CalculationResultForFormula{
		Total: complex(0,0),
		ContributionByTerm: []complex128{},
	}

	for _, term := range waveFormula.Terms {
		termContribution := term.Calculate(zInLatticeCoordinates)
		result.Total += termContribution
		result.ContributionByTerm = append(result.ContributionByTerm, termContribution)
	}
	result.Total *= waveFormula.Multiplier

	return result
}

// NewWaveFormulaFromJSON reads the data and returns a formula term from it.
func NewWaveFormulaFromJSON(data []byte) (*WavePacket, error) {
	return newWaveFormulaFromDatastream(data, json.Unmarshal)
}

// NewWaveFormulaFromYAML reads the data and returns a formula term from it.
func NewWaveFormulaFromYAML(data []byte) (*WavePacket, error) {
	return newWaveFormulaFromDatastream(data, yaml.Unmarshal)
}

//newWaveFormulaFromDatastream consumes a given bytestream and tries to create a new object from it.
func newWaveFormulaFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*WavePacket, error) {
	var unmarshalError error
	var formulaMarshal Marshal
	unmarshalError = unmarshal(data, &formulaMarshal)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	formulaTerm := NewWaveFormulaFromMarshalObject(formulaMarshal)
	return formulaTerm, nil
}

// NewWaveFormulaFromMarshalObject converts the marshaled intermediary object into a usable object.
func NewWaveFormulaFromMarshalObject(marshalObject Marshal) *WavePacket {
	formulaTerms := []*formula.EisensteinFormulaTerm{}
	for _,term := range marshalObject.Terms {
		newEisenstein := formula.NewEisensteinFormulaTermFromMarshalObject(*term)
		formulaTerms = append(formulaTerms, newEisenstein)
	}

	return &WavePacket{
		Terms: 		formulaTerms,
		Multiplier:	complex(marshalObject.Multiplier.Real, marshalObject.Multiplier.Imaginary),
	}
}

// Symmetry tracks the various ways a Wave Pattern formula can have symmetry.
type Symmetry struct {
	P3 bool
	P31m bool
	P3m1 bool
	P6 bool
	P6m bool
	P4 bool
	P4m bool
	P4g bool
}

// GetWavePacketRelationship returns a list of relationships that all of the wave packets conform to.
func GetWavePacketRelationship(wavePacket1, wavePacket2 *WavePacket) []coefficient.Relationship {
	if wavePacket1 == nil || wavePacket2 == nil {
		return []coefficient.Relationship{}
	}

	return formula.GetAllPossibleTermRelationships(wavePacket1.Terms[0], wavePacket2.Terms[0])
}

// ContainsRelationship returns true if the target relationship is in the list of relationships.
func ContainsRelationship(relationships []coefficient.Relationship, target coefficient.Relationship) bool {
	for _, thisRelationship := range relationships {
		if target == thisRelationship {
			return true
		}
	}
	return false
}

// CanWavePacketsBeGroupedAmongCoefficientRelationships returns true if the WavePackets involved satisfy the relationships.
func CanWavePacketsBeGroupedAmongCoefficientRelationships(wavePackets []*WavePacket, desiredRelationships []coefficient.Relationship) bool {
	wavePacketsMatched := []bool{}
	for range wavePackets {
		wavePacketsMatched = append(wavePacketsMatched, false)
	}

	for indexA, wavePacketA := range wavePackets {
		relationshipWasFound := map[coefficient.Relationship]bool{}
		for _, r := range desiredRelationships {
			relationshipWasFound[r] = false
		}

		if wavePacketsMatched[indexA] == true {
			continue
		}

		for offsetB, wavePacketB := range wavePackets[indexA+1:] {
			relationshipsFound := GetWavePacketRelationship(
				wavePacketA,
				wavePacketB,
			)

			for _, relationshipToLookFor := range desiredRelationships {
				if ContainsRelationship(relationshipsFound, relationshipToLookFor) {
					wavePacketsMatched[indexA + offsetB + 1] = true
					relationshipWasFound[relationshipToLookFor] = true
					break
				}
			}
		}

		for _, relationshipFound := range relationshipWasFound {
			if relationshipFound != true {
				return false
			}
		}
		wavePacketsMatched[indexA] = true
	}

	return true
}
