package criterionlib

// Config - Maintains Criterion configuration information.
type Config struct {
	Clusters       []Cluster
	DropRates      []float64
	NumClients     []int64
	RandomFailures bool
}
