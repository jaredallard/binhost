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

package parser

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// encodeColonFormat serializes the given value into the writer using
// the colon format. v must be a pointer to a struct.
func encodeColonFormat(w io.Writer, v any) error {
	var prve reflect.Value
	var prt reflect.Type
	{ // Make sure we don't accidentally use prv since it is unsafe.
		prv := reflect.ValueOf(v)
		if prv.Kind() != reflect.Ptr {
			return fmt.Errorf("expected pointer to struct, got %T", v)
		}

		prve = prv.Elem()
		prt = prve.Type()
	}

	for i := 0; i < prve.NumField(); i++ {
		fv := prve.Field(i)
		ft := prt.Field(i)

		value := fv.Interface()
		if value == nil || fv.IsZero() {
			// Don't encode non-existent fields or fields with the zero value.
			continue
		}

		key := ft.Tag.Get("colon")
		if key == "" {
			// Skip fields without a colon tag.
			continue
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

	// Write a newline to end the document.
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

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
		return nil, scanner.Err()
	}

	// If we ended with a non-empty document, append it to the output now.
	if len(cur) > 0 {
		out = append(out, cur)
	}
	return out, nil
}

// decodeColonFormat deserializes the given document into the value v.
// The doc map should be from parseColonDocuments. v must be a pointer.
func decodeColonFormat(doc map[string]string, v any) error {
	var vrv reflect.Value
	var vrt reflect.Type
	{ // Make sure we don't accidentally use prv since it is unsafe.
		_vrv := reflect.ValueOf(v)
		if _vrv.Kind() != reflect.Ptr {
			return fmt.Errorf("expected pointer to struct, got %T", v)
		}

		vrv = _vrv.Elem()
		vrt = vrv.Type()
	}

	// Create a map of the colon struct tags into the field index for
	// iteration later.
	tagToField := make(map[string]int)
	for i := 0; i < vrt.NumField(); i++ {
		tag := vrt.Field(i).Tag.Get("colon")
		if tag == "" {
			continue
		}

		tagToField[tag] = i
	}

	// Set the fields from the document
	for k, v := range doc {
		fieldIndex, ok := tagToField[k]
		if !ok {
			return fmt.Errorf("unknown field: %s", k)
		}

		field := vrv.Field(fieldIndex)
		switch field.Kind() {
		case reflect.String:
			field.SetString(v)
		case reflect.Slice:
			field.Set(reflect.ValueOf(strings.Split(v, " ")))
		case reflect.Int:
			i, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("invalid integer: %s", v)
			}
			field.SetInt(int64(i))
		default:
			return fmt.Errorf("unsupported field type %s for field %s", field.Kind(), k)
		}

		// Remove the field from the map so we can check for missing
		// fields later.
		delete(doc, k)
	}

	// Check for missing fields
	if len(doc) > 0 {
		return fmt.Errorf("unknown fields: %v", doc)
	}

	return nil
}
