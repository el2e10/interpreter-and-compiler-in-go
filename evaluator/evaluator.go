package evaluator

import (
	"fmt"

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
		return eval_program(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return native_bool_to_boolean_object(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return eval_prefix_expression(node.Operator, right)

	case *ast.InfixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		return eval_infix_expression(node.Operator, left, right)

	case *ast.BlockStatement:
		return eval_block_statement(node)

	case *ast.IfExpression:
		return eval_if_expression(node)

	case *ast.ReturnStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}

	return nil
}

func eval_if_expression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isError(condition) {
		return condition
	}

	if is_truthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func is_truthy(obj object.Object) bool {
	switch obj {

	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true

	}
}

func eval_infix_expression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return eval_integer_infix_expression(operator, left, right)
	case operator == "==":
		return native_bool_to_boolean_object(left == right)
	case operator == "!=":
		return native_bool_to_boolean_object(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func eval_minus_operator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
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

func eval_program(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func eval_block_statement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
