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
		return native_bool_to_boolean_object(node.Value)

	case *ast.PrefixExpression:
		{
			right := Eval(node.Right)
			return eval_prefix_expression(node.Operator, right)
		}
	case *ast.InfixExpression:
		{
			right := Eval(node.Right)
			left := Eval(node.Left)
			return eval_infix_expression(node.Operator, left, right)
		}
	}

	return nil
}

func eval_infix_expression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return eval_integer_infix_expression(operator, left, right)
	case operator == "==":
		return native_bool_to_boolean_object(left == right)
	case operator == "!=":
		return native_bool_to_boolean_object(left != right)
	default:
		return NULL
	}
}

func eval_integer_infix_expression(operator string, left object.Object, right object.Object) object.Object {
	left_value := left.(*object.Integer).Value
	right_value := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: left_value + right_value}

	case "-":
		return &object.Integer{Value: left_value - right_value}

	case "*":
		return &object.Integer{Value: left_value * right_value}

	case "/":
		return &object.Integer{Value: left_value / right_value}

	case "<":
		return native_bool_to_boolean_object(left_value < right_value)

	case ">":
		return native_bool_to_boolean_object(left_value > right_value)

	case "==":
		return native_bool_to_boolean_object(left_value == right_value)

	case "!=":
		return native_bool_to_boolean_object(left_value != right_value)

	default:
		return NULL
	}
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

func native_bool_to_boolean_object(value bool) object.Object {
	if value {
		return TRUE
	} else {
		return FALSE
	}
}
