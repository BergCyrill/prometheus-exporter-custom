package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// GitCloneDuration measures the duration of git clone operations.
	GitCloneDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "git_clone_duration_seconds",
			Help: "Duration of git clone",
		},
		[]string{"name", "repo"},
	)

	// GitLsRemoteDuration measures the duration of git ls-remote operations.
	GitLsRemoteDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "git_ls_remote_duration_seconds",
			Help: "Duration of git ls-remote",
		},
		[]string{"name", "repo"},
	)

	// RESTCallDuration measures the duration of REST API calls.
	RESTCallDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "rest_api_call_duration_seconds",
			Help: "Duration of REST API long polling",
		},
		[]string{"name"},
	)

	// DockerPullDuration measures the duration of docker image pull operations.
	DockerPullDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "docker_pull_duration_seconds",
			Help: "Duration of docker image pull",
		},
		[]string{"name", "image"},
	)

	// GitSuccess indicates whether the last git operation was successful.
	GitSuccess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "git_probe_success",
			Help: "Whether the last git operation succeeded (1 for success, 0 for failure)",
		},
		[]string{"name", "repo", "type"},
	)

	// DockerSuccess indicates whether the last docker pull operation was successful.
	DockerSuccess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "docker_pull_success",
			Help: "Whether the last docker pull succeeded (1 for success, 0 for failure)",
		},
		[]string{"name", "image"},
	)

	// RESTSuccess indicates whether the last REST API call was successful.
	RESTSuccess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rest_api_probe_success",
			Help: "Whether the last REST API call probe succeeded (1 for success, 0 for failure)",
		},
		[]string{"name"},
	)
)

func init() {
	prometheus.MustRegister(GitCloneDuration)
	prometheus.MustRegister(GitLsRemoteDuration)
	prometheus.MustRegister(RESTCallDuration)
	prometheus.MustRegister(DockerPullDuration)
	prometheus.MustRegister(GitSuccess)
	prometheus.MustRegister(DockerSuccess)
	prometheus.MustRegister(RESTSuccess)
}
