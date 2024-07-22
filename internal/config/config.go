package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"ENV"`
	LogLevel string `env:"LOG"`

	GrpcPort int `env:"GRPC_PORT"`
	HttpPort int `env:"HTTP_PORT"`

	PGConfig PostgresConfig

	KafkaBrokers []string `env:"KAFKA_BROKERS"`
	KafkaTopic   string   `env:"KAFKA_TOPIC"`
}

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
	DBName   string `env:"POSTGRES_DBNAME" env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Load() Config {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatal("couldn't bind settings to config")
	}

	return config
}
