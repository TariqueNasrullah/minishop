package config

import (
	"github.com/spf13/viper"
	"time"
)

// AppConfig is the application configuration type.
type AppConfig struct {
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	JwtSecretKey         string        `mapstructure:"jwt_secret_key"`
	Port                 int           `mapstructure:"port"`
}

var app AppConfig

func App() AppConfig {
	return app
}

// PostgresConfig holds Postgresql db configuration
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"db_name"`
	Timezone string `mapstructure:"timezone"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

var postgres PostgresConfig

func Postgres() PostgresConfig {
	return postgres
}

// Load parse the env/configuration file to the config variables.
func Load() error {
	if err := viper.UnmarshalKey("app", &app); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &postgres); err != nil {
		return err
	}
	return nil
}
