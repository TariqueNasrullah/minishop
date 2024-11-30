package domain

import "context"

// LoginRequest Username and Password in required
type LoginRequest struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

// LoginResponse contains the access and refresh token
type LoginResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AuthUsecase interface {
	// Login receives a LoginRequest generates JWT token and return LoginResponse
	Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)

	// CreateUser is used to insert user into the persistence layer
	CreateUser(ctx context.Context, user UserCreateParameters) (User, error)
}
