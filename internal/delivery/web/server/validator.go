package server

import (
	"fmt"
	"net/http"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (v *CustomValidator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		verr := entity.ErrValidation{
			Message: "Invalid format data",
			Err:     err,
		}
		if len(validationErrors) > 0 {
			validationError := validationErrors[0]
			validationField := validationError.Field()
			validationParam := validationError.Param()
			switch validationError.Tag() {
			case "required":
				verr.Message = fmt.Sprintf("%s is required", validationField)
			case "min":
				verr.Message = fmt.Sprintf("Min value of %s is %s", validationField, validationParam)
			case "max":
				verr.Message = fmt.Sprintf("Max value of %s is %s", validationField, validationParam)
			}
		}

		return verr
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
