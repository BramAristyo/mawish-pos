package config

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Addr           string
	Handler        *gin.Engine
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	RunMode        string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxIdleCons     int
	MaxOpenCons     int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	AccessTokenExpireDuration time.Duration
	Secret                    string
}

type R2Config struct {
	AccountID       string
	AccessKey       string
	SecretAccessKey string
	BucketName      string
	PresignExpires  int
}

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	JWT      JWTConfig
	R2       R2Config
}

func GetConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config: " + err.Error())
	}

	viper.SetDefault("SERVER_ADDR", ":8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", "10s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "10s")
	viper.SetDefault("SERVER_MAX_HEADER_BYTES", 1048576)

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "postgres")
	viper.SetDefault("SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", "1h")

	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRE_DURATION", "24h")
	viper.SetDefault("JWT_SECRET_KEY", "changeme")

	return &Config{
		Server: ServerConfig{
			Addr:           viper.GetString("SERVER_ADDR"),
			ReadTimeout:    viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout:   viper.GetDuration("SERVER_WRITE_TIMEOUT"),
			MaxHeaderBytes: viper.GetInt("SERVER_MAX_HEADER_BYTES"),
			RunMode:        viper.GetString("SERVER_MODE"),
		},
		Postgres: PostgresConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetString("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			DbName:          viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("SSL_MODE"),
			MaxIdleCons:     viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenCons:     viper.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		JWT: JWTConfig{
			AccessTokenExpireDuration: viper.GetDuration("JWT_ACCESS_TOKEN_EXPIRE_DURATION"),
			Secret:                    viper.GetString("JWT_SECRET_KEY"),
		},
		R2: R2Config{
			AccountID:       viper.GetString("R2_ACCOUNT_ID"),
			AccessKey:       viper.GetString("R2_ACCESS_KEY"),
			SecretAccessKey: viper.GetString("R2_SECRET_ACCESS_KEY"),
			BucketName:      viper.GetString("R2_BUCKET_NAME"),
			PresignExpires:  viper.GetInt("R2_PRESIGN_EXPIRES"),
		},
	}
}
