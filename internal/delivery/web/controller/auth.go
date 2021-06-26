package controller

import (
	"net/http"
	"os"

	"github.com/ardafirdausr/discuss-server/internal/app"
	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/ardafirdausr/discuss-server/internal/service/auth"
	"github.com/ardafirdausr/discuss-server/internal/service/token"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	ucs *app.Usecases
}

func NewAuthController(ucs *app.Usecases) *AuthController {
	return &AuthController{ucs: ucs}
}

func (ctrl AuthController) Login(c echo.Context) error {
	googleAuth := entity.GoogleAuth{}
	if err := c.Bind(&googleAuth); err != nil {
		c.Logger().Error(err.Error())
		return echo.ErrInternalServerError
	}

	googleSSOClientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleAuthenticator := auth.NewGoogleSSOAuthenticator(googleSSOClientID)
	user, err := ctrl.ucs.AuthUsecase.SSO(googleAuth.TokenID, googleAuthenticator)
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	JWTSecretKey := os.Getenv("JWT_SECRET_KEY")
	JWTToknizer := token.NewJWTTokenizer(JWTSecretKey)
	JWTToken, err := ctrl.ucs.AuthUsecase.GenerateAuthToken(*user, JWTToknizer)
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	response := echo.Map{
		"message": "Login Successful",
		"data":    user,
		"token":   JWTToken,
	}
	return c.JSON(http.StatusOK, response)
}
