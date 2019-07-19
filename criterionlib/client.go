package criterionlib

// Client - A client used for interacting with your distributed services.
type Client interface {
	Run(Cluster, float64, bool)
}
