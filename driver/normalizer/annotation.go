package normalizer

import (
	"errors"

	"github.com/bblfsh/php-driver/driver/normalizer/phpast"

	"gopkg.in/bblfsh/sdk.v1/uast"
	. "gopkg.in/bblfsh/sdk.v1/uast/ann"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer/annotatter"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer/positioner"
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

var someAssignOp = Or(phpast.AssignOpPlus,
			phpast.AssignOpMinus,
			phpast.AssignOpMul,
			phpast.AssignOpDiv,
			phpast.AssignOpMod)

// AnnotationRules for the PHP language
var AnnotationRules = On(Any).Self(
	On(Not(phpast.File)).Error(errors.New("root must be uast.File")),
	On(phpast.File).Roles(uast.File, uast.Module).Descendants(
		// Misc
		On(phpast.Alias).Roles(uast.Statement, uast.Alias),
		On(phpast.Arg).Roles(uast.Argument),
		On(phpast.Array).Roles(uast.Expression, uast.Literal, uast.List),
		On(phpast.ArrayDimFetch).Roles(uast.Expression, uast.List, uast.Value,
			uast.Entry),
		On(phpast.ArrayItem).Roles(uast.Expression, uast.List, uast.Entry),
		On(phpast.Variable).Roles(uast.Identifier, uast.Variable),
		On(phpast.Name).Roles(uast.Expression, uast.Identifier).Self(
			On(HasInternalRole("class")).Roles(uast.Qualified),
		),
		On(phpast.NameRelative).Roles(uast.Expression, uast.Identifier, uast.Qualified, uast.Incomplete),
		On(phpast.Comment).Roles(uast.Noop, uast.Comment),
		On(phpast.Doc).Roles(uast.Noop, uast.Comment, uast.Documentation),
		On(phpast.Nop).Roles(uast.Noop),
		On(phpast.Echo).Roles(uast.Statement, uast.Call),
		On(phpast.Print).Roles(uast.Statement, uast.Call),
		On(phpast.Empty).Roles(uast.Expression, uast.Call),
		On(phpast.Isset).Roles(uast.Expression, uast.Call),
		On(phpast.Unset).Roles(uast.Expression, uast.Call),
		On(HasInternalRole("stmts")).Roles(uast.Expression, uast.Body),
		On(phpast.PropertyFetch).Roles(uast.Expression, uast.Map, uast.Identifier,
			uast.Entry, uast.Value),
		// no static in UAST
		On(phpast.StaticPropertyFetch).Roles(uast.Expression, uast.Map,
			uast.Identifier, uast.Entry, uast.Value, uast.Incomplete),
		// no error supress in UAST
		On(phpast.ErrorSuppress).Roles(uast.Expression, uast.Incomplete),
		On(phpast.Eval).Roles(uast.Expression, uast.Call),
		On(phpast.Exit).Roles(uast.Expression, uast.Call),
		On(phpast.Namespace).Roles(uast.Block),
		// no const in UAST
		On(phpast.Const).Roles(uast.Expression, uast.Variable, uast.Incomplete),
		On(phpast.ConstFetch).Roles(uast.Expression, uast.Variable, uast.Incomplete),
		On(phpast.FullyQualified).Roles(uast.Expression, uast.Variable, uast.Incomplete),
		On(phpast.ClassConstFetch).Roles(uast.Expression, uast.Type, uast.Incomplete),
		On(phpast.Clone).Roles(uast.Expression, uast.Call, uast.Incomplete),
		On(phpast.Param).Roles(uast.Argument),
		On(phpast.Closure).Roles(uast.Function, uast.Declaration, uast.Expression,
			uast.Anonymous),
		On(phpast.ClosureUse).Roles(uast.Visibility, uast.Incomplete),
		On(phpast.Coalesce).Roles(uast.Expression, uast.Incomplete),
	        On(HasInternalRole("cond")).Roles(uast.Condition),
	        On(Or(phpast.Use, phpast.UseUse)).Roles(uast.Alias),
	        On(Or(phpast.Yield, phpast.YieldFrom)).Roles(uast.Return, uast.Incomplete),

		// Control flow
		On(phpast.Break).Roles(uast.Statement, uast.Break),
		On(phpast.Continue).Roles(uast.Statement, uast.Continue),
		On(phpast.Return).Roles(uast.Statement, uast.Return),
		On(phpast.Throw).Roles(uast.Statement, uast.Throw),
		On(phpast.Goto).Roles(uast.Statement, uast.Goto),
		// no UAST role for goto target labels
		On(phpast.Label).Roles(uast.Statement, uast.Goto, uast.Incomplete),

		On(Or(phpast.Assign, someAssignOp)).Roles(uast.Expression, uast.Assignment).Children(
			On(HasInternalRole("var")).Roles(uast.Left),
			On(HasInternalRole("expr")).Roles(uast.Right),
		),

		On(HasInternalRole("left")).Roles(uast.Left),
		On(HasInternalRole("right")).Roles(uast.Right),
		On(phpast.Name).Self(On(HasToken("null")).Roles(uast.Null)),
		// no Nullable/Optional in UAST
		On(phpast.NullableType).Roles(uast.Type, uast.Incomplete),
		On(phpast.Global).Roles(uast.Visibility, uast.World),
		// no Static in UAST
		On(phpast.Static).Roles(uast.Visibility, uast.Type),
		On(phpast.StaticVar).Roles(uast.Expression, uast.Identifier, uast.Variable, uast.Visibility, uast.Type),
		On(phpast.InlineHTML).Roles(uast.String, uast.Literal, uast.Incomplete),
		On(phpast.List).Roles(uast.Call, uast.List),

		// Operators
		On(phpast.BinaryOpPlus).Roles(uast.Expression, uast.Operator, uast.Add),
		On(phpast.BinaryOpMinus).Roles(uast.Expression, uast.Operator, uast.Substract),
		On(phpast.BinaryOpMul).Roles(uast.Expression, uast.Operator, uast.Multiply),
		On(phpast.BinaryOpDiv).Roles(uast.Expression, uast.Operator, uast.Divide),
		On(phpast.BinaryOpMod).Roles(uast.Expression, uast.Operator, uast.Modulo),
		On(phpast.BinaryOpPow).Roles(uast.Expression, uast.Operator, uast.Incomplete),

		On(phpast.AssignOpPlus).Roles(uast.Expression, uast.Operator, uast.Add),
		On(phpast.AssignOpMinus).Roles(uast.Expression, uast.Operator, uast.Substract),
		On(phpast.AssignOpMul).Roles(uast.Expression, uast.Operator, uast.Multiply),
		On(phpast.AssignOpDiv).Roles(uast.Expression, uast.Operator, uast.Divide),
		On(phpast.AssignOpMod).Roles(uast.Expression, uast.Operator, uast.Modulo),

		On(phpast.BitwiseAnd).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Bitwise, uast.And),
		On(phpast.BitwiseOr).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Bitwise, uast.Or),
		On(phpast.BitwiseXor).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Bitwise, uast.Xor),
		On(phpast.BitwiseNot).Roles(uast.Expression, uast.Unary, uast.Operator,
			uast.Bitwise, uast.Not),

		On(phpast.BooleanAnd).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Boolean, uast.And),
		On(phpast.BooleanOr).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Boolean, uast.Or),
		On(phpast.BooleanXor).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Boolean, uast.Xor),
		On(phpast.BooleanNot).Roles(uast.Expression, uast.Operator, uast.Boolean,
			uast.Not),

		On(phpast.UnaryPlus).Roles(uast.Expression, uast.Unary, uast.Incomplete),
		On(phpast.UnaryMinus).Roles(uast.Expression, uast.Unary, uast.Incomplete),
		On(phpast.PreInc).Roles(uast.Expression, uast.Unary, uast.Increment),
		On(phpast.PostInc).Roles(uast.Expression, uast.Unary, uast.Increment, uast.Postfix),
		On(phpast.PreDec).Roles(uast.Expression, uast.Unary, uast.Decrement),
		On(phpast.PostDec).Roles(uast.Expression, uast.Unary, uast.Decrement, uast.Postfix),

	        // no join/concatenation role in UAST
	        On(phpast.Concat).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Add, uast.Incomplete),

		On(phpast.ShiftLeft).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Bitwise, uast.LeftShift),
		On(phpast.ShiftRight).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Bitwise, uast.RightShift),

		On(phpast.Equal).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.Equal),
		On(phpast.Identical).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.Identical),
		On(phpast.NotEqual).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.Not, uast.Equal),
		On(phpast.NotIdentical).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.Not, uast.Identical),
		On(phpast.GreaterOrEqual).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.GreaterThanOrEqual),
		On(phpast.SmallerOrEqual).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.LessThanOrEqual),
		On(phpast.Greater).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.GreaterThan),
		On(phpast.Smaller).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.LessThan),
		On(phpast.Spaceship).Roles(uast.Expression, uast.Binary, uast.Operator,
			uast.Relational, uast.GreaterThanOrEqual, uast.LessThanOrEqual),

		// Scalars
		On(phpast.ScalarString).Roles(uast.Expression, uast.Literal, uast.String),
		On(Or(phpast.ScalarLNumber, phpast.ScalarDNumber)).Roles(uast.Expression,
			uast.Literal, uast.Number),
		// __CLASS__ and similar constants. Also mising a Const role in the UAST.
		On(Or(phpast.ScalarMagicClass,
	              phpast.ScalarMagicDir,
	              phpast.ScalarMagicFile,
	              phpast.ScalarMagicFunction,
	              phpast.ScalarMagicLine,
	              phpast.ScalarMagicMethod,
	              phpast.ScalarMagicNamespace,
	              phpast.ScalarMagicTrait)).Roles(uast.Expression, uast.Literal,
			uast.Incomplete),


		// Switch
		On(phpast.Switch).Roles(uast.Statement, uast.Switch).Children(
			On(HasInternalRole("cond")).Roles(uast.Switch),
		),
		On(phpast.Case).Roles(uast.Statement, uast.Case).Self(
			On(Not(HasChild(HasInternalRole("cond")))).Roles(uast.Default),
		).Children(
			On(HasInternalRole("cond")).Roles(uast.Case),
		),

		// Casts... no Cast in the UAST
		On(phpast.CastArray).Roles(uast.Expression, uast.List, uast.Incomplete),
		On(phpast.CastBool).Roles(uast.Expression, uast.Boolean, uast.Incomplete),
		On(phpast.CastDouble).Roles(uast.Expression, uast.Number, uast.Incomplete),
		On(phpast.CastInt).Roles(uast.Expression, uast.Number, uast.Incomplete),
		On(phpast.CastObject).Roles(uast.Expression, uast.Type, uast.Incomplete),
		On(phpast.CastString).Roles(uast.Expression, uast.String, uast.Incomplete),
		On(phpast.CastUnset).Roles(uast.Expression, uast.Incomplete),

		// TryCatch
		On(phpast.TryCatch).Roles(uast.Statement, uast.Try),
		On(phpast.Catch).Roles(uast.Statement, uast.Catch).Children(
			On(HasInternalRole("types")).Roles(uast.Catch, uast.Type),
		),
		On(phpast.Finally).Roles(uast.Statement, uast.Finally),

		// Class
		On(phpast.Class).Roles(uast.Statement, uast.Declaration, uast.Type),
		On(HasInternalRole("extends")).Roles(uast.Base),
		On(HasInternalRole("implements")).Roles(uast.Implements),
		// FIXME: php-parser doesn't give Visibility information (public, private, etc)
		// no const in UAST
		On(phpast.ClassConst).Roles(uast.Type, uast.Variable, uast.Incomplete),
		// no member role in UAST
		On(Or(phpast.Property, phpast.PropertyProperty)).Roles(uast.Type,
			uast.Variable, uast.Incomplete),
		// ditto
		On(phpast.ClassMethod).Roles(uast.Type, uast.Function),

		// If + Ternary
		On(phpast.Ternary).Roles(uast.Expression, uast.If).Children(
			On(HasInternalRole("if")).Roles(uast.Then),
			On(HasInternalRole("else")).Roles(uast.Else),
		),
		On(phpast.If).Roles(uast.Statement, uast.If),
		On(phpast.ElseIf).Roles(uast.Statement, uast.If, uast.Else),
		On(phpast.Else).Roles(uast.Statement, uast.Else),

		// Declare, we interpret it as an assignment-ish
		On(phpast.Declare).Roles(uast.Expression, uast.Assignment, uast.Incomplete).Children(
			On(HasInternalRole("declares")).Roles(uast.Identifier, uast.Left).Children(
				On(HasInternalRole("value")).Roles(uast.Right),
			),
		),

		// While and DoWhile
		On(phpast.Do).Roles(uast.Statement, uast.DoWhile),
		On(phpast.While).Roles(uast.Statement, uast.While),

		// Encapsed; incomplete: no encapsed/ string varsubst in UAST
		On(phpast.Encapsed).Roles(uast.Expression, uast.Literal, uast.String,
			uast.Incomplete),
		On(phpast.EncapsedStringPart).Roles(uast.Expression, uast.Identifier,
			uast.Value),

		// For
		On(phpast.For).Roles(uast.Statement, uast.For).Children(
			On(HasInternalRole("init")).Roles(uast.Expression, uast.For, uast.Initialization),
			On(HasInternalRole("cond")).Roles(uast.For), // Condition role added elsewhere
			On(HasInternalRole("loop")).Roles(uast.Expression, uast.For, uast.Update),
		),

		// Foreach
		On(phpast.Foreach).Roles(uast.Statement, uast.For, uast.Incomplete).Children(
			On(HasInternalRole("valueVar")).Roles(uast.Iterator),
		),

		// FuncCalls, StaticCalls and MethodCalls
		On(phpast.FuncCall).Roles(uast.Expression, uast.Call).Children(
			On(HasInternalRole("name")).Roles(uast.Function, uast.Name),
		),
		On(phpast.StaticCall).Roles(uast.Expression, uast.Call, uast.Identifier).Children(
			On(HasInternalRole("class")).Roles(uast.Type, uast.Receiver),
		),
		On(phpast.MethodCall).Roles(uast.Expression, uast.Call, uast.Identifier).Children(
			On(HasInternalRole("var")).Roles(uast.Receiver, uast.Receiver, uast.Identifier),
		),

		// Function declarations
		On(phpast.Function).Roles(uast.Function, uast.Declaration).Children(
			On(phpast.Param).Self(
				// No reference/value in the UAST
				On(HasProperty("byRef", "true")).Roles(uast.Incomplete),
				On(HasProperty("variadic", "true")).Roles(uast.ArgsList),
			).Children(
				On(HasInternalRole("default")).Roles(uast.Default),
			),
			On(phpast.FunctionReturn).Roles(uast.Return, uast.Type),
		),

		// Include and require
		On(phpast.Include).Roles(uast.Import).Children(
			On(Any).Roles(uast.Import, uast.Pathname),
		),

		// Instanceof
		On(phpast.Instanceof).Roles(uast.Expression, uast.Call).Children(
			On(Any).Roles(uast.Argument),
			On(HasInternalRole("class")).Roles(uast.Type),
		),

		// Interface
		On(phpast.Interface).Roles(uast.Type, uast.Declaration).Children(
			On(HasInternalRole("extends")).Roles(uast.Type, uast.Base),
		),

		// Traits
		On(phpast.Trait).Roles(uast.Type, uast.Declaration),
		On(phpast.TraitUse).Roles(uast.Base),
		On(phpast.TraitPrecedence).Roles(uast.Base, uast.Alias, uast.Incomplete),
		On(HasInternalRole("insteadof")).Roles(uast.Alias, uast.Incomplete),

		// New
		On(phpast.New).Roles(uast.Expression, uast.Initialization, uast.Call).Children(
			On(HasInternalRole("class")).Roles(uast.Type),
		),
	),
)
