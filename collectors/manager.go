package collectors

import "prometheus-exporter/config"

func RegisterCollectors(cfg *config.Config) {
	handleGit(cfg.Git)
	handleREST(cfg.RESTAPI)
	handleDocker(cfg.Docker)
}
