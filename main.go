package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"wallpaper/entities/command"
	"wallpaper/entities/formula/frieze"
	"wallpaper/entities/formula/rosette"
	"wallpaper/entities/formula/wavepacket"

	_ "image/png"
	"os"
	"wallpaper/entities/mathutility"
)

func main() {
	createWallpaperYAML, err := ioutil.ReadFile("data/formula.yml")
	if err != nil {
		log.Fatal(err)
	}
	wallpaperCommand, err := command.NewCreateWallpaperCommandFromYAML(createWallpaperYAML)
	if err != nil {
		log.Fatal(err)
	}
	sampleSpaceMin := complex(wallpaperCommand.SampleSpace.MinX, wallpaperCommand.SampleSpace.MinY)
	sampleSpaceMax := complex(wallpaperCommand.SampleSpace.MaxX, wallpaperCommand.SampleSpace.MaxY)
	outputWidth := wallpaperCommand.OutputImageSize.Width
	outputHeight := wallpaperCommand.OutputImageSize.Height
	colorSourceFilename := wallpaperCommand.SampleSourceFilename
	outputFilename := wallpaperCommand.OutputFilename
	colorValueBoundMin := complex(wallpaperCommand.ColorValueSpace.MinX, wallpaperCommand.ColorValueSpace.MinY)
	colorValueBoundMax := complex(wallpaperCommand.ColorValueSpace.MaxX, wallpaperCommand.ColorValueSpace.MaxY)

	//sampleSpaceMin := complex(-1e0, -1e0)
	//sampleSpaceMax := complex(1e0, 1e0)
	////outputWidth := 800
	////outputHeight := 450
	//outputWidth := 3840
	//outputHeight := 2160
	//colorSourceFilename := "exampleImage/brownie.png"
	//outputFilename := "exampleImage/newBrownie.png"
	//colorValueBoundMin := complex(-2e5, -2e5)
	//colorValueBoundMax := complex(5e5, 5e5)

	reader, err := os.Open(colorSourceFilename)
	if err != nil {
	  log.Fatal(err)
	}
	defer reader.Close()

	colorSourceImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	destinationBounds := image.Rect(0,0, outputWidth, outputHeight)
	destinationCoordinates := flattenCoordinates(destinationBounds)

	scaledCoordinates := scaleDestinationPixels(
		destinationBounds,
		destinationCoordinates,
		sampleSpaceMin,
		sampleSpaceMax,
	)

	transformedCoordinates := transformCoordinatesForFormula(wallpaperCommand, scaledCoordinates)
	minz, maxz := mathutility.GetBoundingBox(transformedCoordinates)
	println(minz)
	println(maxz)

	// Consider how to give a preview image? What's the picture ration
	outputImage := image.NewNRGBA(image.Rect(0, 0, outputWidth, outputHeight))
	colorDestinationImage(outputImage, colorSourceImage, destinationCoordinates, transformedCoordinates, colorValueBoundMin, colorValueBoundMax)

	outputToFile(outputFilename, outputImage)
}

func outputToFile(outputFilename string, outputImage image.Image) {
	outputImageFile, err := os.Create(outputFilename)
	if err != nil {
		panic(err)
	}
	defer outputImageFile.Close()
	png.Encode(outputImageFile, outputImage)
}

func transformCoordinatesForFormula(command *command.CreateWallpaperCommand, scaledCoordinates []complex128) []complex128 {
	if command.FriezeFormula != nil {
		return transformCoordinatesForFriezeFormula(command.FriezeFormula, scaledCoordinates)
	}
	if command.RosetteFormula != nil {
		return transformCoordinatesForRosetteFormula(command.RosetteFormula, scaledCoordinates)
	}
	if command.HexagonalWallpaperFormula != nil {
		return transformCoordinatesForHexagonalWallpaperFormula(command.HexagonalWallpaperFormula, scaledCoordinates)
	}
	if command.SquareWallpaperFormula != nil {
		return transformCoordinatesForSquareWallpaperFormula(command.SquareWallpaperFormula, scaledCoordinates)
	}
	if command.RhombicWallpaperFormula != nil {
		return transformCoordinatesForRhombicWallpaperFormula(command.RhombicWallpaperFormula, scaledCoordinates)
	}
	if command.RectangularWallpaperFormula != nil {
		return transformCoordinatesForRectangularWallpaperFormula(command.RectangularWallpaperFormula, scaledCoordinates)
	}
	log.Fatal(errors.New("no formula found"))
	return []complex128{}
}

