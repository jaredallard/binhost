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
	assert.Equal(t, "CPV: x11-terms/alacritty-0.12.3\n", buf.String())
}

// TODO: Implement this.
func TestCanReadPackages(t *testing.T) {
	input := `CPV: x11-terms/alacritty-0.12.3` + "\n"

	pkgs, err := parser.ParsePackages(strings.NewReader(input))
	assert.NilError(t, err)
	assert.Equal(t, 1, len(pkgs))
	assert.Equal(t, "x11-terms/alacritty-0.12.3", pkgs[0].CPV)
}
