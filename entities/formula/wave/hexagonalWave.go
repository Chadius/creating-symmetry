package wave

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"math"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/coefficient"
	"wallpaper/entities/utility"
)

// HexagonalWallpaperFormulaMarshalable is a marshalable object that can be turned into a real object.
type HexagonalWallpaperFormulaMarshalable struct {
	WavePackets []*FormulaMarshalable			`json:"wave_packets" yaml:"wave_packets"`
	Multiplier utility.ComplexNumberForMarshal	`json:"multiplier" yaml:"multiplier"`
}

// HexagonalWallpaperFormula uses waves that create a 3 rotation symmetry.
//  Each term will be rotated 3 times and averaged by 1/3.
type HexagonalWallpaperFormula struct {
	WavePackets []*Formula
	Multiplier complex128
}

// SetUp initializes all of the needed wallpaper terms.
func (hexWaveFormula *HexagonalWallpaperFormula) SetUp() {
	for _, wavePacket := range hexWaveFormula.WavePackets {
		baseCoefficientPairing := coefficient.Pairing{
			PowerN: wavePacket.Terms[0].PowerN,
			PowerM: wavePacket.Terms[0].PowerM,
		}

		newPairings := baseCoefficientPairing.GenerateCoefficientSets([]coefficient.Relationship{
			coefficient.PlusMMinusSumNAndM,
			coefficient.MinusSumNAndMPlusN,
		})

		baseXVector := complex(1, 0)
		baseYVector := complex(-0.5, math.Sqrt(3.0)/2.0)
		wavePacket.Terms[0].XLatticeVector = complex(real(baseXVector), imag(baseXVector))
		wavePacket.Terms[0].YLatticeVector = complex(real(baseYVector), imag(baseYVector))

		for _, newCoefficientPair := range newPairings {
			xVector := baseXVector
			yVector := baseYVector

			if newCoefficientPair.NegateMultiplier == true {
				xVector *= -1
				yVector *= -1
			}
			newEisenstein := &formula.EisensteinFormulaTerm{
				XLatticeVector: xVector,
				YLatticeVector: yVector,
				PowerN:         newCoefficientPair.PowerN,
				PowerM:         newCoefficientPair.PowerM,
			}
			wavePacket.Terms = append(wavePacket.Terms, newEisenstein)
		}
	}
}

// Calculate takes the complex number z and processes it using the mathematical terms.
func (hexWaveFormula *HexagonalWallpaperFormula) Calculate(z complex128) *formula.CalculationResultForFormula {

	result := &formula.CalculationResultForFormula{
		Total: complex(0,0),
		ContributionByTerm: []complex128{},
	}

	for _, wavePacket := range hexWaveFormula.WavePackets {
		termContribution := wavePacket.Calculate(z)
		result.Total += termContribution.Total * 1.0/3.0
		result.ContributionByTerm = append(result.ContributionByTerm, termContribution.Total)
	}
	result.Total *= hexWaveFormula.Multiplier

	return result
}

// FindSymmetries returns an object with a bunch of symmetries.
func (hexWaveFormula *HexagonalWallpaperFormula) FindSymmetries() *Symmetry {
	foundSymmetries := Symmetry{
		P3: true,
	}

	symmetryFound := FindWaveRelationships(hexWaveFormula.WavePackets)
	if symmetryFound == PlusMPlusN {
		foundSymmetries.P31m = true
	}
	if symmetryFound == MinusMMinusN {
		foundSymmetries.P3m1 = true
	}
	if symmetryFound == MinusNMinusM {
		foundSymmetries.P6 = true
	}
	if symmetryFound == MinusNMinusMPlusMPlusNMinusMMinusN {
		foundSymmetries.P6m = true
	}
	return &foundSymmetries
}

// NewHexagonalWallpaperFormulaFromJSON reads the data and returns a formula term from it.
func NewHexagonalWallpaperFormulaFromJSON(data []byte) (*HexagonalWallpaperFormula, error) {
	return newHexagonalWallpaperFormulaFromDatastream(data, json.Unmarshal)
}

// NewHexagonalWallpaperFormulaFromYAML reads the data and returns a formula term from it.
func NewHexagonalWallpaperFormulaFromYAML(data []byte) (*HexagonalWallpaperFormula, error) {
	return newHexagonalWallpaperFormulaFromDatastream(data, yaml.Unmarshal)
}

//newHexagonalWallpaperFormulaFromDatastream consumes a given bytestream and tries to create a new object from it.
func newHexagonalWallpaperFormulaFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*HexagonalWallpaperFormula, error) {
	var unmarshalError error
	var formulaMarshal HexagonalWallpaperFormulaMarshalable
	unmarshalError = unmarshal(data, &formulaMarshal)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	formulaTerm := NewHexagonalWallpaperFormulaFromMarshalObject(formulaMarshal)
	return formulaTerm, nil
}

// NewHexagonalWallpaperFormulaFromMarshalObject uses a marshalled object to create a new object.
func NewHexagonalWallpaperFormulaFromMarshalObject(marshalObject HexagonalWallpaperFormulaMarshalable) *HexagonalWallpaperFormula {
	wavePackets := []*Formula{}
	for _,packet := range marshalObject.WavePackets {
		newWavePacket := NewWaveFormulaFromMarshalObject(*packet)
		wavePackets = append(wavePackets, newWavePacket)
	}

	return &HexagonalWallpaperFormula{
		WavePackets: wavePackets,
		Multiplier:	complex(marshalObject.Multiplier.Real, marshalObject.Multiplier.Imaginary),
	}
}
