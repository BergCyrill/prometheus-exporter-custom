package collectors

import "prometheus-exporter/config"

// RegisterCollectors initializes and starts all collectors based on the provided configuration.
func RegisterCollectors(cfg *config.Config) {
	handleGit(cfg.Git)
	handleREST(cfg.RESTAPI)
	handleDocker(cfg.Docker)
}
