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
	"strconv"
	"strings"

	"github.com/jaredallard/binhost/internal/archive"
	"github.com/jaredallard/binhost/internal/parser"
)

// supportedCompressionExtensions is a list of supported compression
// extensions that this package can handle.
var supportedCompressionExtensions = []string{"xz", "gz", "bz2"}

// Package represents a Gentoo gpkg (xpkg is not supported).
type Package struct {
	// Metadata contains package fields from a metadata.tar of a gpkg.
	Metadata

	// Fields below are populated from the extracted contents of the gpkg.

	// path is the path to the extracted gpkg on disk.
	path string
}

// Delete removes the package from disk.
func (p *Package) Delete() error {
	return os.RemoveAll(p.path)
}

// Metadata is the representation of a metadata.tar from a gpkg.
type Metadata struct {
	parser.PackageCommon

	// Name is the name of the package as calculated from the PF.
	Name string
	// Version is the version of the package as calculated from the PF.
	Version string

	CBuild        string
	CFlags        string
	CHost         string
	CXX           string
	CXXFlags      string
	Category      string
	Description   string
	Features      []string
	IUseEffective string
	Inherited     []string
	LDFLAGS       string
	PF            string
}

// New creates a new Package from the provided [io.ReadCloser]. The
// provided ReadCloser should be streaming the raw contents of a Gentoo
// package (gpkg).
//
// The package will be stored on disk in a temporary directory due to
// the nature of gpkgs being usually a large tarball.
func New(r io.Reader) (*Package, error) {
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

	// Move te files out of the sub dir by finding the first dir in the
	// temp dir.
	contents, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read temporary directory: %w", err)
	}

	var subDir string
	for _, f := range contents {
		if f.IsDir() {
			subDir = f.Name()
			break
		}
	}
	// If there's no sub dir, optimistically assume the contents are at
	// the root.
	if subDir != "" {
		contents, err := os.ReadDir(filepath.Join(tmpDir, subDir))
		if err != nil {
			return nil, fmt.Errorf("failed to read temp directory sub dir: %w", err)
		}

		for _, f := range contents {
			if err := os.Rename(filepath.Join(tmpDir, subDir, f.Name()), filepath.Join(tmpDir, f.Name())); err != nil {
				return nil, fmt.Errorf("failed to move file out of sub dir: %w", err)
			}
		}

		if err := os.Remove(filepath.Join(tmpDir, subDir)); err != nil {
			return nil, fmt.Errorf("failed to remove sub dir: %w", err)
		}
	}

	p, err := packageFromDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create package from extracted contents: %w", err)
	}
	p.path = tmpDir

	keepTempDir = true
	return p, nil
}

// metadataFromDir creates a package manifest out of the contents of an
// extracted manifest.tar file.
func metadataFromDir(dir string) (*Metadata, error) {
	md := &Metadata{}

	filesToFields := map[string]any{
		"BUILD_ID":   &md.BuildID,
		"BUILD_TIME": &md.BuildTime,
		"PF":         &md.PF,

		"CATEGORY":       &md.Category,
		"CBUILD":         &md.CBuild,
		"CFLAGS":         &md.CFlags,
		"CHOST":          &md.CHost,
		"CXXFLAGS":       &md.CXXFlags,
		"DEFINED_PHASES": &md.DefinedPhases,
		"DESCRIPTION":    &md.Description,
		"EAPI":           &md.EAPI,
		"FEATURES":       &md.Features,
		"INHERITED":      &md.Inherited,
		"IUSE":           &md.IUse,
		"IUSE_EFFECTIVE": &md.IUseEffective,
		"KEYWORDS":       &md.Keywords,
		"LDFLAGS":        &md.LDFLAGS,
		"LICENSE":        &md.Licenses,
		"repository":     &md.Repo,
		"SIZE":           &md.Size,
		"SLOT":           &md.Slot,
		"USE":            &md.Use,
	}

	for file, field := range filesToFields {
		data, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			// Can't read? It's likely just not set.
			continue
		}

		strData := strings.TrimSuffix(string(data), "\n")

		switch fv := field.(type) {
		case *int:
			i, err := strconv.Atoi(strings.TrimSpace(strData))
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s's value (%s) as int", file, string(data))
			}
			*fv = i
		case *string:
			*fv = strings.TrimSuffix(strData, "\n")
		case *[]string:
			*fv = strings.Split(strData, " ")
		default:
			panic(fmt.Errorf("unable to parse %s into %t", file, field))
		}
	}

	// Calculate name and version from PF. The version is the last part of
	// the string split by the '-' character.
	//
	// TODO(jaredallard): This will break on packages with - in the
	// version. I'm not sure how we're supposed to get the version then.
	parts := strings.Split(md.PF, "-")
	md.Name = strings.Join(parts[:len(parts)-1], "-")
	md.Version = parts[len(parts)-1]
	return md, nil
}

// packageFromDir creates a Package from the extracted contents of a
// gpkg tar. The supplied directory's Manifest is used to validate the
// contents of the package.
//
// TODO(jaredallard): We don't currently validate the Manifest.
func packageFromDir(dir string) (*Package, error) {
	expectedFiles := []string{"Manifest", "gpkg-1"}
	expectedArchives := []string{"image", "metadata"}

	for _, name := range expectedFiles {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			return nil, fmt.Errorf("package missing required file: %s", name)
		}
	}

	for _, name := range expectedArchives {
		var found bool
		for _, ext := range supportedCompressionExtensions {
			archiveName := name + ".tar." + ext
			if _, err := os.Stat(filepath.Join(dir, archiveName)); err != nil {
				continue
			}
			found = true

			// Extract the archive
			if err := archive.Extract(archive.ExtractOptions{
				Path: filepath.Join(dir, archiveName)}, dir,
			); err != nil {
				return nil, fmt.Errorf("failed to extract archive %s: %w", archiveName, err)
			}

			// Ensure we extracted to a directory with the same name as the archive.
			if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
				return nil, fmt.Errorf("failed to extract archive %s: %w", archiveName, err)
			}

			// We're done.
			break
		}
		if !found {
			return nil, fmt.Errorf("package missing required archive: %s", name)
		}
	}

	mf, err := metadataFromDir(filepath.Join(dir, "metadata"))
	if err != nil {
		return nil, fmt.Errorf("failed to create manifest from directory: %w", err)
	}

	// Create manifest from the extracted manifest.
	return &Package{
		Metadata: *mf,
	}, nil
}
