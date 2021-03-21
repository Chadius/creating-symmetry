package formula_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math"
	"wallpaper/entities/formula"
)

var _ = Describe("Common formula formats", func() {
	Context("Terms that involve e^(inz) * e^(-imzConj)", func() {
		It("Can calculate a formula that uses Euler and complex numbers", func() {
			form := formula.EulerFormulaTerm{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
			}
			result := form.Calculate(complex(math.Pi / 6.0,1))
			Expect(real(result)).To(BeNumerically("~", 3 * math.Exp(-2) * 1.0 / 2.0))
			Expect(imag(result)).To(BeNumerically("~", 3 * math.Exp(-2) * math.Sqrt(3.0) / 2.0))
		})
		It("Can calculate a formula that uses locked coefficient pairs", func() {
			form := formula.EulerFormulaTerm{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 0,
				IgnoreComplexConjugate: true,
				CoefficientPairs: formula.LockedCoefficientPair{
					Multiplier: 1,
					OtherCoefficientRelationships: []formula.CoefficientRelationship{
						formula.PlusMPlusN,
					},
				},
			}
			result := form.Calculate(complex(math.Pi / 6.0,1))
			Expect(real(result)).To(BeNumerically("~", 3 * ((math.Exp(-2) * 1.0 / 2.0) + 1.0)))
			Expect(imag(result)).To(BeNumerically("~", 3 * math.Exp(-2) * math.Sqrt(3.0) / 2.0))
		})
		It("Can calculate a formula that uses the complex conjugate", func() {
			form := formula.EulerFormulaTerm{
				Scale:                  complex(3, 0),
				PowerN:                 2,
				PowerM:                 1,
				IgnoreComplexConjugate: false,
			}
			result := form.Calculate(complex(math.Pi / 6.0,2))
			Expect(real(result)).To(BeNumerically("~", 3 * math.Exp(-6) * math.Sqrt(3.0) / 2.0))
			Expect(imag(result)).To(BeNumerically("~", 3 * math.Exp(-6) * 1.0 / 2.0))
		})
	})

	It("Can calculate a Frieze formula", func() {
		friezeFormula := formula.FriezeFormula{
			Terms: []*formula.EulerFormulaTerm{
				{
					Scale: complex(2, 0),
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
		result := friezeFormula.Calculate(complex(math.Pi/6, 1))
		total := result.Total

		expectedResult := complex(math.Exp(-1), 0) * complex(math.Sqrt(3) * 2, 0)
		Expect(real(total)).To(BeNumerically("~", real(expectedResult)))
		Expect(imag(total)).To(BeNumerically("~", imag(expectedResult)))
	})

	Context("Analyze Friezes for symmetry", func() {
		It("Knows when a pattern is p211", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.MinusNMinusM,
							},
						},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P211).To(BeTrue())
		})
		It("Knows when a pattern is p1m1", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
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
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P1m1).To(BeTrue())
		})
		It("Knows when a pattern is p11m", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.MinusMMinusN,
							},
						},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P11m).To(BeTrue())
		})
		It("Knows when a pattern is p11g", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 1,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.MinusMMinusNMaybeFlipScale,
							},
						},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P11g).To(BeTrue())
		})
		It("Knows when a pattern is p11m if a p11g pattern has even sum powers", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.MinusMMinusNMaybeFlipScale,
							},
						},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P11m).To(BeTrue())
		})
		It("Knows when a pattern is p2mm", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.MinusNMinusM,
								formula.PlusMPlusN,
								formula.MinusMMinusN,
							},
						},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P2mm).To(BeTrue())
		})
		It("Knows when a pattern is p2mg", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: -1,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.MinusNMinusM,
								formula.PlusMPlusNMaybeFlipScale,
								formula.MinusMMinusNMaybeFlipScale,
							},
						},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P2mg).To(BeTrue())
		})
		It("Knows when a pattern is p2mm if a p2mg pattern has even sum powers", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.MinusNMinusM,
								formula.PlusMPlusNMaybeFlipScale,
								formula.MinusMMinusNMaybeFlipScale,
							},
						},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P2mm).To(BeTrue())
		})
		It("Knows when a pattern is p111", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 0,
						IgnoreComplexConjugate: false,
						CoefficientPairs: formula.LockedCoefficientPair{},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P111).To(BeTrue())
		})
		It("Knows when a pattern is p111 because complex coordinates are ignored", func() {
			friezeFormula := formula.FriezeFormula{
				Terms: []*formula.EulerFormulaTerm{
					{
						Scale: complex(1, 0),
						PowerN: 2,
						PowerM: 0,
						IgnoreComplexConjugate: true,
						CoefficientPairs: formula.LockedCoefficientPair{
							Multiplier: 1,
							OtherCoefficientRelationships: []formula.CoefficientRelationship{
								formula.MinusNMinusM,
							},
						},
					},
				},
			}
			symmetriesDetected := friezeFormula.AnalyzeForSymmetry()
			Expect(symmetriesDetected.P111).To(BeTrue())
			Expect(symmetriesDetected.P211).To(BeFalse())
		})
	})
	It("Can determine the contribution by each term of a Frieze formula", func() {
		friezeFormula := formula.FriezeFormula{
			Terms: []*formula.EulerFormulaTerm{
				{
					Scale: complex(2, 0),
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
		result := friezeFormula.Calculate(complex(math.Pi/6, 1))

		Expect(result.ContributionByTerm).To(HaveLen(1))
		contributionByFirstTerm := result.ContributionByTerm[0]

		expectedResult := complex(math.Exp(-1), 0) * complex(math.Sqrt(3) * 2, 0)
		Expect(real(contributionByFirstTerm)).To(BeNumerically("~", real(expectedResult)))
		Expect(imag(contributionByFirstTerm)).To(BeNumerically("~", imag(expectedResult)))
	})
})
