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

// Package archive implements a helper for extracting archives without
// needing to configure the extraction process.
package archive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Configures extractors supported by this package and values
// initialized by the init function.
var (
	extractors = []Extractor{&tarExtractor{}}
	extensions = map[string]Extractor{}
)

// init initializes calls all extractors to register their supported
// extensions.
func init() {
	for i := range extractors {
		for _, ext := range extractors[i].Extensions() {
			extensions[ext] = extractors[i]
		}
	}
}

// ExtractOptions contains the options for extracting an archive.
type ExtractOptions struct {
	// Reader is the io.Reader to read the archive from. Either [Reader]
	// or [Path] must be provided.
	Reader io.Reader

	// Extension is the extension of the archive to extract. This
	// overrides the extension detection from [Path] if provided. This is
	// required if [Reader] is provided.
	Extension string

	// Path is the path to the archive to extract. Either [Reader] or
	// [Path] must be provided.
	Path string
}

// Extract extracts an archive to the provided destination.
func Extract(opts ExtractOptions, dest string) error {
	if opts.Reader == nil && opts.Path == "" {
		return fmt.Errorf("either reader or path must be provided")
	}

	if opts.Reader != nil && opts.Path != "" {
		return fmt.Errorf("only one of reader or path can be provided")
	}

	ext := opts.Extension
	if opts.Reader != nil {
		if ext == "" {
			return fmt.Errorf("extension must be provided when using a reader (set opts.Extension)")
		}
	} else if opts.Path != "" && ext == "" {
		// If not set, default to the extension of the provided path.
		ext = filepath.Ext(opts.Path)
	}
	ext = strings.TrimPrefix(ext, ".")

	// Read the file from disk if Path is provided.
	if opts.Path != "" {
		r, err := os.Open(opts.Path)
		if err != nil {
			return fmt.Errorf("failed to open archive: %w", err)
		}
		defer r.Close()

		opts.Reader = r
	}

	for eext, extractor := range extensions {
		if ext == eext {
			return extractor.Extract(opts.Reader, ext, dest)
		}
	}

	return fmt.Errorf("unsupported archive extension: %s", ext)
}

// Extractor is an interface for extracting archives.
type Extractor interface {
	// Extract extracts all files from the provided reader to the
	// destination.
	Extract(r io.Reader, ext, dest string) error

	// Extensions should return a list of supported extensions for this
	// extractor.
	Extensions() []string
}
