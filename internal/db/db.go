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

// TODO(jaredallard): This probably should implement a "store" for also
// storing the packages themselves. Maybe S3 or something else backed?

// Package db implements a thin wrapper around a DB for the purpose of
// storing package information.
package db

// Client contains a DB client.
type Client struct{}

// New creates a new DB client.
func New() (*Client, error) { return nil, nil }

// TODO(jaredallard): package struct from somewhere?
// NewPackage creates a new package and stores it.
func (c *Client) NewPackage() error { return nil }

func (c *Client) DeletePackage(id string) error { return nil }

func (c *Client) CreateTarget(name string) error { return nil }

func (p *Client) DeleteTarget(name string) error { return nil }
