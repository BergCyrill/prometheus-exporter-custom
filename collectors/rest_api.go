package collectors

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"

	"prometheus-exporter/config"
	"prometheus-exporter/metrics"
	"prometheus-exporter/secrets"

	"github.com/tidwall/gjson"
)

func monitorFollowUp(cfg config.RESTAPIConfig) {
	client := &http.Client{}

	for {
		var authHeader string
		var authKey string
		if cfg.Auth != nil {
			token, err := secrets.ReadSecret(cfg.Auth.SecretPath)
			if err != nil {
				log.Printf("[%s] Failed to load auth token: %v", cfg.Name, err)
				log.Printf("[%s] Will rerun in %v seconds", cfg.Name, cfg.IntervalSeconds)
				time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
				continue
			}
			authHeader = token
			authKey = cfg.Auth.HeaderName
		}
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

func monitorRequest(cfg config.RESTAPIConfig) {
	client := &http.Client{}
	var authHeader, authKey string

	for {
		// Always re-read the token from the secret file every interval
		if cfg.Auth != nil {
			token, err := secrets.ReadSecret(cfg.Auth.SecretPath)
			if err != nil {
				log.Printf("[%s] Failed to load auth token: %v", cfg.Name, err)
				log.Printf("[%s] Will rerun in %v seconds", cfg.Name, cfg.IntervalSeconds)
				time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
				continue
			} else {
				authHeader = token
				authKey = cfg.Auth.HeaderName
			}
		}
		start := time.Now()
		// Create the request
		req, _ := http.NewRequest(cfg.Method, cfg.URL, nil)

		if authHeader != "" && authKey != "" {
			req.Header.Set(authKey, authHeader)
		}

		resp, err := client.Do(req)
		if err != nil {
			metrics.RESTSuccess.WithLabelValues(cfg.Name).Set(0)
			log.Printf("[%s] Request failed: %v", cfg.Name, err)
			time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		content := string(body)

		match := false
		if cfg.Filter.JSONQuery != "" {
			value := gjson.Get(content, cfg.Filter.JSONQuery).String()
			match = value == cfg.Filter.Match
		} else if cfg.Filter.Regex != "" {
			re := regexp.MustCompile(cfg.Filter.Regex)
			match = re.MatchString(content) && strings.Contains(content, cfg.Filter.Match)
		}

		if match {
			metrics.RESTSuccess.WithLabelValues(cfg.Name).Set(1)
		} else {
			metrics.RESTSuccess.WithLabelValues(cfg.Name).Set(0)
			log.Printf("[%s] Response did not match filter", cfg.Name)
		}
		metrics.RESTCallDuration.WithLabelValues(cfg.Name).Observe(time.Since(start).Seconds())
		time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
	}
}

func handleREST(cfgs []config.RESTAPIConfig) {
	for _, cfg := range cfgs {
		if cfg.Type == "follow_up" {
			go monitorFollowUp(cfg)
		} else if cfg.Type == "request" || cfg.Type == "" {
			go monitorRequest(cfg)
		}
	}
}
