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
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if validationErrors, ok := err.(validator.ValidationErrors); ok {
		verr := &entity.ErrValidation{
			Message: "Invalid format data",
			Err:     err,
		}
		for _, validationError := range validationErrors {
			errorField := map[string]string{
				"field": validationError.Field(),
			}
			switch validationError.Tag() {
			case "required":
				errorField["message"] = fmt.Sprintf("%s is required", validationError.Field())
			case "min":
				errorField["message"] = fmt.Sprintf("Min value of %s is %s",
					validationError.Field(), validationError.Param())
			case "max":
				errorField["message"] = fmt.Sprintf("Max value of %s is %s",
					validationError.Field(), validationError.Param())
			}
			verr.ValidationErrors = append(verr.ValidationErrors, errorField)
		}
		return verr
	}

	return nil
}
