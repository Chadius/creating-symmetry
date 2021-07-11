package wavepacket

import (
	//"encoding/json"
	//"gopkg.in/yaml.v2"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/coefficient"
	"wallpaper/entities/utility"

	//"wallpaper/entities/utility"
)

// RectangularWallpaperFormulaMarshalled can be marshalled into a Rectangular formula
type RectangularWallpaperFormulaMarshalled struct {
	Formula *WallpaperFormulaMarshalled `json:"formula" yaml:"formula"`
	LatticeHeight float64  `json:"lattice_height" yaml:"lattice_height"`
}

// RectangularWallpaperFormula helps transform one point to a 2D wallpaper pattern that uses the Rectangular lattice.
//   The underlying lattice forms a rectangle. Horizontal is 1 unit, Vertical LatticeHeight cannot be 0.
type RectangularWallpaperFormula struct {
	Formula *WallpaperFormula
	LatticeHeight float64
}

// SetUp will create the Rectangular RectangularWallpaperFormula using the given LatticeHeight.
func (Rectangular *RectangularWallpaperFormula) SetUp() error {

	Rectangular.Formula.Lattice = &formula.LatticeVectorPair{
		XLatticeVector: complex(1, 0),
		YLatticeVector: complex(0, Rectangular.LatticeHeight),
	}

	err := Rectangular.Formula.Lattice.Validate()
	if err != nil {
		return err
	}

	for _, wavePacket := range Rectangular.Formula.WavePackets {
		for _, term := range wavePacket.Terms {
			if real(term.Multiplier) == 0 && imag(term.Multiplier) == 0 {
				term.Multiplier = complex(1, 0)
			}
		}
	}

	return nil
}

// Calculate applies the formula to the complex number z.
//  It modifies the formula's result to track the contribution per term
//  As well as the final numerical result.
func (Rectangular *RectangularWallpaperFormula) Calculate(z complex128) *formula.CalculationResultForFormula {
	return Rectangular.Formula.Calculate(z)
}

// HasSymmetry returns true if the WavePackets involved form symmetry.
func (Rectangular *RectangularWallpaperFormula) HasSymmetry(desiredSymmetry Symmetry) bool {
	return HasSymmetry(Rectangular.Formula.WavePackets, desiredSymmetry, map[Symmetry][]coefficient.Relationship {
		Pm: {coefficient.PlusNMinusM},
		Pg: {coefficient.PlusNMinusMNegateMultiplierIfOddPowerN},
		Pmm: {
			coefficient.PlusNMinusM,
			coefficient.MinusNMinusM,
			coefficient.MinusNPlusM,
		},
		Pmg: {
			coefficient.MinusNMinusM,
			coefficient.PlusNMinusMNegateMultiplierIfOddPowerN,
			coefficient.MinusNPlusMNegateMultiplierIfOddPowerN,
		},
		Pgg: {
			coefficient.MinusNMinusM,
			coefficient.PlusNMinusMNegateMultiplierIfOddPowerSum,
			coefficient.MinusNPlusMNegateMultiplierIfOddPowerSum,
		},
	})
}

// NewRectangularWallpaperFormulaFromJSON reads the data and returns a formula term from it.
func NewRectangularWallpaperFormulaFromJSON(data []byte) (*RectangularWallpaperFormula, error) {
	return newRectangularWallpaperFormulaFromDatastream(data, json.Unmarshal)
}

// NewRectangularWallpaperFormulaFromYAML reads the data and returns a formula term from it.
func NewRectangularWallpaperFormulaFromYAML(data []byte) (*RectangularWallpaperFormula, error) {
	return newRectangularWallpaperFormulaFromDatastream(data, yaml.Unmarshal)
}

//newRectangularWallpaperFormulaFromDatastream consumes a given bytestream and tries to create a new object from it.
func newRectangularWallpaperFormulaFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*RectangularWallpaperFormula, error) {
	var unmarshalError error
	var RectangularFormulaMarshalled RectangularWallpaperFormulaMarshalled
	unmarshalError = unmarshal(data, &RectangularFormulaMarshalled)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	return NewRectangularWallpaperFormulaFromMarshalObject(RectangularFormulaMarshalled), nil
}

// NewRectangularWallpaperFormulaFromMarshalObject uses a marshalled object to create a new object.
func NewRectangularWallpaperFormulaFromMarshalObject(marshalObject RectangularWallpaperFormulaMarshalled) *RectangularWallpaperFormula {
	formula := NewWallpaperFormulaFromMarshalObject(*marshalObject.Formula)

	if marshalObject.Formula.DesiredSymmetry != "" {
		wallpaper, err := NewRectangularWallpaperFormulaWithSymmetry(
			formula.WavePackets[0].Terms,
			formula.Multiplier,
			marshalObject.LatticeHeight,
			Symmetry(marshalObject.Formula.DesiredSymmetry),
		)

		if err != nil {
			return nil
		}
		return wallpaper
	}

	return &RectangularWallpaperFormula{
		Formula:       formula,
		LatticeHeight: marshalObject.LatticeHeight,
	}
}

// NewRectangularWallpaperFormulaWithSymmetry will try to create a new RectangularWallpaperFormula WavePacket
//   with the desired Terms, Multiplier and Symmetry.
func NewRectangularWallpaperFormulaWithSymmetry(terms []*formula.EisensteinFormulaTerm, wallpaperMultiplier complex128, latticeHeight float64, desiredSymmetry Symmetry) (*RectangularWallpaperFormula, error) {
	newWavePackets := []*WavePacket{}
	for _, term := range terms {
		if real(term.Multiplier) == 0 && imag(term.Multiplier) == 0 {
			term.Multiplier = complex(1, 0)
		}

		newWavePackets = append(
			newWavePackets,
			&WavePacket{
				Terms:      []*formula.EisensteinFormulaTerm{term},
				Multiplier: term.Multiplier,
			},
		)

		newWavePackets = addNewWavePacketsBasedOnSymmetry(term, desiredSymmetry, newWavePackets)
	}

	newBaseWallpaper := &RectangularWallpaperFormula{
		Formula: &WallpaperFormula{
			WavePackets: newWavePackets,
			Multiplier:  wallpaperMultiplier,
		},
		LatticeHeight: latticeHeight,
	}
	newBaseWallpaper.SetUp()
	return newBaseWallpaper, nil
}
