package controller

import (
	"github.com/labstack/echo/v4"
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
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, auth)
}

func (a *AuthController) logout(c echo.Context) error {
	time.Sleep(time.Second * 8)
	return c.JSON(http.StatusOK, "ok!")
}
