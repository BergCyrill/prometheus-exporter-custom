package collectors

import (
	"io"
	"net/http"
	"prometheus-exporter/config"
	"strings"
	"testing"
)

// mockRoundTripper returns canned responses to simulate HTTP calls
type mockRoundTripper struct {
	response string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(m.response)),
		Header:     make(http.Header),
	}, nil
}

func TestTryFollowUp(t *testing.T) {
	client := &http.Client{
		Transport: &mockRoundTripper{
			response: `{"status": "completed"}`,
		},
	}

	success := tryFollowUp(client,
		"http://fakeurl",
		"", "", // no auth
		mockRestConfig("status", "completed"))

	if !success {
		t.Error("Expected follow-up call to succeed, but it failed")
	}
}

func mockRestConfig(key, val string) config.RESTAPIConfig {
	return config.RESTAPIConfig{
		Name: "test",
		FollowUp: config.FollowConfig{
			StopCondition: config.KeyValue{
				Key:   key,
				Value: val,
			},
		},
	}
}
