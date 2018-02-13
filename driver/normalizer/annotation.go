package normalizer

import (
	"errors"

	"github.com/bblfsh/php-driver/driver/normalizer/phpast"

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

var someAssignOp = On(Or(phpast.AssignOpPlus,
                         phpast.AssignOpMinus,
                         phpast.AssignOpMul,
                         phpast.AssignOpDiv,
                         phpast.AssignOpMod))

var AnnotationRules = On(Any).Self(
	On(Not(phpast.File)).Error(errors.New("root must be uast.File")),
	On(phpast.File).Roles(uast.File, uast.Module).Descendants(
		On(phpast.Alias).Roles(uast.Statement, uast.Alias),
		On(phpast.Arg).Roles(uast.Argument),
		On(phpast.Array).Roles(uast.Expression, uast.Literal, uast.List),
		On(phpast.ArrayDimFetch).Roles(uast.Expression, uast.List, uast.Value, uast.Incomplete),
		On(phpast.ArrayItem).Roles(uast.Expression, uast.List, uast.Value),

		On(Or(phpast.Assign, someAssignOp).roles(uast.Expression, uast.Assignment).Roles(uast.Expression,
			uast.Assignment).Children(
				On(HasInternalRole("var").Roles(uast.Left),
				On(HasInternalRole("expr").Roles(uast.Right),
			),
		),
		On(phpast.AssignOpPlus).Roles(uast.Operator, uast.Add),
		On(phpast.AssignOpMinus).Roles(uast.Operator, uast.Substract),
		On(phpast.AssignOpMul).Roles(uast.Operator, uast.Multiply),
		On(phpast.AssignOpDiv).Roles(uast.Operator, uast.Divide),
		On(phpast.AssignOpMod).Roles(uast.Operator, uast.Modulo),
	),
)
