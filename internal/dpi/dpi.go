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

// Package dpi creates dependencies for usage across the binhost server.
package dpi

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"log/slog"

	"entgo.io/ent/dialect"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/jackc/pgx/v5/stdlib" // Used by ent.
	"github.com/jaredallard/binhost/internal/config"
	"github.com/jaredallard/binhost/internal/ent"
)

// Dependencies contains dependencies for the binhost server that is
// passed around to various components.
type Dependencies struct {
	// DB is a database client
	DB *ent.Client

	// S3 is a S3 client
	S3 *minio.Client

	// Conf is the configuration for the binhost server.
	Conf *config.Config

	// Log is a configured logger that should be used for logging purposes.
	Log *slog.Logger
}

// New creates a new dependencies struct with all of the required
// clients and configuration.
func New(ctx context.Context, log *slog.Logger) (*Dependencies, error) {
	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.With("error", err).Error("failed to load configuration")
		os.Exit(1)
	}

	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPass),
		Host:   cfg.DBHost + ":" + cfg.DBPort,
		Path:   cfg.DBName,
	}
	sanitizedURL := *connURL
	sanitizedURL.User = url.UserPassword(cfg.DBUser, "REDACTED")

	log.Info("connecting to postgres", "url", sanitizedURL.String())
	db, err := sql.Open("pgx", connURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	if err := client.Schema.Create(ctx, schema.WithDropColumn(true), schema.WithDropIndex(true)); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	log.Info("connecting to S3", "endpoint", cfg.S3Endpoint)
	s3, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		DB:   client,
		S3:   s3,
		Conf: cfg,
		Log:  log,
	}, nil
}
