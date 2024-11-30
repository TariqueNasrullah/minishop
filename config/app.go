package config

import (
	"github.com/spf13/viper"
	"time"
)

type AppConfig struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

var app AppConfig

func App() AppConfig {
	return app
}

func Load() error {
	return viper.UnmarshalKey("app", &app)
}
