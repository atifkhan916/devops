package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// filePath is the path to the file containing the metrics data
const filePath = "data/metrics_from_special_app.txt"

// cacheTTL is the amount of time the metrics data should be cached in memory
const cacheTTL = 30 * time.Second

var (
	// lastCacheTime is the time the cache was last updated
	lastCacheTime time.Time
	// metricsCache contains the cached metrics data
	metricsCache string
)

// getMetrics reads the metrics data from the file specified in filePath.
// It uses caching to limit disk I/O by storing the metrics data in memory
// and only re-reading the file if the cache is stale (based on the cacheTTL constant).
func getMetrics() (string, error) {
	now := time.Now()
	if now.Sub(lastCacheTime) < cacheTTL {
		return metricsCache, nil
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	metrics := string(data)
	metricsCache = metrics
	lastCacheTime = now
	return metrics, nil
}

// metricsHandler is the handler function for GET requests on the path "/metrics".
// It reads the metrics data using the getMetrics function and returns it as a string
// with key-value fields separated by line breaks.
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/metrics" {
		http.NotFound(w, r)
		return
	}
	metrics, err := getMetrics()
	if err != nil {
		http.Error(w, "Failed to read metrics data", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, metrics)
}

func main() {
	// Register the metricsHandler function to handle requests on the path "/metrics".
	http.HandleFunc("/metrics", metricsHandler)
	// Start the server and listen on port 12345.
	http.ListenAndServe(":12345", nil)
}
