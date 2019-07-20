package criterionlib

// EmptyClient - A client that just prints.
type EmptyClient struct{}

// Run - Print the arguments.
func (EmptyClient) Run(cluster Cluster, dropRate float64, randomFailures bool) Metric {
	return NewMetric("test", 1)
}
