package main

import (
	"os"
	"prometheus-exporter/collectors"
	"prometheus-exporter/config"
	"testing"
	"time"
)

func TestIntegrationStartup(t *testing.T) {
	configYaml := `
git:
  - name: "test"
    repo_url: "https://github.com/kubernetes/kubernetes.git"
    type: "ls-remote"
    interval_seconds: 5
rest_api: []
docker: []
`
	os.MkdirAll("/tmp/config", 0755)
	fpath := "/tmp/config/test.yaml"
	os.WriteFile(fpath, []byte(configYaml), 0644)

	cfg, err := config.LoadConfig(fpath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	collectors.RegisterCollectors(cfg)

	// Simulate exporter running
	go main() // starts /metrics endpoint

	time.Sleep(3 * time.Second) // give it a moment to run

	// We're only verifying that startup works
}
