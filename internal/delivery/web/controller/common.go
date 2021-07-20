package controller

import "github.com/labstack/echo/v4"

func jsonResponse(c echo.Context, code int, message string, data interface{}) error {
	payload := echo.Map{"message": message}
	if data != nil {
		payload["data"] = data
	}
	return c.JSON(code, payload)
}