func transformCoordinatesForFriezeFormula(friezeFormula *frieze.Formula, scaledCoordinates []complex128) []complex128 {
	symmetryAnalysis := friezeFormula.AnalyzeForSymmetry()
	if symmetryAnalysis.P111 {
		println("Has these symmetries: p111")
	}
	if symmetryAnalysis.P211 {
		println("  P211")
	}
	if symmetryAnalysis.P1m1 {
		println("  P1m1")
	}
	if symmetryAnalysis.P11g {
		println("  P11g")
	}
	if symmetryAnalysis.P11m {
		println("  P11m")
	}
	if symmetryAnalysis.P2mm {
		println("  P2mm")
	}
	if symmetryAnalysis.P2mg {
		println("  P2mg")
	}

	transformedCoordinates := []complex128{}
	resultsByTerm := [][]complex128{}
	for range friezeFormula.Terms {
		resultsByTerm = append(resultsByTerm, []complex128{})
	}

	for _, complexCoordinate := range scaledCoordinates {
		friezeResults := friezeFormula.Calculate(complexCoordinate)
		for index, formulaResult := range friezeResults.ContributionByTerm {
			resultsByTerm[index] = append(resultsByTerm[index], formulaResult)
		}

		transformedCoordinate := friezeResults.Total
		transformedCoordinates = append(transformedCoordinates, transformedCoordinate)
	}

	println("Min/Max ranges, by Term")
	for index, results := range resultsByTerm {
		minz, maxz := mathutility.GetBoundingBox(results)
		fmt.Printf("%d: %e - %e\n", index, minz, maxz)
	}
	return transformedCoordinates
}

func transformCoordinatesForRosetteFormula(rosetteFormula *rosette.Formula, scaledCoordinates []complex128) []complex128 {
	transformedCoordinates := []complex128{}
	resultsByTerm := [][]complex128{}
	for range rosetteFormula.Terms {
		resultsByTerm = append(resultsByTerm, []complex128{})
	}

	for _, complexCoordinate := range scaledCoordinates {
		rosetteResults := rosetteFormula.Calculate(complexCoordinate)
		for index, formulaResult := range rosetteResults.ContributionByTerm {
			resultsByTerm[index] = append(resultsByTerm[index], formulaResult)
		}

		transformedCoordinate := rosetteResults.Total
		transformedCoordinates = append(transformedCoordinates, transformedCoordinate)
	}

	println("Min/Max ranges, by Term")
	for index, results := range resultsByTerm {
		minz, maxz := mathutility.GetBoundingBox(results)
		fmt.Printf("%d: %e - %e\n", index, minz, maxz)
	}
	return transformedCoordinates
}

func transformCoordinatesForHexagonalWallpaperFormula(wallpaperFormula *wavepacket.HexagonalWallpaperFormula, scaledCoordinates []complex128) []complex128 {
	wallpaperFormula.SetUp()

	transformedCoordinates := []complex128{}
	resultsByTerm := [][]complex128{}
	for range wallpaperFormula.Formula.WavePackets {
		resultsByTerm = append(resultsByTerm, []complex128{})
	}

	for _, complexCoordinate := range scaledCoordinates {
		hexagonalWallpaperResults := wallpaperFormula.Calculate(complexCoordinate)
		for index, formulaResult := range hexagonalWallpaperResults.ContributionByTerm {
			resultsByTerm[index] = append(resultsByTerm[index], formulaResult)
		}

		transformedCoordinate := hexagonalWallpaperResults.Total
		transformedCoordinates = append(transformedCoordinates, transformedCoordinate)
	}

	println("Min/Max ranges, by Term")
	for index, results := range resultsByTerm {
		minz, maxz := mathutility.GetBoundingBox(results)
		fmt.Printf("%d: %e - %e\n", index, minz, maxz)
	}

	println("Symmetries found:")

	if wallpaperFormula.HasSymmetry(wavepacket.P31m) {
		println("  p31m")
	}
	if wallpaperFormula.HasSymmetry(wavepacket.P3m1) {
		println("  p3m1")
	}
	if wallpaperFormula.HasSymmetry(wavepacket.P6) {
		println("  p6")
	}
	if wallpaperFormula.HasSymmetry(wavepacket.P6m) {
		println("  p6m")
	}
	if wallpaperFormula.HasSymmetry(wavepacket.P3) {
		println("  p3")
	}

	return transformedCoordinates
}

