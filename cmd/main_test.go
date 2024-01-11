package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_redirectHandler(t *testing.T) {
	//fixme

	go main()

	time.Sleep(3 * time.Second)
	url := "http://localhost:8082/ad"

	_, err := http.Get(url)
	if strings.Contains(err.Error(), "lookup hot-news.local: no such host") {
	} else if err != nil {
		fmt.Println("Ошибка при выполнении GET-запроса:", err)
		return
	}

	fmt.Println(err)

}

//
//	// Создание HTTP клиента для выполнения запросов к тестовому серверу
//	client := &http.Client{}
//
//	// Отправляем GET запрос на тестовый сервер
//	response, err := client.Get(server.URL)
//	if err != nil {
//		t.Fatalf("Failed to send GET request: %v", err)
//	}
//
//	// Проверяем, что код ответа равен 307 (Temporary Redirect)
//	assert.Equal(t, http.StatusTemporaryRedirect, response.StatusCode)
//
//	// Проверяем, что заголовок Location установлен и содержит целевой URL редиректа
//	location, err := response.Location()
//	if err != nil {
//		t.Fatalf("Failed to parse Location header: %v", err)
//	}
//
//	expectedLocation := "https://example.com/new-location"
//	assert.Equal(t, expectedLocation, location.String())
//}
