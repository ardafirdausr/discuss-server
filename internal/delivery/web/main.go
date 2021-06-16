package web

import (
	"github.com/ardafirdausr/discuss-server/internal/app"
	"github.com/ardafirdausr/discuss-server/internal/delivery/web/controller"
	"github.com/ardafirdausr/discuss-server/internal/delivery/web/server"
)

func Start(app *app.App) {
	web := server.New()

	authController := controller.NewAuthController(app.Usecases)
	authGroup := web.Group("/auth")
	authGroup.POST("/login", authController.Login)

	server.Start(web)
}
