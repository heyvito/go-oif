package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	src := `package main

import "io"

import (
	"github.com/test/proj-name/importc"
	"net"
	"github.com/test/proj-name/importb"
	"github.com/test/proj-name/importj"
	"os"
	ig "github.com/test/proj-name/importg"
	ni "github.com/test/proj-name/namedimport"
	_ "github.com/lib/side_effects_import"
	"github.com/foo/barA"
	"fmt"
	"github.com/test/proj-name/importi"

	"github.com/test/proj-name/importa"
	"github.com/foo/barC"
	"github.com/test/proj-name/importf"
	"github.com/test/proj-name/importh"
	"github.com/test/proj-name/importk"
	"github.com/foo/barB"
	"github.com/test/proj-name/importd"
	_ "github.com/test/proj-name/side-effects-import"
	"context"

	annotatedImport "github.com/foo/barB"
	"github.com/foo/barD"
	"log"
	"github.com/test/proj-name/importe"
)`

	str := FormatImports("github.com/test/proj-name", src)
	exp := `package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/foo/barA"
	"github.com/foo/barB"
	annotatedImport "github.com/foo/barB"
	"github.com/foo/barC"
	"github.com/foo/barD"
	_ "github.com/lib/side_effects_import"

	"github.com/test/proj-name/importa"
	"github.com/test/proj-name/importb"
	"github.com/test/proj-name/importc"
	"github.com/test/proj-name/importd"
	"github.com/test/proj-name/importe"
	"github.com/test/proj-name/importf"
	ig "github.com/test/proj-name/importg"
	"github.com/test/proj-name/importh"
	"github.com/test/proj-name/importi"
	"github.com/test/proj-name/importj"
	"github.com/test/proj-name/importk"
	ni "github.com/test/proj-name/namedimport"
	_ "github.com/test/proj-name/side-effects-import"
)
`

	assert.Equal(t, exp, str)
}
