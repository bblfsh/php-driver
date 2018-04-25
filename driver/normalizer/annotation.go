package normalizer

import (
	php "github.com/bblfsh/php-driver/driver/normalizer/phpast"

	"gopkg.in/bblfsh/sdk.v1/uast"
	"gopkg.in/bblfsh/sdk.v1/uast/role"
	. "gopkg.in/bblfsh/sdk.v1/uast/transformer"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer/positioner"
)

var Native = Transformers([][]Transformer{
	{
		ResponseMetadata{
			TopLevelIsRootNode: false,
		},
	},
	{Mappings(Annotations...)},
	{RolesDedup()},
}...)

var Code = []CodeTransformer{
	positioner.NewFillOffsetFromLineCol(),
}

// FIXME: move to the SDK and remove from here and the python driver
func annotateTypeToken(typ, token string, roles ...role.Role) Mapping {
	return AnnotateTypeFields(typ,
		FieldRoles{
			uast.KeyToken: {Add: true, Op: String(token)},
		}, roles...)
}

func mapInternalProperty(key string, roles ...role.Role) Mapping {
	return Map(key,
		Part("other", Obj{
			key: ObjectRoles(key),
		}),
		Part("other", Obj{
			key: ObjectRoles(key, roles...),
		}),
	)
}

func annAssign(typ string, opRoles ...role.Role) Mapping {
	return AnnotateType(typ, ObjRoles{
		"var": {role.Left},
		"expr": {role.Right},
	}, opRoles...)
}

// XXX Missing:
// Callbacks from the ToNoder object.
// Name with parts ["Null"] should get role.Null
// Function declaration
// Case
// Tokens

