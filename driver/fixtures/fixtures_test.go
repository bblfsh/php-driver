package fixtures

import (
	"path/filepath"
	"testing"

	"github.com/bblfsh/php-driver/driver/normalizer"
	"github.com/bblfsh/sdk/v3/driver"
	"github.com/bblfsh/sdk/v3/driver/fixtures"
	"github.com/bblfsh/sdk/v3/driver/native"
)

const projectRoot = "../../"

var Suite = &fixtures.Suite{
	Lang: "php",
	Ext:  ".php",
	Path: filepath.Join(projectRoot, fixtures.Dir),
	NewDriver: func() driver.Native {
		return native.NewDriverAt(filepath.Join(projectRoot, "build/bin/native"), native.UTF8)
	},
	Transforms: normalizer.Transforms,
	BenchName:  "complex",
	Semantic: fixtures.SemanticConfig{
		BlacklistTypes: []string{
			"Name",
			"Scalar_String",
			"Scalar_EncapsedStringPart",
			"Comment",
			"Comment_Doc",
			"Expr_Include",
			"Stmt_Use",
			"Stmt_UseUse",
			"Stmt_GroupUse",
			"Param",
			"Stmt_Function",
		},
	},
}

func TestPHPDriver(t *testing.T) {
	Suite.RunTests(t)
}

func BenchmarkPHPDriver(b *testing.B) {
	Suite.RunBenchmarks(b)
}
