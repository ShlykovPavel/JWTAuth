package requests

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestMakePostRequest Тест основной функции выполнения запроса клиентом
func TestMakePostRequest(t *testing.T) {
	// Создание тестового сервера
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch string(body) {
		case `{"test": 1}`:
			w.WriteHeader(http.StatusOK)
			return
		case `{"test": 2}`:
			w.WriteHeader(http.StatusBadRequest)
			return
		case `{"test": 3}`:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case `{"test": 4}`:
			time.Sleep(time.Second * 11)

		default:
			w.WriteHeader(http.StatusTeapot)
		}

	}))
	defer testServer.Close()
	//Сетапим тест кейсы
	tests := []struct {
		testName      string
		requestBody   string
		responseBody  string
		responseCode  int
		wantError     bool
		errorContains string
	}{
		{"PositivePostRequest", `{"test": 1}`, "", http.StatusOK, false, ""},
		{"NegativePostRequest", `{"test": 2}`, "", http.StatusBadRequest, false, ""},
		{"ServerErrorTest", `{"test": 3}`, "", http.StatusInternalServerError, false, ""},
		{"TimeoutTest", `{"test": 4}`, "", http.StatusRequestTimeout, true, "timeout"},
	}
	var logOutput bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logOutput, nil))
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Запуск подтеста
			log.Info("launch test", slog.String("Test", tt.testName),
				slog.String("request body from test case", tt.requestBody))
			result, err := makePostRequest(testServer.URL, []byte(tt.requestBody), log)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				log.Info("Got error", slog.String("error", err.Error()))
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Fatalf("expected error containing %q, got %q", tt.errorContains, logOutput.String())
				}
				return
			}
			if err != nil {
				t.Error("Error while making post request: ", err)

			}
			if result.StatusCode != tt.responseCode {
				t.Errorf("got %v status code, want %v status code", result, tt.responseCode)
			}
		})

	}
}

func TestLoginOrRefreshInService(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch string(body) {
		case `{"accessKey":"test","secretKey":"password"}`:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"accessToken": "valid", "refreshToken": "valid"}`))
			return
		case `{"accessToken":"access","refreshToken":"refresh"}`:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"accessToken": "valid", "refreshToken": "valid"}`))
			return
		case `{"accessKey":"test2","secretKey":"password2"}`:
			w.WriteHeader(http.StatusBadRequest)
			return
		case `{"accessToken":"access2","refreshToken":"refresh2"}`:
			w.WriteHeader(http.StatusBadRequest)
			return
		case `{"accessKey":"test3","secretKey":"password3"}`:
			time.Sleep(time.Second * 11)
			return
		case `{"accessToken":"access3","refreshToken":"refresh3"}`:
			time.Sleep(time.Second * 11)
			return

		default:
			w.WriteHeader(http.StatusTeapot)
		}
	}))
	defer testServer.Close()
	tests := []struct {
		testName      string
		requestBody   any
		responseBody  string
		responseCode  int
		wantError     bool
		errorContains string
		retryCount    int
	}{
		{"PositivePostRequestLogin",
			Credentials{Username: "test", Password: "password"},
			`{"accessToken": "valid", "refreshToken": "valid"}`,
			http.StatusOK,
			false,
			"", 1},
		{"PositivePostRequestRefresh",
			Tokens{AccessToken: "access", RefreshToken: "refresh"},
			`{"accessToken": "valid", "refreshToken": "valid"}`,
			http.StatusOK,
			false,
			"", 1},
		{"Negative400PostRequestLogin",
			Credentials{Username: "test2", Password: "password2"},
			`{"accessToken": "valid", "refreshToken": "valid"}`,
			http.StatusBadRequest,
			true,
			"after 1 attempts login failed", 1},
		{"Negative400PostRequestRefresh",
			Tokens{AccessToken: "access2", RefreshToken: "refresh2"},
			`{"accessToken": "valid", "refreshToken": "valid"}`,
			http.StatusBadRequest,
			true,
			"after 1 attempts login failed", 1},
		{"Negative3RetryCountPostRequestLogin",
			Credentials{Username: "test2", Password: "password2"},
			`{"accessToken": "valid", "refreshToken": "valid"}`,
			http.StatusBadRequest,
			true,
			"after 3 attempts login failed", 3},
		{"Negative3RetryCountPostRequestRefresh",
			Tokens{AccessToken: "access2", RefreshToken: "refresh2"},
			`{"accessToken": "valid", "refreshToken": "valid"}`,
			http.StatusBadRequest,
			true,
			"after 3 attempts login failed", 3},
		{"TimeoutPostRequestLogin",
			Credentials{Username: "test3", Password: "password3"},
			`{"accessToken": "valid", "refreshToken": "valid"}`,
			http.StatusBadRequest,
			true,
			"timeout", 1},
		{"TimeoutPostRequestRefresh",
			Tokens{AccessToken: "access3", RefreshToken: "refresh3"},
			`{"accessToken": "valid", "refreshToken": "valid"}`,
			http.StatusBadRequest,
			true,
			"timeout", 1},
	}
	var logOutput bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logOutput, nil))
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			var result *Tokens
			var err error
			// Определяем тип requestBody и вызываем функцию с правильным дженериком
			switch body := tt.requestBody.(type) {
			case Credentials:
				result, err = LoginOrRefreshInService(testServer.URL, body, log, tt.retryCount)
			case Tokens:
				result, err = LoginOrRefreshInService(testServer.URL, body, log, tt.retryCount)
			default:
				t.Fatal("Unsupported request body type")
			}
			//Проверка ожидаемых ошибок
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Fatalf("expected error containing %q, got %q", tt.errorContains, logOutput.String())

				}
				return
			}
			//Ошибка во время выполнения POST запроса
			if err != nil {
				t.Error("Error while making post request: ", err)
			}
			//Проверка пустого результата
			if result == nil {
				t.Fatal("Expected tokens, got nil")
			}
		})
	}
}
