package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig() *Config {
	// Go по умолчанию читает переменные окружения (os.Getenv),
	// но НЕ читает .env файлы автоматически.
	// godotenv.Load() загружает переменные из .env в окружение,
	// после чего они становятся доступны через os.Getenv()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	return &Config{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:   getEnv("POSTGRES_NAME", "geno-tree-db"),
		SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return defaultValue
}

func CreateDSN(c *Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