func transformCoordinatesForSquareWallpaperFormula(wallpaperFormula *wavepacket.SquareWallpaperFormula, scaledCoordinates []complex128) []complex128 {
	wallpaperFormula.SetUp()

	transformedCoordinates := []complex128{}
	resultsByTerm := [][]complex128{}
	for range wallpaperFormula.Formula.WavePackets {
		resultsByTerm = append(resultsByTerm, []complex128{})
	}

	for _, complexCoordinate := range scaledCoordinates {
		hexagonalWallpaperResults := wallpaperFormula.Calculate(complexCoordinate)
		for index, formulaResult := range hexagonalWallpaperResults.ContributionByTerm {
			resultsByTerm[index] = append(resultsByTerm[index], formulaResult)
		}

		transformedCoordinate := hexagonalWallpaperResults.Total
		transformedCoordinates = append(transformedCoordinates, transformedCoordinate)
	}

	println("Min/Max ranges, by Term")
	for index, results := range resultsByTerm {
		minz, maxz := mathutility.GetBoundingBox(results)
		fmt.Printf("%d: %e - %e\n", index, minz, maxz)
	}

	println("Symmetries found:")
	if wallpaperFormula.HasSymmetry(wavepacket.P4) {
		println("  p4")
	}
	if wallpaperFormula.HasSymmetry(wavepacket.P4m) {
		println("  p4m")
	}
	if wallpaperFormula.HasSymmetry(wavepacket.P4g) {
		println("  p4g")
	}

	return transformedCoordinates
}

func transformCoordinatesForRhombicWallpaperFormula(wallpaperFormula *wavepacket.RhombicWallpaperFormula, scaledCoordinates []complex128) []complex128 {
	err := wallpaperFormula.SetUp()
	if err != nil {
		println(err.Error())
	}

	transformedCoordinates := []complex128{}
	resultsByTerm := [][]complex128{}
	for range wallpaperFormula.Formula.WavePackets {
		resultsByTerm = append(resultsByTerm, []complex128{})
	}

	for _, complexCoordinate := range scaledCoordinates {
		hexagonalWallpaperResults := wallpaperFormula.Calculate(complexCoordinate)
		for index, formulaResult := range hexagonalWallpaperResults.ContributionByTerm {
			resultsByTerm[index] = append(resultsByTerm[index], formulaResult)
		}

		transformedCoordinate := hexagonalWallpaperResults.Total
		transformedCoordinates = append(transformedCoordinates, transformedCoordinate)
	}

	println("Min/Max ranges, by Term")
	for index, results := range resultsByTerm {
		minz, maxz := mathutility.GetBoundingBox(results)
		fmt.Printf("%d: %e - %e\n", index, minz, maxz)
	}

	println("Symmetries found:")
	if wallpaperFormula.HasSymmetry(wavepacket.Cm) {
		println("  cm")
	}
	if wallpaperFormula.HasSymmetry(wavepacket.Cmm) {
		println("  cmm")
	}

	return transformedCoordinates
}

func transformCoordinatesForRectangularWallpaperFormula(wallpaperFormula *wavepacket.RectangularWallpaperFormula, scaledCoordinates []complex128) []complex128 {
	err := wallpaperFormula.SetUp()
	if err != nil {
		println(err.Error())
	}

	transformedCoordinates := []complex128{}
	resultsByTerm := [][]complex128{}
	for range wallpaperFormula.Formula.WavePackets {
		resultsByTerm = append(resultsByTerm, []complex128{})
	}

	for _, complexCoordinate := range scaledCoordinates {
		hexagonalWallpaperResults := wallpaperFormula.Calculate(complexCoordinate)
		for index, formulaResult := range hexagonalWallpaperResults.ContributionByTerm {
			resultsByTerm[index] = append(resultsByTerm[index], formulaResult)
		}

		transformedCoordinate := hexagonalWallpaperResults.Total
		transformedCoordinates = append(transformedCoordinates, transformedCoordinate)
	}

	println("Min/Max ranges, by Term")
	for index, results := range resultsByTerm {
		minz, maxz := mathutility.GetBoundingBox(results)
		fmt.Printf("%d: %e - %e\n", index, minz, maxz)
	}

	println("Symmetries found:")
	hasSymmetry := false
	if wallpaperFormula.HasSymmetry(wavepacket.Pm) {
		println("  pm")
		hasSymmetry = true
	}
	if wallpaperFormula.HasSymmetry(wavepacket.Pg) {
		println("  pg")
		hasSymmetry = true
	}
	if wallpaperFormula.HasSymmetry(wavepacket.Pmm) {
		println("  pmm")
		hasSymmetry = true
	}
	if wallpaperFormula.HasSymmetry(wavepacket.Pmg) {
		println("  pmg")
		hasSymmetry = true
	}
	if wallpaperFormula.HasSymmetry(wavepacket.Pgg) {
		println("  pgg")
		hasSymmetry = true
	}
	if hasSymmetry == false {
		println("  none found")
	}
	return transformedCoordinates
}

