package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		// пока сравнивать не будем, а просто выведем ответы
		// удалите потом этот вывод
		fmt.Println(response.Body.String())
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "tula"
	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList[city])},
	}

	for _, v := range requests {
		url := fmt.Sprintf("/cafe?city=%s&count=%d", city, v.count)
		req := httptest.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())

		if len(body) == 0 {
			assert.Equal(t, v.want, 0)
		} else {
			cafes := strings.Split(body, ",")
			assert.Equal(t, v.want, len(cafes))
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"
	requests := []struct {
		search string
		want   int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		url := fmt.Sprintf("/cafe?city=%s&search=%s", city, v.search)
		req := httptest.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())

		if len(body) == 0 {
			assert.Equal(t, v.want, 0)
		} else {
			cafes := strings.Split(body, ",")
			assert.Equal(t, v.want, len(cafes))
		}
	}
}
