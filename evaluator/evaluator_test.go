package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"25", 25},
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