func flattenCoordinates(destinationBounds image.Rectangle) []complex128 {
	flattenedCoordinates := []complex128{}
	for destinationY := destinationBounds.Min.Y ; destinationY < destinationBounds.Max.Y; destinationY++ {
		for destinationX := destinationBounds.Min.X; destinationX < destinationBounds.Max.X; destinationX++ {
			flattenedCoordinates = append(flattenedCoordinates, complex(float64(destinationX), float64(destinationY)))
		}
	}
	return flattenedCoordinates
}

func scaleDestinationPixels(destinationBounds image.Rectangle, destinationCoordinates []complex128, viewPortMin complex128, viewPortMax complex128) []complex128 {
	scaledCoordinates := []complex128{}
	for _, destinationCoordinate := range destinationCoordinates {
		destinationScaledX := mathutility.ScaleValueBetweenTwoRanges(
			float64(real(destinationCoordinate)),
			float64(destinationBounds.Min.X),
			float64(destinationBounds.Max.X),
			real(viewPortMin),
			real(viewPortMax),
		)
		destinationScaledY := mathutility.ScaleValueBetweenTwoRanges(
			float64(imag(destinationCoordinate)),
			float64(destinationBounds.Min.Y),
			float64(destinationBounds.Max.Y),
			imag(viewPortMin),
			imag(viewPortMax),
		)
		scaledCoordinates = append(scaledCoordinates, complex(destinationScaledX, destinationScaledY))
	}
	return scaledCoordinates
}

func colorDestinationImage(
	destinationImage *image.NRGBA,
	sourceImage image.Image,
	destinationCoordinates []complex128,
	transformedCoordinates []complex128,
	colorValueBoundMin complex128,
	colorValueBoundMax complex128,
	) {
	sourceImageBounds := sourceImage.Bounds()
	for index, transformedCoordinate := range transformedCoordinates {
		var sourceColorR, sourceColorG, sourceColorB, sourceColorA uint32

		if real(transformedCoordinate) < real(colorValueBoundMin) ||
			imag(transformedCoordinate) < imag(colorValueBoundMin) ||
		real(transformedCoordinate) > real(colorValueBoundMax) ||
			imag(transformedCoordinate) > imag(colorValueBoundMax) {
			sourceColorR,sourceColorG,sourceColorB,sourceColorA = 0,0,0,0
		} else {
			sourceImagePixelX := int(mathutility.ScaleValueBetweenTwoRanges(
				float64(real(transformedCoordinate)),
				real(colorValueBoundMin),
				real(colorValueBoundMax),
				float64(sourceImageBounds.Min.X),
				float64(sourceImageBounds.Max.X),
			))
			sourceImagePixelY := int(mathutility.ScaleValueBetweenTwoRanges(
				float64(imag(transformedCoordinate)),
				imag(colorValueBoundMin),
				imag(colorValueBoundMax),
				float64(sourceImageBounds.Min.Y),
				float64(sourceImageBounds.Max.Y),
			))
			sourceColorR, sourceColorG, sourceColorB, sourceColorA = sourceImage.At(sourceImagePixelX, sourceImagePixelY).RGBA()
		}

		destinationPixelX := int(real(destinationCoordinates[index]))
		destinationPixelY := int(imag(destinationCoordinates[index]))

		destinationImage.Set(
			destinationPixelX,
			destinationPixelY,
			color.NRGBA{
				R: uint8(sourceColorR>>8),
				G: uint8(sourceColorG>>8),
				B: uint8(sourceColorB>>8),
				A: uint8(sourceColorA>>8),
			},
		)
	}
}
