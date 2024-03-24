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

// Package main contains a CLI for starting the binhost server. See the
// README in the root of the repository for more information.
package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaredallard/binhost/internal/packages"
)

// main runs the binhost server.
func main() {
	f, err := os.Open("onepassword-cli-0-1.gpkg.tar")
	if err != nil {
		fmt.Println("failed to open file:", err)
		os.Exit(1)
	}

	pkg, err := packages.New(f)
	if err != nil {
		panic(err)
	}

	spew.Dump(pkg)
}
