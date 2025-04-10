package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func Load(filename string) Config {
	if filename == "" {
		filename = ".env"
	}
	if err := LoadEnv(filename); err != nil {
		fmt.Printf("Error loading %s file: %v\n", filename, err)
	}
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

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			err := os.Setenv(key, value)
			if err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
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
