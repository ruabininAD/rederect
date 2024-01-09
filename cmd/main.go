package main

import (
	"fmt"
	"go.uber.org/zap"
	"math"
	"net/http"
	"net/url"
	"redirect/internal/config"
	"redirect/internal/metrics"
	"redirect/internal/storage"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var db storage.DB
var (
	Version string
	Build   string
)

func main() {

	config.Init()
	config.Log.Info("run redirect", zap.String("Version", Version), zap.String("Build", Build))

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

	originalURL := r.URL
	path := strings.TrimLeft(originalURL.Path, "/")
	metrics.RequestCounter.WithLabelValues(r.Host).Inc()

	// Если весь путь - цифровой
	if path != "" {
		matched, err := regexp.MatchString(`^([0-9]){1,8}$`, path)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if matched {
			// То ищем новость с таким ID и формируем адрес
			newPath, err := getNewsPath(path)
			if err != nil {
				code := http.StatusInternalServerError
				http.Error(w, "news not found", code)
				config.Log.Warn("news " + path + " not found")
				metrics.DigitalPathCounter.WithLabelValues(r.Host, strconv.Itoa(code)).Inc()
				return
			}
			metrics.DigitalPathCounter.WithLabelValues(r.Host, "200").Inc()
			path = newPath
		}
	}

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
	metrics.ResponseCounter.WithLabelValues(domain).Inc()

	newUrl := url.URL{
		Scheme:   "https",
		Host:     domain,
		Path:     originalURL.Path,
		RawQuery: originalURL.RawQuery,
	}
	config.Log.Debug(r.Host + " to " + newUrl.String())

	// Выставление заголовков редиректа и выполнение редиректа
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Expires", "Sat, 26 Jul 1997 05:00:00 GMT")
	w.Header().Set("Location", newUrl.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func getNewsPath(id string) (string, error) {

	// Проверка ID на соответствие формату
	matched, err := regexp.MatchString(`^([0-9]){1,8}$`, id)
	if err != nil {
		return "", err
	}
	if matched {
		// получения пути новости по ID
		newPath, err := db.PathId(id)
		if err != nil {
			return "", err
		}
		return newPath, nil
	}
	return "", fmt.Errorf("Invalid news ID")
}
