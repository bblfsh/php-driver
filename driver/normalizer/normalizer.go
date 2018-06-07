package normalizer

import (
	"gopkg.in/bblfsh/sdk.v2/uast"
	. "gopkg.in/bblfsh/sdk.v2/uast/transformer"
)

var Preprocess = Transformers([][]Transformer{
	{Mappings(Preprocessors...)},
}...)

var Normalize = Transformers([][]Transformer{
	{Mappings(Normalizers...)},
}...)

// Preprocessors is a block of AST preprocessing rules rules.
var Preprocessors = []Mapping{
	MapPart("root", ObjMap{ // name field as string value
		"name": Map(
			Check(isString{}, Var("name")),

			Obj{
				uast.KeyType:  String("Name"),
				uast.KeyToken: Var("name"),
				uast.KeyPos:   UASTType(uast.Positions{}, Obj{}),
			},
		),
	}),

	//The native AST puts positions and comments inside an "attribute" node. Here
	//we reparent them to the current node.
	MapPart("root", MapObj(
		Obj{
			"attributes": Part("attrs", Fields{
				// Ignore those because they're wrong in the native AST; we instead
				// compute line and col from the offset

				{Name: "startLine", Op: AnyVal(nil)},     // sline
				{Name: "endLine", Op: AnyVal(nil)},       // eline
				{Name: "startTokenPos", Op: AnyVal(nil)}, // stoken
				{Name: "endTokenPos", Op: AnyVal(nil)},   // etoken
				{Name: "startFilePos", Op: Var("sfile")},
				{Name: "endFilePos", Op: Var("efile")},
				{Name: "comments", Op: Var("comments"), Optional: "comments_exists"},
			}),
		},

		Fields{
			{Name: uast.KeyPos, Op: UASTType(uast.Positions{}, Obj{
				uast.KeyStart: UASTType(uast.Position{}, Obj{
					//uast.KeyPosLine: Var("sline"),
					//uast.KeyPosCol:  Var("stoken"),
					uast.KeyPosOff: Var("sfile"),
				}),
				uast.KeyEnd: UASTType(uast.Position{}, Obj{
					//uast.KeyPosLine: Var("eline"),
					//uast.KeyPosCol:  Var("etoken"),
					uast.KeyPosOff: Var("efile"),
				}),
			})},
			{Name: "attributes", Op: Var("attrs")},
			{Name: "comments", Op: Var("comments"), Optional: "comments_exists"},
		},
	)),

	MapPart("root", MapObj(Obj{
		"filePos": Var("fp"),
		"line":    AnyVal(nil),
	}, Obj{
		uast.KeyPos: UASTType(uast.Positions{}, Obj{
			uast.KeyStart: UASTType(uast.Position{}, Obj{
				uast.KeyPosOff: Var("fp"),
			}),
		}),
	})),

	Map( // trim attributes if it's empty
		Part("_", Obj{
			"attributes": EmptyObj(),
		}),
		Part("_", Obj{}),
	),

	ObjectToNode{
		InternalTypeKey: "nodeType",
	}.Mapping(),
}

