package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input         string
	position      int  // current position or the index of 'current_char'
	read_position int  // next position to read
	current_char  byte // current character that is getting analyzed
}

func New(code string) *Lexer {
	lexer := &Lexer{input: code}
	lexer.read_char()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var tkn token.Token

	lexer.skip_whitespace()
	lexer.skipComment()

	switch lexer.current_char {
	case '=':
		if lexer.peek_next_char() == '=' {
			tkn.Type = token.EQ
			tkn.Literal = "=="
			lexer.read_char()
		} else {
			tkn = new_token(token.ASSIGN, lexer.current_char)
		}
	case '+':
		tkn = new_token(token.PLUS, lexer.current_char)
	case '-':
		tkn = new_token(token.MINUS, lexer.current_char)
	case '(':
		tkn = new_token(token.LPAREN, lexer.current_char)
	case ')':
		tkn = new_token(token.RPAREN, lexer.current_char)
	case '{':
		tkn = new_token(token.LBRACE, lexer.current_char)
	case '}':
		tkn = new_token(token.RBRACE, lexer.current_char)
	case ',':
		tkn = new_token(token.COMMA, lexer.current_char)
	case ';':
		tkn = new_token(token.SEMICOLON, lexer.current_char)
	case ':':
		tkn = new_token(token.COLON, lexer.current_char)
	case '!':
		if lexer.peek_next_char() == '=' {
			tkn.Type = token.NOT_EQ
			tkn.Literal = "!="
			lexer.read_char()
		} else {
			tkn = new_token(token.BANG, lexer.current_char)
		}
	case '*':
		tkn = new_token(token.ASTERISK, lexer.current_char)
	case '/':
		tkn = new_token(token.SLASH, lexer.current_char)
	case '>':
		tkn = new_token(token.GT, lexer.current_char)
	case '<':
		tkn = new_token(token.LT, lexer.current_char)
	case '"':
		s := lexer.readString()
		tkn = token.Token{Type: token.STRING, Literal: s}
	case 0:
		tkn.Literal = ""
		tkn.Type = token.EOF
	case '[':
		tkn = new_token(token.LBRACKET, lexer.current_char)
	case ']':
		tkn = new_token(token.RBRACKET, lexer.current_char)
	default:
		if is_letter(lexer.current_char) {
			tkn.Literal = lexer.read_identifier()
			tkn.Type = token.LookupIdentifier(tkn.Literal)
			return tkn
		} else if is_digit(lexer.current_char) {
			tkn.Type = token.INT
			tkn.Literal = lexer.read_digit()
			return tkn
		} else {
			tkn = new_token(token.ILLEGAL, lexer.current_char)
		}
	}

	lexer.read_char()
	return tkn
}

func (lexer *Lexer) skipComment() {
	if lexer.current_char != '#' {
		return
	}
	for lexer.current_char != 0 {
		lexer.read_char()
	}
}

func (lexer *Lexer) readString() string {
	start_position := lexer.read_position
	for {
		lexer.read_char()
		if lexer.current_char == '"' || lexer.current_char == 0 {
			break
		}
	}
	return lexer.input[start_position:lexer.position]
}

func (lexer *Lexer) skip_whitespace() {
	for lexer.current_char == ' ' || lexer.current_char == '\t' || lexer.current_char == '\n' || lexer.current_char == '\r' {
		lexer.read_char()
	}
}

func (lexer *Lexer) read_identifier() string {
	start_position := lexer.position
	for is_letter(lexer.current_char) {
		lexer.read_char()
	}
	return lexer.input[start_position:lexer.position]
}

func (lexer *Lexer) read_digit() string {
	start_position := lexer.position
	for is_digit(lexer.current_char) {
		lexer.read_char()
	}
	return lexer.input[start_position:lexer.position]
}

func is_letter(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func is_digit(char byte) bool {
	return char >= '0' && char <= '9'
}

func new_token(token_type token.TokenType, char byte) token.Token {
	return token.Token{Type: token_type, Literal: string(char)}
}

func (lexer *Lexer) read_char() {
	if lexer.read_position >= len(lexer.input) {
		lexer.current_char = 0
	} else {
		lexer.current_char = lexer.input[lexer.read_position]
	}
	lexer.position = lexer.read_position
	lexer.read_position += 1
}

func (lexer *Lexer) peek_next_char() byte {
	if lexer.read_position >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.read_position]
	}
}
