package main

import (
	"log"
	"net/http"

	"prometheus-exporter/collectors"
	"prometheus-exporter/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.LoadConfig("/etc/config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	collectors.RegisterCollectors(cfg)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter running on :2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
