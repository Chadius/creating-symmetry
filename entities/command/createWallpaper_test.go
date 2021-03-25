package command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"wallpaper/entities/command"
)

var _ = Describe("Commands that create wallpaper", func() {
	Context("Create commands from YAML", func() {
		It("Can create a command using YAML", func() {
			yamlByteStream := []byte(`sample_source_filename: input.png
output_filename: output.png
output_size:
  width: 800
  height: 600
sample_space:
  minx: 0
  miny: 0
  maxx: 1e-10
  maxy: 3e5
color_value_space:
  minx: -50
  miny: 9001
  maxx: -1e-1
  maxy: 2e10
rosette_formula:
  terms:
    -
      multiplier:
        real: -1.0
        imaginary: 2e-2
      power_n: 3
      power_m: 0
      coefficient_pairs:
        multiplier: 1
        relationships:
        - -M-N
        - "+M+NF"
    -
      multiplier:
        real: 1e-10
        imaginary: 0
      power_n: 1
      power_m: 1
      coefficient_pairs:
        multiplier: 1
        relationships:
        - -M-NF
`)
			wallpaperCommand, err := command.NewCreateWallpaperCommandFromYAML(yamlByteStream)
			Expect(err).To(BeNil())
			Expect(wallpaperCommand.SampleSourceFilename).To(Equal("input.png"))

			Expect(wallpaperCommand.OutputFilename).To(Equal("output.png"))
			Expect(wallpaperCommand.OutputImageSize.Width).To(Equal(800))
			Expect(wallpaperCommand.OutputImageSize.Height).To(Equal(600))

			Expect(wallpaperCommand.SampleSpace.MinX).To(BeNumerically("~", 0))
			Expect(wallpaperCommand.SampleSpace.MinY).To(BeNumerically("~", 0))
			Expect(wallpaperCommand.SampleSpace.MaxX).To(BeNumerically("~", 1e-10))
			Expect(wallpaperCommand.SampleSpace.MaxY).To(BeNumerically("~", 3e5))

			Expect(wallpaperCommand.ColorValueSpace.MinX).To(BeNumerically("~", -50))
			Expect(wallpaperCommand.ColorValueSpace.MinY).To(BeNumerically("~", 9001))
			Expect(wallpaperCommand.ColorValueSpace.MaxX).To(BeNumerically("~", -1e-1))
			Expect(wallpaperCommand.ColorValueSpace.MaxY).To(BeNumerically("~", 2e10))

			Expect(wallpaperCommand.RosetteFormula.Terms).To(HaveLen(2))
		})
		It("Can create a command using JSON", func() {
			jsonByteStream := []byte(`{
				"sample_source_filename": "input.png",
				"output_filename": "output.png",
				"output_size": {
				  "width": 800,
				  "height": 600
				},
				"sample_space": {
				  "minx": 0,
				  "miny": 0,
				  "maxx": 1e-10,
				  "maxy": 3e5
				},
				"color_value_space": {
				  "minx": -50,
				  "miny": 9001,
				  "maxx": -1e-1,
				  "maxy": 2e10
				}
			}`)
			wallpaperCommand, err := command.NewCreateWallpaperCommandFromJSON(jsonByteStream)
			Expect(err).To(BeNil())
			Expect(wallpaperCommand.SampleSourceFilename).To(Equal("input.png"))

			Expect(wallpaperCommand.OutputFilename).To(Equal("output.png"))
			Expect(wallpaperCommand.OutputImageSize.Width).To(Equal(800))
			Expect(wallpaperCommand.OutputImageSize.Height).To(Equal(600))

			Expect(wallpaperCommand.SampleSpace.MinX).To(BeNumerically("~", 0))
			Expect(wallpaperCommand.SampleSpace.MinY).To(BeNumerically("~", 0))
			Expect(wallpaperCommand.SampleSpace.MaxX).To(BeNumerically("~", 1e-10))
			Expect(wallpaperCommand.SampleSpace.MaxY).To(BeNumerically("~", 3e5))

			Expect(wallpaperCommand.ColorValueSpace.MinX).To(BeNumerically("~", -50))
			Expect(wallpaperCommand.ColorValueSpace.MinY).To(BeNumerically("~", 9001))
			Expect(wallpaperCommand.ColorValueSpace.MaxX).To(BeNumerically("~", -1e-1))
			Expect(wallpaperCommand.ColorValueSpace.MaxY).To(BeNumerically("~", 2e10))
		})
	})
})