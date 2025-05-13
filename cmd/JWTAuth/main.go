package main

import (
	"github.com/ShlykovPavel/JWTAuth/auth"
	"github.com/ShlykovPavel/JWTAuth/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "production"
)

func main() {
	cfg := config.LoadConfig()
	//log.Default().Println("cfg:", cfg)
	log := setupLogger(cfg.Env)
	log.Info("Starting application")
	log.Debug("Debug messages enabled")

	jwtauth := auth.NewJwtAuth(
		"https://unuk-admin-stage.devol.xyz/api/accounts/login",
		"https://unuk-admin-stage.devol.xyz/api/accounts/refresh-tokens",
		cfg.Username,
		cfg.Password,
		cfg.RetryCount,
		log)
	err := jwtauth.Start()
	if err != nil {
		log.Error("Error starting jwtauth", err.Error())
	}
	token, err := jwtauth.GetToken()
	if err != nil {
		log.Error("failed to get token", "error", err)
		return
	}

	log.Info("successfully got token", "token", token)
	select {}
	//time.Sleep(time.Minute * 10)
}

// setupLogger
//
// Configures and initializes a structured logger (slog.Logger) tailored to the specified runtime environment.
// The logger outputs log messages in JSON format, ensuring consistency and ease of parsing across different environments.
//
// Parameters:
// - env (string): The runtime environment. Supported values:
//   - envLocal: Local development environment. Enables Debug-level logging for detailed debugging.
//   - envDev: Shared development or staging environment. Also enables Debug-level logging.
//   - envProd: Production environment. Restricts logging to Info level to reduce noise and focus on critical events.
//
// Behavior:
//   - In local (`envLocal`) and development (`envDev`) environments, the logger is configured with the `Debug` level.
//     This ensures that all log messages, including debug-level information, are captured and output in JSON format.
//     This is particularly useful for troubleshooting and detailed analysis during development.
//   - In the production environment (`envProd`), the logger is configured with the `Info` level. This ensures that only
//     informational and higher-priority logs (e.g., warnings, errors) are captured, reducing verbosity and focusing on
//     critical operational data.
//
// Returns:
// - *slog.Logger: A configured logger instance ready for use in the specified environment.
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return log
}
