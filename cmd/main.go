package main

import (
	"fmt"
	"net/http"
	"net/url"
	"rederect/internal/config"
	"rederect/internal/metrics"
	"rederect/internal/storage"
	"regexp"
	"strings"
	//	"zap"
	//"net/http"
)

// zap
// net/http
// .env
// config

var db storage.DB

func main() {
	config.Init()

	go metrics.InitMetrics()

	db = &storage.MariaDBS{}
	db.Connect()

	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {}) //fixme
	http.ListenAndServe(":"+config.Cfg.Port, nil)

}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("redirectHandler")
	originalURL := r.URL
	path := strings.TrimLeft(originalURL.Path, "/")
	metrics.RequestCounter.WithLabelValues(r.Host).Inc()

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
				http.Error(w, "news not found", http.StatusInternalServerError)
				config.Log.Debug("news " + path + " not found")
				return
			}
			metrics.DigitalPathCounter.WithLabelValues(domain).Inc()
			path = newPath
		}
	}

	newUrl := url.URL{
		Scheme:   originalURL.Scheme,
		Host:     domain,
		Path:     originalURL.Path,
		RawQuery: originalURL.RawQuery,
	}

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
		// Выполнение запроса к базе данных для получения пути новости по ID
		newPath, err := db.PathId(id)
		if err != nil {
			return "", err
		}
		return newPath, nil
	}
	config.Log.Debug("Invalid news ID")
	return "", fmt.Errorf("Invalid news ID")
}
