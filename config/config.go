package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config Top-level config structure
type Config struct {
	Git     []GitConfig     `yaml:"git"`
	RESTAPI []RESTAPIConfig `yaml:"rest_api"`
	Docker  []DockerConfig  `yaml:"docker"`
}

// AuthConfig Used for mounted secret + HTTP header auth
type AuthConfig struct {
	HeaderName string `yaml:"header_name"` // e.g. "Authorization"
	SecretPath string `yaml:"secret_path"` // e.g. "/etc/secrets/github-token"
}

// GitConfig Git config (clone or ls-remote)
type GitConfig struct {
	Name            string      `yaml:"name"`
	RepoURL         string      `yaml:"repo_url"`
	Type            string      `yaml:"type"` // "clone" or "ls-remote"
	IntervalSeconds int         `yaml:"interval_seconds"`
	Auth            *AuthConfig `yaml:"auth,omitempty"`
}

// RESTAPIConfig REST API config with chained call and conditional polling
type RESTAPIConfig struct {
	Name            string            `yaml:"name"`
	Type            string            `yaml:"type"` // "follow_up" or "request"
	FirstURL        string            `yaml:"first_url,omitempty"`
	Method          string            `yaml:"method,omitempty"`
	URL             string            `yaml:"url,omitempty"`
	IntervalSeconds int               `yaml:"interval_seconds"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	FollowUp        FollowConfig      `yaml:"follow_up,omitempty"`
	Filter          FilterConfig      `yaml:"filter,omitempty"`
	Auth            *AuthConfig       `yaml:"auth,omitempty"`
}

// FilterConfig Config for filtering JSON response
type FilterConfig struct {
	JSONQuery string `yaml:"json_query,omitempty"`
	Regex     string `yaml:"regex,omitempty"`
	Match     string `yaml:"match"`
}

// FollowConfig Config for the follow-up URL
type FollowConfig struct {
	URLTemplate         string            `yaml:"url_template"`
	Replace             map[string]string `yaml:"replace"`
	StopCondition       KeyValue          `yaml:"stop_condition"`
	InitialDelaySeconds int               `yaml:"initial_delay_seconds"` // <-- new
	IntervalSeconds     int               `yaml:"interval_seconds"`
	TimeoutSeconds      int               `yaml:"timeout_seconds"`
}

// KeyValue represents a key-value pair for the stop condition
type KeyValue struct {
	Key   string `yaml:"key"`   // e.g. "status"
	Value string `yaml:"value"` // e.g. "completed"
}

// DockerConfig Docker image pull config
type DockerConfig struct {
	Name            string `yaml:"name"`
	Image           string `yaml:"image"`
	Registry        string `yaml:"registry"`
	IntervalSeconds int    `yaml:"interval_seconds"`
}

// LoadConfig reads and parses the YAML config file from a given path
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
