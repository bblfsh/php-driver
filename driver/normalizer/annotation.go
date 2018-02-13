package normalizer

import (
	"errors"

	"github.com/bblfsh/python-driver/driver/normalizer/pyast"

	"gopkg.in/bblfsh/sdk.v1/uast/transformer/positioner"
	"gopkg.in/bblfsh/sdk.v1/uast"
	. "gopkg.in/bblfsh/sdk.v1/uast/ann"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer/annotatter"
)

// Transformers is the list of `transformer.Transfomer` to apply to a UAST, to
// learn more about the Transformers and the available ones take a look to:
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast/transformers
var Transformers = []transformer.Tranformer{
	annotatter.NewAnnotatter(AnnotationRules),
	positioner.NewFillLineColFromOffset(),
}

// AnnotationRules describes how a UAST should be annotated with `uast.Role`.
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast/ann
var AnnotationRules = On(Any).Self(
	On(Not(phpast.File)).Error(errors.New("root must be uast.File")),
	On(phpast.File).Roles(uast.File, uast.Module).Descendants(
	)
)
