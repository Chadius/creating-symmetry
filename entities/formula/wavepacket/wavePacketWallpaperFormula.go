package wavepacket

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/coefficient"
	"wallpaper/entities/utility"
)

// WallpaperFormulaMarshalled can be marshalled into Wave Packet formulas
type WallpaperFormulaMarshalled struct {
	WavePackets []*Marshal                     `json:"wave_packets" yaml:"wave_packets"`
	Multiplier utility.ComplexNumberForMarshal `json:"multiplier" yaml:"multiplier"`
	Lattice *formula.LatticeVectorPairMarshal  `json:"lattice" yaml:"lattice"`
}

// WallpaperFormula uses wave packets that enforce rotation symmetry.
type WallpaperFormula struct {
	WavePackets []*WavePacket
	Multiplier complex128
	Lattice *formula.LatticeVectorPair
}

// SetUp adds the locked Eisenstein terms to the formula.
//   SetUp assumes only the first term of each Wave Packet was given.
//   The type of wallpaper supplies the paired coefficient relationships
//     as well as the base vectors.
func (wallpaperFormula *WallpaperFormula) SetUp(
	lockedRelationships []coefficient.Relationship,
	) {
	for _, wavePacket := range wallpaperFormula.WavePackets {
		baseCoefficientPairing := coefficient.Pairing{
			PowerN: wavePacket.Terms[0].PowerN,
			PowerM: wavePacket.Terms[0].PowerM,
		}

		newPairings := baseCoefficientPairing.GenerateCoefficientSets(lockedRelationships)

		if real(wavePacket.Terms[0].Multiplier) == 0 && imag(wavePacket.Terms[0].Multiplier) == 0 {
			wavePacket.Terms[0].Multiplier = complex(1, 0)
		}

		for _, newCoefficientPair := range newPairings {
			baseMultiplier := complex(real(wavePacket.Terms[0].Multiplier), imag(wavePacket.Terms[0].Multiplier))

			if newCoefficientPair.NegateMultiplier == true {
				baseMultiplier = complex(-1 * real(wavePacket.Terms[0].Multiplier), -1 * imag(wavePacket.Terms[0].Multiplier))
			}
			newEisenstein := &formula.EisensteinFormulaTerm{
				PowerN:         newCoefficientPair.PowerN,
				PowerM:         newCoefficientPair.PowerM,
				Multiplier: baseMultiplier,
			}
			wavePacket.Terms = append(wavePacket.Terms, newEisenstein)
		}
	}
}

// Calculate takes the complex number z and processes it using the mathematical terms.
func (wallpaperFormula *WallpaperFormula) Calculate(z complex128) *formula.CalculationResultForFormula {

	result := &formula.CalculationResultForFormula{
		Total: complex(0,0),
		ContributionByTerm: []complex128{},
	}

	zInLatticeCoordinates := wallpaperFormula.Lattice.ConvertToLatticeCoordinates(z)

	for _, wavePacket := range wallpaperFormula.WavePackets {
		termContribution := wavePacket.Calculate(zInLatticeCoordinates)
		result.Total += termContribution.Total / complex(float64(len(wavePacket.Terms)), 0)
		result.ContributionByTerm = append(result.ContributionByTerm, termContribution.Total)
	}
	result.Total *= wallpaperFormula.Multiplier

	return result
}

// NewWallpaperFormulaFromYAML reads the data and returns a formula from it.
func NewWallpaperFormulaFromYAML(data []byte) (*WallpaperFormula, error) {
	return newWallpaperFormulaFromDatastream(data, yaml.Unmarshal)
}

// NewWallpaperFormulaFromJSON reads the data and returns a formula from it.
func NewWallpaperFormulaFromJSON(data []byte) (*WallpaperFormula, error) {
	return newWallpaperFormulaFromDatastream(data, json.Unmarshal)
}

//newWallpaperFormulaFromDatastream consumes a given bytestream and tries to create a new object from it.
func newWallpaperFormulaFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*WallpaperFormula, error) {
	var unmarshalError error
	var formulaMarshal WallpaperFormulaMarshalled
	unmarshalError = unmarshal(data, &formulaMarshal)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	formula := NewWallpaperFormulaFromMarshalObject(formulaMarshal)
	return formula, nil
}

// NewWallpaperFormulaFromMarshalObject uses a marshalled object to create a new Wave WavePacket.
func NewWallpaperFormulaFromMarshalObject(marshalObject WallpaperFormulaMarshalled) *WallpaperFormula {
	wavePackets := []*WavePacket{}
	for _,packet := range marshalObject.WavePackets {
		newWavePacket := NewWaveFormulaFromMarshalObject(*packet)
		wavePackets = append(wavePackets, newWavePacket)
	}

	return &WallpaperFormula{
		WavePackets: wavePackets,
		Multiplier:  complex(marshalObject.Multiplier.Real, marshalObject.Multiplier.Imaginary),
	}
}
