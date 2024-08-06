package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return eval_statement(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Value {
			return TRUE
		} else {
			return FALSE
		}

	case *ast.PrefixExpression:
		{
			right := Eval(node.Right)
			return eval_prefix_expression(node.Operator, right)
		}
	}

	return nil
}

func eval_prefix_expression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		{
			return eval_bang_operator(right)
		}
	case "-":
		{
			return eval_minus_operator(right)
		}
	default:
		return NULL
	}
}

func eval_minus_operator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func eval_bang_operator(obj object.Object) object.Object {
	switch obj {

	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return FALSE
	default:
		return FALSE

	}
}

func eval_statement(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}
