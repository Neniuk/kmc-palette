package kmeans

import (
	"log"

	"github.com/Neniuk/kmc-terminal-palette/internal/models"
	"github.com/Neniuk/kmc-terminal-palette/internal/utils"
)

func AssignPixelsToClusters(pixels []models.Pixel, means []models.Pixel, numberOfClusters int) models.Clusters {
	clusters := make(models.Clusters, numberOfClusters)

	for _, pixel := range pixels {
		distances := make([]float64, numberOfClusters)
		for j, mean := range means {
			distances[j] = utils.CalculateEuclideanDistance(pixel, mean)
		}
		nearestMean := utils.FindMinimumDistanceIndex(distances)
		clusters[nearestMean] = append(clusters[nearestMean], pixel)
	}

	return clusters
}

func CalculateNewMeans(clusters models.Clusters, numberOfClusters int) []models.Pixel {
	newMeans := make([]models.Pixel, numberOfClusters)
	for i, cluster := range clusters {
		newMeans[i] = utils.CalculateMeanPixel(cluster)
	}
	return newMeans
}

func HasConverged(oldMeans, newMeans []models.Pixel) bool {
	for i := range oldMeans {
		if oldMeans[i] != newMeans[i] {
			return false
		}
	}
	return true
}

func Cluster(pixels []models.Pixel, numberOfClusters int, maxIterations int) models.MeansAndClusters {
	means := utils.InitializeRandomCentroids(pixels, numberOfClusters)
	log.Printf("Means: %v", means)

	for i := 0; i < maxIterations; i++ {
		log.Printf("Iteration %d: Assigning pixels to clusters...", i)

		clusters := AssignPixelsToClusters(pixels, means, numberOfClusters)
		newMeans := CalculateNewMeans(clusters, numberOfClusters)
		log.Printf("Iteration %d: means = %v", i, newMeans)

		if HasConverged(means, newMeans) {
			log.Printf("Means have converged.")
			return models.MeansAndClusters{Means: newMeans, Clusters: clusters}
		}

		means = newMeans
	}

	// Return the result after maxIterations
	log.Printf("Means have not converged, exceeded maximum number of iterations.")
	clusters := AssignPixelsToClusters(pixels, means, numberOfClusters)
	return models.MeansAndClusters{Means: means, Clusters: clusters}
}

func KMeansClustering(numberOfClusters int, maxIterations int, pixels []models.Pixel) []models.Pixel {
	meansAndClusters := Cluster(pixels, numberOfClusters, maxIterations)
	meanPixels := meansAndClusters.Means
	return meanPixels
}