// FIXME XXX: also migrate the callbacks in tonode.go
var Annotations = []Mapping{
	ObjectToNode{
		InternalTypeKey: "nodeType",
	}.Mapping(),
	ObjectToNode{
		LineKey:   "attributes.startLine",
		ColumnKey: "attributes.startTokenPos",
		OffsetKey: "attributes.startFilePos",
	}.Mapping(),
	ObjectToNode{
		EndLineKey:   "attributes.endLine",
		EndColumnKey: "attributes.endTokenPos",
		EndOffsetKey: "attributes.endFilePos",
	}.Mapping(),

	// FIXME: remove this
	mapInternalProperty("attributes", role.Incomplete),

	// FIXME: doesnt work for stmts (neither with mapInternalProperty)
	// check after other parts annotate this
	Map("statements",
		Part("m", Obj{
			"stmts": ObjectRoles("stmts"),
		}),
		Part("m", Obj{
			"stmts": ObjectRoles("stmts", role.Body),
		}),
	),
	//mapInternalProperty("stmts", role.Expression, role.Body),
	mapInternalProperty("left", role.Left),
	mapInternalProperty("right", role.Right),

	AnnotateType(php.File, nil, role.File),

	AnnotateType(php.Name, ObjRoles{
		"class": {role.Qualified},
	}, role.Expression, role.Identifier),

	annAssign(php.Assign, role.Expression, role.Assignment),
	annAssign(php.AssignOpMinus, role.Expression, role.Assignment, role.Operator, role.Substract),
	annAssign(php.AssignOpPlus, role.Expression, role.Assignment, role.Operator, role.Add),
	annAssign(php.AssignOpMul, role.Expression, role.Assignment, role.Operator, role.Multiply),
	annAssign(php.AssignOpDiv, role.Expression, role.Assignment, role.Operator, role.Divide),
	annAssign(php.AssignOpMod, role.Expression, role.Assignment, role.Operator, role.Modulo),

	// __CLASS__ and similar constants. Also mising a Const role in the UAST.
	AnnotateType(php.ScalarMagicClass, nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType(php.ScalarMagicDir, nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType(php.ScalarMagicFile, nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType(php.ScalarMagicFunction, nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType(php.ScalarMagicLine, nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType(php.ScalarMagicMethod, nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType(php.ScalarMagicNamespace, nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType(php.ScalarMagicTrait, nil, role.Expression, role.Literal, role.Incomplete),
	AnnotateType(php.Alias, nil, role.Statement, role.Alias),
	AnnotateType(php.Arg, nil, role.Argument),
	AnnotateType(php.Array, nil, role.Expression, role.Literal, role.List),
	AnnotateType(php.ArrayDimFetch, nil, role.Expression, role.List, role.Value, role.Entry),
	AnnotateType(php.ArrayItem, nil, role.Expression, role.List, role.Entry),
	AnnotateType(php.Variable, nil, role.Identifier, role.Variable),
	AnnotateType(php.NameRelative, nil, role.Expression, role.Identifier, role.Qualified, role.Incomplete),
	AnnotateType(php.Comment, nil, role.Noop, role.Comment),
	AnnotateType(php.Doc, nil, role.Noop, role.Comment, role.Documentation),
	AnnotateType(php.Nop, nil, role.Noop),
	AnnotateType(php.Echo, nil, role.Statement, role.Call),
	AnnotateType(php.Print, nil, role.Statement, role.Call),
	AnnotateType(php.Empty, nil, role.Expression, role.Call),
	AnnotateType(php.Isset, nil, role.Expression, role.Call),
	AnnotateType(php.Unset, nil, role.Expression, role.Call),
	AnnotateType(php.PropertyFetch, nil, role.Expression, role.Map, role.Identifier, role.Entry, role.Value),

	// no static in UAST
	AnnotateType(php.StaticPropertyFetch, nil, role.Expression, role.Map, role.Identifier,
		role.Entry, role.Value, role.Incomplete),

	// no error supress in UAST
	AnnotateType(php.ErrorSuppress, nil, role.Expression, role.Incomplete),
	AnnotateType(php.Eval, nil, role.Expression, role.Call),
	AnnotateType(php.Exit, nil, role.Expression, role.Call),
	AnnotateType(php.Namespace, nil, role.Block),
	// no const in UAST
	AnnotateType(php.Const, nil, role.Expression, role.Variable, role.Incomplete),
	AnnotateType(php.ConstFetch, nil, role.Expression, role.Variable, role.Incomplete),
	AnnotateType(php.FullyQualified, nil, role.Expression, role.Variable, role.Incomplete),
	AnnotateType(php.ClassConstFetch, nil, role.Expression, role.Type, role.Incomplete),
	AnnotateType(php.Clone, nil, role.Expression, role.Call, role.Incomplete),
	AnnotateType(php.Param, nil, role.Argument),
	AnnotateType(php.Closure, nil, role.Function, role.Declaration, role.Expression, role.Anonymous),
	AnnotateType(php.ClosureUse, nil, role.Visibility, role.Incomplete),
	AnnotateType(php.Coalesce, nil, role.Expression, role.Incomplete),
	AnnotateType(php.Use, nil, role.Alias),
	AnnotateType(php.UseUse, nil, role.Alias),
	AnnotateType(php.Yield, nil, role.Return, role.Incomplete),
	AnnotateType(php.YieldFrom, nil, role.Return, role.Incomplete),

	// Control flow
	AnnotateType(php.Break, nil, role.Statement, role.Break),
	AnnotateType(php.Continue, nil, role.Statement, role.Continue),
	AnnotateType(php.Return, nil, role.Statement, role.Return),
	AnnotateType(php.Throw, nil, role.Statement, role.Throw),
	AnnotateType(php.Goto, nil, role.Statement, role.Goto),

	// no role role for goto target labels
	AnnotateType(php.Label, nil, role.Statement, role.Goto, role.Incomplete),

	// XXX Implement this (native annotations were wrong too), must
	// have the role.Null if has a parts key and the single element in the
	// array is the string "NULL".
	//AnnotateType(php.Name, nil, role(HasToken("null")).Roles(role.Null)),

	// no Nullable/Optional in UAST
	AnnotateType(php.NullableType, nil, role.Type, role.Incomplete),
	AnnotateType(php.Global, nil, role.Visibility, role.World),

	// no Static in UAST
	AnnotateType(php.Static, nil, role.Visibility, role.Type),
	AnnotateType(php.StaticVar, nil, role.Expression, role.Identifier, role.Variable,
		role.Visibility, role.Type),
	AnnotateType(php.InlineHTML, nil, role.String, role.Literal, role.Incomplete),
	AnnotateType(php.List, nil, role.Call, role.List),

	// Operators
	AnnotateType(php.BinaryOpPlus, nil, role.Expression, role.Operator, role.Add),
	AnnotateType(php.BinaryOpMinus, nil, role.Expression, role.Operator, role.Substract),
	AnnotateType(php.BinaryOpMul, nil, role.Expression, role.Operator, role.Multiply),
	AnnotateType(php.BinaryOpDiv, nil, role.Expression, role.Operator, role.Divide),
	AnnotateType(php.BinaryOpMod, nil, role.Expression, role.Operator, role.Modulo),
	AnnotateType(php.BinaryOpPow, nil, role.Expression, role.Operator, role.Incomplete),

	AnnotateType(php.AssignOpPlus, nil, role.Expression, role.Operator, role.Add),
	AnnotateType(php.AssignOpMinus, nil, role.Expression, role.Operator, role.Substract),
	AnnotateType(php.AssignOpMul, nil, role.Expression, role.Operator, role.Multiply),
	AnnotateType(php.AssignOpDiv, nil, role.Expression, role.Operator, role.Divide),
	AnnotateType(php.AssignOpMod, nil, role.Expression, role.Operator, role.Modulo),

	AnnotateType(php.BitwiseAnd, nil, role.Expression, role.Binary, role.Operator, role.Bitwise, role.And),
	AnnotateType(php.BitwiseOr, nil, role.Expression, role.Binary, role.Operator, role.Bitwise, role.Or),
	AnnotateType(php.BitwiseXor, nil, role.Expression, role.Binary, role.Operator, role.Bitwise, role.Xor),
	AnnotateType(php.BitwiseNot, nil, role.Expression, role.Unary, role.Operator, role.Bitwise, role.Not),

	AnnotateType(php.BooleanAnd, nil, role.Expression, role.Binary, role.Operator, role.Boolean, role.And),
	AnnotateType(php.BooleanOr, nil, role.Expression, role.Binary, role.Operator, role.Boolean, role.Or),
	AnnotateType(php.BooleanXor, nil, role.Expression, role.Binary, role.Operator, role.Boolean, role.Xor),
	AnnotateType(php.BooleanNot, nil, role.Expression, role.Operator, role.Boolean, role.Not),

	AnnotateType(php.UnaryPlus, nil, role.Expression, role.Unary, role.Incomplete),
	AnnotateType(php.UnaryMinus, nil, role.Expression, role.Unary, role.Incomplete),
	AnnotateType(php.PreInc, nil, role.Expression, role.Unary, role.Increment),
	AnnotateType(php.PostInc, nil, role.Expression, role.Unary, role.Increment, role.Postfix),
	AnnotateType(php.PreDec, nil, role.Expression, role.Unary, role.Decrement),
	AnnotateType(php.PostDec, nil, role.Expression, role.Unary, role.Decrement, role.Postfix),

	// no join/concatenation role in UAST
	AnnotateType(php.Concat, nil, role.Expression, role.Binary, role.Operator, role.Add, role.Incomplete),

	AnnotateType(php.ShiftLeft, nil, role.Expression, role.Binary, role.Operator, role.Bitwise, role.LeftShift),
	AnnotateType(php.ShiftRight, nil, role.Expression, role.Binary, role.Operator, role.Bitwise, role.RightShift),

	AnnotateType(php.Equal, nil, role.Expression, role.Binary, role.Operator, role.Relational, role.Equal),
	AnnotateType(php.Identical, nil, role.Expression, role.Binary, role.Operator, role.Relational, role.Identical),
	AnnotateType(php.NotEqual, nil, role.Expression, role.Binary, role.Operator, role.Relational, role.Not, role.Equal),
	AnnotateType(php.NotIdentical, nil, role.Expression, role.Binary, role.Operator, role.Relational, role.Not, role.Identical),
	AnnotateType(php.GreaterOrEqual, nil, role.Expression, role.Binary, role.Operator, role.Relational, role.GreaterThanOrEqual),
	AnnotateType(php.SmallerOrEqual, nil, role.Expression, role.Binary, role.Operator, role.Relational, role.LessThanOrEqual),
	AnnotateType(php.Greater, nil, role.Expression, role.Binary, role.Operator, role.Relational, role.GreaterThan),
	AnnotateType(php.Smaller, nil, role.Expression, role.Binary, role.Operator, role.Relational, role.LessThan),
	AnnotateType(php.Spaceship, nil, role.Expression, role.Binary, role.Operator, role.Relational,
		role.GreaterThanOrEqual, role.LessThanOrEqual),

	// Scalars
	AnnotateType(php.ScalarString, nil, role.Expression, role.Literal, role.String),

	AnnotateType(php.ScalarLNumber, nil, role.Expression, role.Literal, role.Number),
	AnnotateType(php.ScalarDNumber, nil, role.Expression, role.Literal, role.Number),

	// Casts... no Cast in the UAST
	AnnotateType(php.CastArray, nil, role.Expression, role.List, role.Incomplete),
	AnnotateType(php.CastBool, nil, role.Expression, role.Boolean, role.Incomplete),
	AnnotateType(php.CastDouble, nil, role.Expression, role.Number, role.Incomplete),
	AnnotateType(php.CastInt, nil, role.Expression, role.Number, role.Incomplete),
	AnnotateType(php.CastObject, nil, role.Expression, role.Type, role.Incomplete),
	AnnotateType(php.CastString, nil, role.Expression, role.String, role.Incomplete),
	AnnotateType(php.CastUnset, nil, role.Expression, role.Incomplete),

	// TryCatch
	AnnotateType(php.TryCatch, nil, role.Statement, role.Try),
	AnnotateType(php.Catch, ObjRoles{
		"types": {role.Catch, role.Type},
	}, role.Catch, role.Type),

	AnnotateType(php.Finally, nil, role.Statement, role.Finally),

	// Class
	// FIXME: php-parser doesn't give Visibility information (public, private, etc)
	AnnotateTypeFields(php.Class, FieldRoles{
		"extends":    {Arr: true, Roles: role.Roles{role.Base}},
		"implements": {Arr: true, Roles: role.Roles{role.Implements}},
		"name":       {Rename: uast.KeyToken},
		"stmts":      {Arr: true, Roles: role.Roles{role.Type, role.Body}},
	}, role.Statement, role.Declaration, role.Type),

	// plus no const in UAST
	AnnotateType(php.ClassConst, nil, role.Type, role.Variable, role.Incomplete),

	// no member role in UAST
	AnnotateType(php.Property, nil, role.Type, role.Variable, role.Incomplete),
	AnnotateType(php.PropertyProperty, nil, role.Type, role.Variable, role.Incomplete),

	// ditto
	AnnotateType(php.ClassMethod, nil, role.Type, role.Function),

	// If + Ternary
	AnnotateType(php.Ternary, ObjRoles{
		"if":   {role.Then},
		"else": {role.Else},
		"cond": {role.If, role.Condition},
	}, role.Expression, role.If),

	AnnotateType(php.If, nil, role.Statement, role.If),
	AnnotateType(php.ElseIf, nil, role.Statement, role.If, role.Else),
	AnnotateType(php.Else, nil, role.Statement, role.Else),

	// Declare, we interpret it as an assignment-ish
	MapAST(php.Declare, Obj{
		"stmts":      Var("body_stmts"),
		"declares":   Var("declares"),
		"attributes": Var("attributes"),
	}, Obj{
		"stmts": Obj{
			uast.KeyType:  String("Declare.stmts"),
			uast.KeyRoles: Roles(role.Assignment, role.Body),
			"stmts":       Var("body_stmts"),
		},
		"declares":   Var("declares"),
		"attributes": Var("attributes"),
	}, role.Expression, role.Assignment, role.Incomplete),

	AnnotateTypeFields(php.DeclareDeclare, FieldRoles{
		"key":   {Rename: uast.KeyToken},
		"value": {Roles: role.Roles{role.Right}},
	}, role.Identifier, role.Left),

	// While and DoWhile
	AnnotateType(php.Do, nil, role.Statement, role.DoWhile),
	AnnotateType(php.While, nil, role.Statement, role.While),

	// Encapsed; incomplete: no encapsed/ string varsubst in UAST
	AnnotateType(php.Encapsed, nil, role.Expression, role.Literal, role.String, role.Incomplete),
	AnnotateType(php.EncapsedStringPart, nil, role.Expression, role.Identifier, role.Value),

	// For
	AnnotateTypeFields(php.For, FieldRoles{
		"init":  {Arr: true, Roles: role.Roles{role.Expression, role.For, role.Initialization}},
		"cond":  {Arr: true, Roles: role.Roles{role.For, role.Condition}},
		"loop":  {Arr: true, Roles: role.Roles{role.Expression, role.For, role.Update}},
		"stmts": {Arr: true, Roles: role.Roles{role.For, role.Body}},
	}, role.Statement, role.For),

	// Foreach
	AnnotateType(php.Foreach, ObjRoles{
		"valueVar": {role.Iterator},
	}, role.Statement, role.For, role.Incomplete),

	// FuncCalls, StaticCalls and MethodCalls
	AnnotateType(php.FuncCall, ObjRoles{
		"name": {role.Function, role.Name},
	}, role.Expression, role.Call),

	AnnotateType(php.StaticCall, ObjRoles{
		"class": {role.Type, role.Receiver},
	}, role.Expression, role.Call, role.Identifier),

	AnnotateType(php.MethodCall, ObjRoles{
		"var": {role.Receiver, role.Identifier},
	}, role.Expression, role.Call, role.Identifier),

	// Function declarations
	// XXX
	/*
		On(phpast.Function).Roles(uast.Function, uast.Declaration).Children(
			On(phpast.Param).Self(
				//// No reference/value in the UAST
				On(HasProperty("byRef", "true")).Roles(uast.Incomplete),
				On(HasProperty("variadic", "true")).Roles(uast.ArgsList),
			).Children(
					On(HasInternalRole("default")).Roles(uast.Default),
			),
			On(phpast.FunctionReturn).Roles(uast.Return, uast.Type),
		),
	*/

	// Include and require
	AnnotateType(php.Include, ObjRoles{
		"expr": {role.Import, role.Pathname},
	}, role.Import),

	// Instanceof
	AnnotateTypeFields(php.Instanceof, FieldRoles{
		"class":       {Roles: role.Roles{role.Call, role.Argument, role.Type, role.Identifier}},
		"expr":        {Roles: role.Roles{role.Call, role.Argument, role.Type, role.Identifier}},
		uast.KeyToken: {Add: true, Op: String("instanceof")},
	}, role.Expression, role.Call),

	// Interface
	AnnotateType(php.Interface, ObjRoles{
		"extends": {role.Receiver, role.Identifier},
	}, role.Type, role.Declaration),

	// Traits
	AnnotateType(php.Trait, nil, role.Type, role.Declaration),
	AnnotateType(php.TraitUse, nil, role.Base),
	AnnotateType(php.TraitPrecedence, ObjRoles{
		"insteadof": {role.Alias, role.Incomplete},
	}, role.Base, role.Alias, role.Incomplete),

	// New
	AnnotateType(php.New, ObjRoles{
		"class": {role.Type},
	}, role.Expression, role.Initialization, role.Call),

	// Switch
	AnnotateType(php.Switch, ObjRoles{
		"cond": {role.Switch, role.Condition},
	}, role.Statement, role.Switch),

	// XXX
	/*
		On(phpast.Case).Roles(uast.Statement, uast.Case).Self(
			On(Not(HasChild(HasInternalRole("cond")))).Roles(uast.Default),
		).Children(
			On(HasInternalRole("cond")).Roles(uast.Case),
		),
	*/
}
