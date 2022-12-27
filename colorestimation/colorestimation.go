package colorestimation

import (
	"image/color"

	"github.com/mattn/go-ciede2000"
)

type cs struct {
	Name   string
	colors []color.RGBA
}

// Each colour set has been represnted as a basket
// which contains a one or more colours.
// We pick the colour set which has the element
// that is closest to the input colour according
// to the CIED2000 algorithm.
// Finetuning using this approach involves
// Playing with the colours in each colour set
// until satisfactory results are obtained
var (
	// red     = []color.RGBA{{255, 0, 0, 255}}
	// green   = []color.RGBA{{0, 255, 0, 255}, {34, 139, 34, 255}, {72, 160, 84, 255}}
	// blue    = []color.RGBA{{0, 0, 205, 255}, {135, 206, 250, 255}}
	// yellow  = []color.RGBA{{255, 255, 0, 25}}
	// pink    = []color.RGBA{{255, 0, 255, 25}, {230, 170, 190, 255}}
	// bw      = []color.RGBA{{255, 255, 255, 255}, {0, 0, 0, 255}}
	// orange  = []color.RGBA{{255, 69, 0, 255}}
	red     = []color.RGBA{{255, 0, 0, 255}}
	green   = []color.RGBA{{0, 255, 0, 255}, {0, 150, 0, 255}}
	blue    = []color.RGBA{{0, 0, 255, 255}, {0, 255, 255, 255}}
	yellow  = []color.RGBA{{255, 255, 0, 255}}
	pink    = []color.RGBA{{255, 0, 255, 255}, {230, 170, 190, 255}}
	bw      = []color.RGBA{{255, 255, 255, 255}, {0, 0, 0, 255}}
	orange  = []color.RGBA{{255, 69, 0, 255}}
	colours = []cs{
		{"red", red},
		{"green", green},
		{"blue", blue},
		{"yellow", yellow},
		{"pink/purple", pink},
		{"bw", bw},
		{"orange", orange},
	}
)

// Find closest distance using CIED2000 algorithm
func FindClosest(inputColor color.RGBA) (closest string) {
	closest = "bw"
	diff := ciede2000.Diff(inputColor, bw[0])
	for _, colorset := range colours {
		for _, currColor := range colorset.colors {
			currDiff := ciede2000.Diff(currColor, inputColor)
			if currDiff < diff {
				closest = colorset.Name
				diff = currDiff
			}
		}
	}
	return
}
