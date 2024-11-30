package config

import (
	"github.com/spf13/viper"
	"time"
)

type AppConfig struct {
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	JwtSecretKey         string        `mapstructure:"jwt_secret_key"`
}

var app AppConfig

func App() AppConfig {
	return app
}

func Load() error {
	return viper.UnmarshalKey("app", &app)
}
