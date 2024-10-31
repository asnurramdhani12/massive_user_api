package config

import (
	"context"
	"os"
	"user_api/utils/logger"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Configuration struct {
	AppName string `mapstructure:"APP_NAME" validate:"required"`
	Env     string `mapstructure:"ENVIRONMENT" validate:"required"`
	Host    string `mapstructure:"HOST" validate:"required"`
	Port    string `mapstructure:"PORT" validate:"required"`

	JWTSecret string `mapstructure:"JWT_SECRET" validate:"required"`

	PostgreSql PostgreSql `mapstructure:",squash"`
}

type PostgreSql struct {
	Host     string `mapstructure:"POSTGRES_HOST" validate:"required"`
	Port     string `mapstructure:"POSTGRES_PORT" validate:"required"`
	Username string `mapstructure:"POSTGRES_USER" validate:"required"`
	Password string `mapstructure:"POSTGRES_PASSWORD" validate:"required"`
	Database string `mapstructure:"POSTGRES_DATABASE" validate:"required"`
}

func NewConfig(ctx context.Context) (*Configuration, error) {
	var cfg Configuration

	viper.SetConfigType("env")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if !os.IsNotExist(err) {
		viper.SetConfigFile(envFile)

		if err := viper.ReadInConfig(); err != nil {
			logger.GetLogger(ctx).Errorf("Error reading config file, %v", err)
			return nil, err
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logger.GetLogger(ctx).Errorf("config error: %+v", err)
		}
		logger.GetLogger(ctx).Errorf("config error: %+v", err)
		return nil, err
	}

	logger.GetLogger(ctx).Infof("Config Loaded: %+v", cfg)
	return &cfg, nil
}
