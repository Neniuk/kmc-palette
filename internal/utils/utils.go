package utils

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"sort"

	"github.com/Neniuk/kmc-terminal-palette/internal/models"
)

func GetPixels(image image.Image) []models.Pixel {
	bounds := image.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	pixels := make([]models.Pixel, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			red, green, blue, _ := image.At(x, y).RGBA()
			pixels[y*width+x] = models.Pixel{
				R: uint8(red >> 8),
				G: uint8(green >> 8),
				B: uint8(blue >> 8),
			}
		}
	}

	return pixels
}

func CalculateEuclideanDistance(x models.Pixel, y models.Pixel) float64 {
	return math.Sqrt(
		math.Pow(float64(x.R)-float64(y.R), 2) +
			math.Pow(float64(x.G)-float64(y.G), 2) +
			math.Pow(float64(x.B)-float64(y.B), 2),
	)
}

func FindMinimumDistanceIndex(distances []float64) int {
	minimumDistanceIndex := 0
	for i, distance := range distances {
		if distance < distances[minimumDistanceIndex] {
			minimumDistanceIndex = i
		}
	}
	return minimumDistanceIndex
}

func CalculateMeanPixel(pixels []models.Pixel) models.Pixel {
	length := len(pixels)
	if length == 0 {
		return models.Pixel{R: 0, G: 0, B: 0}
	}

	var sumRed, sumGreen, sumBlue uint32
	for _, pixel := range pixels {
		sumRed += uint32(pixel.R)
		sumGreen += uint32(pixel.G)
		sumBlue += uint32(pixel.B)
	}

	return models.Pixel{
		R: uint8(sumRed / uint32(length)),
		G: uint8(sumGreen / uint32(length)),
		B: uint8(sumBlue / uint32(length)),
	}
}

func GenerateLighterPixelVariant(pixel models.Pixel) models.Pixel {
	return models.Pixel{
		R: uint8(math.Min(float64(pixel.R)+30, 255)),
		G: uint8(math.Min(float64(pixel.G)+30, 255)),
		B: uint8(math.Min(float64(pixel.B)+30, 255)),
	}
}

func GenerateDarkerPixelVariant(pixel models.Pixel) models.Pixel {
	return models.Pixel{
		R: uint8(math.Max(float64(pixel.R)-30, 0)),
		G: uint8(math.Max(float64(pixel.G)-30, 0)),
		B: uint8(math.Max(float64(pixel.B)-30, 0)),
	}
}

// AddColorVariantsToPalette adds lighter or darker variants of the colors to the palette based on brightness.
func AddColorVariantsToPalette(palette *[]models.Pixel) {
	originalPalette := *palette
	newPalette := make([]models.Pixel, 0, len(originalPalette)*2)

	for _, pixel := range originalPalette {
		brightness := CalculatePixelSum(pixel)

		if brightness < 128 {
			// Original color is dark, generate a lighter variant
			lighterPixel := GenerateLighterPixelVariant(pixel)
			// Add the original and lighter variant to the new palette
			newPalette = append(newPalette, pixel, lighterPixel)
		} else {
			// Original color is light, generate a darker variant
			darkerPixel := GenerateDarkerPixelVariant(pixel)
			// Add the original and darker variant to the new palette
			newPalette = append(newPalette, darkerPixel, pixel)
		}
	}
	*palette = newPalette
}

func InitializeRandomCentroids(pixels []models.Pixel, numberOfClusters int) []models.Pixel {
	centroids := make([]models.Pixel, numberOfClusters)
	for i := 0; i < numberOfClusters; i++ {
		randomIndex := rand.Intn(len(pixels))
		centroids[i] = pixels[randomIndex]
	}
	return centroids
}

func CalculatePixelSum(pixel models.Pixel) uint32 {
	return uint32(pixel.R) + uint32(pixel.G) + uint32(pixel.B)
}

func SortPixelsByBrightness(pixels []models.Pixel) {
	sort.Slice(pixels, func(i, j int) bool {
		return CalculatePixelSum(pixels[i]) < CalculatePixelSum(pixels[j])
	})
}

func ConvertRgbToHex(pixel models.Pixel) string {
	return fmt.Sprintf("#%02x%02x%02x", pixel.R, pixel.G, pixel.B)
}

func ConvertRgbToAnsiBackground(pixel models.Pixel) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", pixel.R, pixel.G, pixel.B)
}
