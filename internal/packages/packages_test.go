package packages_test

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaredallard/binhost/internal/packages"
	"gotest.tools/v3/assert"
)

func TestCanParseGpkg(t *testing.T) {
	f, err := os.Open("testdata/onepassword-cli-0-1.gpkg.tar")
	assert.NilError(t, err)

	pkg, err := packages.New(f)
	assert.NilError(t, err)

	spew.Dump(pkg)

	// Check that one field is set. Maybe one day check them all (I'm
	// lazy)
	assert.Equal(t, "onepassword-cli-0", pkg.PF)
}
