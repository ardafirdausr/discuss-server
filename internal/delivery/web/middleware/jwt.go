package middleware

import (
	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWT(JWTSecretKey string) echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:       &entity.JWTPayload{},
		SigningKey:   []byte(JWTSecretKey),
		ErrorHandler: customJWTErrorHandler,
	}
	return middleware.JWTWithConfig(config)
}

func customJWTErrorHandler(err error) error {
	return echo.ErrUnauthorized
}
