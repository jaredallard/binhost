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

// Package server contains the code for creating an HTTP server for
// serving a binhost.
package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	_ "github.com/jackc/pgx/v5/stdlib" // Used by ent.
	"github.com/jaredallard/binhost/internal/config"
	"github.com/jaredallard/binhost/internal/dpi"
	"github.com/jaredallard/binhost/internal/ent"
	"github.com/jaredallard/binhost/internal/ent/target"
	"github.com/jaredallard/binhost/internal/packages"
)

// New creates a new Activity.
func New(deps *dpi.Dependencies) *Activity {
	return &Activity{&Server{deps}, deps.Conf}
}

// Activity is a service activity that spawns an HTTP server to serve
// the binhost application. Create using the New() function.
type Activity struct {
	srv *Server
	cfg *config.Config
}

type Server struct {
	deps *dpi.Dependencies
}

func (s *Server) listTargets(c fiber.Ctx) error {
	targets, err := s.deps.DB.Target.Query().All(c.Context())
	if err != nil {
		return fmt.Errorf("failed querying targets: %w", err)
	}

	return c.Status(fiber.StatusOK).JSON(targets)
}

func (s *Server) createTarget(c fiber.Ctx) error {
	targetName := c.Params("target")
	_, err := s.deps.DB.Target.Create().SetName(targetName).Save(c.Context())
	if err != nil {
		if ent.IsConstraintError(err) {
			return c.SendStatus(fiber.StatusConflict)
		}

		return fmt.Errorf("failed creating target: %w", err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (s *Server) uploadPackage(c fiber.Ctx) error {
	targetName := c.Params("target")
	if targetName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("missing target")
	}

	// Ensure the target exists
	t, err := s.deps.DB.Target.Query().Where(target.NameEQ(targetName)).First(c.Context())
	if t == nil || err != nil {
		return c.Status(fiber.StatusNotFound).SendString("target not found")
	}

	pkg, err := packages.New(c.Request().BodyStream())
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}
	defer c.Request().CloseBodyStream() //nolint:errcheck // Why: Best effort close body.
	defer pkg.Delete()                  //nolint:errcheck // Why: Best effort delete.

	// name suitable for logging
	logName := pkg.Category + "/" + pkg.Name + "-" + pkg.Version + "::" + pkg.Repository

	s.deps.Log.Info("uploading package", "package", logName, "target", t.Name)

	if err := s.deps.DB.Pkg.Create().
		SetName(pkg.Name).
		SetCategory(pkg.Category).
		SetRepository(pkg.Repository).
		SetTarget(t).
		SetVersion(pkg.Version).
		Exec(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (s *Server) getPackage(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (s *Server) getTargetPackageIndex(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Run starts the HTTP service activity. Blocks until the provided
// context is cancelled.
func (a *Activity) Run(ctx context.Context) error {
	app := fiber.New(fiber.Config{StreamRequestBody: true})

	app.Use(logger.New(logger.Config{
		LoggerFunc: func(c fiber.Ctx, data *logger.Data, cfg logger.Config) error {
			a.srv.deps.Log.Info("http request", "method", c.Method(), "path", c.OriginalURL(), "status", c.Response().StatusCode(), "duration", data.Stop.Sub(data.Start).String())
			return nil
		},
	})).Name("logger")

	app.Get("/v1/targets", a.srv.listTargets).Name("list targets")
	app.Post("/v1/targets/:target", a.srv.createTarget).Name("create target")
	app.Post("/v1/targets/:target/upload", a.srv.uploadPackage).Name("upload package")

	// Gentoo Paths
	app.Get("/t/:target/Packages", a.srv.getPackage)
	app.Get("/t/:target/*", a.srv.getTargetPackageIndex)

	return app.Listen(":8080", fiber.ListenConfig{
		GracefulContext:       ctx,
		DisableStartupMessage: a.cfg.LogLevel != "debug",
		EnablePrintRoutes:     a.cfg.LogLevel == "debug",
	})
}
