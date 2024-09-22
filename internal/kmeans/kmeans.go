package kmeans

import (
	"fmt"
	"log/slog"

	"github.com/Neniuk/kmc-palette/internal/models"
	"github.com/Neniuk/kmc-palette/internal/utils"
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
	slog.Info("Initial means", slog.Any("means", means))

	for iteration := range maxIterations {
		fmt.Printf("[%d] Assigning pixels to clusters...\n", iteration)

		clusters := AssignPixelsToClusters(pixels, means, numberOfClusters)
		newMeans := CalculateNewMeans(clusters, numberOfClusters)

		slog.Info("Iteration means", slog.Int("iteration", iteration), slog.Any("means", newMeans))

		if HasConverged(means, newMeans) {
			slog.Info("Means have converged.")

			return models.MeansAndClusters{Means: newMeans, Clusters: clusters}
		}

		means = newMeans
	}

	// Return the result after maxIterations
	slog.Info("Means have not converged, exceeded maximum number of iterations.")
	clusters := AssignPixelsToClusters(pixels, means, numberOfClusters)

	return models.MeansAndClusters{Means: means, Clusters: clusters}
}

func KMeansClustering(numberOfClusters int, maxIterations int, pixels []models.Pixel) []models.Pixel {
	meansAndClusters := Cluster(pixels, numberOfClusters, maxIterations)
	meanPixels := meansAndClusters.Means

	return meanPixels
}
