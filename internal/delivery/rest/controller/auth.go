package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	httpError "github.com/minishop/internal/delivery/rest/errors"
	"github.com/minishop/internal/delivery/rest/middleware"
	"github.com/minishop/internal/domain"
	"net/http"
	"time"
)

type AuthController struct {
	authUsecase domain.AuthUsecase
}

func NewAuthController(e *echo.Group, authUsecase domain.AuthUsecase, authMiddleware *middleware.Auth) *AuthController {
	controller := &AuthController{authUsecase: authUsecase}

	e.POST("/login", controller.login)
	e.GET("/logout", authMiddleware.AuthRequired(controller.logout))

	return controller
}

func (a *AuthController) login(c echo.Context) error {
	var loginRequest domain.LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	auth, err := a.authUsecase.Login(c.Request().Context(), &loginRequest)
	if err != nil {
		if errors.Is(err, domain.BadRequestError) || errors.Is(err, domain.NotFoundError) {
			return c.JSON(http.StatusBadRequest, httpError.HTTPError{Message: "The user credentials were incorrect.", Type: "error", Code: http.StatusBadRequest})
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, auth)
}

func (a *AuthController) logout(c echo.Context) error {
	time.Sleep(time.Second * 8)
	return c.JSON(http.StatusOK, "ok!")
}
