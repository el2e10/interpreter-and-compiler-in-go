package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node) object.Object {

	switch node := node.(type) {

	case *ast.Program:
		return eval_statement(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func eval_statement(stmts []ast.Statement) object.Object {

	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}
