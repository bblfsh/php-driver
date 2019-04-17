package normalizer

import (
	"github.com/bblfsh/sdk/v3/uast"
	"github.com/bblfsh/sdk/v3/uast/nodes"
	. "github.com/bblfsh/sdk/v3/uast/transformer"
	"github.com/bblfsh/sdk/v3/uast/transformer/positioner"
)

var Preprocess = Transformers([][]Transformer{
	{Mappings(Preprocessors...)},
}...)

var Normalize = Transformers([][]Transformer{
	{Mappings(PreNormilizers...)},
	{Mappings(Normalizers...)},
}...)

var PreprocessCode = []CodeTransformer{
	positioner.FromOffset(),
}

// Preprocessors is a block of AST preprocessing rules rules.
var Preprocessors = []Mapping{
	MapPart("root", ObjMap{ // name field as string value
		"name": Map(
			VarKind("name", nodes.KindString),

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
					uast.KeyPosOff: opAdd{op: Var("efile"), n: +1},
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

var PreNormilizers = []Mapping{
	Map(
		splitUse{vr: "uses"},
		VarKind("uses", nodes.KindArray),
	),
}

// Normalizers is the main block of normalization rules to convert native AST to semantic UAST.
var Normalizers = []Mapping{
	MapSemantic("Name", uast.Identifier{}, MapObj(
		Fields{
			{Name: uast.KeyToken, Op: Var("name")},
			{Name: "comments", Drop: true, Op: Any()}, // FIXME(dennwc): handle comments
		},
		Obj{
			"Name": Var("name"),
		},
	)),
	MapSemantic("Name", uast.Identifier{}, MapObj(
		Fields{
			{Name: "parts", Op: One(Var("name"))},
			{Name: "comments", Drop: true, Op: Any()}, // FIXME(dennwc): handle comments
		},
		Obj{
			"Name": Var("name"),
		},
	)),
	MapSemantic("Name", uast.QualifiedIdentifier{}, MapObj(
		Fields{
			{Name: "parts", Op: Each("names", Var("name"))},
			{Name: "comments", Drop: true, Op: Any()}, // FIXME(dennwc): handle comments
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
			"text": CommentText([2]string{"#", "\n"}, "text"),
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
				"name": Cases("name_case",
					Check(HasType(uast.Identifier{}), Var("name")),
					UASTTypePart("path", uast.QualifiedIdentifier{}, Obj{
						"Names": Append(Var("path_pref"), One(Var("name"))),
					}),
				),
			}),
		},
		CasesObj("name_case",
			// common
			Obj{
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
			Objs{
				{
					"Path": Is(nil),
				},
				{
					"Path": UASTTypePart("path", uast.QualifiedIdentifier{}, Obj{
						"Names": Var("path_pref"),
					}),
				},
			},
		),
	)),

	MapSemantic("Param", uast.Argument{}, MapObj(
		Obj{
			"byRef": Cases("by_ref",
				Bool(false),
				Bool(true),
			),
			"default":  Var("init"),
			"name":     Var("name"),
			"type":     typeCaseLeft("typ"),
			"variadic": Var("variadic"),
		},
		Obj{
			"Name": Var("name"),
			"Type": Cases("by_ref",
				typeCaseRight("typ"),
				Obj{
					uast.KeyType: String("ByRef"),
					"Type":       typeCaseRight("typ"),
				},
			),
			"Init":     Var("init"),
			"Variadic": Var("variadic"),
		},
	)),
	MapSemantic("Stmt_Function", uast.FunctionGroup{}, MapObj(
		Fields{
			{Name: "byRef", Op: Cases("by_ref",
				Bool(false),
				Bool(true),
			)},
			{Name: "name", Op: Var("name")},
			{Name: "params", Op: Var("params")},
			{Name: "returnType", Op: typeCaseLeft("return")},
			{Name: "stmts", Op: Var("body")},
			{Name: "comments", Drop: true, Op: Any()}, // FIXME(dennwc): handle comments
		},
		Obj{
			"Nodes": Arr(
				UASTType(uast.Alias{}, Obj{
					"Name": Var("name"),
					"Node": UASTType(uast.Function{}, Obj{
						"Type": UASTType(uast.FunctionType{}, Obj{
							"Arguments": Var("params"),
							"Returns": One(UASTType(uast.Argument{}, Obj{
								"Type": Cases("by_ref",
									// by val
									typeCaseRight("return"),
									// by ref
									Obj{
										uast.KeyType: String("ByRef"),
										"Type":       typeCaseRight("return"),
									},
								),
							})),
						}),
						"Body": UASTType(uast.Block{}, Obj{
							"Statements": Var("body"),
						}),
					}),
				}),
			),
		},
	)),
}

func typeCaseLeft(vr string) Op {
	return Cases(vr+"_case",
		Is(nil),
		VarKind(vr, nodes.KindString),
		VarKind(vr, nodes.KindObject|nodes.KindArray),
	)
}

func typeCaseRight(vr string) Op {
	return Cases(vr+"_case",
		Is(nil),
		UASTType(uast.Identifier{}, Obj{"Name": Var(vr)}),
		VarKind(vr, nodes.KindObject|nodes.KindArray),
	)
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

type splitUse struct {
	vr string
}

func (op splitUse) Kinds() nodes.Kind {
	return nodes.KindArray
}

func (op splitUse) Check(st *State, n nodes.Node) (bool, error) {
	arr, ok := n.(nodes.Array)
	if !ok {
		return false, nil
	}
	contains := false
	for _, s := range arr {
		obj, ok := s.(nodes.Object)
		if !ok && s != nil {
			return false, nil
		}
		if uast.TypeOf(obj) == "Stmt_Use" {
			uses, _ := obj["uses"].(nodes.Array)
			if len(uses) > 1 {
				contains = true
				break
			}
		}
	}
	if !contains {
		return false, nil
	}
	arr = arr.CloneList()
	for i := 0; i < len(arr); i++ {
		s := arr[i]
		obj, ok := s.(nodes.Object)
		if !ok || uast.TypeOf(obj) != "Stmt_Use" {
			continue
		}
		uses, _ := obj["uses"].(nodes.Array)
		if len(uses) < 2 {
			continue
		}
		sub := make(nodes.Array, 0, len(uses))
		for _, u := range uses {
			use := obj.CloneObject()
			use["uses"] = nodes.Array{u}
			sub = append(sub, use)
		}
		arr = append(arr[:i], append(sub, arr[i+1:]...)...)
		i += len(uses) - 1
	}
	err := st.SetVar(op.vr, arr)
	return err == nil, err
}

func (op splitUse) Construct(st *State, n nodes.Node) (nodes.Node, error) {
	// TODO: add some info to join nodes back
	return st.MustGetVar(op.vr)
}

type opAdd struct {
	op Op
	n  int
}

func (op opAdd) Kinds() nodes.Kind {
	return nodes.KindInt | nodes.KindUint | nodes.KindFloat
}

func (op opAdd) Check(st *State, n nodes.Node) (bool, error) {
	switch v := n.(type) {
	case nodes.Float:
		v -= nodes.Float(op.n)
		n = v
	case nodes.Int:
		v -= nodes.Int(op.n)
		n = v
	case nodes.Uint:
		v -= nodes.Uint(op.n)
		n = v
	default:
		return false, nil
	}
	return op.op.Check(st, n)
}

func (op opAdd) Construct(st *State, n nodes.Node) (nodes.Node, error) {
	n, err := op.op.Construct(st, n)
	if err != nil {
		return nil, err
	}
	switch v := n.(type) {
	case nodes.Float:
		v += nodes.Float(op.n)
		n = v
	case nodes.Int:
		v += nodes.Int(op.n)
		n = v
	case nodes.Uint:
		v += nodes.Uint(op.n)
		n = v
	default:
		return nil, ErrUnexpectedType.New(n)
	}
	return n, nil
}
