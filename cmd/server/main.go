package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type MemStorage struct {
	sync.Mutex
	metrics map[string]map[string]interface{}
}

func (m *MemStorage) gauge(name string, value float64) {
	m.Lock()
	defer m.Unlock()

	if m.metrics["gauge"] == nil {
		m.metrics["gauge"] = make(map[string]interface{})
	}

	m.metrics["gauge"][name] = value
}

func (m *MemStorage) counter(name string, value int64) {
	m.Lock()
	defer m.Unlock()

	if m.metrics["counter"] == nil {
		m.metrics["counter"] = make(map[string]interface{})
	}

	existingValue, ok := m.metrics["counter"][name].(int64)
	if ok {
		value += existingValue
	}

	m.metrics["counter"][name] = value
}

func (m *MemStorage) getMetrics() map[string]map[string]interface{} {
	m.Lock()
	defer m.Unlock()

	return m.metrics
}

var storage MemStorage

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	path := r.URL.Path[len("/update/"):]
	params := strings.Split(path, "/")

	if len(params) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := params[0]
	metricName := params[1]
	metricValue := params[2]

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.gauge(metricName, value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.counter(metricName, value)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	storage.metrics = make(map[string]map[string]interface{})

	http.HandleFunc("/update/", metricsHandler)
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := storage.getMetrics()
		fmt.Fprint(w, metrics)
	})

	http.ListenAndServe(":8080", nil)
}
