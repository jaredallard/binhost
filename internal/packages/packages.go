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
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jaredallard/binhost/internal/archive"
)

// supportedCompressionExtensions is a list of supported compression
// extensions that this package can handle.
var supportedCompressionExtensions = []string{"xz", "gz"}

// Package represents a Gentoo gpkg (xpkg is not supported).
type Package struct {
	// Name is the name of the package.

	// Fields below are populated from the extracted contents of the gpkg.

	// imagePath contains the path on disk to the extracted image archive.
	imagePath string
	// metadataPath contains the path on disk to the extracted metadata archive.
	metadataPath string
}

// New creates a new Package from the provided [io.ReadCloser]. The
// provided ReadCloser should be streaming the raw contents of a Gentoo
// package (gpkg).
//
// The package will be stored on disk in a temporary directory due to
// the nature of gpkgs being usually a large tarball.
func New(r io.ReadCloser) (*Package, error) {
	tmpDir, err := os.MkdirTemp("", "binhost-extract-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Cleanup the temp directory if we fail.
	var keepTempDir bool
	defer func() {
		if !keepTempDir {
			os.RemoveAll(tmpDir)
		}
	}()

	if err := archive.Extract(archive.ExtractOptions{
		Reader:    r,
		Extension: "tar", // gpkg files are tar archives.
	}, tmpDir); err != nil {
		return nil, fmt.Errorf("failed to extract gpkg: %w", err)
	}

	p, err := packageFromDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create package from extracted contents: %w", err)
	}

	keepTempDir = true
	return p, nil
}

// packageFromDir creates a Package from the extracted contents of a
// gpkg tar. The supplied directory's Manifest is used to validate the
// contents of the package.
func packageFromDir(dir string) (*Package, error) {
	expectedFiles := []string{"Manifest", "gpkg-1"}
	expectedArchives := []string{"image", "metadata"}

	for _, name := range expectedFiles {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			return nil, fmt.Errorf("package missing required file: %s", name)
		}
	}

	archives := make(map[string]string)
	for _, name := range expectedArchives {
		for _, ext := range supportedCompressionExtensions {
			archiveName := name + ".tar." + ext
			if _, err := os.Stat(filepath.Join(dir, archiveName)); err == nil {
				archives[name] = archiveName
				break
			}

			// Extract the archive
			if err := archive.Extract(archive.ExtractOptions{
				Path: filepath.Join(dir, archiveName)}, filepath.Join(dir, name),
			); err != nil {
				return nil, fmt.Errorf("failed to extract archive %s: %w", archiveName, err)
			}
		}
		if _, ok := archives[name]; !ok {
			return nil, fmt.Errorf("package missing required archive: %s", name)
		}
	}

}
