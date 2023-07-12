package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Metrics struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

func main() {
	pollInterval := 10
	go sendMetrics(pollInterval)
	for {
		metrics := collectMetrics()
		metrics.Gauges["RandomValue"] = float64(rand.Intn(100))
		sendMetricsToServer(metrics)
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func collectMetrics() Metrics {
	metrics := Metrics{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	metrics.Gauges["Alloc"] = float64(memStats.Alloc)
	metrics.Gauges["BuckHashSys"] = float64(memStats.BuckHashSys)
	metrics.Gauges["Frees"] = float64(memStats.Frees)
	metrics.Gauges["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	metrics.Gauges["GCSys"] = float64(memStats.GCSys)
	metrics.Gauges["HeapAlloc"] = float64(memStats.HeapAlloc)
	metrics.Gauges["HeapIdle"] = float64(memStats.HeapIdle)
	metrics.Gauges["HeapInuse"] = float64(memStats.HeapInuse)
	metrics.Gauges["HeapObjects"] = float64(memStats.HeapObjects)
	metrics.Gauges["HeapReleased"] = float64(memStats.HeapReleased)
	metrics.Gauges["HeapSys"] = float64(memStats.HeapSys)
	metrics.Gauges["LastGC"] = float64(memStats.LastGC)
	metrics.Gauges["Lookups"] = float64(memStats.Lookups)
	metrics.Gauges["MCacheInuse"] = float64(memStats.MCacheInuse)
	metrics.Gauges["MCacheSys"] = float64(memStats.MCacheSys)
	metrics.Gauges["MSpanInuse"] = float64(memStats.MSpanInuse)
	metrics.Gauges["MSpanSys"] = float64(memStats.MSpanSys)
	metrics.Gauges["Mallocs"] = float64(memStats.Mallocs)
	metrics.Gauges["NextGC"] = float64(memStats.NextGC)
	metrics.Gauges["NumForcedGC"] = float64(memStats.NumForcedGC)
	metrics.Gauges["NumGC"] = float64(memStats.NumGC)
	metrics.Gauges["OtherSys"] = float64(memStats.OtherSys)
	metrics.Gauges["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	metrics.Gauges["StackInuse"] = float64(memStats.StackInuse)
	metrics.Gauges["StackSys"] = float64(memStats.StackSys)
	metrics.Gauges["Sys"] = float64(memStats.Sys)
	metrics.Gauges["TotalAlloc"] = float64(memStats.TotalAlloc)

	return metrics
}

func sendMetricsToServer(metrics Metrics) {

	_, err := json.Marshal(metrics)
	if err != nil {
		log.Println("Ошибка при преобразовании метрик в JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080", "application/json", nil)
	if err != nil {
		log.Println("Ошибка при отправке метрик на сервер:", err)
		return
	}
	defer resp.Body.Close()
}

func sendMetrics(pollInterval int) {

	for {

		metrics := collectMetrics()

		metrics.Counters["PollCount"]++

		sendMetricsToServer(metrics)

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
