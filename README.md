# Go JWT Auth Library

Библиотека для управления JWT-аутентификацией с автоматическим обновлением токенов.

## Установка

```bash
go get github.com/ShlykovPavel/JWTAuth
```

## Быстрый старт

```go
package main

import (
	"github.com/yourusername/go-jwtauth/auth"
	"github.com/yourusername/go-jwtauth/config"
	"log"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()
	
	// Инициализация логгера (Использовать пакет slog)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	
	// Создание клиента JWT аутентификации
	jwtauth := auth.NewJwtAuth(
		"https://example.com/api/login",       // URL для входа
		"https://example.com/api/refresh",     // URL для обновления токена
		"your_username",                       // Логин
		"your_password",                       // Пароль
		3,                                     // Количество попыток повторных запросов
		logger,                                // Логгер (реализующий ваш интерфейс)
	)
	
	// Запуск клиента (начинает процесс аутентификации и обновления токенов)
	if err := jwtauth.Start(); err != nil {
		logger.Fatal("Failed to start JWT auth:", err)
	}
	
	// Получение текущего токена
	token, err := jwtauth.GetToken()
	if err != nil {
		logger.Fatal("Failed to get token:", err)
	}
	
	logger.Println("Successfully authenticated. Token:", token)
	
	// Блокировка main (или работа вашего приложения)
	select {}
}
```

## Конфигурация

Библиотека ожидает конфигурацию в следующем формате:

```.env
AUTH_USERNAME=your_username
AUTH_PASSWORD=your_password
ENV=your_env (Необязательно)
```

Или же в добавить эти данные в переменные окружиения при запуске. 
Сначала пакет ищет переменные окружения, а потом файл .env

## Основные методы

### `NewJwtAuth(authURL, refreshURL, username, password string, retryCount int, logger LoggerInterface) *JwtAuth`

Создает новый экземпляр клиента JWT аутентификации.

Параметры:
- `authURL` - URL для первоначальной аутентификации
- `refreshURL` - URL для обновления токена
- `username` - логин пользователя
- `password` - пароль пользователя
- `retryCount` - количество попыток повторных запросов при ошибках
- `logger` - интерфейс логгера

### `(j *JwtAuth) Start() error`

Запускает процесс аутентификации и начинает автоматическое обновление токенов.

### `(j *JwtAuth) GetToken() (string, error)`

Возвращает текущий JWT токен. Если токен недействителен или отсутствует, пытается получить новый.

### `(j *JwtAuth) Stop()`

Останавливает автоматическое обновление токенов.

## Логирование

Библиотека ожидает, что переданный логгер реализует следующий интерфейс:

```go
type LoggerInterface interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}
```

Необходимо использовать логгер из библиотеки slog
