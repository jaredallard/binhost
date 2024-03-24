// Copyright (C) 2024 Jared Allard
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package config stores configuration for the binhost server.
package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config contains configuration for the binhost server.
type Config struct {
	// ListenAddress is the address to listen on.
	ListenAddress string `env:"LISTEN_ADDRESS" envDefault:":5100"`

	// LogLevel is the log level to use.
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// DBHost is the host of the database to connect to.
	DBHost string `env:"DB_HOST" envDefault:"localhost"`

	// DBPort is the port of the database to connect to.
	DBPort string `env:"DB_PORT" envDefault:"5432"`

	// DBUser is the user to connect to the database as.
	DBUser string `env:"DB_USER"`

	// DBPass is the password to connect to the database with.
	DBPass string `env:"DB_PASS"`

	// DBName is the name of the database to connect to.
	DBName string `env:"DB_NAME" envDefault:"binhost"`

	// DBSSLMode is the SSL mode to use when connecting to the database.
	DBSSLMode string `env:"DB_SSL_MODE" envDefault:"disable"`

	// S3Endpoint is the endpoint to connect to the S3 server at.
	S3Endpoint string `env:"S3_ENDPOINT" envDefault:"localhost:9000"`

	// S3AccessKey is the access key to use when connecting to the S3
	// server.
	S3AccessKey string `env:"S3_ACCESS_KEY"`

	// S3SecretKey is the secret key to use when connecting to the S3
	// server.
	S3SecretKey string `env:"S3_SECRET_KEY"`

	// S3Bucket is the bucket to store files in.
	S3Bucket string `env:"S3_BUCKET"`
}

// LoadConfig loads configuration from the environment and returns a
// Config struct.
func LoadConfig(log *slog.Logger) (*Config, error) {
	environment := strings.ToLower(os.Getenv("ENV"))

	var envFile string
	switch environment {
	case "dev", "development":
		envFile = ".env.development"
	case "prod", "production":
		envFile = ".env.production"
	default:
		envFile = ".env"
	}

	// If there's an environment file, load it.
	if envFile != "" {
		log.Info("loading environment file", "file", envFile)
		if err := godotenv.Load(envFile); err != nil {
			return nil, err
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
