package web

import (
	"os"

	"github.com/ardafirdausr/discuss-server/internal/app"
	"github.com/ardafirdausr/discuss-server/internal/delivery/web/controller"
	"github.com/ardafirdausr/discuss-server/internal/delivery/web/middleware"
	"github.com/ardafirdausr/discuss-server/internal/delivery/web/server"
)

func Start(app *app.App) {
	web := server.New()

	authController := controller.NewAuthController(app.Usecases)
	authGroup := web.Group("/auth")
	authGroup.POST("/login", authController.Login)

	JWTSecretKey := os.Getenv("JWT_SECRET_KEY")
	JWTmiddleware := middleware.JWT(JWTSecretKey)

	discussionController := controller.NewDiscussionController(app.Usecases)
	discussionGroup := web.Group("/discussions", JWTmiddleware)
	discussionGroup.POST("", discussionController.CreateDiscussion)
	discussionGroup.PUT("/:discussionId", discussionController.UpdateDiscussion)
	discussionGroup.DELETE("/:discussionId", discussionController.DeleteDiscussion)

	server.Start(web)
}
