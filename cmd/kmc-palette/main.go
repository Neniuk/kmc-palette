package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/Neniuk/kmc-palette/internal/kmeans"
	"github.com/Neniuk/kmc-palette/internal/models"
	"github.com/Neniuk/kmc-palette/internal/utils"
)

const (
	defaultNumberOfIterations = 1
	defaultNumberOfClusters   = 8
	resetAnsiCode             = "\033[0m"
)

var errUsage = errors.New("usage: kmc-terminal-palette [-i <number>] <image>")

func printPalette(palette []models.Pixel) {
	fmt.Println("Palette:")

	for _, pixel := range palette {
		rgbStr := fmt.Sprintf("rgb(%3d, %3d, %3d)", pixel.R, pixel.G, pixel.B)
		hexStr := utils.ConvertRgbToHex(pixel)
		colorBlock := utils.ConvertRgbToAnsiBackground(pixel) + "    " + resetAnsiCode
		fmt.Printf("%-20s %-10s %s\n", rgbStr, hexStr, colorBlock)
	}
}

func printHorizontalPalette(palette []models.Pixel) {
	darkerVariants, lighterVariants := separatePalette(palette)

	fmt.Println()

	for _, pixel := range darkerVariants {
		colorBlock := utils.ConvertRgbToAnsiBackground(pixel) + "  " + resetAnsiCode
		fmt.Print(colorBlock)
	}

	fmt.Println()

	for _, pixel := range lighterVariants {
		colorBlock := utils.ConvertRgbToAnsiBackground(pixel) + "  " + resetAnsiCode
		fmt.Print(colorBlock)
	}

	fmt.Println()
}

func separatePalette(palette []models.Pixel) ([]models.Pixel, []models.Pixel) {
	darkerVariants := make([]models.Pixel, 0)
	lighterVariants := make([]models.Pixel, 0)

	for i, pixel := range palette {
		if i%2 == 0 {
			darkerVariants = append(darkerVariants, pixel)
		} else {
			lighterVariants = append(lighterVariants, pixel)
		}
	}

	return darkerVariants, lighterVariants
}

func run() error {
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
		return errUsage
	}

	log.Printf("Number of iterations set to %d", iterations)
	log.Printf("Number of clusters set to %d", numberOfClusters)

	// Open the image file for reading
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()
	log.Printf("Opened file: %s", flag.Args()[0])

	// Decode the image
	image, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	log.Printf("Decoded image: %s", flag.Args()[0])

	// Get the pixels from the image
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

	// Print the palette
	printPalette(palette)

	if !hideHorizontalPalette {
		printHorizontalPalette(palette)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
