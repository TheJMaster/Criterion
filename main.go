package main

// TODO: Fix any relative imports.
import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"../Criterion/criterionlib"
)

var (
	configFilePath = "config.json"
)

func main() {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("could not open config file: %s\n", err.Error())
	}

	configJSON, err := ioutil.ReadAll(bufio.NewReader(configFile))
	if err != nil {
		log.Fatalf("could not read config file: %s\n", err.Error())
	}

	var config criterionlib.Config
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		log.Fatalf("could not unmarshal json from config file: %s\n", err.Error())
	}

	_ = os.Mkdir("./out", 0777)
	t := time.Now()
	outFile := fmt.Sprintf("out/%s.csv", t.Format("2006-01-02_15:04:05"))
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("could not create output file: %s\n", err.Error())
	}

	var wg sync.WaitGroup
	client := criterionlib.EmptyClient{}
	for i := int64(1); i <= config.Iterations; i++ {
		for _, cluster := range config.Clusters {
			for _, dropRate := range config.DropRates {
				for _, numClients := range config.NumClients {
					var metrics []criterionlib.Metric
					for i := int64(1); i <= numClients; i++ {
						wg.Add(1)
						go func() {
							defer wg.Done()
							metric := client.Run(cluster, dropRate, false)
							metrics = append(metrics, metric)
						}()
					}
					wg.Wait()
					err := writeMetric(f, cluster.Name, dropRate, numClients, false, averageMetric(metrics))
					if err != nil {
						log.Fatalf("error writing metric: %s\n", err.Error())
					}

					metrics = nil
					if config.RandomFailures {
						for i := int64(1); i <= numClients; i++ {
							wg.Add(1)
							go func() {
								defer wg.Done()
								metric := client.Run(cluster, dropRate, true)
								metrics = append(metrics, metric)
							}()
						}
						wg.Wait()
						err := writeMetric(f, cluster.Name, dropRate, numClients, true, averageMetric(metrics))
						if err != nil {
							log.Fatalf("error writing metric: %s\n", err.Error())
						}
					}
				}
			}
		}
	}
}

func writeMetric(f *os.File, clusterName string, dropRate float64, numClients int64, randomFailures bool, metric criterionlib.Metric) error {
	metricStr := fmt.Sprintf("%s,%f,%d,%v,%s,%f\n", clusterName, dropRate, numClients, randomFailures, metric.Name(), metric.Value())
	n, err := f.WriteString(metricStr)
	if err != nil {
		return err
	} else if n != len(metricStr) {
		return errors.New("did not write entire metric string")
	}
	return nil
}

func averageMetric(metrics []criterionlib.Metric) criterionlib.Metric {
	total := 0.0
	for _, metric := range metrics {
		total += metric.Value()
	}
	return criterionlib.NewMetric(metrics[0].Name(), total/float64(len(metrics)))
}
