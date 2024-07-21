package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatemnts(t *testing.T) {
	input := ` 
	let x = 5;
   	let y = 10;
   	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	check_parser_errors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program statements doesnt include all the three statements got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

	}

}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {

	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let' got %q", s.TokenLiteral())
		return false
	}

	let_stmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement. got=%T", s)
		return false
	}

	if let_stmt.Name.Value != name {
		t.Errorf("let_stmt.Name.Value not %s. got %s", name, let_stmt.Name.Value)
		return false
	}
	if let_stmt.Name.TokenLiteral() != name {
		t.Errorf("let_stmt.Name.TokenLiteral() not '%s'. got=%s",
			name, let_stmt.Name.TokenLiteral())
		return false
	}

	return true

}

func check_parser_errors(t *testing.T, p *Parser) {

	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()

}

func TestReturnStatement(t *testing.T) {

	input := `
	return 5;
	return 10;
	return 993322;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	check_parser_errors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statement doesn't contain 3 statements got %d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		// Type assertion, basically checking if the statement is off correct type
		return_statement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement is not a return statement")
		}

		if return_statement.TokenLiteral() != "return" {

			t.Errorf("return_statement.TokenLiter() didn't return 'return' got %s instead", return_statement.TokenLiteral())

		}

	}
}

func TestIdentifierExprssion(t *testing.T) {

	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	check_parser_errors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has wrong number of statements. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statement[0] is not an expression it's an %T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expression not an IntegerLiteral got %T", statement.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("identifier not %s, got=%d", "foobar", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not %s, got=%s", "5", literal.TokenLiteral())

	}

}

func TestParsingPrefixExpressions(t *testing.T) {

	prefix_tests := []struct {
		input         string
		operator      string
		integer_value int64
	}{{"!5", "!", 5}, {"!5", "!", 5}}

	for _, tt := range prefix_tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		check_parser_errors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program contained incorrect number of statements expected was %d got=%d", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Wrong type of statement was created got=%T", program.Statements[0])
		}

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("statement was not a Prefix expression got=%T", statement.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Exp.Operator is not '%s'  got=%s", tt.operator, exp.Operator)
		}

		if !test_integer_literal(t, exp.Right, tt.integer_value) {
			return
		}
	}

}

func test_integer_literal(t *testing.T, il ast.Expression, value int64) bool {

	integer, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integer.TokenLiteral())
		return false
	}
	return true
}
