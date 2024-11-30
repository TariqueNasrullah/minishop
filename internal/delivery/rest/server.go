package rest

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/minishop/config"
	"github.com/minishop/internal/delivery/rest/controller"
	"github.com/minishop/internal/delivery/rest/middleware"
	minishopJwt "github.com/minishop/internal/pkg/jwt"
	postgresRepo "github.com/minishop/internal/repository/postgres"
	"github.com/minishop/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"net/http"
	"os"
	"time"
)

func New(ctx context.Context) *http.Server {
	dsn := "host=localhost user=minishop password=supersecretpasswd dbname=minishop port=5432 sslmode=disable TimeZone=Asia/Dhaka"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		os.Exit(1)
	}

	uRepo := postgresRepo.NewUserRepository(db)

	jwtService := minishopJwt.NewTokenService([]byte(config.App().JwtSecretKey))
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	authUcase := usecase.NewAuthUsecase(uRepo, jwtService)

	e := echo.New()
	v1Router := e.Group("/api/v1")

	controller.NewAuthController(v1Router, authUcase, authMiddleware)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      e,
		ReadTimeout:  time.Second * 100,
		WriteTimeout: time.Second * 100,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}
	return &srv
}
