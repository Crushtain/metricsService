package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSendMetricsToServer(t *testing.T) {
	// Имитация метрик
	metrics := Metrics{
		Gauges: map[string]float64{
			"Alloc":       100,
			"BuckHashSys": 200,
			// ...
		},
		Counters: map[string]int64{
			"PollCount": 10,
		},
	}

	// Создание тестового сервера
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка содержимого запроса
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Парсинг полученных метрик
		var receivedMetrics Metrics
		err := json.NewDecoder(r.Body).Decode(&receivedMetrics)
		assert.NoError(t, err)

		// Сравнение полученных метрик с ожидаемыми значениями
		assert.Equal(t, metrics.Gauges, receivedMetrics.Gauges)
		assert.Equal(t, metrics.Counters, receivedMetrics.Counters)

		// Отправка ответа
		w.WriteHeader(http.StatusOK)
	}))

	// Подмена адреса сервера для отправки метрик
	oldURL := server.URL + "/metrics"
	newURL := "http://localhost:8080/metrics"
	server.URL = newURL

	// Запуск тестирования
	sendMetricsToServer(metrics)

	// Восстановление оригинального адреса сервера
	server.URL = oldURL
}

func TestSendMetrics(t *testing.T) {
	// Задержка между сбором метрик в секундах
	pollInterval := 1

	// Создание тестового сервера
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Подмена адреса сервера для отправки метрик
	oldURL := server.URL + "/metrics"
	newURL := "http://localhost:8080/metrics"
	server.URL = newURL

	// Запуск отправки метрик на сервер через заданный интервал
	go sendMetrics(pollInterval)

	// Задержка для выполнения нескольких итераций
	time.Sleep(3 * time.Second)

	// Проверка количества отправленных метрик
	resp, err := http.Get(server.URL)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var receivedMetrics Metrics
		err := json.NewDecoder(resp.Body).Decode(&receivedMetrics)
		if assert.NoError(t, err) {
			assert.True(t, receivedMetrics.Counters["PollCount"] > 0)
		}
		resp.Body.Close()
	}

	// Остановка отправки метрик

	// Восстановление оригинального адреса сервера
	server.URL = oldURL
}

func TestMain(m *testing.M) {
	// Запуск основной функции
	go main()

	// Задержка для выполнения нескольких итераций
	time.Sleep(3 * time.Second)

	// Остановка работы главной функции

}
