package fixtures

import (
	"path/filepath"
	"testing"

	"github.com/bblfsh/php-driver/driver/normalizer"
	"gopkg.in/bblfsh/sdk.v2/sdk/driver"
	"gopkg.in/bblfsh/sdk.v2/sdk/driver/fixtures"
)

const projectRoot = "../../"

var Suite = &fixtures.Suite{
	Lang: "php",
	Ext:  ".php",
	WriteYML: true,
	Path: filepath.Join(projectRoot, fixtures.Dir),
	NewDriver: func() driver.BaseDriver {
		return driver.NewExecDriverAt(filepath.Join(projectRoot, "build/bin/native"))
	},
	Transforms: driver.Transforms{
		Native: normalizer.Native,
		Code:   normalizer.Code,
	},
	BenchName: "complex",
}

func TestXxxDriver(t *testing.T) {
	Suite.RunTests(t)
}

func BenchmarkXxxDriver(b *testing.B) {
	Suite.RunBenchmarks(b)
}
