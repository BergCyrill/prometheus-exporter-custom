package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	GitCloneDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "git_clone_duration_seconds",
			Help: "Duration of git clone",
		},
		[]string{"name", "repo"},
	)

	GitLsRemoteDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "git_ls_remote_duration_seconds",
			Help: "Duration of git ls-remote",
		},
		[]string{"name", "repo"},
	)

	RESTCallDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "rest_api_call_duration_seconds",
			Help: "Duration of REST API long polling",
		},
		[]string{"name"},
	)

	DockerPullDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "docker_pull_duration_seconds",
			Help: "Duration of docker image pull",
		},
		[]string{"name", "image"},
	)

	GitSuccess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "git_probe_success",
			Help: "Whether the last git operation succeeded (1 for success, 0 for failure)",
		},
		[]string{"name", "repo", "type"},
	)

	DockerSuccess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "docker_pull_success",
			Help: "Whether the last docker pull succeeded (1 for success, 0 for failure)",
		},
		[]string{"name", "image"},
	)

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
