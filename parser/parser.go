package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l             *lexer.Lexer
	current_token token.Token
	peek_token    token.Token
	errors        []string
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{l: l, errors: []string{}}
	p.next_token()
	p.next_token()
	return p
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
		return nil
	}
}

func (p *Parser) parse_return_statement() *ast.ReturnStatement {
	/*
	   Checks if the statement is a valid return statement
	   'return <expression>'
	*/

	statement := &ast.ReturnStatement{Token: p.current_token}

	for !p.current_token_is(token.SEMICOLON){
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
