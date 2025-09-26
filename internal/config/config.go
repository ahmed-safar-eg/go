package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func loadEnv() {
	// حاول تحميل من ملف .env أولاً
	err := godotenv.Load()
	if err != nil {
		// إذا لم يوجد ملف .env، استخدم المتغيرات البيئية النظامية
		log.Println("No .env file found, using system environment variables")
	}
}
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
    Mongo    MongoConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type MongoConfig struct {
    URI    string
    DBName string
}

type JWTConfig struct {
	Secret string
	Expiry int
}

func LoadConfig() *Config {
	loadEnv()
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8090"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "angazny"),
			Password: getEnv("DB_PASSWORD", "Angazny@123"),
			Name:     getEnv("DB_NAME", "angazny"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "secret"),
			Expiry: getEnvAsInt("JWT_EXPIRY", 24),
		},
        Mongo: MongoConfig{
            URI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
            DBName: getEnv("MONGO_DB", "appdb"),
        },
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}