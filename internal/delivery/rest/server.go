package rest

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/minishop/config"
	"github.com/minishop/internal/delivery/rest/controller"
	"github.com/minishop/internal/delivery/rest/middleware"
	"github.com/minishop/internal/domain"
	minishopJwt "github.com/minishop/internal/pkg/jwt"
	postgresRepo "github.com/minishop/internal/repository/postgres"
	"github.com/minishop/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net"
	"net/http"
	"os"
	"time"
)

func NewServer(ctx context.Context) (*http.Server, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Postgres().Host,
		config.Postgres().User,
		config.Postgres().Password,
		config.Postgres().DbName,
		config.Postgres().Port,
		config.Postgres().SSLMode,
		config.Postgres().Timezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		os.Exit(1)
	}

	// Service Or Other Package initializations
	jwtService := minishopJwt.NewTokenService([]byte(config.App().JwtSecretKey))
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Repository initializations
	uRepo := postgresRepo.NewUserRepository(db)
	orderRepo := postgresRepo.NewOrderRepository(db)

	// Usecase Initializations
	authUsecase := usecase.NewAuthUsecase(uRepo, jwtService)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	// Run Db Migrations and Seed Some data
	if err = uRepo.AutoMigrate(); err != nil {
		return nil, err
	}

	authUsecase.CreateUser(context.Background(), domain.UserCreateParameters{Username: "01901901901@mailinator.com", Password: "321dsaf"})

	if err = orderRepo.AutoMigrate(); err != nil {
		return nil, err
	}

	e := echo.New()
	v1Router := e.Group("/api/v1")

	controller.NewAuthController(v1Router, authUsecase, authMiddleware)
	controller.NewOrderController(v1Router, orderUsecase, authMiddleware)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", config.App().Port),
		Handler:      e,
		ReadTimeout:  time.Second * 100,
		WriteTimeout: time.Second * 100,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}
	return &srv, nil
}
