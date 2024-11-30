package usecase

import (
	"context"
	"fmt"
	"github.com/minishop/internal/domain"
	minishipJwt "github.com/minishop/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo     domain.UserRepository
	tokenService *minishipJwt.TokenService
}

func (a *authUsecase) CreateUser(ctx context.Context, user domain.UserCreateParameters) (domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	user.Password = string(hashedPassword)
	cUser, err := a.userRepo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return cUser, nil
}

func (a *authUsecase) Login(ctx context.Context, request *domain.LoginRequest) (*domain.LoginResponse, error) {
	// validation
	user, err := a.userRepo.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, err
	}

	// generate token
	token, err := a.tokenService.Generate(ctx, minishipJwt.Payload{Aud: fmt.Sprint(user.ID), Name: user.Username})
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func NewAuthUsecase(userRepo domain.UserRepository, tokenService *minishipJwt.TokenService) domain.AuthUsecase {
	return &authUsecase{userRepo: userRepo, tokenService: tokenService}
}
