package lexer

import (
	"github.com/vita-dounai/Firework/token"
)

type Lexer struct {
	Input        string
	Position     int
	ReadPosition int
	Ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{Input: input, Position: 0, ReadPosition: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.ReadPosition >= len(l.Input) {
		l.Ch = 0
	} else {
		l.Ch = l.Input[l.ReadPosition]
	}
	l.Position = l.ReadPosition
	l.ReadPosition++
}

func (l *Lexer) peekChar() byte {
	if l.ReadPosition >= len(l.Input) {
		return 0
	}
	return l.Input[l.ReadPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.Position
	for isLetter(l.Ch) {
		l.readChar()
	}
	return l.Input[position:l.Position]
}

func (l *Lexer) readNumber() string {
	position := l.Position
	for isDigit(l.Ch) {
		l.readChar()
	}
	return l.Input[position:l.Position]
}

func (l *Lexer) skipWhitespace() {
	for l.Ch == ' ' || l.Ch == '\t' || l.Ch == '\n' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.Ch {
	case '=':
		nextCh := l.peekChar()
		if nextCh == '=' {
			tok = token.Token{Type: token.EQ, Literal: token.EQ}
			l.readChar()
		} else {
			tok = newToken(token.ASSIGN, l.Ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.Ch)
	case '(':
		tok = newToken(token.LPAREN, l.Ch)
	case ')':
		tok = newToken(token.RPAREN, l.Ch)
	case ',':
		tok = newToken(token.COMMA, l.Ch)
	case '+':
		tok = newToken(token.PLUS, l.Ch)
	case '{':
		tok = newToken(token.LBRACE, l.Ch)
	case '}':
		tok = newToken(token.RBRACE, l.Ch)
	case '-':
		tok = newToken(token.MINUS, l.Ch)
	case '!':
		nextCh := l.peekChar()
		if nextCh == '=' {
			tok = token.Token{Type: token.NOT_EQ, Literal: token.NOT_EQ}
			l.readChar()
		} else {
			tok = newToken(token.EXCLAMATION, l.Ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.Ch)
	case '/':
		tok = newToken(token.SLASH, l.Ch)
	case '<':
		tok = newToken(token.LT, l.Ch)
	case '>':
		tok = newToken(token.GT, l.Ch)
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.Ch) {
			identifier := l.readIdentifier()
			tok = token.Token{Type: token.LookupIdentifier(identifier), Literal: identifier}
			return tok
		} else if isDigit(l.Ch) {
			number := l.readNumber()
			tok = token.Token{Type: token.INT, Literal: number}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.Ch)
		}
	}
	l.readChar()
	return tok
}
