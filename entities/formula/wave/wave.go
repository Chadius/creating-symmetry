package wave

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"wallpaper/entities/formula"
	"wallpaper/entities/utility"
)

// FormulaMarshalable can be marshaled and converted to a EisensteinFormulaTerm
type FormulaMarshalable struct {
	Terms []*formula.EisensteinFormulaTermMarshalable `json:"terms" yaml:"terms"`
	Multiplier utility.ComplexNumberForMarshal	`json:"multiplier" yaml:"multiplier"`
}

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

// NewWaveFormulaFromJSON reads the data and returns a formula term from it.
func NewWaveFormulaFromJSON(data []byte) (*Formula, error) {
	return newWaveFormulaFromDatastream(data, json.Unmarshal)
}

// NewWaveFormulaFromYAML reads the data and returns a formula term from it.
func NewWaveFormulaFromYAML(data []byte) (*Formula, error) {
	return newWaveFormulaFromDatastream(data, yaml.Unmarshal)
}

//newWaveFormulaFromDatastream consumes a given bytestream and tries to create a new object from it.
func newWaveFormulaFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*Formula, error) {
	var unmarshalError error
	var formulaMarshal FormulaMarshalable
	unmarshalError = unmarshal(data, &formulaMarshal)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	formulaTerm := NewWaveFormulaFromMarshalObject(formulaMarshal)
	return formulaTerm, nil
}

// NewWaveFormulaFromMarshalObject converts the marshaled intermediary object into a usable object.
func NewWaveFormulaFromMarshalObject(marshalObject FormulaMarshalable) *Formula {
	formulaTerms := []*formula.EisensteinFormulaTerm{}
	for _,term := range marshalObject.Terms {
		newEisenstein := formula.NewEisensteinFormulaTermFromMarshalObject(*term)
		formulaTerms = append(formulaTerms, newEisenstein)
	}

	return &Formula{
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
}