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
	Path: filepath.Join(projectRoot, fixtures.Dir),
	NewDriver: func() driver.BaseDriver {
		return driver.NewExecDriverAt(filepath.Join(projectRoot, "build/bin/native"))
	},
	//UpdateNative:true,
	//UpdateUAST:true,
	Transforms: driver.Transforms{
		Preprocess: normalizer.Preprocess,
		Normalize:  normalizer.Normalize,
		Native:     normalizer.Native,
		Code:       normalizer.Code,
	},
	BenchName: "complex",
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
