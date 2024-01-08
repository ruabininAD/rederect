package main

import (
	"fmt"
	"net/http"
	"net/url"
	"rederect/internal/config"
	"rederect/internal/db"
	"rederect/internal/metrics"
	"sync"
	//	"zap"
	//"net/http"
)

// zap
// net/http
// .env
// config

var wg sync.WaitGroup

func main() {
	config.Init()
	wg.Add(1)
	go metrics.InitMetrics()

	mar := db.MariaDBS{}
	mar.Connect()
	wg.Add(1)
	go f()
	wg.Wait()
}
func f() {
	http.HandleFunc("/social/new-school-opened-in-vladovostok", redirectHandler)
	http.HandleFunc("/1", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte{12, 213, 43, 23, 43, 54, 12, 23, 34, 45, 56})
	})
	http.ListenAndServe(":8000", nil)
}
func redirectHandler(w http.ResponseWriter, r *http.Request) {

	//originalURL := r.URL.String()
	parsedURL, err := url.Parse("https://localhost:8000/social/new-school-opened-in-vladovostok?utm_source=smi2&utm_content=23415125")
	if err != nil {
		fmt.Println("Ошибка при разборе URL:", err)
		return
	}

	fmt.Println("Схема (Scheme):", parsedURL.Scheme)
	fmt.Println("Хост (Host):", parsedURL.Host)
	fmt.Println("Путь (Path):", parsedURL.Path)
	fmt.Println("Сырой Query параметры (Raw Query):", parsedURL.RawQuery)

	http.Redirect(w, r, "https://www.youtube.com/", http.StatusTemporaryRedirect)
}
