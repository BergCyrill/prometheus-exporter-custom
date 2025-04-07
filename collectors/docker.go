package collectors

import (
	"log"
	"os/exec"
	"time"

	"prometheus-exporter/config"
	"prometheus-exporter/metrics"
)

func monitorDockerPull(cfg config.DockerConfig) {
	for {
		start := time.Now()

		cmd := exec.Command("docker", "pull", cfg.Image)
		err := cmd.Run()
		duration := time.Since(start).Seconds()

		if err != nil {
			metrics.DockerSuccess.WithLabelValues(cfg.Name, cfg.Image).Set(0)
			log.Printf("Docker pull failed for %s: %v", cfg.Image, err)
		} else {
			metrics.DockerSuccess.WithLabelValues(cfg.Name, cfg.Image).Set(1)
			metrics.DockerPullDuration.WithLabelValues(cfg.Name, cfg.Image).Observe(duration)
		}

		_ = exec.Command("docker", "rmi", "-f", cfg.Image).Run()

		time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
	}
}

func handleDocker(cfgs []config.DockerConfig) {
	for _, cfg := range cfgs {
		go monitorDockerPull(cfg)
	}
}
