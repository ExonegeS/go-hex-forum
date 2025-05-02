package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type (
	Config struct {
		Server        Server
		DataBase      DataBase
		SessionConfig SessionConfig
		Storage       Storage
	}

	Server struct {
		Address string
		Port    string
	}

	DataBase struct {
		DBUser     string
		DBPassword string
		DBHost     string
		DBPort     string
		DBName     string
	}

	SessionConfig struct {
		DefaultTTL    time.Duration
		MaxNameLength int64
	}

	Storage struct {
		Host          string
		Port          string
		MaxNameLength int64
	}
)

func NewConfig() *Config {
	return &Config{
		Server{
			Address: getEnvStr("ADDRESS", ""),
			Port:    getEnvStr("PORT", "8080"),
		},
		DataBase{
			DBUser:     getEnvStr("DB_USER", "hacker"),
			DBPassword: getEnvStr("DB_PASSWORD", "0000"),
			DBHost:     getEnvStr("DB_HOST", "localhost"),
			DBPort:     getEnvStr("DB_PORT", "5432"),
			DBName:     getEnvStr("DB_NAME", "forum"),
		},
		SessionConfig{
			DefaultTTL:    time.Duration(getEnvInt64("SESSION_TTL", 7*24*60*60)) * time.Second,
			MaxNameLength: getEnvInt64("SESSION_TOKEN_LENGTH", 10),
		},
		Storage{
			Host:          getEnvStr("STORAGE_HOST", "http://localhost"),
			Port:          getEnvStr("STORAGE_PORT", "6969"),
			MaxNameLength: getEnvInt64("STORAGE_CODE_LENGTH", 6),
		},
	}
}

func (d *DataBase) MakeConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.DBHost, d.DBPort, d.DBUser, d.DBPassword, d.DBName,
	)
}

func (s *Storage) MakeAddressString() string {
	return fmt.Sprintf(
		"%s:%s",
		s.Host,
		s.Port,
	)
}

func getEnvStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
