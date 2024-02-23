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
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// parseColonDocuments parses fields from a colon separated file and
// returns the documents as a slice of maps. Empty newlines are
// considered to be the end of a document.
func parseColonDocuments(r io.Reader) ([]map[string]string, error) {
	out := make([]map[string]string, 0)
	cur := make(map[string]string)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// if the line is empty, we've reached the end of the document
		if line == "" {
			out = append(out, cur)
			cur = make(map[string]string)
			continue
		}

		// split the line by the first colon
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		// set the key and value
		cur[parts[0]] = parts[1]
	}
	if scanner.Err() != nil {
		if errors.Is(scanner.Err(), io.EOF) {
			// Last entry, if cur is not empty (already had a newline?) append
			// it.
			if len(cur) > 0 {
				out = append(out, cur)
			}
			return out, nil
		}

		return nil, scanner.Err()
	}

	return out, nil
}

// Packages is a Gentoo binhost packages file.
type Packages []Package

// EncodeInto serializes the packages into the given writer using the
// standard Packages format.
//
// See: https://wiki.gentoo.org/wiki/Binary_package_guide
func (pkgs Packages) EncodeInto(w io.Writer) error {
	for _, pkg := range pkgs {
		if err := pkg.EncodeInto(w); err != nil {
			return err
		}

		// Write a newline to separate the packages
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

// Package is a Gentoo binhost package.
type Package struct {
	BDepends      string   `colon:"BDEPEND"`
	BuildID       string   `colon:"BUILD_ID"`
	BuildTime     string   `colon:"BUILD_TIME"`
	CPV           string   `colon:"CPV"`
	DefinedPhases string   `colon:"DEFINED_PHASES"`
	Depends       string   `colon:"DEPEND"`
	EAPI          int      `colon:"EAPI"`
	IUse          string   `colon:"IUSE"`
	Keywords      []string `colon:"KEYWORDS"`
	Licenses      []string `colon:"LICENSE"`
	Path          string   `colon:"PATH"`
	RDepends      string   `colon:"RDEPEND"`
	Requires      string   `colon:"REQUIRES"`
	SHA1          string   `colon:"SHA1"`
	Size          int      `colon:"SIZE"`
	Use           string   `colon:"USE"`
	ModifiedTime  int      `colon:"MTIME"`
	Repo          string   `colon:"REPO"`
}

// EncodeInto serializes the package into the given writer using the
// standard Package format.
func (pkg *Package) EncodeInto(w io.Writer) error {
	prv := reflect.ValueOf(pkg).Elem()
	prt := prv.Type()

	for i := 0; i < prv.NumField(); i++ {
		fv := prv.Field(i)
		ft := prt.Field(i)

		value := fv.Interface()
		if value == nil || fv.IsZero() {
			continue
		}

		key := ft.Tag.Get("colon")
		if key == "" {
			// default to the field name as it is in the struct
			key = ft.Name
		}
		switch v := value.(type) {
		case string:
			value = v
		case []string:
			value = strings.Join(v, " ")
		default:
			value = fmt.Sprintf("%v", v)
		}

		if _, err := w.Write([]byte(fmt.Sprintf("%s: %s\n", key, value))); err != nil {
			return err
		}
	}

	return nil
}

func ParsePackages(r io.Reader) (Packages, error) {
	docs, err := parseColonDocuments(r)
	if err != nil {
		return nil, err
	}

	spew.Dump(docs)

	pkgs := make(Packages, 0, len(docs))
	for _, doc := range docs {
		var pkg Package
		pkgT := reflect.TypeOf(pkg)
		pkgV := reflect.ValueOf(&pkg).Elem()

		// Create a map of the struct tags to the field index
		tagToField := make(map[string]int)
		for i := 0; i < pkgT.NumField(); i++ {
			tag := pkgT.Field(i).Tag.Get("colon")
			if tag == "" {
				tag = pkgT.Field(i).Name
			}
			tagToField[tag] = i
		}

		// Set the fields from the document
		for k, v := range doc {
			fieldIndex, ok := tagToField[k]
			if !ok {
				return nil, fmt.Errorf("unknown field: %s", k)
			}

			field := pkgV.Field(fieldIndex)
			switch field.Kind() {
			case reflect.String:
				field.SetString(v)
			case reflect.Slice:
				field.Set(reflect.ValueOf(strings.Fields(v)))
			case reflect.Int:
				i, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("invalid integer: %s", v)
				}
				field.SetInt(int64(i))
			default:
				return nil, fmt.Errorf("unsupported field type: %s", field.Kind())
			}

			// Remove the field from the map so we can check for missing
			// fields later.
			delete(tagToField, k)

		}

		// Check for missing fields
		if len(tagToField) > 0 {
			return nil, fmt.Errorf("unknown fields: %v", tagToField)
		}

		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}
