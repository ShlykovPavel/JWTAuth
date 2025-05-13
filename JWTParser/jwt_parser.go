package JWTParser

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

func ParseUnverified(tokenString string, log *slog.Logger) (jwt.MapClaims, error) {
	const op = "requests.ParseUnverified"
	log = slog.With(
		slog.String("token", tokenString),
		slog.String("operation", op))
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Error("parsing error", err)
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}

func IsTokenExpired(claims jwt.MapClaims) bool {
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return true
	}
	return exp.Before(time.Now())
}

// GetExpirationTime возвращает время истечения токена и флаг валидности
func GetExpirationTime(claims jwt.MapClaims, log *slog.Logger) (time.Time, error) {
	const op = "JWTParser.GetExpirationTime"
	log = slog.With(
		slog.String("operation", op),
		slog.String("claims", fmt.Sprintf("%v", claims)))

	exp, err := claims.GetExpirationTime()
	if err != nil {
		log.Error("failed to get expiration time", "error", err)
		return time.Time{}, err
	}

	if exp == nil {
		err := fmt.Errorf("token has no expiration claim")
		log.Warn("token has no expiration claim")
		return time.Time{}, err
	}

	return exp.Time, nil

}
