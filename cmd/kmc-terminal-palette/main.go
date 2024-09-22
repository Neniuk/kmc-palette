package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/Neniuk/kmc-terminal-palette/internal/kmeans"
	"github.com/Neniuk/kmc-terminal-palette/internal/models"
	"github.com/Neniuk/kmc-terminal-palette/internal/utils"
)

const (
	defaultNumberOfIterations = 1
	defaultNumberOfClusters   = 8
)

func main() {
	var iterations int
	var numberOfClusters int
	var skipColorVariants bool
	var hideHorizontalPalette bool

	flag.IntVar(&iterations, "i", defaultNumberOfIterations, "number of iterations for K-means clustering")
	flag.IntVar(&numberOfClusters, "k", defaultNumberOfClusters, "number of clusters for K-means clustering")
	flag.BoolVar(&skipColorVariants, "scv", false, "skip adding color variants to the palette")
	flag.BoolVar(&hideHorizontalPalette, "hhp", false, "hide horizontal palette")
	flag.Parse()

	// Check if the user has provided the correct number of arguments
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: kmc-terminal-palette [-i <number>] <image>")
		os.Exit(1)
	}

	log.Printf("Number of iterations set to %d", iterations)
	log.Printf("Number of clusters set to %d", numberOfClusters)

	// Open the image file for reading
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()
	log.Printf("Opened file: %s", flag.Args()[0])

	// Decode the image
	image, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		os.Exit(1)
	}
	log.Printf("Decoded image: %s", flag.Args()[0])

	pixels := utils.GetPixels(image)
	log.Printf("Pixels fetched: %d", len(pixels))

	// Perform K-means clustering
	palette := kmeans.KMeansClustering(numberOfClusters, iterations, pixels)

	// Sort the palette in ascending order based on the sum of RGB values
	log.Printf("Palette ready, sorting...")
	utils.SortPixelsByBrightness(palette)

	if !skipColorVariants {
		// Add color variations to the palette
		log.Printf("Adding color variations...")
		utils.AddColorVariantsToPalette(&palette)
	} else {
		log.Printf("Skipping adding color variations...")
	}

	fmt.Println("Palette:")
	for _, pixel := range palette {
		rgbStr := fmt.Sprintf("rgb(%3d, %3d, %3d)", pixel.R, pixel.G, pixel.B)
		hexStr := utils.ConvertRgbToHex(pixel)
		colorBlock := utils.ConvertRgbToAnsiBackground(pixel) + "    " + "\033[0m"
		fmt.Printf("%-20s %-10s %s\n", rgbStr, hexStr, colorBlock)
	}

	if !hideHorizontalPalette {
		// Separate the palette into darker and lighter variants
		darkerVariants := make([]models.Pixel, 0)
		lighterVariants := make([]models.Pixel, 0)
		for i, pixel := range palette {
			if i%2 == 0 {
				darkerVariants = append(darkerVariants, pixel)
				lighterVariants = append(lighterVariants, utils.GenerateLighterPixelVariant(pixel))
			} else {
				darkerVariants = append(darkerVariants, utils.GenerateDarkerPixelVariant(pixel))
				lighterVariants = append(lighterVariants, pixel)
			}
		}

		// Print color blocks horizontally
		fmt.Println()
		for _, pixel := range darkerVariants {
			colorBlock := utils.ConvertRgbToAnsiBackground(pixel) + "  " + "\033[0m"
			fmt.Print(colorBlock)
		}
		fmt.Println()
		for _, pixel := range lighterVariants {
			colorBlock := utils.ConvertRgbToAnsiBackground(pixel) + "  " + "\033[0m"
			fmt.Print(colorBlock)
		}
		fmt.Println()
	}
}
