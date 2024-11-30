/*
Copyright © 2024 TARIQUE M NASRULLAH nasrullahtarique@gmail.com
*/

package cmd

import (
	"context"
	"github.com/minishop/config"
	"github.com/minishop/internal/domain"
	minishopJwt "github.com/minishop/internal/pkg/jwt"
	postgresRepo "github.com/minishop/internal/repository/postgres"
	"github.com/minishop/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/spf13/cobra"
)

// seekCmd represents the seek command
var seekCmd = &cobra.Command{
	Use:   "seek",
	Short: "Initialize db with dummy data",
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := "host=localhost user=minishop password=supersecretpasswd dbname=minishop port=5432 sslmode=disable TimeZone=Asia/Dhaka"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}

		uRepo := postgresRepo.NewUserRepository(db)

		jwtService := minishopJwt.NewTokenService([]byte(config.App().JwtSecretKey))
		authUseCase := usecase.NewAuthUsecase(uRepo, jwtService)

		if err = uRepo.AutoMigrate(); err != nil {
			return err
		}

		if _, err = authUseCase.CreateUser(context.Background(), domain.UserCreateParameters{
			Username: "01901901901@mailinator.com",
			Password: "321dsaf",
		}); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(seekCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seekCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seekCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
