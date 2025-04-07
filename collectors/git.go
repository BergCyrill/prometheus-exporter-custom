package collectors

import (
	"log"
	"os"
	"os/exec"
	"time"

	"prometheus-exporter/config"
	"prometheus-exporter/metrics"
	"prometheus-exporter/secrets"
)

func startGitClone(cfg config.GitConfig) {
	for {
		dir := "/tmp/" + cfg.Name
		start := time.Now()

		cmdArgs := []string{"clone", "--depth=1", cfg.RepoURL, dir}
		if cfg.Auth != nil {
			token, err := secrets.ReadSecret(cfg.Auth.SecretPath)
			if err == nil {
				// Inject token into repo URL (https://<token>@github.com/org/repo.git)
				urlWithToken := injectTokenIntoURL(cfg.RepoURL, token)
				cmdArgs = []string{"clone", "--depth=1", urlWithToken, dir}
			}
		}

		cmd := exec.Command("git", cmdArgs...)
		err := cmd.Run()
		duration := time.Since(start).Seconds()

		if err != nil {
			metrics.GitSuccess.WithLabelValues(cfg.Name, cfg.RepoURL, "clone").Set(0)
			log.Printf("Git clone failed: %v", err)
		} else {
			metrics.GitSuccess.WithLabelValues(cfg.Name, cfg.RepoURL, "clone").Set(1)
			metrics.GitCloneDuration.WithLabelValues(cfg.Name, cfg.RepoURL).Observe(duration)
		}

		_ = os.RemoveAll(dir)
		time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
	}
}

func startGitLsRemote(cfg config.GitConfig) {
	for {
		start := time.Now()
		repoURL := cfg.RepoURL

		if cfg.Auth != nil {
			token, err := secrets.ReadSecret(cfg.Auth.SecretPath)
			if err == nil {
				repoURL = injectTokenIntoURL(cfg.RepoURL, token)
			}
		}

		cmd := exec.Command("git", "ls-remote", repoURL)
		err := cmd.Run()
		duration := time.Since(start).Seconds()

		if err != nil {
			metrics.GitSuccess.WithLabelValues(cfg.Name, cfg.RepoURL, "ls-remote").Set(0)
			log.Printf("Git ls-remote failed: %v", err)
		} else {
			metrics.GitSuccess.WithLabelValues(cfg.Name, cfg.RepoURL, "ls-remote").Set(1)
			metrics.GitLsRemoteDuration.WithLabelValues(cfg.Name, cfg.RepoURL).Observe(duration)
		}

		time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
	}
}

func injectTokenIntoURL(url, token string) string {
	// Basic assumption: https://github.com/... or https://gitlab.com/...
	if url[:8] == "https://" {
		return "https://" + token + "@" + url[8:]
	}
	return url
}

func handleGit(cfgs []config.GitConfig) {
	for _, cfg := range cfgs {
		if cfg.Type == "clone" {
			go startGitClone(cfg)
		} else if cfg.Type == "ls-remote" {
			go startGitLsRemote(cfg)
		}
	}
}
