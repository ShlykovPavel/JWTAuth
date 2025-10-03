package auth

import (
	"errors"
	"github.com/ShlykovPavel/JWTAuth/JWTParser"
	"github.com/ShlykovPavel/JWTAuth/http-server/requests"
	"github.com/ShlykovPavel/JWTAuth/scheduler"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type JWTAuth struct {
	loginURL    string
	refreshURL  string
	credentials *requests.Credentials
	retryCount  int
	logger      *slog.Logger
	scheduler   *scheduler.Scheduler
	tokens      *requests.Tokens
	httpClient  *http.Client
	mu          sync.RWMutex // Используем RWMutex для оптимизации чтения
}

func NewJwtAuth(loginURL, refreshURL, username, password string, retryCount int, logger *slog.Logger) *JWTAuth {
	return &JWTAuth{
		loginURL:    loginURL,
		refreshURL:  refreshURL,
		credentials: &requests.Credentials{Username: username, Password: password},
		retryCount:  retryCount,
		logger:      logger,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *JWTAuth) Start() error {
	// Первоначальный логин
	tokens, err := requests.LoginOrRefreshInService(
		a.loginURL,
		*a.credentials,
		a.logger,
		a.retryCount,
	)
	if err != nil {
		return err
	}

	a.mu.Lock()
	a.tokens = tokens
	a.mu.Unlock()

	// Инициализация планировщика
	a.scheduler = scheduler.NewScheduler(a.handleRefresh, a.logger)

	// Планируем обновление
	if err := a.scheduleNextRefresh(); err != nil {
		return err
	}

	return nil
}

func (a *JWTAuth) handleRefresh() {
	a.mu.Lock()
	defer a.mu.Unlock()

	newTokens, err := requests.LoginOrRefreshInService(
		a.refreshURL,
		*a.tokens,
		a.logger,
		a.retryCount,
	)
	if err != nil {
		a.logger.Error("refresh failed", "error", err)
		return
	}
	a.tokens = newTokens

	// Планируем следующее обновление (без блокировки, так как мьютекс уже заблокирован)
	if err := a.scheduleNextRefreshUnlocked(); err != nil {
		a.logger.Error("failed to schedule next refresh", "error", err)
	}
}

func (a *JWTAuth) scheduleNextRefresh() error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.scheduleNextRefreshUnlocked()
}

// scheduleNextRefreshUnlocked - внутренний метод без блокировок
func (a *JWTAuth) scheduleNextRefreshUnlocked() error {
	claims, err := JWTParser.ParseUnverified(a.tokens.AccessToken, a.logger)
	if err != nil {
		return err
	}

	expiry, err := JWTParser.GetExpirationTime(claims, a.logger)
	if err != nil {
		return err
	}

	a.scheduler.ScheduleRefresh(expiry)
	return nil
}

func (a *JWTAuth) GetToken() (string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.tokens == nil {
		return "", errors.New("not authenticated")
	}
	return a.tokens.AccessToken, nil
}

func (a *JWTAuth) Stop() {
	if a.scheduler != nil {
		a.scheduler.Stop()
	}
}
