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
	"net/http"
	"os"

	"github.com/jaredallard/binhost/internal/parser"
)

// main runs the binhost server.
func main() {
	req, err := http.NewRequest("GET", "https://gentoo.rgst.io/t/arm64/asahi/Packages", nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	pkgs, err := parser.ParsePackages(resp.Body)
	if err != nil {
		panic(fmt.Errorf("failed to parse packages: %w", err))
	}

	pkgs.EncodeInto(os.Stdout)
}
