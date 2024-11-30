package config

import "time"

type AppConfig struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

var app AppConfig

func App() AppConfig {
	return app
}
