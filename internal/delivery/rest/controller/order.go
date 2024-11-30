package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	minishopHttpError "github.com/minishop/internal/delivery/rest/errors"
	"github.com/minishop/internal/delivery/rest/middleware"
	"github.com/minishop/internal/domain"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderUsecase domain.OrderUsecase
}

func NewOrderController(e *echo.Group, orderUsecase domain.OrderUsecase, authMiddleware *middleware.Auth) *OrderController {
	controller := &OrderController{orderUsecase: orderUsecase}

	e.POST("/orders", authMiddleware.AuthRequired(controller.createOrder))
	return controller
}

func (o *OrderController) createOrder(c echo.Context) error {

	var (
		audStr string
		aud    uint64
		ok     bool
		err    error
	)

	audStr, ok = c.Get("aud").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, minishopHttpError.Unauthrized)
	}
	if aud, err = strconv.ParseUint(audStr, 10, 64); err != nil {
		return c.JSON(http.StatusBadRequest, minishopHttpError.Unauthrized)
	}

	var orderRequest domain.OrderCreateParameters
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	orderRequest.CreatedBy = aud

	ord, err := o.orderUsecase.Create(c.Request().Context(), orderRequest)
	if err != nil {
		var validationErr *domain.ValidationError
		if errors.As(err, &validationErr) {
			return c.JSON(http.StatusUnprocessableEntity, minishopHttpError.HTTPError{
				Message: "Please fix the given errors",
				Type:    "error",
				Code:    http.StatusUnprocessableEntity,
				Errors:  validationErr.ErrorMap,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, ord)
}