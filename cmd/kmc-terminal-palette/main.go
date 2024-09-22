package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"sort"
)

// Types
type Pixel struct {
	r uint8
	g uint8
	b uint8
}

type Clusters [][]Pixel

type MeansAndClusters struct {
	means    []Pixel
	clusters Clusters
}

type SmallestTotalClustersVariance struct {
	variance         float64
	meansAndClusters MeansAndClusters
}

// Utility
func getPixels(img image.Image) []Pixel {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	pixels := make([]Pixel, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels[y*width+x] = Pixel{
				r: uint8(r >> 8),
				g: uint8(g >> 8),
				b: uint8(b >> 8),
			}
		}
	}

	return pixels
}

func getEuclideanDistance(x Pixel, y Pixel) float64 {
	distance := math.Sqrt(
		math.Pow(float64(x.r)-float64(y.r), 2) +
			math.Pow(float64(x.g)-float64(y.g), 2) +
			math.Pow(float64(x.b)-float64(y.b), 2),
	)
	return distance
}

func getMinimumDistance(distances []float64) int {
	length := len(distances)
	minimumDistanceIndex := 0
	for i := 0; i < length; i++ {
		if distances[i] < distances[minimumDistanceIndex] {
			minimumDistanceIndex = i
		}
	}
	return minimumDistanceIndex
}

func getPixelMean(pixels []Pixel) Pixel {
	length := len(pixels)
	if length == 0 {
		// Return a default Pixel if the cluster is empty
		return Pixel{0, 0, 0}
	}

	var sumRed uint32
	var sumGreen uint32
	var sumBlue uint32
	for _, pixel := range pixels {
		sumRed += uint32(pixel.r)
		sumGreen += uint32(pixel.g)
		sumBlue += uint32(pixel.b)
	}

	meanRed := uint8(sumRed / uint32(length))
	meanGreen := uint8(sumGreen / uint32(length))
	meanBlue := uint8(sumBlue / uint32(length))

	meanPixel := Pixel{
		r: meanRed,
		g: meanGreen,
		b: meanBlue,
	}

	return meanPixel
}

func initializeCentroids(pixels []Pixel, numberOfClusters int) []Pixel {
	centroids := make([]Pixel, numberOfClusters)
	for i := 0; i < numberOfClusters; i++ {
		randomIndex := rand.Intn(len(pixels))
		centroids[i] = pixels[randomIndex]
	}
	return centroids
}

func assignPixelsToClusters(pixels []Pixel, means []Pixel, numberOfClusters int) Clusters {
	clusters := make(Clusters, numberOfClusters)

	for _, pixel := range pixels {
		distances := make([]float64, numberOfClusters)
		for j, mean := range means {
			distances[j] = getEuclideanDistance(pixel, mean)
		}
		nearestMean := getMinimumDistance(distances)
		clusters[nearestMean] = append(clusters[nearestMean], pixel)
	}

	return clusters
}

func calculateNewMeans(clusters Clusters, numberOfClusters int) []Pixel {
	newMeans := make([]Pixel, numberOfClusters)
	for i, cluster := range clusters {
		newMeans[i] = getPixelMean(cluster)
	}
	return newMeans
}

func hasConverged(oldMeans, newMeans []Pixel) bool {
	for i := range oldMeans {
		if oldMeans[i] != newMeans[i] {
			return false
		}
	}
	return true
}

func getClusterVariance(cluster []Pixel) float64 {
	length := len(cluster)
	if length == 0 {
		// Return 0 variance if the cluster is empty
		return 0.0
	}

	meanPixel := getPixelMean(cluster)
	var sumSquaredDistances float64

	for _, pixel := range cluster {
		distance := getEuclideanDistance(pixel, meanPixel)
		sumSquaredDistances += distance * distance
	}

	variance := sumSquaredDistances / float64(length)
	return variance
}

func sumPixel(pixel Pixel) uint32 {
	return uint32(pixel.r) + uint32(pixel.g) + uint32(pixel.b)
}

func sortAscending(pixels []Pixel) {
	sort.Slice(pixels, func(i, j int) bool {
		return sumPixel(pixels[i]) < sumPixel(pixels[j])
	})
}

func rgbToHex(pixel Pixel) string {
	return fmt.Sprintf("#%02x%02x%02x", pixel.r, pixel.g, pixel.b)
}

func rgbToAnsiBackground(pixel Pixel) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", pixel.r, pixel.g, pixel.b)
}

// Core
func cluster(pixels []Pixel, numberOfClusters int) MeansAndClusters {
	means := initializeCentroids(pixels, numberOfClusters)

	for {
		clusters := assignPixelsToClusters(pixels, means, numberOfClusters)
		newMeans := calculateNewMeans(clusters, numberOfClusters)

		if hasConverged(means, newMeans) {
			return MeansAndClusters{means: newMeans, clusters: clusters}
		}

		means = newMeans
	}
}

func kMeansClustering(numberOfClusters int, attempts int, pixels []Pixel) []Pixel {
	initialMeansAndClusters := cluster(pixels, numberOfClusters)

	var totalInitialVariance float64
	for _, cluster := range initialMeansAndClusters.clusters {
		totalInitialVariance += getClusterVariance(cluster)
	}

	var smallestTotalClustersVariance = SmallestTotalClustersVariance{
		variance:         totalInitialVariance,
		meansAndClusters: initialMeansAndClusters,
	}

	for attempt := 0; attempt < attempts; attempt++ {
		meansAndClusters := cluster(pixels, numberOfClusters)

		var totalVariance float64
		for _, cluster := range meansAndClusters.clusters {
			totalVariance += getClusterVariance(cluster)
		}

		if totalVariance < smallestTotalClustersVariance.variance {
			smallestTotalClustersVariance.variance = totalVariance
			smallestTotalClustersVariance.meansAndClusters = meansAndClusters
		}
	}

	meanPixels := smallestTotalClustersVariance.meansAndClusters.means
	return meanPixels
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kmc-terminal-palette <image>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		os.Exit(1)
	}

	pixels := getPixels(image)
	palette := kMeansClustering(8, 1, pixels)

	// Sort the palette in ascending order based on the sum of RGB values
	sortAscending(palette)

	fmt.Println("Palette:")
	for _, pixel := range palette {
		rgbStr := fmt.Sprintf("rgb(%3d, %3d, %3d)", pixel.r, pixel.g, pixel.b)
		hexStr := rgbToHex(pixel)
		colorBlock := rgbToAnsiBackground(pixel) + "    " + "\033[0m"
		fmt.Printf("%-20s %-10s %s\n", rgbStr, hexStr, colorBlock)
	}
}
