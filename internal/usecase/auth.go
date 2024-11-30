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
	// Basic validation. username and password in required. Instead of validating through validation package
	// If checker is faster.
	if request.Username == "" || request.Password == "" {
		return nil, domain.BadRequestError
	}

	// Db Query with the Username
	user, err := a.userRepo.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	// Compare request.Password with Stored Hashed Password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, domain.BadRequestError
	}

	// Token Generation process
	token, err := a.tokenService.Generate(ctx, minishipJwt.Payload{Aud: fmt.Sprint(user.ID), Name: user.Username})
	if err != nil {
		return nil, domain.InternalServerError
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
