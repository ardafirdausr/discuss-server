package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Logger() echo.MiddlewareFunc {
	config := middleware.LoggerConfig{
		Format:           "${time_custom} ${method} ${uri} (${status}) - ${latency_human} \t${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}
	return middleware.LoggerWithConfig(config)
}
