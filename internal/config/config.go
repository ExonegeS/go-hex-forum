package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server Server
	DB     DataBase
}

type Server struct {
	Address string
	Port    string
}

type DataBase struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func Load() Config {
	return Config{
		Server{
			Address: getEnv("ADDRESS", ""),
			Port:    getEnv("PORT", "8080"),
		},
		DataBase{
			DBUser:     getEnv("DB_USER", "test"),
			DBPassword: getEnv("DB_PASSWORD", "test"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5444"),
			DBName:     getEnv("DB_NAME", "test"),
		},
	}
}

func (d *DataBase) MakeConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.DBHost, d.DBPort, d.DBUser, d.DBPassword, d.DBName,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
