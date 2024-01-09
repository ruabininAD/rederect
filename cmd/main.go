package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"rederect/internal/config"
	"rederect/internal/metrics"
	"rederect/internal/storage"
	"regexp"
	"strings"
	"sync"
	//	"zap"
	//"net/http"
)

// zap
// net/http
// .env
// config

var wg sync.WaitGroup
var db storage.DB

func main() {
	config.Init()
	wg.Add(1)
	go metrics.InitMetrics()

	db = &storage.MariaDBS{}
	db.Connect()

	wg.Add(1)
	go f()
	wg.Wait()
}

func f() {
	http.HandleFunc("/g", getHandler)
	http.HandleFunc("/u", updateHandler)
	http.HandleFunc("/", redirectHandler)
	http.ListenAndServe(":"+config.Cfg.Port, nil)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL
	path := strings.TrimLeft(originalURL.Path, "/")

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
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
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
	// Пример подключения к базе данных MySQL
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database_name")
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Проверка ID на соответствие формату
	matched, err := regexp.MatchString(`^([0-9]){1,8}$`, id)
	if err != nil {
		return "", err
	}
	if matched {
		// Выполнение запроса к базе данных для получения пути новости по ID
		var newPath string
		err := db.QueryRow("SELECT CONCAT(`cat_url`, '/', `item_url`) FROM `news_posts` "+
			"LEFT JOIN `news_posts_cat` ON `news_posts_cat`.`cat_id` = `news_posts`.`item_cat` "+
			"WHERE `item_id` = ?", id).Scan(&newPath)
		if err != nil {
			return "", err
		}
		return newPath, nil
	}
	return "", fmt.Errorf("Invalid news ID")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	lastDomain, err := db.GetLast()
	if err != nil {
		config.Log.Info("error updateHandler: " + err.Error())
	}
	println(lastDomain)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	err := db.Update(r.Host)
	fmt.Println("update URL" + r.URL.String())
	if err != nil {
		config.Log.Info("error updateHandler: " + err.Error())
	}
}
