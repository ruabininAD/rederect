package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"math"
	"net/http"
	"net/url"
	"redirect/internal/config"
	"redirect/internal/metrics"
	"redirect/internal/storage"
	"time"
)

var db storage.DB
var (
	Version string
	Build   string
)

func main() {

	config.Init()
	config.Log.Info("run redirect", zap.String("Version", Version), zap.String("Build", Build), zap.String("Port", config.Cfg.Port))

	go metrics.InitMetrics()

	db = &storage.MariaDBS{}
	db.Connect()

	http.HandleFunc("/", measureTime(redirectHandler))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	err := http.ListenAndServe(":"+config.Cfg.Port, nil)
	if err != nil {
		config.Log.Fatal(err.Error())
		return
	}

}

func measureTime(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		f(w, r)
		elapsed := time.Since(start)
		metrics.ResponseTimeHistogram.WithLabelValues(r.Host).Observe(float64(elapsed) / math.Pow(10, 6))
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {

	metrics.RequestCounter.WithLabelValues(r.Host).Inc()

	tracer := otel.Tracer("redirect-tracer")
	ctx, span := tracer.Start(r.Context(), "redirectHandler")
	defer span.End()

	// Получаем последний домен
	domain, err := db.GetLast()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		config.Log.Fatal(err.Error()) //fixme
	}
	// Обновляем инфу о домене в списках доменов
	err = db.Update(domain)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		config.Log.Fatal(err.Error()) //fixme
	}

	newUrl := url.URL{
		Scheme:   "https",
		Host:     domain,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}
	config.Log.Debug(
		r.Host+" to "+newUrl.String(),
		zap.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
		zap.String("span_id", trace.SpanContextFromContext(ctx).SpanID().String()),
	)

	// Выставление заголовков редиректа и выполнение редиректа
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Expires", "Sat, 26 Jul 1997 05:00:00 GMT")
	w.Header().Set("Location", newUrl.String())
	w.WriteHeader(http.StatusTemporaryRedirect)

	metrics.ResponseCounter.WithLabelValues(domain).Inc()
}
