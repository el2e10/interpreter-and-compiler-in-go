package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func (p *Parser) peek_precedence() int {

	if p, ok := precedences[p.peek_token.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) cur_precedence() int {

	if p, ok := precedences[p.current_token.Type]; ok {
		return p
	}
	return LOWEST
}

type (
	prefix_parse_fn func() ast.Expression
	infix_parse_fn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l             *lexer.Lexer
	current_token token.Token
	peek_token    token.Token
	errors        []string

	prefix_parse_fns map[token.TokenType]prefix_parse_fn
	infix_parse_fns  map[token.TokenType]infix_parse_fn
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{l: l, errors: []string{}}
	p.next_token()
	p.next_token()

	p.prefix_parse_fns = make(map[token.TokenType]prefix_parse_fn)
	p.register_prefix(token.IDENT, p.parse_identifier)
	p.register_prefix(token.INT, p.parse_integer_literal)
	p.register_prefix(token.BANG, p.parse_prefix_expression)
	p.register_prefix(token.MINUS, p.parse_prefix_expression)

	p.infix_parse_fns = make(map[token.TokenType]infix_parse_fn)
	p.register_infix(token.PLUS, p.parse_infix_expression)
	p.register_infix(token.MINUS, p.parse_infix_expression)
	p.register_infix(token.SLASH, p.parse_infix_expression)
	p.register_infix(token.ASTERISK, p.parse_infix_expression)
	p.register_infix(token.EQ, p.parse_infix_expression)
	p.register_infix(token.NOT_EQ, p.parse_infix_expression)
	p.register_infix(token.LT, p.parse_infix_expression)
	p.register_infix(token.GT, p.parse_infix_expression)

	return p
}

func (p *Parser) parse_infix_expression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.current_token,
		Operator: p.current_token.Literal,
		Left:     left,
	}

	precedence := p.cur_precedence()
	p.next_token()
	expression.Right = p.parse_expression(precedence)

	return expression
}

func (p *Parser) parse_identifier() ast.Expression {
	return &ast.Identifier{Token: p.current_token, Value: p.current_token.Literal}
}

func (p *Parser) next_token() {

	p.current_token = p.peek_token
	p.peek_token = p.l.NextToken()

}

func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.current_token_is(token.EOF) {

		statement := p.parse_statement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.next_token()
	}

	return program
}

func (p *Parser) parse_statement() ast.Statement {

	switch p.current_token.Type {
	case token.LET:
		return p.parse_let_statement()
	case token.RETURN:
		return p.parse_return_statement()
	default:
		return p.parse_expression_statement()
	}
}

func (p *Parser) parse_expression_statement() *ast.ExpressionStatement {

	statement := &ast.ExpressionStatement{Token: p.current_token}
	statement.Expression = p.parse_expression(LOWEST)

	if p.peek_token_is(token.SEMICOLON) {
		p.next_token()
	}

	return statement
}

func (p *Parser) parse_integer_literal() ast.Expression {

	literal := &ast.IntegerLiteral{Token: p.current_token}

	value, error := strconv.ParseInt(p.current_token.Literal, 0, 64)
	if error != nil {
		msg := fmt.Sprintf("Could not parser %q as integer", p.current_token.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value
	return literal

}

func (p *Parser) parse_expression(precedence int) ast.Expression {

	prefix := p.prefix_parse_fns[p.current_token.Type]
	if prefix == nil {
		p.no_prefix_parse_fn_error(p.current_token.Type)
		return nil
	}

	left_exp := prefix()

	for !p.peek_token_is(token.SEMICOLON) && precedence < p.peek_precedence() {

		infix := p.infix_parse_fns[p.peek_token.Type]

		if infix == nil {
			return left_exp
		}
		p.next_token()
		left_exp = infix(left_exp)
	}

	return left_exp

}

func (p *Parser) parse_prefix_expression() ast.Expression {

	expression := &ast.PrefixExpression{
		Token:    p.current_token,
		Operator: p.current_token.Literal,
	}

	p.next_token()
	expression.Right = p.parse_expression(PREFIX)

	return expression

}

func (p *Parser) no_prefix_parse_fn_error(t token.TokenType) {

	msg := fmt.Sprintf("No prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)

}

func (p *Parser) parse_return_statement() *ast.ReturnStatement {
	/*
	   Checks if the statement is a valid return statement
	   'return <expression>'
	*/

	statement := &ast.ReturnStatement{Token: p.current_token}

	for !p.current_token_is(token.SEMICOLON) {
		p.next_token()
	}

	return statement

}

func (p *Parser) parse_let_statement() *ast.LetStatement {

	/*
		Checks if the statement is of the for 'let x = 5;'
	*/
	statement := &ast.LetStatement{Token: p.current_token}

	if !p.expect_peek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.current_token, Value: p.current_token.Literal}

	if !p.expect_peek(token.ASSIGN) {
		return nil
	}

	for !p.current_token_is(token.SEMICOLON) {
		p.next_token()
	}

	return statement

}

func (p *Parser) current_token_is(t token.TokenType) bool {
	return p.current_token.Type == t
}

func (p *Parser) peek_token_is(t token.TokenType) bool {
	return p.peek_token.Type == t
}

func (p *Parser) expect_peek(tkn token.TokenType) bool {
	if p.peek_token_is(tkn) {
		p.next_token()
		return true
	} else {
		p.peek_error(tkn)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors

}

func (p *Parser) peek_error(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s but got %s instead", t, p.peek_token.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) register_prefix(tokenType token.TokenType, fn prefix_parse_fn) {
	p.prefix_parse_fns[tokenType] = fn
}

func (p *Parser) register_infix(tokenType token.TokenType, fn infix_parse_fn) {
	p.infix_parse_fns[tokenType] = fn
}
