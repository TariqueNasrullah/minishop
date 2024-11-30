/*
Copyright © 2024 TARIQUE M NASRULLAH nasrullahtarique@gmail.com
*/

package cmd

import (
	"context"
	"fmt"
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
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			config.Postgres().Host,
			config.Postgres().User,
			config.Postgres().Password,
			config.Postgres().DbName,
			config.Postgres().Port,
			config.Postgres().SSLMode,
			config.Postgres().Timezone,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}

		uRepo := postgresRepo.NewUserRepository(db)
		orderRepo := postgresRepo.NewOrderRepository(db)

		jwtService := minishopJwt.NewTokenService([]byte(config.App().JwtSecretKey))
		authUseCase := usecase.NewAuthUsecase(uRepo, jwtService)

		if err = uRepo.AutoMigrate(); err != nil {
			return err
		}

		if _, err = authUseCase.CreateUser(context.Background(), domain.UserCreateParameters{
			Username: "01901901901@mailinator.com",
			Password: "321dsaf",
		}); err != nil {
			fmt.Println("could not insert dummy user: ", err)
		}

		if err = orderRepo.AutoMigrate(); err != nil {
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
