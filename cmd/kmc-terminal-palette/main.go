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
)

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

	var sumRed uint8
	var sumGreen uint8
	var sumBlue uint8
	for _, pixel := range pixels {
		sumRed += pixel.r
		sumGreen += pixel.g
		sumBlue += pixel.b
	}

	meanRed := sumRed / uint8(length)
	meanGreen := sumGreen / uint8(length)
	meanBlue := sumBlue / uint8(length)

	meanPixel := Pixel{
		r: meanRed,
		g: meanGreen,
		b: meanBlue,
	}

	return meanPixel
}

func cluster(pixels []Pixel, numberOfClusters int) MeansAndClusters {
	var dataPoints []Pixel
	length := len(pixels)

	clusters := make(Clusters, numberOfClusters)
	means := make([]Pixel, numberOfClusters)

	// Get starting data points
	for i := 0; i < numberOfClusters; i++ {
		// Generate a random number between 0 and pixels length
		randomNumber := rand.Intn(length)
		dataPoints[i] = pixels[randomNumber]
	}

	// Sort into clusters by initial data points
	for i := 0; i < length; i++ {
		var distances []float64
		for j := 0; j < numberOfClusters; j++ {
			distances = append(distances, getEuclideanDistance(pixels[i], dataPoints[j]))
		}

		nearestDataPoint := getMinimumDistance(distances)
		clusters[nearestDataPoint] = append(clusters[nearestDataPoint], pixels[i])
	}

	// Get mean pixel/vector
	for index, cluster := range clusters {
		means[index] = getPixelMean(cluster)
	}

	meanClusters := make(Clusters, numberOfClusters)
	// Sort into clusters by means
	for i := 0; i < length; i++ {
		var distances []float64
		for j := 0; j < numberOfClusters; j++ {
			distances = append(distances, getEuclideanDistance(pixels[i], means[j]))
		}

		nearestMean := getMinimumDistance(distances)
		meanClusters[nearestMean] = append(meanClusters[nearestMean], pixels[i])
	}

	var meansAndClusters = MeansAndClusters{
		means:    means,
		clusters: meanClusters,
	}

	return meansAndClusters
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

	fmt.Println("Palette:")
	for _, pixel := range palette {
		fmt.Printf("rgb(%d, %d, %d)\n", pixel.r, pixel.g, pixel.b)
	}
}
