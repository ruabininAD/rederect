package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"redirect/internal/config"
)

var (
	RequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "redirect",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "redirect http requests total",
	}, []string{"host"})

	ResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "redirect",
		Subsystem: "http",
		Name:      "response_total",
		Help:      "response http requests total",
	}, []string{"host"})

	ResponseTimeHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "redirect",
			Subsystem: "http",
			Name:      "request_duration_milliseconds",
			Help:      "redirect http request duration (ms)",
			Buckets: []float64{
				2, 6, 10, 14, 18, 22, 26, 30,
			},
		}, []string{"host"})
)

func InitMetrics() {

	config.Log.Debug("prometheus metrics init http://localhost:" + config.Cfg.MetricPort + "/metrics")

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+config.Cfg.MetricPort, nil)
	if err != nil {
		config.Log.Warn(err.Error())
		return
	}
}
