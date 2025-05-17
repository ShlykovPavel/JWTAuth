package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Tokens struct for access and refresh tokens after auth
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

// LoginOrRefreshInService выполняет аутентификацию или обновление токена.
// Поддерживает типы Credentials (для логина) и Tokens (для refresh).
// Возвращает новые токены или ошибку.
func LoginOrRefreshInService[T Credentials | Tokens](URL string, body T, log *slog.Logger, retryCount int) (*Tokens, error) {
	const op = "requests.LoginOrRefreshInService"
	var operation string
	switch any(body).(type) {
	case Credentials:
		operation = "login"
	case Tokens:
		operation = "refresh"
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Ошибка при маршалинге JSON: %s\n", err)
		return nil, err
	}
	log = log.With(
		slog.String("operation", op),
		slog.String("auth_type", operation),
		slog.String("url", URL),
		slog.String("body", string(jsonData)),
	)

	log.Debug("request body", slog.String("data", string(jsonData)))
	for attempt := 0; attempt <= retryCount; attempt++ {
		resp, err := makePostRequest(URL, jsonData, log)
		if err != nil {
			log.Error("Error in request: ", slog.String("error", err.Error()))
			if attempt == retryCount {
				return nil, err
			}
			continue
		}
		defer resp.Body.Close()

		//Выход из функции при положительном результате
		if resp.StatusCode == http.StatusOK {
			var tokens Tokens
			err = json.NewDecoder(resp.Body).Decode(&tokens)
			if err != nil {
				return nil, fmt.Errorf("decode tokens: %w", err)
			}

			return &tokens, nil
		}
		//Читаем тело ошибки и логируем
		respBody, err := io.ReadAll(resp.Body)
		log.Warn("server error",
			"attempt", attempt,
			"status", resp.StatusCode,
			"body", string(respBody))
		if err != nil {
			log.Error("Error while reading response body", slog.String("error", err.Error()))
			return nil, err
		}
		//Небольшая задержка перед следующей попыткой
		time.Sleep(time.Duration(attempt) * time.Second)

	}
	return nil, fmt.Errorf("after %d attempts login failed", retryCount)

}

// makePostRequest Выполняет post запрос к api
//
// Возвращает ответ или ошибку
func makePostRequest(URL string, data []byte, log *slog.Logger) (*http.Response, error) {
	const op = "requests.makePostRequest"
	resp, err := defaultClient.Post(URL, "application/json", bytes.NewBuffer(data))
	// Обработка если произошла ошибка сети
	log = log.With(
		slog.String("operation", op),
		slog.String("url", URL))
	if err != nil {
		//Обработка ошибки по таймауту
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, fmt.Errorf("request timeout: %w", err)
		}
		//Общая логика ошибок
		log.Error("HTTP request failed: ", slog.String("error", err.Error()))
		return nil, err

	}
	return resp, nil
}
