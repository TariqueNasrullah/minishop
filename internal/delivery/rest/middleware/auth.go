package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	deliverError "github.com/minishop/internal/delivery/rest/errors"
	minishopJwt "github.com/minishop/internal/pkg/jwt"
	"net/http"
	"strings"
)

type Auth struct {
	tokenService *minishopJwt.TokenService
}

func NewAuthMiddleware(tokenService *minishopJwt.TokenService) *Auth {
	return &Auth{tokenService: tokenService}
}

func lookupBearerToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")

	if token == "" {
		return "", errors.New("no bearer token found")
	}

	if strings.HasPrefix(token, "Bearer ") {
		tokenPart := strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))
		if tokenPart != "" {
			return tokenPart, nil
		}
	}

	return "", errors.New("invalid token")
}

func (a *Auth) AuthRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := lookupBearerToken(c.Request())
		if err != nil {
			return c.JSON(http.StatusUnauthorized, deliverError.Unauthrized)
		}

		tokenPayload, err := a.tokenService.Parse(c.Request().Context(), token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, deliverError.Unauthrized)
		}

		c.Set("aud", tokenPayload.Aud)
		c.Set("username", tokenPayload.Name)
		return next(c)
	}
}
