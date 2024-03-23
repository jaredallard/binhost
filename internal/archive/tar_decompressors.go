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
	"compress/bzip2"
	"compress/gzip"
	"io"

	"github.com/jamespfennell/xz"
)

// newGzipReader creates a new gzip reader from the provided reader.
func newGzipReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

// newXZReader creates a new xz reader from the provided reader.
func newXZReader(r io.Reader) io.ReadCloser {
	return xz.NewReader(r)
}

// newBzip2Reader creates a new bzip2 reader from the provided reader.
func newBzip2Reader(r io.Reader) io.ReadCloser {
	return io.NopCloser(bzip2.NewReader(r))
}
