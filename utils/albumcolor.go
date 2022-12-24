package utils

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/zmb3/spotify/v2"
)

// Gets the prominent colour from an album cover
func GetAlbumColour(album spotify.SimpleAlbum) (colorstring string, colorVal color.RGBA) {
	var c *os.File
	c, _ = os.CreateTemp("", "album.jpg")
	defer c.Close()
	defer os.Remove(c.Name())
	// Images is an array of 3 images of different
	// sizes - from 300x300 to 64x64.
	album.Images[0].Download(c)
	t, err := loadImage(c.Name())
	if err != nil {
		log.Fatalln(err)
	}
	// Perform Kmeans without cropping to find promiment colour
	cols, err := prominentcolor.KmeansWithArgs(prominentcolor.ArgumentNoCropping, t)
	if err != nil {
		// Donda moment.
		// A fully black picture has no non-alpha pixels and returns an error
		colorstring = "#000000"
		colorVal = color.RGBA{0, 0, 0, 255}
	} else {
		domcolour := cols[0]
		colorstring = "#" + domcolour.AsString()
		colorVal = color.RGBA{
			uint8(domcolour.Color.R),
			uint8(domcolour.Color.G),
			uint8(domcolour.Color.B),
			255,
		}
	}
	return colorstring, colorVal
}

// Load image from disk
func loadImage(fileInput string) (image.Image, error) {
	f, err := os.Open(fileInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	if err != nil {
		log.Println("File not found:", fileInput)
		return nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
