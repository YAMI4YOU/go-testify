package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMainHandlerWhenCountMoreThanTotal проверяет, что если в параметре count указано больше,
// чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil) // здесь нужно создать запрос к сервису

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// проверяем код ответа 200, если нет, то завершаем тест
	require.Equal(t, http.StatusOK, responseRecorder.Code)

	// проверяем, что тело ответа не пустое
	assert.NotEmpty(t, responseRecorder.Body.String())

	// проверяем количество кафе, должно быть 4
	returnedCafes := strings.Split(responseRecorder.Body.String(), ",")
	assert.Len(t, returnedCafes, totalCount)

	// проверяем все ли кафе из списка есть в ответе
	assert.Equal(t, returnedCafes, cafeList["moscow"])
}

// TestMainHandlerRequestIsCorrect проверяет, что запрос сформирован корректно,
// сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerRequestIsCorrect(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// проверяем код ответа 200
	require.Equal(t, http.StatusOK, responseRecorder.Code)

	// проверяем, что тело ответа не пустое
	assert.NotEmpty(t, responseRecorder.Body.String())

	// проверяем, что вернулись 2 кафе, как мы и запросили
	returnedCafes := strings.Split(responseRecorder.Body.String(), ",")
	assert.Len(t, returnedCafes, 2)
}

// TestMainHandlerWhenCityIsNotCorrect проверяет, что город, который передаётся в параметре city,
// не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhenCityIsNotCorrect(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=wrong", nil) // запрашиваем 2 кафе
	// из заведомо неправильного города

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// проверяем, что возвращается код ответа 400
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	// проверяем описание ошибки в теле ответа
	assert.Equal(t, "wrong city value", responseRecorder.Body.String())
}
