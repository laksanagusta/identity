package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Logger   LoggerConfig
	JWT      JWTconfig
}

type AppConfig struct {
	Name    string
	Version string
	Port    string
	Env     string
	Key     string
	Domain  string
	StoreId string
}

type LoggerConfig struct {
	Mode  string
	Level string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTconfig struct {
	SecretKey string
}

func LoadConfig(env string) (Config, error) {
	v := viper.New()

	v.SetConfigName(fmt.Sprintf("config/config-%s", env))
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, errors.New("config file not found")
		}
		return Config{}, err
	}

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return Config{}, err
	}

	return c, nil
}

func LoadConfigV2() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	config := Config{
		App: AppConfig{
			Name:    os.Getenv("APP_NAME"),
			Version: os.Getenv("APP_VERSION"),
			Port:    os.Getenv("APP_PORT"),
			Env:     os.Getenv("APP_ENV"),
			Key:     os.Getenv("APP_KEY"),
			Domain:  os.Getenv("APP_DOMAIN"),
			StoreId: os.Getenv("APP_STOREID"),
		},
		Postgres: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
		},
		Logger: LoggerConfig{
			Mode:  os.Getenv("LOGGER_MODE"),
			Level: os.Getenv("LOGGER_LEVEL"),
		},
		JWT: JWTconfig{
			SecretKey: os.Getenv("JWT_SECRETKEY"),
		},
	}

	return config, nil
}
