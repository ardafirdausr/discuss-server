package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ardafirdausr/discuss-server/internal/delivery/web/middleware"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// New instantates new Echo server
func New() *echo.Echo {
	debug := os.Getenv("DEBUG")
	isDebuging, _ := strconv.ParseBool(debug)

	e := echo.New()
	e.Debug = isDebuging
	e.Server.ReadTimeout = 30 * time.Second
	e.Server.WriteTimeout = 30 * time.Second

	validator := &CustomValidator{validator: validator.New()}
	e.Validator = validator

	SentryDsn := os.Getenv("SENTRY_DSN")
	sentryMiddleware := middleware.Sentry(SentryDsn, isDebuging)
	e.Use(sentryMiddleware)

	errorHandler := &CustomHTTPErrorHandler{debug: isDebuging, logger: e.Logger}
	e.HTTPErrorHandler = errorHandler.Handler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echoMiddleware.Secure())
	e.Use(echoMiddleware.CORS())
	return e
}

func Start(e *echo.Echo) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Start server
	go func() {
		e.Logger.Info("Starting server...")
		if err := e.Start(host + ":" + port); err != nil {
			e.Logger.Info("Shutting down the server. error: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
