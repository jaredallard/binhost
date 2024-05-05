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

// Package parser implements a parser for parsing binhost files,
// primarily Packages and gpkg files.
package parser

import (
	"io"
)

// Index is a binhost index file (Packages).
type Index struct {
	AcceptLicense         string   `colon:"ACCEPT_LICENSE"`
	AcceptProperties      string   `colon:"ACCEPT_PROPERTIES"`
	AcceptRestrict        string   `colon:"ACCEPT_RESTRICT"`
	AcceptKeywords        []string `colon:"ACCEPT_KEYWORDS"`
	Arch                  string   `colon:"ARCH"`
	CBuild                string   `colon:"CBUILD"`
	CHost                 string   `colon:"CHOST"`
	ConfigProtect         []string `colon:"CONFIG_PROTECT"`
	ConfigProtectMask     []string `colon:"CONFIG_PROTECT_MASK"`
	ELibc                 string   `colon:"ELIBC"`
	Features              []string `colon:"FEATURES"`
	GentooMirrors         string   `colon:"GENTOO_MIRRORS"`
	IUseImplicit          string   `colon:"IUSE_IMPLICIT"`
	Kernel                string   `colon:"KERNEL"`
	Packages              int      `colon:"PACKAGES"`
	Profile               string   `colon:"PROFILE"`
	Timestamp             int      `colon:"TIMESTAMP"`
	Use                   string   `colon:"USE"`
	UseExpand             []string `colon:"USE_EXPAND"`
	UseExpandHidden       []string `colon:"USE_EXPAND_HIDDEN"`
	UseExpandImplicit     []string `colon:"USE_EXPAND_IMPLICIT"`
	UseExpandUnprefixed   []string `colon:"USE_EXPAND_UNPREFIXED"`
	UseExpandValuesArch   []string `colon:"USE_EXPAND_VALUES_ARCH"`
	UseExpandValuesELibc  []string `colon:"USE_EXPAND_VALUES_ELIBC"`
	UseExpandValuesKernel []string `colon:"USE_EXPAND_VALUES_KERNEL"`
	Version               int      `colon:"VERSION"`

	// PackageEntries is a slice of packages contained within the index.
	PackageEntries []Package
}

// EncodeInto serializes the packages into the given writer using the
// standard Packages format.
//
// See: https://wiki.gentoo.org/wiki/Binary_package_guide
func (index Index) EncodeInto(w io.Writer) error {
	// Write the header first.
	if err := encodeColonFormat(w, &index); err != nil {
		return err
	}

	for _, pkg := range index.PackageEntries {
		if err := pkg.EncodeInto(w); err != nil {
			return err
		}
	}

	return nil
}

// PackageCommon is data shared between both a Packages index and the
// actual package's metadata.tar.xz distributed with it.
type PackageCommon struct {
	BDepends      string   `colon:"BDEPEND"`
	BuildID       string   `colon:"BUILD_ID"`
	BuildTime     string   `colon:"BUILD_TIME"`
	DefinedPhases []string `colon:"DEFINED_PHASES"`
	EAPI          int      `colon:"EAPI"`
	ELibc         string   `colon:"ELIBC"`
	IUse          string   `colon:"IUSE"`
	Use           string   `colon:"USE"`
	Keywords      []string `colon:"KEYWORDS"`
	Licenses      []string `colon:"LICENSE"`
	Slot          string   `colon:"SLOT"`
	Depends       string   `colon:"DEPEND"`
	IDepend       string   `colon:"IDEPEND"`
	PDepends      string   `colon:"PDEPEND"`
	RDepends      string   `colon:"RDEPEND"`
	Requires      []string `colon:"REQUIRES"`
	Restrict      string   `colon:"RESTRICT"`
	Provides      []string `colon:"PROVIDES"`
	Size          int      `colon:"SIZE"`
	Repo          string   `colon:"REPO"`
}

// Package is a Gentoo binhost package.
type Package struct {
	PackageCommon

	CPV          string `colon:"CPV"`
	Path         string `colon:"PATH"`
	SHA1         string `colon:"SHA1"`
	MD5          string `colon:"MD5"`
	ModifiedTime int    `colon:"MTIME"`
}

// EncodeInto serializes the package into the given writer using the
// standard Package format.
func (pkg *Package) EncodeInto(w io.Writer) error {
	return encodeColonFormat(w, pkg)
}

// ParsePackages parses the provided reader into a the Index type.
func ParsePackages(r io.Reader) (*Index, error) {
	docs, err := parseColonDocuments(r)
	if err != nil {
		return nil, err
	}

	// Read the first document into the index.
	if len(docs) == 0 {
		return nil, nil
	}

	var index Index
	if err := decodeColonFormat(docs[0], &index); err != nil {
		return nil, err
	}

	// No package entries, return the index.
	if len(docs) == 1 {
		return &index, nil
	}

	index.PackageEntries = make([]Package, len(docs)-1)

	// Read the rest of the documents into the package entries.
	for i, doc := range docs[1:] {
		var pkg Package
		if err := decodeColonFormat(doc, &pkg); err != nil {
			return nil, err
		}
		index.PackageEntries[i] = pkg
	}

	return &index, nil
}
