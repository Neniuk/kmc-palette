package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/Neniuk/kmc-terminal-palette/internal/kmeans"
	"github.com/Neniuk/kmc-terminal-palette/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kmc-terminal-palette <image>")
		os.Exit(1)
	}

	log.Printf("Correct number of args")

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	log.Printf("Opened file")

	image, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		os.Exit(1)
	}

	log.Printf("Decoded image")

	pixels := utils.GetPixels(image)
	log.Printf("Pixels fetched: %d", len(pixels))

	palette := kmeans.KMeansClustering(8, 2, pixels)

	// Sort the palette in ascending order based on the sum of RGB values
	log.Printf("Palette ready, sorting...")
	utils.SortPixelsByBrightness(palette)

	fmt.Println("Palette:")
	for _, pixel := range palette {
		rgbStr := fmt.Sprintf("rgb(%3d, %3d, %3d)", pixel.R, pixel.G, pixel.B)
		hexStr := utils.ConvertRgbToHex(pixel)
		colorBlock := utils.ConvertRgbToAnsiBackground(pixel) + "    " + "\033[0m"
		fmt.Printf("%-20s %-10s %s\n", rgbStr, hexStr, colorBlock)
	}
}
