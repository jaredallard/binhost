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

// Package main contains a CLI for starting the binhost server. See the
// README in the root of the repository for more information.
package main

import (
	"context"
	"os/signal"

	"log/slog"
	"os"

	charmlog "github.com/charmbracelet/log"

	_ "github.com/jackc/pgx/v5/stdlib" // Used by ent.

	"github.com/jaredallard/binhost/internal/dpi"
	"github.com/jaredallard/binhost/internal/server"
)

// main runs the binhost server.
func main() {
	handler := charmlog.New(os.Stderr)
	log := slog.New(handler)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	defer log.Info("shutting down")

	deps, err := dpi.New(ctx, log)
	if err != nil {
		log.With("error", err).Error("failed to create dependencies")
		os.Exit(1)
	}
	defer deps.DB.Close()

	log.Info("starting server")
	if err := server.New(deps).Run(ctx); err != nil {
		log.With("error", err).Error("failed to run server")
		os.Exit(1)
	}
}
