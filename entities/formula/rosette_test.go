package formula_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"wallpaper/entities/formula"
)

var _ = Describe("Rosette formulas", func() {
	It("Can calculate a Rosette formula", func() {
		rosetteFormula := formula.RosetteFormula{
			Terms: []*formula.ZExponentialFormulaTerm{
				{
					Scale: complex(3, 0),
					PowerN: 1,
					PowerM: 0,
					IgnoreComplexConjugate: false,
					CoefficientPairs: formula.LockedCoefficientPair{
						Multiplier: 1,
						OtherCoefficientRelationships: []formula.CoefficientRelationship{
							formula.PlusMPlusN,
						},
					},
				},
			},
		}
		result := rosetteFormula.Calculate(complex(2,1))
		total := result.Total
		Expect(real(total)).To(BeNumerically("~", 12))
		Expect(imag(total)).To(BeNumerically("~", 0))
	})
	Context("Analyze Rosettes for symmetry", func() {
		It("Can determine there is multifold symmetry with 1 term", func() {
			rosetteFormula := formula.RosetteFormula{
				Terms: []*formula.ZExponentialFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 6,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.PlusMPlusN,
							},
						},
					},
				},
			}
			symmetriesDetected := rosetteFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.Multifold).To(Equal(6))
		})
		It("Multifold symmetry is always a positive value", func() {
			rosetteFormula := formula.RosetteFormula{
				Terms: []*formula.ZExponentialFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: -6,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.PlusMPlusN,
							},
						},
					},
				},
			}
			symmetriesDetected := rosetteFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.Multifold).To(Equal(6))
		})
		It("Multifold symmetry uses the greatest common denominator of all elements", func() {
			rosetteFormula := formula.RosetteFormula{
				Terms: []*formula.ZExponentialFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: -6,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.PlusMPlusN,
							},
						},
					},
					{
						Scale: complex(1, 0),
						PowerN: -8,
						PowerM: 4,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.PlusMPlusN,
							},
						},
					},
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.PlusMPlusN,
							},
						},
					},
				},
			}
			symmetriesDetected := rosetteFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.Multifold).To(Equal(2))
		})
	})
	It("Can determine the contribution by each term of a Rosette formula", func() {
		rosetteFormula := formula.RosetteFormula{
			Terms: []*formula.ZExponentialFormulaTerm{
				{
					Scale: complex(3, 0),
					PowerN: 1,
					PowerM: 0,
					IgnoreComplexConjugate: false,
					CoefficientPairs: formula.LockedCoefficientPair{
						Multiplier: 1,
						OtherCoefficientRelationships: []formula.CoefficientRelationship{
							formula.PlusMPlusN,
						},
					},
				},
			},
		}
		result := rosetteFormula.Calculate(complex(2,1))
		Expect(result.ContributionByTerm).To(HaveLen(1))
		contributionByFirstTerm := result.ContributionByTerm[0]
		Expect(real(contributionByFirstTerm)).To(BeNumerically("~", 12))
		Expect(imag(contributionByFirstTerm)).To(BeNumerically("~", 0))
	})

	Context("Terms that involve z^n * zConj^m", func() {
		It("Can make a z to the n exponential formula", func() {
			form := formula.ZExponentialFormulaTerm{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
			}
			total := form.Calculate(complex(3,2))
			Expect(real(total)).To(BeNumerically("~", 15))
			Expect(imag(total)).To(BeNumerically("~", 36))
		})
		It("Can make a z to the n exponential formula using locked pairs", func() {
			form := formula.ZExponentialFormulaTerm{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
				CoefficientPairs: formula.LockedCoefficientPair{
					Multiplier: -1,
					OtherCoefficientRelationships: []formula.CoefficientRelationship{
						formula.PlusMPlusN,
					},
				},
			}
			total := form.Calculate(complex(3,2))
			Expect(real(total)).To(BeNumerically("~", 12))
			Expect(imag(total)).To(BeNumerically("~", 36))
		})
		It("Can make a z to the n exponential formula using a complex conjugate", func() {
			form := formula.ZExponentialFormulaTerm{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 1,
				IgnoreComplexConjugate: false,
			}
			total := form.Calculate(complex(3,2))
			Expect(real(total)).To(BeNumerically("~", 117))
			Expect(imag(total)).To(BeNumerically("~", 78))
		})
	})

	Context("Coefficient Relationships", func() {
		It("Can keep coefficients in same order", func() {
			form := formula.ZExponentialFormulaTerm{
				Scale:                  complex(1, 0),
				PowerN:                 1,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
				CoefficientPairs: formula.LockedCoefficientPair{
					Multiplier: -1,
					OtherCoefficientRelationships: []formula.CoefficientRelationship{
						formula.PlusNPlusM,
					},
				},
			}
			total := form.Calculate(complex(3,2))
			Expect(real(total)).To(BeNumerically("~", 0))
			Expect(imag(total)).To(BeNumerically("~", 0))
		})
		It("Can swap coefficients", func() {
			form := formula.ZExponentialFormulaTerm{
				Scale:                  complex(1, 0),
				PowerN:                 1,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
				CoefficientPairs: formula.LockedCoefficientPair{
					Multiplier: -1,
					OtherCoefficientRelationships: []formula.CoefficientRelationship{
						formula.PlusMPlusN,
					},
				},
			}
			total := form.Calculate(complex(3,2))
			Expect(real(total)).To(BeNumerically("~", 2))
			Expect(imag(total)).To(BeNumerically("~", 2))
		})
	})
})
