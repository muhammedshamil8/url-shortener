package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var URLsCreated = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "snippy",
		Name:      "urls_created_total",
		Help:      "Total number of shortened URLs created.",
	},
)

var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "snippy",
		Name:      "request_duration_seconds",
		Help:      "HTTP request duration in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "endpoint", "status"},
)

var Redirects = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "snippy",
		Name:      "redirects_total",
		Help:      "Total redirect attempts grouped by result.",
	},
	[]string{"result"},
)

var Cache = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "snippy",
		Name:      "cache_total",
		Help:      "Cache operations grouped by result.",
	},
	[]string{"result"},
)

func Register() {
	prometheus.MustRegister(URLsCreated)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(Redirects)
	prometheus.MustRegister(Cache)
}

func URLCreated() {
	URLsCreated.Inc()
}

func Redirect(result string) {
	Redirects.WithLabelValues(result).Inc()
}

func CacheResult(result string) {

	Cache.WithLabelValues(result).Inc()
}
