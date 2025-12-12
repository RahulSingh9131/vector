package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary Primary `koanf:"primary" validate:"required"`
	Server ServerConfig `koanf:"server" validate:"required"`
	Database DatabaseConfig `koanf:"database" validate:"required"`
	Auth AuthConfig `koanf:"auth" validate:"required"`
	Redis RedisConfig `koanf:"redis" validate:"required"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerConfig struct {
	Port string `koanf:"port" validate:"required"`
	ReadTimeout string `koanf:"read_timeout" validate:"required"`
	WriteTimeout string `koanf:"write_timeout" validate:"required"`
	IdleTimeout string `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host string `koanf:"host" validate:"required"`
	Port string `koanf:"port" validate:"required"`
	User string `koanf:"user" validate:"required"`
	Password string `koanf:"password" validate:"required"`
	Name string `koanf:"name" validate:"required"`
	SSLMode string `koanf:"sslmode" validate:"required"`
	MaxOpenConns int `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns int `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime string `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime string `koanf:"conn_max_idle_time" validate:"required"`
}

type RedisConfig struct {
	Address string `koanf:"address" validate:"required"`
}

type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validate:"required"`
}

func LoadConfig() (*Config , error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	k := koanf.New(".")

	err := k.Load(env.Provider("VECTOR_", "." ,func (s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "VECTOR_"))
	}), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load environment variables")
	}

	mainconfig := &Config{}
	err = k.Unmarshal("", mainconfig)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to unmarshal environment variables")
	}

	validate := validator.New()
	if err := validate.Struct(mainconfig); err != nil {
		logger.Fatal().Err(err).Msg("Failed to validate environment variables")
	}

	return mainconfig, nil
}