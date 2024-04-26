package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HttpRequestsCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "runners_app_http_requests",
		Help: "Total number of HTTP requests",
	},
)

var HttpResponsesCounter = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "runners_app_http_responses",
		Help: "Total number of HTTP responses",
	},
	[]string{"status"},
)
