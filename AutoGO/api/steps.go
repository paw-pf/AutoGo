package api

import (
	"autogo/config"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// LoginResponse — обычный логин для вызова апи тестов
type LoginResponse struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}

// Login выполняет аутентификацию и возвращает токен
func (c *Client) Login() (string, error) {
	body, resp := c.POST("/Account/v1/GenerateToken", `{
		"userName": "`+config.Username+`",
		"password": "`+config.Password+`"
	}`)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed: %d, response: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", fmt.Errorf("failed to parse login response: %w", err)
	}

	return loginResp.Token, nil
}

// AuthenticatedClient создаёт клиент API и автоматически авторизует его.
func AuthenticatedClient(t *testing.T, baseURL string) *Client {
	t.Helper()

	client := NewClient(t, baseURL)
	token, err := client.Login()
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	return client.WithToken(token)
}

func (c *Client) InfoUser() (string, error) {
	body, resp := c.POST("/Account/v1/User", `{
		"userName": "`+config.Username+`",
		"password": "`+config.Password+`"
	}`)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("info user failed: %d, response: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", fmt.Errorf("failed to parse login response: %w", err)
	}

	return loginResp.UserID, nil
}
