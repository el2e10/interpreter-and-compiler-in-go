package evaluator

import (
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":   object.GetBuildinByName("len"),
	"first": object.GetBuildinByName("first"),
	"last":  object.GetBuildinByName("last"),
	"push":  object.GetBuildinByName("push"),
	"rest":  object.GetBuildinByName("rest"),
	"puts":   object.GetBuildinByName("put"),
}
