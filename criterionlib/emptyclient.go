package criterionlib

import (
	"fmt"
)

// EmptyClient - A client that just prints.
type EmptyClient struct{}

// Run - Print the arguments.
func (EmptyClient) Run(cluster Cluster, dropRate float64, randomFailures bool) {
	fmt.Printf("%s/%f/%v\n", cluster.Name, dropRate, randomFailures)
}
