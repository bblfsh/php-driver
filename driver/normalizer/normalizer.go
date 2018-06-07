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
	Map( // name field as string value
		Part("_", Obj{
			"name": Check(isString{}, Var("name")),
		}),
		Part("_", Obj{
			"name": Obj{
				uast.KeyType:  String("Name"),
				uast.KeyToken: Var("name"),
			},
		}),
	),

	//The native AST puts positions and comments inside an "attribute" node. Here
	//we reparent them to the current node.
	Map(
		Part("root", Obj{
			"attributes": Part("attrs", Fields{
				// Ignore those because they're wrong in the native AST; we instead
				// compute line and col from the offset

				//{Name: "startLine", Op: Var("sline")},
				//{Name: "endLine", Op: Var("eline")},
				//{Name: "startTokenPos", Op: Var("stoken")},
				//{Name: "endTokenPos", Op: Var("etoken")},
				{Name: "startFilePos", Op: Var("sfile")},
				{Name: "endFilePos", Op: Var("efile")},
				{Name: "comments", Op: Var("comments"), Optional: "comments_exists"},
			}),
		}),

		Part("root", Fields{
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
			{Name: "comments", Op: Var("comments"), Optional: "comments_exists"},
		}),
	),

	ObjectToNode{
		InternalTypeKey: "nodeType",
	}.Mapping(),
}

// Normalizers is the main block of normalization rules to convert native AST to semantic UAST.
var Normalizers = []Mapping{}
