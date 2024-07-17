package parser

import (
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

		if return_statement.TokenLiteral() != "return"{

			t.Errorf("return_statement.TokenLiter() didn't return 'return' got %s instead", return_statement.TokenLiteral())

		}

	}
}
