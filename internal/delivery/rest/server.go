package rest

import (
	"context"
	"fmt"
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Postgres().Host,
		config.Postgres().User,
		config.Postgres().Password,
		config.Postgres().DbName,
		config.Postgres().Port,
		config.Postgres().SSLMode,
		config.Postgres().Timezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		os.Exit(1)
	}

	uRepo := postgresRepo.NewUserRepository(db)
	orderRepo := postgresRepo.NewOrderRepository(db)

	jwtService := minishopJwt.NewTokenService([]byte(config.App().JwtSecretKey))
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	authUcase := usecase.NewAuthUsecase(uRepo, jwtService)
	orderUcase := usecase.NewOrderUsecase(orderRepo)

	e := echo.New()
	v1Router := e.Group("/api/v1")

	controller.NewAuthController(v1Router, authUcase, authMiddleware)
	controller.NewOrderController(v1Router, orderUcase, authMiddleware)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", config.App().Port),
		Handler:      e,
		ReadTimeout:  time.Second * 100,
		WriteTimeout: time.Second * 100,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}
	return &srv
}
