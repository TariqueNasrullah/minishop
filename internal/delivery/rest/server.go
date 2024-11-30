package rest

import (
	"context"
	"github.com/minishop/internal/delivery/rest/controller"
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
	authUcase := usecase.NewAuthUsecase(uRepo)

	authController := controller.NewAuthController(authUcase)

	router := http.NewServeMux()

	router.HandleFunc("/version", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
	})

	router.Handle("/v1/", http.StripPrefix("/v1", authController.Handle()))

	srv := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  time.Second * 100,
		WriteTimeout: time.Second * 100,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}
	return &srv
}
