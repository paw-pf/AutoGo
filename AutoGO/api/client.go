// package api

package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Client — HTTP-клиент для API-запросов (только для тестов!)
type Client struct {
	BaseURL      string
	HTTPClient   *http.Client
	Token        string
	T            *testing.T
	lastResponse []byte // ← добавлено: хранит тело последнего ответа
}

// NewClient создаёт новый API-клиент
func NewClient(t *testing.T, baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		T:          t,
	}
}

// WithToken устанавливает Bearer-токен
func (c *Client) WithToken(token string) *Client {
	c.Token = token
	return c
}

// LastResponse возвращает тело последнего ответа (для аттачмента в Allure)
func (c *Client) LastResponse() []byte {
	return c.lastResponse
}

// SendRequest выполняет HTTP-запрос
func (c *Client) SendRequest(method, endpoint, body string) ([]byte, *http.Response) {
	url := c.BaseURL + endpoint

	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	assert.NoError(c.T, err, "не удалось создать HTTP-запрос")

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	assert.NoError(c.T, err, "ошибка при выполнении HTTP-запроса")
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	assert.NoError(c.T, err, "ошибка при чтении тела ответа")

	// Сохраняем тело для последующего использования в отчётах
	c.lastResponse = rawBody

	return rawBody, resp
}

// GET — выполняет GET-запрос
func (c *Client) GET(endpoint string) ([]byte, *http.Response) {
	return c.SendRequest("GET", endpoint, "")
}

// POST — выполняет POST-запрос
func (c *Client) POST(endpoint, body string) ([]byte, *http.Response) {
	return c.SendRequest("POST", endpoint, body)
}

// PUT — выполняет PUT-запрос
func (c *Client) PUT(endpoint, body string) ([]byte, *http.Response) {
	return c.SendRequest("PUT", endpoint, body)
}

// DELETE — выполняет DELETE-запрос
func (c *Client) DELETE(endpoint string) ([]byte, *http.Response) {
	return c.SendRequest("DELETE", endpoint, "")
}
