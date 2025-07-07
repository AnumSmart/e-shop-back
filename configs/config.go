package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	SecretAcc       string
	SecretRef       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

const (
	timeExpAccessToken  = time.Minute * 15
	timeExpRefreshToken = time.Hour * 24
)

func LoadConfig() *Config {
	err := godotenv.Load("c:\\Son_Alex\\GO_projects\\e-commerce_proj\\simple_gin_server\\.env")
	if err != nil {
		fmt.Println("Error loading .env file, using default config", err.Error())
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			SecretAcc:       os.Getenv("JWT_ACC_SECRET"),
			SecretRef:       os.Getenv("JWT_REF_SECRET"),
			AccessTokenExp:  timeExpAccessToken,
			RefreshTokenExp: timeExpRefreshToken,
		},
	}
}
