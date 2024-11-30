package domain

import "context"

type LoginRequest struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type LoginResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AuthUsecase interface {
	Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
	CreateUser(ctx context.Context, user UserCreateParameters) (User, error)
}
