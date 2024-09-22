package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
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

var errUsage = errors.New("usage: ./kmc-palette [-i <number>] [-k <number>] [-scv] [-hhp] <filepath-to-image>")

func printPalette(palette []models.Pixel) {
	fmt.Println("\nPalette:")

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
	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelError,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

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

	slog.Info("Number of iterations set", slog.Int("iterations", iterations))
	slog.Info("Number of clusters set", slog.Int("clusters", numberOfClusters))

	// Open the image file for reading
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		slog.Error("Error opening file", slog.String("error", err.Error()))
		return fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()
	fmt.Printf("Opened file %s\n", flag.Args()[0])

	// Decode the image
	image, _, err := image.Decode(file)
	if err != nil {
		slog.Error("Error decoding image", slog.String("error", err.Error()))
		return fmt.Errorf("error decoding image: %w", err)
	}

	fmt.Println("Image decoded")

	// Get the pixels from the image
	pixels := utils.GetPixels(image)
	fmt.Printf("Pixels fetched: %d\n", len(pixels))

	// Perform K-means clustering
	palette := kmeans.KMeansClustering(numberOfClusters, iterations, pixels)

	// Sort the palette in ascending order based on the sum of RGB values
	fmt.Println("Palette ready, sorting...")
	utils.SortPixelsByBrightness(palette)

	if !skipColorVariants {
		// Add color variations to the palette
		fmt.Println("Adding color variations...")
		utils.AddColorVariantsToPalette(&palette)
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
