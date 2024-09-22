package models

type Pixel struct {
	R uint8
	G uint8
	B uint8
}

type Clusters [][]Pixel

type MeansAndClusters struct {
	Means    []Pixel
	Clusters Clusters
}

type SmallestTotalClustersVariance struct {
	Variance         float64
	MeansAndClusters MeansAndClusters
}
