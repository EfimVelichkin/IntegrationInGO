package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	// Создаем экземпляр структуры servie
	srv := servie{store: make(map[string]*User), file: "../users/users.json"}

	// Создаем тестовый сервер
	reqBody := `{"name":"John", "age":30}`
	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа сервера
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.Create) // Используем метод Create через экземпляр srv

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус ответа
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Проверяем, что в ответе содержится ожидаемая строка
	expected := "User was created John\nYour ID "
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
