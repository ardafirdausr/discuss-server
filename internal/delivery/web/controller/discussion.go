package controller

import (
	"log"
	"net/http"

	"github.com/ardafirdausr/discuss-server/internal/app"
	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type DiscussionController struct {
	ucs *app.Usecases
}

func NewDiscussionController(ucs *app.Usecases) *DiscussionController {
	return &DiscussionController{ucs: ucs}
}

func (ctrl DiscussionController) GetUserDiscussions(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*entity.JWTPayload)
	userID := claims.ID

	discussions, err := ctrl.ucs.DiscussionUsecase.GetAllUserDiscussions(userID)
	if err != nil {
		return err
	}

	return jsonResponse(c, http.StatusOK, "Success", discussions)
}

func (ctrl DiscussionController) CreateDiscussion(c echo.Context) error {
	var param entity.CreateDiscussionParam
	if err := c.Bind(&param); err != nil {
		return echo.ErrInternalServerError
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*entity.JWTPayload)
	param.CreatorID = claims.ID

	if err := c.Validate(&param); err != nil {
		log.Println(err.Error())
		return err
	}

	discussion, err := ctrl.ucs.DiscussionUsecase.Create(param)
	if err != nil {
		return err
	}

	return jsonResponse(c, http.StatusCreated, "Success", discussion)
}

func (ctrl DiscussionController) JoinDiscussion(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*entity.JWTPayload)
	userID := claims.ID
	discussionID := c.Param("discussionId")

	discussion, err := ctrl.ucs.DiscussionUsecase.JoinDiscussion(discussionID, userID)
	if err != nil {
		return err
	}

	return jsonResponse(c, http.StatusOK, "Success", discussion)
}

func (ctrl DiscussionController) LeaveDiscussion(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*entity.JWTPayload)
	userID := claims.ID
	discussionID := c.Param("discussionId")

	err := ctrl.ucs.DiscussionUsecase.LeaveDiscussion(discussionID, userID)
	if err != nil {
		return err
	}

	return jsonResponse(c, http.StatusOK, "Success", nil)
}

func (ctrl DiscussionController) UpdateDiscussion(c echo.Context) error {
	discussionID := c.Param("discussionId")
	discussion, err := ctrl.ucs.DiscussionUsecase.GetDiscussionByID(discussionID)
	if _, isErrNF := err.(entity.ErrNotFound); isErrNF {
		log.Println(err.Error())
		return echo.ErrNotFound
	}

	if err != nil {
		log.Println(err.Error())
		return err
	}

	param := entity.UpdateDiscussionParam{
		Code:        discussion.Code,
		Name:        discussion.Name,
		Description: discussion.Description,
	}
	if err := c.Bind(&param); err != nil {
		log.Println(err.Error())
		return echo.ErrInternalServerError
	}

	if err := c.Validate(&param); err != nil {
		log.Println(err.Error())
		return err
	}

	if err != nil {
		return err
	}

	if err = ctrl.ucs.DiscussionUsecase.Update(discussionID, param); err != nil {
		return err
	}

	discussion.Code = param.Code
	discussion.Name = param.Name
	discussion.Description = param.Description
	return jsonResponse(c, http.StatusOK, "Success", discussion)
}

func (ctrl DiscussionController) UpdateDiscussionPhoto(c echo.Context) error {
	return echo.ErrInternalServerError
}

func (ctrl DiscussionController) UpdateDiscussionPassword(c echo.Context) error {
	var param entity.UpdateDiscussionPassword
	if err := c.Bind(&param); err != nil {
		log.Println(err.Error())
		return err
	}

	if err := c.Validate(&param); err != nil {
		log.Println(err.Error())
		return err
	}

	discussionId := c.Param("discussionId")
	if err := ctrl.ucs.DiscussionUsecase.UpdatePassword(discussionId, param); err != nil {
		return err
	}

	return jsonResponse(c, http.StatusOK, "Success", nil)
}

func (ctrl DiscussionController) DeleteDiscussion(c echo.Context) error {
	discussionID := c.Param("discussionId")
	if err := ctrl.ucs.DiscussionUsecase.Delete(discussionID); err != nil {
		return err
	}

	return c.JSON(http.StatusNoContent, nil)
}
