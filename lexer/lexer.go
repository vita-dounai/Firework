package lexer

import (
	"bytes"

	"github.com/vita-dounai/Firework/token"
)

type Lexer struct {
	Input        string
	Position     int
	ReadPosition int
	Ch           byte
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{Input: input, Position: 0, ReadPosition: 0, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.ReadPosition >= len(l.Input) {
		l.Ch = 0
	} else {
		l.column++
		l.Ch = l.Input[l.ReadPosition]
	}

	if l.Ch == '\n' {
		l.line++
		l.column = 0
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
	l.readChar()
	for isLetter(l.Ch) || isDigit(l.Ch) || l.Ch == '_' {
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

func (l *Lexer) readString() string {
	var str bytes.Buffer
	for {
		l.readChar()
		if l.Ch == '\\' {
			l.readChar()
			switch l.Ch {
			case 'n':
				str.WriteByte('\n')
			case 't':
				str.WriteByte('\t')
			case '"':
				str.WriteByte('"')
			default:
				str.Write([]byte{'\\', l.Ch})
			}
		} else {
			if l.Ch == '"' {
				break
			}

			str.WriteByte(l.Ch)
		}
	}
	return str.String()
}

func (l *Lexer) skipWhitespace() {
	for l.Ch == ' ' || l.Ch == '\t' || l.Ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) newToken(tokenType token.TokenType, literal string, startLine int, startColumn int) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    startLine,
		Column:  startColumn,
	}
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
		startColumn := l.column
		nextCh := l.peekChar()
		if nextCh == '=' {
			tok = l.newToken(token.EQ, token.EQ, l.line, startColumn)
			l.readChar()
		} else {
			tok = l.newToken(token.ASSIGN, string(l.Ch), l.line, startColumn)
		}
	case ';':
		tok = l.newToken(token.SEMICOLON, string(l.Ch), l.line, l.column)
	case '(':
		tok = l.newToken(token.LPAREN, string(l.Ch), l.line, l.column)
	case ')':
		tok = l.newToken(token.RPAREN, string(l.Ch), l.line, l.column)
	case ',':
		tok = l.newToken(token.COMMA, string(l.Ch), l.line, l.column)
	case '+':
		tok = l.newToken(token.PLUS, string(l.Ch), l.line, l.column)
	case '{':
		tok = l.newToken(token.LBRACE, string(l.Ch), l.line, l.column)
	case '}':
		tok = l.newToken(token.RBRACE, string(l.Ch), l.line, l.column)
	case '[':
		tok = l.newToken(token.LBRACKET, string(l.Ch), l.line, l.column)
	case ']':
		tok = l.newToken(token.RBRACKET, string(l.Ch), l.line, l.column)
	case ':':
		tok = l.newToken(token.COLON, string(l.Ch), l.line, l.column)
	case '-':
		tok = l.newToken(token.MINUS, string(l.Ch), l.line, l.column)
	case '!':
		startColumn := l.column
		nextCh := l.peekChar()
		if nextCh == '=' {
			tok = l.newToken(token.NOT_EQ, token.NOT_EQ, l.line, startColumn)
			l.readChar()
		} else {
			tok = l.newToken(token.EXCLAMATION, string(l.Ch), l.line, startColumn)
		}
	case '*':
		startColumn := l.column
		nextCh := l.peekChar()
		if nextCh == '*' {
			tok = l.newToken(token.EXP, token.EXP, l.line, startColumn)
			l.readChar()
		} else {
			tok = l.newToken(token.ASTERISK, string(l.Ch), l.line, startColumn)
		}
	case '/':
		tok = l.newToken(token.SLASH, string(l.Ch), l.line, l.column)
	case '<':
		startColumn := l.column
		nextCh := l.peekChar()
		if nextCh == '=' {
			tok = l.newToken(token.LTE, token.LTE, l.line, startColumn)
			l.readChar()
		} else {
			tok = l.newToken(token.LT, string(l.Ch), l.line, startColumn)
		}
	case '>':
		startColumn := l.column
		nextCh := l.peekChar()
		if nextCh == '=' {
			tok = l.newToken(token.GTE, token.GTE, l.line, startColumn)
			l.readChar()
		} else {
			tok = l.newToken(token.GT, string(l.Ch), l.line, startColumn)
		}
	case '|':
		tok = l.newToken(token.VERTICAL, string(l.Ch), l.line, l.column)
	case '%':
		tok = l.newToken(token.PERCENT, string(l.Ch), l.line, l.column)
	case '"':
		startLine := l.line
		startColumn := l.column
		tok = l.newToken(token.STRING, l.readString(), startLine, startColumn)
	case 0:
		tok = l.newToken(token.EOF, "", l.line, l.column)
	default:
		startLine := l.line
		startColumn := l.column
		if isLetter(l.Ch) || l.Ch == '_' {
			identifier := l.readIdentifier()
			tok = l.newToken(token.LookupIdentifier(identifier), identifier, startLine, startColumn)
			return tok
		} else if isDigit(l.Ch) {
			number := l.readNumber()
			tok = l.newToken(token.INT, number, startLine, startColumn)
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.Ch), startLine, startColumn)
		}
	}
	l.readChar()
	return tok
}