// Normalizers is the main block of normalization rules to convert native AST to semantic UAST.
var Normalizers = []Mapping{
	MapSemantic("Name", uast.Identifier{}, MapObj(
		Obj{
			uast.KeyToken: Var("name"),
		},
		Obj{
			"Name": Var("name"),
		},
	)),
	MapSemantic("Name", uast.Identifier{}, MapObj(
		Obj{
			"parts": One(Var("name")),
		},
		Obj{
			"Name": Var("name"),
		},
	)),
	MapSemantic("Name", uast.QualifiedIdentifier{}, MapObj(
		Obj{
			"parts": Each("names", Var("name")),
		},
		Obj{
			"Names": Each("names", UASTType(uast.Identifier{}, Obj{
				"Name": Var("name"),
			})),
		},
	)),

	MapSemantic("Scalar_String", uast.String{}, MapObj(
		Obj{
			"value": Var("val"),
			"attributes": Obj{"kind": Cases("kind",
				Int(1), // raw string
				Int(2), // escaped string
			)},
		},
		Obj{
			"Value": Var("val"),
			"Format": Cases("kind",
				String("raw"),
				String(""),
			),
		},
	)),

	MapSemantic("Scalar_String", uast.String{}, MapObj(
		Obj{
			"value": Var("val"),
			"attributes": Obj{
				"kind":     Int(3),
				"docLabel": AnyVal(nil), // TODO: store it
			},
		},
		Obj{
			"Value":  Var("val"),
			"Format": String("raw_custom"),
		},
	)),
	MapSemantic("Scalar_String", uast.String{}, MapObj(
		Obj{
			"value": Var("val"),
		},
		Obj{
			"Value":  Var("val"),
			"Format": String(""),
		},
	)),
	MapSemantic("Scalar_EncapsedStringPart", uast.String{}, MapObj(
		Obj{
			"value": Var("val"),
		},
		Obj{
			"Value":  Var("val"),
			"Format": String("encapsed"),
		},
	)),
	MapSemantic("Comment", uast.Comment{}, MapObj(
		Obj{
			"text": CommentText([2]string{"//", "\n"}, "text"),
		},
		CommentNode(false, "text", nil),
	)),
	MapSemantic("Comment", uast.Comment{}, MapObj(
		Obj{
			"text": CommentText([2]string{"/*", "*/"}, "text"),
		},
		CommentNode(true, "text", nil),
	)),
	MapSemantic("Comment_Doc", uast.Comment{}, MapObj(
		Obj{ // FIXME: doc should write additional flag
			"text": CommentText([2]string{"/**", "*/"}, "text"),
		},
		CommentNode(true, "text", nil),
	)),

	convertBlock("Stmt_If", ""),
	convertBlock("Stmt_ElseIf", ""),
	convertBlock("Stmt_Else", ""),
	convertBlock("Stmt_For", ""),
	convertBlock("Stmt_Foreach", ""),
	convertBlock("Stmt_While", ""),
	convertBlock("Stmt_Class", ""),
	convertBlock("Stmt_Do", ""),
	convertBlock("Stmt_ClassMethod", ""),
	convertBlock("Expr_Closure", ""),
	convertBlock("Stmt_Namespace", ""),
	convertBlock("Stmt_Interface", ""),
	convertBlock("Stmt_Catch", ""),
	convertBlock("Stmt_Finally", ""),
	convertBlock("Stmt_TryCatch", ""),
	convertBlock("Stmt_Case", ""),
	convertBlock("Stmt_Switch", ""),
	convertBlock("Stmt_Trait", ""),
	convertBlock("Stmt_Declare", "body_stmts"),

	MapSemantic("Expr_Include", uast.RuntimeReImport{}, MapObj(
		Obj{
			"expr": Var("path"),
			"type": Cases("typ",
				Int(1), // include
				Int(3), // require
			),
		},
		Obj{
			"Path": Var("path"),
			"All": Cases("typ",
				Bool(true), // include
				Bool(true), // require
			),
		},
	)),

	MapSemantic("Expr_Include", uast.RuntimeImport{}, MapObj(
		Obj{
			"expr": Var("path"),
			"type": Cases("typ",
				Int(2), // include_once
				Int(4), // require_once
			),
		},
		Obj{
			"Path": Var("path"),
			"All": Cases("typ",
				Bool(true), // include
				Bool(true), // require
			),
		},
	)),

	MapSemantic("Stmt_GroupUse", uast.RuntimeImport{}, MapObj(
		Obj{
			"prefix": Var("path"),
			"type":   Int(0),
			"uses": Each("names", Obj{
				uast.KeyType: String("Stmt_UseUse"),
				uast.KeyPos:  Var("name_pos"),
				"type":       Int(1),
				"alias":      Var("alias"),
				"name":       Var("name"),
			}),
		},
		Obj{
			"Path": Var("path"),
			"All":  Bool(false),
			"Names": Each("names", UASTType(uast.Alias{}, Obj{
				uast.KeyPos: Var("name_pos"),
				"Name": UASTType(uast.Identifier{}, Obj{
					"Name": Var("alias"),
				}),
				"Node": Var("name"),
			})),
		},
	)),
	MapSemantic("Stmt_Use", uast.RuntimeImport{}, MapObj(
		Obj{
			"type": Cases("typ",
				// TODO: do we care?
				Int(1), // use
				Int(2), // use function
				Int(3), // use const
			),
			"uses": One(Obj{
				uast.KeyType: String("Stmt_UseUse"),
				uast.KeyPos:  Var("name_pos"),
				"type":       Int(0),
				"alias":      Var("alias"),
				"name": UASTTypePart("path", uast.QualifiedIdentifier{}, Obj{
					"Names": Append(Var("path_pref"), One(Var("name"))),
				}),
			}),
		},
		Obj{
			"Path": UASTTypePart("path", uast.QualifiedIdentifier{}, Obj{
				"Names": Var("path_pref"),
			}),
			"All": Cases("typ",
				Bool(false), // use
				Bool(false), // use function
				Bool(false), // use const
			),
			"Names": One(UASTType(uast.Alias{}, Obj{
				uast.KeyPos: Var("name_pos"),
				"Name": UASTType(uast.Identifier{}, Obj{
					"Name": Var("alias"),
				}),
				"Node": Var("name"),
			})),
		},
	)),
}

func convertBlock(typ, field string) Mapping {
	if field == "" {
		field = "stmts"
	}
	return MapPart("root", ObjMap{
		uast.KeyType: String(typ),
		"stmts": Map(
			Var("body"),

			UASTType(uast.Block{}, Obj{
				"Statements": Var("body"),
			}),
		),
	})
}
