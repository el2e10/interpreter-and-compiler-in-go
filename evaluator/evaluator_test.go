package evaluator

import (
	"testing"

	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input  string
		output int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, test := range tests {
		test_integer_object(t, test_eval(test.input), test.output)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			` if (10 > 1) {
  		if (10 > 1) {
    		return true + false;
		}
		return 1; }`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"foobar", "identifier not found: foobar"},
	}

	for index, tt := range tests {

		evaluated := test_eval(tt.input)
		error_obj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("%d No error object returned. got=%T (%+v)", index, evaluated, evaluated)
			continue
		}

		if error_obj.Message != tt.expectedMessage {
			t.Errorf("Wrong error message was returned expected=%q, got=%q", tt.expectedMessage, error_obj.Message)
		}

	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := test_eval(tt.input)
		test_integer_object(t, evaluated, tt.expected)

	}
}

func TestIfElseExpression(t *testing.T) {
	test := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range test {
		evaluated := test_eval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			test_integer_object(t, evaluated, int64(integer))
		} else {
			test_null_object(t, evaluated)
		}
	}
}

func test_null_object(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	} else {
		return true
	}
}

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

	for _, tt := range tests {
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
	env := object.NewEnvironment()

	return Eval(prgm, env)
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
