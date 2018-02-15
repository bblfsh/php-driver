package normalizer

import (
	"strings"

	"gopkg.in/bblfsh/sdk.v1/uast"
)

// ToNode is an instance of `uast.ObjectToNode`, defining how to transform an
// into a UAST (`uast.Node`).
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast#ObjectToNode
var ToNode = &uast.ObjectToNode{
	InternalTypeKey: "nodeType",
	LineKey:         "attributes.startLine",
	ColumnKey:       "attributes.startTokenPos",
	OffsetKey:       "attributes.startFilePos",
	EndLineKey:      "attributes.endLine",
	EndColumnKey:    "attributes.endTokenPos",
	SyntheticTokens: map[string]string{
		"Expr_Clone":      "clone",
		"Expr_Empty":      "empty",
		"Expr_Isset":      "isset",
		"Stmt_Echo":       "echo",
		"Stmt_Print":      "print",
		"Stmt_Unset":      "unset",
		"Expr_Eval":       "eval",
		"Expr_Exit":       "exit",
		"Expr_Instanceof": "instanceof",
		"Expr_List":       "list",
		"Expr_New":        "new",
	},
	EndOffsetKey: "attributes.endFilePos",
	TokenKeys: map[string]bool{
		"name":    true,
		"text":    true, // for comments
		"value":   true, // Scalars
		"var":     true, // catch list
		"key":     true, // declare
		"newName": true, // trait alias
	},
	PromotedPropertyStrings: map[string]map[string]bool{
		"Stmt_Function": {"returnType": true},
	},
	// PHP AST includes a map called attributes with several properties, should
	// be ignored, otherwise fake nodes are created.
	IsNode: func(v map[string]interface{}) bool {
		_, ok := v["nodeType"]
		return ok
	},
	// The parser returns multiple nodes instead of a single node, a fake node
	// root node with the type "File" is created.
	OnToNode: func(n interface{}) (interface{}, error) {
		return map[string]interface{}{
			"root": map[string]interface{}{
				"nodeType": "File",
				"children": n,
			},
		}, nil
	},
	Modifier: func(n map[string]interface{}) error {
		// Sometimes, if the name includes namespaces, it's given as an array in
		// several parts. The parts are imploded into the name key.
		if parts, ok := n["parts"].([]interface{}); ok {
			deleteParts := false
			n["name"], deleteParts = sliceInterfaceToString(parts, "\\")
			if deleteParts {
				delete(n, "parts")
			}
		}

		// Positions in comments don't follow the same schema as the other
		// nodes, the position info is moved to the same place.
		if pos, ok := n["filePos"].(float64); ok {
			n["attributes.startFilePos"] = pos
			n["attributes.endFilePos"] = pos + float64(len(n["text"].(string)))
			n["attributes.startLine"] = n["line"]
			delete(n, "filePos")
			delete(n, "line")
		}

		return nil
	},
}

func sliceInterfaceToString(s []interface{}, sep string) (string, bool) {
	deleteParts := false
	l := make([]string, len(s))
	for i, v := range s {
		if part, ok := v.(string); ok {
			l[i] = part
			deleteParts = true
		}
	}

	return strings.Join(l, sep), deleteParts
}
