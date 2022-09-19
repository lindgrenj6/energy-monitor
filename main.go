package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// probably should set up DNS for this eventually
const metricsUrl = "http://10.0.0.11/index.php/realtimedata/old_power_graph"

// data of the last scrape
var last *RealTimeData

func main() {
	go func() {
		// first time populating the data
		last = fetchMetrics()

		// once a minute, update the current metrics in-memory
		for range time.Tick(time.Minute) {
			log.Println("Getting metrics...")
			last = fetchMetrics()
			log.Printf("Got metrics, current power: %v\n", last.Power[len(last.Power)-1].EachSystemPower)
		}
	}()

	log.Println("Starting server on :9394")
	http.Handle("/metrics", http.HandlerFunc(currentMetrics))
	http.ListenAndServe(":9394", nil)
}

func currentMetrics(w http.ResponseWriter, r *http.Request) {
	if last == nil {
		http.Error(w, "Data not ready yet", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(last.ToMetrics()))
}

// pull the data from the ECU and parse into the RealTimeData struct
func fetchMetrics() *RealTimeData {
	resp := must(http.Get(metricsUrl))
	if resp.StatusCode != 200 {
		panic("bad status code from metrics: " + resp.Status)
	}
	body := must(io.ReadAll(resp.Body))

	var current RealTimeData
	err := json.Unmarshal(body, &current)
	if err != nil {
		log.Fatalf("bad data came back from ECU: %v\n", err)
	}

	return &current
}
