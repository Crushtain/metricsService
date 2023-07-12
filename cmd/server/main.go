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

func handleValue(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Не тот метод", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/value/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(res, "Не полный запрос", http.StatusBadRequest)
		return
	}

	metricType := parts[0]
	metricName := parts[1]

	value, ok := storage.metrics[metricName]
	if !ok {
		http.Error(res, "Не найдено", http.StatusNotFound)
	}

	response := fmt.Sprintf("Current value of %s metric (%s): %s", metricType, metricName, value)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(response))
}

func handleMain(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body := `<html>
<head>
<title>Metric List</title>
</head>
<body>
<h1>Metric List</h1>
<table>
<tr>
<th>Name</th>
<th>Value</th>
</tr>`

	for metricType, metricMap := range storage.metrics {
		for metricName, metricValue := range metricMap {
			row := fmt.Sprintf(`<tr><td>%s</td><td>%v</td></tr>`, metricType+"/"+metricName, metricValue)
			body += row
		}
	}

	body += `</table>
</body>
</html>`

	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(body))
}
func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Обработка добавления новой метрики

	response := "New metric added"

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func main() {
	storage.metrics = make(map[string]map[string]interface{})
	http.HandleFunc("/add/", handleAdd)
	http.HandleFunc("/value/counter/123", handleMain)
	http.HandleFunc("/value/", handleValue)
	http.HandleFunc("/update/", metricsHandler)
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := storage.getMetrics()
		fmt.Fprint(w, metrics)
	})

	http.ListenAndServe("localhost:8080", nil)
}
