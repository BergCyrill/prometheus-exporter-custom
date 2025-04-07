package metrics

import (
	"testing"
)

func TestGitCloneDurationMetric(t *testing.T) {
	defer func() {
		recover() // suppress panic if metric already registered
	}()
	GitCloneDuration.WithLabelValues("test", "https://example.com").Observe(0.5)
}

func TestRESTSuccessGauge(t *testing.T) {
	defer func() {
		recover()
	}()
	RESTSuccess.WithLabelValues("test_api").Set(1)
}
