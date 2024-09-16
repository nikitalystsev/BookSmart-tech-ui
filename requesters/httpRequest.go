package requesters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTPRequest - структура для представления HTTP-запроса
type HTTPRequest struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        interface{}
	QueryParams map[string]string
	Timeout     time.Duration
}

// HTTPResponse - структура для представления HTTP-ответа
type HTTPResponse struct {
	Status     string
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

// SendRequest - универсальная функция для отправки HTTP-запроса
func SendRequest(req HTTPRequest) (*HTTPResponse, error) {
	// Создаем URL с параметрами запроса
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return nil, err
	}

	// Добавляем query параметры
	if len(req.QueryParams) > 0 {
		q := parsedURL.Query()
		for key, value := range req.QueryParams {
			q.Add(key, value)
		}
		parsedURL.RawQuery = q.Encode()
	}

	// Преобразуем тело запроса в JSON, если оно не nil
	var body io.Reader
	if req.Body != nil {
		var jsonBody []byte
		if jsonBody, err = json.Marshal(req.Body); err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonBody)
	}

	// Создаем новый HTTP-запрос
	httpReq, err := http.NewRequest(req.Method, parsedURL.String(), body)
	if err != nil {
		return nil, err
	}

	// Добавляем заголовки
	for key, value := range req.Headers {
		httpReq.Header.Add(key, value)
	}

	// Устанавливаем таймаут для HTTP-клиента
	client := &http.Client{
		Timeout: req.Timeout,
	}

	// Отправляем запрос
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(resp.Body)

	// Читаем тело ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Формируем и возвращаем HTTP-ответ
	httpResp := &HTTPResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       responseBody,
	}

	return httpResp, nil
}
