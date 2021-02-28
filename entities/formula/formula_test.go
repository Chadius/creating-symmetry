package formula_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"wallpaper/entities/formula"
)

var _ = Describe("Common formula formats", func() {
	It("Can calulate a Rosette formula", func() {
		rosetteFormula := formula.RosetteFormula{
			Elements: []*formula.ZExponentialFormulaElement{
				{
					Scale: complex(3, 0),
					PowerN: 1,
					PowerM: 0,
					IgnoreComplexConjugate: false,
					LockedCoefficientPairs: []*formula.LockedCoefficientPair{
						{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.PlusMPlusN,
							},
						},
					},
				},
			},
		}
		result := rosetteFormula.Calculate(complex(2,1))
		Expect(real(result)).To(BeNumerically("~", 12))
		Expect(imag(result)).To(BeNumerically("~", 0))
	})

	Context("Terms that involve z^n * zConj^m", func() {
		It("Can make a z to the n exponential formula", func() {

			form := formula.ZExponentialFormulaElement{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
			}
			result := form.Calculate(complex(3,2))
			Expect(real(result)).To(BeNumerically("~", 15))
			Expect(imag(result)).To(BeNumerically("~", 36))
		})
		It("Can make a z to the n exponential formula using locked pairs", func() {
			form := formula.ZExponentialFormulaElement{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
				LockedCoefficientPairs: []*formula.LockedCoefficientPair{
					{
						Multiplier: -1,
						OtherCoefficientRelationships: []formula.CoefficientRelationship{
							formula.PlusMPlusN,
						},
					},
				},
			}
			result := form.Calculate(complex(3,2))
			Expect(real(result)).To(BeNumerically("~", 12))
			Expect(imag(result)).To(BeNumerically("~", 36))
		})
		It("Can make a z to the n exponential formula using a complex conjugate", func() {
			form := formula.ZExponentialFormulaElement{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 1,
				IgnoreComplexConjugate: false,
			}
			result := form.Calculate(complex(3,2))
			Expect(real(result)).To(BeNumerically("~", 117))
			Expect(imag(result)).To(BeNumerically("~", 78))
		})
	})

	Context("Coefficient Relationships", func() {
		It("Can keep coefficients in same order", func() {
			form := formula.ZExponentialFormulaElement{
				Scale:                  complex(1, 0),
				PowerN:                 1,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
				LockedCoefficientPairs: []*formula.LockedCoefficientPair{
					{
						Multiplier: -1,
						OtherCoefficientRelationships: []formula.CoefficientRelationship{
							formula.PlusNPlusM,
						},
					},
				},
			}
			result := form.Calculate(complex(3,2))
			Expect(real(result)).To(BeNumerically("~", 0))
			Expect(imag(result)).To(BeNumerically("~", 0))
		})
		It("Can swap coefficients", func() {
			form := formula.ZExponentialFormulaElement{
				Scale:                  complex(1, 0),
				PowerN:                 1,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
				LockedCoefficientPairs: []*formula.LockedCoefficientPair{
					{
						Multiplier: -1,
						OtherCoefficientRelationships: []formula.CoefficientRelationship{
							formula.PlusMPlusN,
						},
					},
				},
			}
			result := form.Calculate(complex(3,2))
			Expect(real(result)).To(BeNumerically("~", 2))
			Expect(imag(result)).To(BeNumerically("~", 2))
		})
	})
})