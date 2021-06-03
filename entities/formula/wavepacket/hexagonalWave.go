package wavepacket

import (
	"errors"
	"math"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/coefficient"
)

// HexagonalWallpaperFormula uses waves that create a 3 rotation symmetry.
//  Each term will be rotated 3 times and averaged by 1/3.
type HexagonalWallpaperFormula struct {
	Formula *WallpaperFormula
}

// SetUp initializes all of the needed wallpaper terms.
func (hexWaveFormula *HexagonalWallpaperFormula) SetUp() {
	hexWaveFormula.Formula.Lattice = &formula.LatticeVectorPair{
		XLatticeVector: complex(1, 0),
		YLatticeVector: complex(-0.5, math.Sqrt(3.0)/2.0),
	}
	hexWaveFormula.Formula.SetUp(
		[]coefficient.Relationship{
			coefficient.PlusMMinusSumNAndM,
			coefficient.MinusSumNAndMPlusN,
		},
	)
}

// Calculate applies the formula to the complex number z.
//  It modifies the formula's result to track the contribution per term
//  As well as the final numerical result.
func (hexWaveFormula *HexagonalWallpaperFormula) Calculate(z complex128) *formula.CalculationResultForFormula {
	return hexWaveFormula.Formula.Calculate(z)
}

// HasSymmetry returns true if the WavePackets involved form symmetry.
func (hexWaveFormula *HexagonalWallpaperFormula) HasSymmetry(desiredSymmetry SymmetryType) bool {
	if desiredSymmetry == P3 {
		return true
	}

	numberOfWavePackets := len(hexWaveFormula.Formula.WavePackets)
	if numberOfWavePackets < 2 || numberOfWavePackets % 2 == 1 {
		return false
	}

	desiredSymmetryToCoefficients := map[SymmetryType][]coefficient.Relationship {
		P31m: {coefficient.PlusMPlusN},
		P3m1: {coefficient.MinusMMinusN},
		P6: {coefficient.MinusNMinusM},
		P6m: {
			coefficient.MinusNMinusM,
			coefficient.MinusMMinusN,
			coefficient.PlusMPlusN,
		},
	}

	coefficientsToFind := desiredSymmetryToCoefficients[desiredSymmetry]

	if coefficientsToFind == nil {
		return false
	}

	return CanWavePacketsBeGroupedAmongCoefficientRelationships(hexWaveFormula.Formula.WavePackets, coefficientsToFind)
}

// NewHexagonalWallpaperFormulaFromJSON reads the data and returns a formula term from it.
func NewHexagonalWallpaperFormulaFromJSON(data []byte) (*HexagonalWallpaperFormula, error) {
	formula, err := NewWallpaperFormulaFromJSON(data)

	if err != nil {
		return nil, err
	}

	return &HexagonalWallpaperFormula{
		Formula: formula,
	}, nil
}

// NewHexagonalWallpaperFormulaFromYAML reads the data and returns a formula term from it.
func NewHexagonalWallpaperFormulaFromYAML(data []byte) (*HexagonalWallpaperFormula, error) {
	formula, err := NewWallpaperFormulaFromYAML(data)

	if err != nil {
		return nil, err
	}

	return &HexagonalWallpaperFormula{
		Formula: formula,
	}, nil
}

// NewHexagonalWallpaperFormulaFromMarshalObject uses a marshalled object to create a new object.
func NewHexagonalWallpaperFormulaFromMarshalObject(marshalObject WallpaperFormulaMarshalled) *HexagonalWallpaperFormula {
	return &HexagonalWallpaperFormula{
		Formula: NewWallpaperFormulaFromMarshalObject(marshalObject),
	}
}

// NewHexagonalWallpaperFormulaWithSymmetry will try to create a new Hexagonal Wallpaper WavePacket
//   with the desired Terms, Multiplier and Symmetry.
func NewHexagonalWallpaperFormulaWithSymmetry(terms []*formula.EisensteinFormulaTerm, wallpaperMultiplier complex128, desiredSymmetry *Symmetry) (*HexagonalWallpaperFormula, error) {
	err := checkForIncompatibleHexagonalSymmetries(desiredSymmetry)
	if err != nil {
		return nil, err
	}

	newWavePackets := []*WavePacket{}
	for _, term := range terms {
		newWavePackets = append(
			newWavePackets,
			&WavePacket{
				Terms:      []*formula.EisensteinFormulaTerm{term},
				Multiplier: term.Multiplier,
			},
		)

		newWavePackets = addNewWavePacketsBasedOnSymmetry(term, desiredSymmetry, newWavePackets)
	}

	newBaseWallpaper := &HexagonalWallpaperFormula{
		Formula: &WallpaperFormula{
			WavePackets: newWavePackets,
			Multiplier:  wallpaperMultiplier,
		},
	}
	newBaseWallpaper.SetUp()
	return newBaseWallpaper, nil
}

func checkForIncompatibleHexagonalSymmetries(desiredSymmetry *Symmetry) error {
	desiredSymmetryAlreadySet := false
	if desiredSymmetry.P3m1 {
		if desiredSymmetryAlreadySet {
			return errors.New("invalid desired symmetry")
		}
		desiredSymmetryAlreadySet = true
	}
	if desiredSymmetry.P31m {
		if desiredSymmetryAlreadySet {
			return errors.New("invalid desired symmetry")
		}
		desiredSymmetryAlreadySet = true
	}
	if desiredSymmetry.P6 {
		if desiredSymmetryAlreadySet {
			return errors.New("invalid desired symmetry")
		}
		desiredSymmetryAlreadySet = true
	}
	if desiredSymmetry.P6m {
		if desiredSymmetryAlreadySet {
			return errors.New("invalid desired symmetry")
		}
		desiredSymmetryAlreadySet = true
	}
	return nil
}