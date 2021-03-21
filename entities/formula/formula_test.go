package formula_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"wallpaper/entities/formula"
)

var _ = Describe("Common formula features", func() {
	Context("Can determine coefficients and scale based on coefficient relationship", func() {
		It("PlusNPlusM reuses given coefficients", func() {
			power1, power2, scale := formula.SetCoefficientsBasedOnRelationship(1, 2, complex(3, 4), formula.PlusNPlusM)
			Expect(power1).To(Equal(1))
			Expect(power2).To(Equal(2))
			Expect(real(scale)).To(BeNumerically("~", 3))
			Expect(imag(scale)).To(BeNumerically("~", 4))
		})
		It("PlusMPlusN swaps given coefficients", func() {
			power1, power2, scale := formula.SetCoefficientsBasedOnRelationship(1, 2, complex(3, 4), formula.PlusMPlusN)
			Expect(power1).To(Equal(2))
			Expect(power2).To(Equal(1))
			Expect(real(scale)).To(BeNumerically("~", 3))
			Expect(imag(scale)).To(BeNumerically("~", 4))
		})
		It("MinusNMinusM reuses given coefficients", func() {
			power1, power2, scale := formula.SetCoefficientsBasedOnRelationship(1, 2, complex(3, 4), formula.MinusNMinusM)
			Expect(power1).To(Equal(-1))
			Expect(power2).To(Equal(-2))
			Expect(real(scale)).To(BeNumerically("~", 3))
			Expect(imag(scale)).To(BeNumerically("~", 4))
		})
		It("MinusMMinusN swaps given coefficients", func() {
			power1, power2, scale := formula.SetCoefficientsBasedOnRelationship(1, 2, complex(3, 4), formula.MinusMMinusN)
			Expect(power1).To(Equal(-2))
			Expect(power2).To(Equal(-1))
			Expect(real(scale)).To(BeNumerically("~", 3))
			Expect(imag(scale)).To(BeNumerically("~", 4))
		})
		It("PlusMPlusNMaybeFlipScale changes the scale of odd summed powers ", func() {
			power1, power2, scale := formula.SetCoefficientsBasedOnRelationship(1, 2, complex(3, 4), formula.PlusMPlusNMaybeFlipScale)
			Expect(power1).To(Equal(2))
			Expect(power2).To(Equal(1))
			Expect(real(scale)).To(BeNumerically("~", -3))
			Expect(imag(scale)).To(BeNumerically("~", -4))
		})
		It("PlusMPlusNMaybeFlipScale does not change the scale of even summed powers ", func() {
			power1, power2, scale := formula.SetCoefficientsBasedOnRelationship(0, 2, complex(3, 4), formula.PlusMPlusNMaybeFlipScale)
			Expect(power1).To(Equal(2))
			Expect(power2).To(Equal(0))
			Expect(real(scale)).To(BeNumerically("~", 3))
			Expect(imag(scale)).To(BeNumerically("~", 4))
		})
		It("MinusMMinusNMaybeFlipScale changes the scale of odd summed powers ", func() {
			power1, power2, scale := formula.SetCoefficientsBasedOnRelationship(1, 2, complex(3, 4), formula.MinusMMinusNMaybeFlipScale)
			Expect(power1).To(Equal(-2))
			Expect(power2).To(Equal(-1))
			Expect(real(scale)).To(BeNumerically("~", -3))
			Expect(imag(scale)).To(BeNumerically("~", -4))
		})
		It("MinusMMinusNMaybeFlipScale does not change the scale of even summed powers ", func() {
			power1, power2, scale := formula.SetCoefficientsBasedOnRelationship(0, 2, complex(3, 4), formula.MinusMMinusNMaybeFlipScale)
			Expect(power1).To(Equal(-2))
			Expect(power2).To(Equal(0))
			Expect(real(scale)).To(BeNumerically("~", 3))
			Expect(imag(scale)).To(BeNumerically("~", 4))
		})
	})

	Context("Create LockedCoefficientPair objects", func() {
		It("Can create from YAML", func() {
			yamlByteStream := []byte(`multiplier: 1
relationships:
  - "+N+M"
  - "-M-NF"
`)
			coefficientPair, err := formula.NewLockedCoefficientPairFromYAML(yamlByteStream)
			Expect(err).To(BeNil())
			Expect(coefficientPair.Multiplier).To(BeNumerically("~", 1))
			Expect(coefficientPair.OtherCoefficientRelationships).To(HaveLen(2))
			Expect(coefficientPair.OtherCoefficientRelationships[0]).To(Equal(formula.CoefficientRelationship(formula.PlusNPlusM)))
			Expect(coefficientPair.OtherCoefficientRelationships[1]).To(Equal(formula.CoefficientRelationship(formula.MinusMMinusNMaybeFlipScale)))
		})
	})
})