package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	cpuTemp prometheus.Gauge
}

func NewMetrics(reg *prometheus.Registry) *metrics {
	m := &metrics{
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temp",
			Help: "CPU temperature",
		}),
	}
	reg.MustRegister(m.cpuTemp)
	return m
}

func main() {

	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)
	// Set values for the new created metrics.
	m.cpuTemp.Set(65.3)

	//expose default go metrics
	// http.Handle("/metrics", promhttp.Handler())
	// http.ListenAndServe(":8080", nil)

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
