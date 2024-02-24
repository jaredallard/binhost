package parser_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaredallard/binhost/internal/parser"
	"gotest.tools/v3/assert"
)

func TestCanEncodeAPackage(t *testing.T) {
	pkg := parser.Package{
		CPV: "x11-terms/alacritty-0.12.3",
	}

	var buf bytes.Buffer
	assert.NilError(t, pkg.EncodeInto(&buf))
	assert.Equal(t, "CPV: x11-terms/alacritty-0.12.3\n\n", buf.String())
}

func TestCanReadPackages(t *testing.T) {
	input := `ARCH: arm64` + "\n" + "\n" + `CPV: x11-terms/alacritty-0.12.3` + "\n"

	index, err := parser.ParsePackages(strings.NewReader(input))
	assert.NilError(t, err)
	assert.Equal(t, 1, len(index.PackageEntries))
	assert.Equal(t, "arm64", index.Arch)
	assert.Equal(t, "x11-terms/alacritty-0.12.3", index.PackageEntries[0].CPV)
}

func TestCanEncodePackages(t *testing.T) {
	index := parser.Index{
		Arch: "arm64",
		PackageEntries: []parser.Package{
			{
				CPV: "x11-terms/alacritty-0.12.3",
			},
		},
	}

	var buf bytes.Buffer
	assert.NilError(t, index.EncodeInto(&buf))
	assert.Equal(t, "ARCH: arm64\n\nCPV: x11-terms/alacritty-0.12.3\n\n", buf.String())
}
