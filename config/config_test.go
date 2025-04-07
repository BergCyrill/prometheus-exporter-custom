package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testYAML := `
git:
  - name: "test"
    repo_url: "https://example.com"
    type: "clone"
    interval_seconds: 60
rest_api: []
docker: []
`
	tmpFile := "test_config.yaml"
	os.WriteFile(tmpFile, []byte(testYAML), 0644)
	defer os.Remove(tmpFile)

	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if len(cfg.Git) != 1 || cfg.Git[0].Name != "test" {
		t.Errorf("unexpected git config: %+v", cfg.Git)
	}
}
