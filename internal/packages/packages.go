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

// Packages contains utilities for interacting with Gentoo packages.
package packages

import (
	"archive/tar"
	"io"
)

type Package struct {
	// TODO(jaredallard): Add fields for a package here.
}

// New creates a new Package from the provided [io.ReadCloser].
func New(r io.ReadCloser) (*Package, error) {
	var p Package

	t := tar.NewReader(r)

	// TODO(jaredallard): Once we have internet, we should actually do
	// something with the tar reader.
	_ = t
	return &p, nil
}
