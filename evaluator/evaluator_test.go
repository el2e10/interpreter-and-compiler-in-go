package evaluator

import (
	"fmt"
	"testing"

	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!!true", true},
		{"!5", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		test_boolean_object(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true == true", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == false", true},
	}

	for index, tt := range tests {
		fmt.Println(index, tt.input)
		evaluated := test_eval(tt.input)
		test_boolean_object(t, evaluated, tt.expected)
	}
}

func test_boolean_object(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value got=%t expected=%t", result.Value, obj)
		return false
	}

	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		test_integer_object(t, evaluated, tt.expected)
	}
}

func test_eval(input string) object.Object {
	l := lexer.New(input)
	prsr := parser.New(l)
	prgm := prsr.ParseProgram()

	return Eval(prgm)
}

func test_integer_object(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object is not an Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Object has the wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}
