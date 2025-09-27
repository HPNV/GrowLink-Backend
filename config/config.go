package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	CFG Cold
)

type Cold struct {
	DB     DatabaseConfig
	Server ServerConfig
}

type DatabaseConfig struct {
	Host      string `env:"DB_HOST"`
	Port      int    `env:"DB_PORT"`
	User      string `env:"DB_USER"`
	Password  string `env:"DB_PASSWORD"`
	Name      string `env:"DB_NAME"`
	Migration bool   `env:"DB_MIGRATION"`
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT"`
	Mode string `env:"GIN_MODE"`
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := envconfig.Process("", &CFG); err != nil {
		log.Fatal("Error processing environment variables:", err)
	}

	log.Printf("Configuration loaded: %+v\n", CFG)
}
