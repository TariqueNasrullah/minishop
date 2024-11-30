package jwt

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minishop/config"
	"time"
)

type TokenService struct {
	key []byte
}

type Payload struct {
	Aud  string `json:"aud"`
	Name string `json:"name"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Jti          string `json:"jti"`
}

func (t *TokenService) Generate(ctx context.Context, payload Payload) (Token, error) {
	jti := uuid.New().String()

	accessTokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":        payload.Aud,
		"name":       payload.Name,
		"jti":        jti,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(config.App().AccessTokenDuration).Unix(),
		"token_type": "access",
	})

	refreshTokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":        payload.Aud,
		"name":       payload.Name,
		"jti":        jti,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(config.App().RefreshTokenDuration).Unix(),
		"token_type": "refresh",
	})

	accessToken, err := accessTokenGenerator.SignedString(t.key)
	if err != nil {
		return Token{}, err
	}

	refreshToken, err := refreshTokenGenerator.SignedString(t.key)
	if err != nil {
		return Token{}, err
	}

	// TODO: Blacklist previous token

	return Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(config.App().RefreshTokenDuration),
		Jti:          jti,
	}, nil
}

func NewTokenService(key []byte) *TokenService {
	return &TokenService{key: key}
}
