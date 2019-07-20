package main

// TODO: Fix any relative imports.
import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"

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

	var wg sync.WaitGroup
	client := criterionlib.EmptyClient{}
	var metrics []criterionlib.Metric
	for i := int64(1); i <= config.Iterations; i++ {
		for _, cluster := range config.Clusters {
			for _, dropRate := range config.DropRates {
				for _, numClients := range config.NumClients {
					for i := int64(1); i <= numClients; i++ {
						wg.Add(1)
						go func() {
							defer wg.Done()
							metric := client.Run(cluster, dropRate, false)
							metrics = append(metrics, metric)
						}()
					}
					if config.RandomFailures {
						for i := int64(1); i <= numClients; i++ {
							wg.Add(1)
							go func() {
								defer wg.Done()
								metric := client.Run(cluster, dropRate, true)
								metrics = append(metrics, metric)
							}()
						}
					}
				}
			}
		}
	}
	wg.Wait()
}
