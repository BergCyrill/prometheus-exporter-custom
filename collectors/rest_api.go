package collectors

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"prometheus-exporter/config"
	"prometheus-exporter/metrics"
	"prometheus-exporter/secrets"

	"github.com/tidwall/gjson"
)

func monitorREST(cfg config.RESTAPIConfig) {
	var authHeader string
	var authKey string
	if cfg.Auth != nil {
		token, err := secrets.ReadSecret(cfg.Auth.SecretPath)
		if err == nil {
			authHeader = token
			authKey = cfg.Auth.HeaderName
		} else {
			log.Printf("[%s] Failed to load auth token: %v", cfg.Name, err)
		}
	}

	client := &http.Client{}

	for {
		start := time.Now()
		// First request
		req, _ := http.NewRequest("GET", cfg.FirstURL, nil)
		if authHeader != "" {
			req.Header.Set(authKey, authHeader)
		}

		resp, err := client.Do(req)
		if err != nil {
			metrics.RESTSuccess.WithLabelValues(cfg.Name).Set(0)
			log.Printf("[%s] First API call failed: %v", cfg.Name, err)
			time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		values := map[string]string{}
		for placeholder, jsonPath := range cfg.FollowUp.Replace {
			values[placeholder] = gjson.GetBytes(body, jsonPath).String()
		}

		tpl, err := template.New("url").Parse(cfg.FollowUp.URLTemplate)
		if err != nil {
			log.Printf("[%s] Failed to parse follow-up template: %v", cfg.Name, err)
			metrics.RESTSuccess.WithLabelValues(cfg.Name).Set(0)
			time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
			continue
		}

		var urlBuffer bytes.Buffer
		_ = tpl.Execute(&urlBuffer, values)
		followURL := urlBuffer.String()

		// Delay before first follow-up attempt
		time.Sleep(time.Duration(cfg.FollowUp.InitialDelaySeconds) * time.Second)

		timeout := time.After(time.Duration(cfg.FollowUp.TimeoutSeconds) * time.Second)
		ticker := time.NewTicker(time.Duration(cfg.FollowUp.IntervalSeconds) * time.Second)
		success := false

		// First call manually
		if tryFollowUp(client, followURL, authKey, authHeader, cfg) {
			metrics.RESTCallDuration.WithLabelValues(cfg.Name).Observe(time.Since(start).Seconds())
			metrics.RESTSuccess.WithLabelValues(cfg.Name).Set(1)
			success = true
		} else {
			// Then loop using ticker
		probeLoop:
			for {
				select {
				case <-ticker.C:
					if tryFollowUp(client, followURL, authKey, authHeader, cfg) {
						metrics.RESTCallDuration.WithLabelValues(cfg.Name).Observe(time.Since(start).Seconds())
						metrics.RESTSuccess.WithLabelValues(cfg.Name).Set(1)
						success = true
						break probeLoop
					}
				case <-timeout:
					break probeLoop
				}
			}
		}

		if !success {
			metrics.RESTSuccess.WithLabelValues(cfg.Name).Set(0)
		}

		time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
	}
}

func tryFollowUp(client *http.Client, url, authKey, authHeader string, cfg config.RESTAPIConfig) bool {
	req, _ := http.NewRequest("GET", url, nil)
	if authHeader != "" {
		req.Header.Set(authKey, authHeader)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("[%s] Follow-up request error: %v", cfg.Name, err)
		return false
	}

	b, _ := io.ReadAll(res.Body)
	res.Body.Close()

	match := gjson.GetBytes(b, cfg.FollowUp.StopCondition.Key).String() == cfg.FollowUp.StopCondition.Value
	// if match {
	// 	log.Printf("[%s] Follow-up condition met", cfg.Name)
	// }
	return match
}

func handleREST(cfgs []config.RESTAPIConfig) {
	for _, cfg := range cfgs {
		go monitorREST(cfg)
	}
}
