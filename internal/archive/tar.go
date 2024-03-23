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

package archive

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// _ ensures that tarExtractor implements the Extractor interface.
var _ Extractor = (&tarExtractor{})

// TarExtractor implements the Extractor interface for tar archives.
type tarExtractor struct{}

func (t *tarExtractor) Extensions() []string {
	return []string{"tar", "tgz", "tar.gz", "tar.xz", "txz", "tar.bz2"}
}

func (t *tarExtractor) Extract(r io.Reader, ext, dest string) error {
	var container io.ReadCloser
	switch ext {
	case "tar":
		container = io.NopCloser(r)
	case "tgz", "tar.gz":
		var err error
		container, err = newGzipReader(r)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
	case "tar.bz2":
		container = newBzip2Reader(r)
	case "tar.xz", "txz":
		container = newXZReader(r)
	default:
		// This only happens if we're missing a case in the switch statement.
		return fmt.Errorf("unsupported tar extension: %s", ext)
	}
	defer container.Close()

	tr := tar.NewReader(container)
	for {
		h, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("failed to read tar header: %w", err)
		}

		path := filepath.Join(dest, h.Name)
		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			f, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				_ = f.Close() //nolint:errcheck // Why: Best effort to close the file.
				return fmt.Errorf("failed to copy file contents: %w", err)
			}

			if err := f.Close(); err != nil {
				return fmt.Errorf("failed to close file: %w", err)
			}
		default:
			return fmt.Errorf("unsupported file type in package (%s: %v)", h.Name, h.Typeflag)
		}

		if err := os.Chmod(path, os.FileMode(h.Mode)); err != nil {
			return fmt.Errorf("failed to set file permissions: %w", err)
		}

		if err := os.Chtimes(path, h.AccessTime, h.ModTime); err != nil {
			return fmt.Errorf("failed to set file times: %w", err)
		}

		// TODO(jaredallard): Symlinks, ownership information, etc...
	}

	return nil
}
