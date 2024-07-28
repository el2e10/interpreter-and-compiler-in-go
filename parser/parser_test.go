package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	check_parser_errors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}
	if !test_identifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	test_literal_expression(t, exp.Arguments[0], 1)
	test_infix_expression(t, exp.Arguments[1], 2, "*", 3)
	test_infix_expression(t, exp.Arguments[2], 4, "+", 5)

}

func TestFunctionParameterParsing(t *testing.T) {

	test := []struct {
		input           string
		expected_params []string
	}{
		{input: "fn() {};", expected_params: []string{}},
		{input: "fn(x) {};", expected_params: []string{"x"}},
		{input: "fn(x, y) {};", expected_params: []string{"x", "y"}},
	}

	for _, tt := range test {

		l := lexer.New(tt.input)
		p := New(l)

		programs := p.ParseProgram()
		check_parser_errors(t, p)

		stmt := programs.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expected_params) {

			t.Errorf("Length of parameters are wrong expected %d, got=%d", len(tt.expected_params), len(function.Parameters))
		}

		for i, ident := range tt.expected_params {
			test_literal_expression(t, function.Parameters[i], ident)
		}

	}

}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) {x + y}`

	lxr := lexer.New(input)
	prsr := New(lxr)

	program := prsr.ParseProgram()
	check_parser_errors(t, prsr)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}
	test_literal_expression(t, function.Parameters[0], "x")
	test_literal_expression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}
	body_stmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}
	test_infix_expression(t, body_stmt.Expression, "x", "+", "y")

}

func TestIfStatement(t *testing.T) {

	input := `if (x<y) { x } else { y }`

	l := lexer.New(input)
	prsr := New(l)

	program := prsr.ParseProgram()
	check_parser_errors(t, prsr)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !test_infix_expression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !test_identifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if exp.Alternative == nil {
		t.Errorf("exp.Alternative.Statements was nil. expect %s", "else")
	}

	if !test_identifier(t, alternative.Expression, "y") {
		return
	}

}

func test_identifier(t *testing.T, exp ast.Expression, value string) bool {

	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s, got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s, got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func test_literal_expression(t *testing.T,
	exp ast.Expression,
	expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return test_integer_literal(t, exp, int64(v))
	case int64:
		return test_integer_literal(t, exp, v)
	case string:
		return test_identifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func test_infix_expression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !test_literal_expression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !test_literal_expression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestBooleanExpression(t *testing.T) {

	input := "let foobar = true;"
	lxr := lexer.New(input)
	psr := New(lxr)

	program := psr.ParseProgram()
	check_parser_errors(t, psr)

	if len(program.Statements) != 1 {
		t.Fatalf("Wrong number of statements returned expected %d, but got=%d", 1, len(program.Statements))
	}

	//statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	statement, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("Wrong program.Statement[0] type, got=%T", program.Statements[0])
	}

	expression, ok := statement.ReturnValue.(*ast.Boolean)
	if !ok {
		t.Fatalf("Wrong Expression type expected Boolean got=%T %q", statement.ReturnValue, statement)
	}

	if expression.Value != true {
		t.Fatalf("Wrong value expected %t, but got=%t", true, expression.Value)
	}

}

func TestParsingInfixExpressions(t *testing.T) {
	infix_tests := []struct {
		input       string
		left_value  int64
		operator    string
		right_value int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infix_tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		check_parser_errors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if !test_integer_literal(t, exp.Left, tt.left_value) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !test_integer_literal(t, exp.Right, tt.right_value) {
			return
		}
	}

}

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

		fmt.Print(program.Statements[0])
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
