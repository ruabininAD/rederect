package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"rederect/internal/config"
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

	DigitalPathCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "redirect",
		Subsystem: "http",
		Name:      "requests_digital_path_total",
		Help:      "redirect http requests digitalPath total",
	}, []string{"host", "code"})
)

func InitMetrics() {

	config.Log.Debug("prometheus metrics init")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
