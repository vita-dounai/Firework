package lexer

import (
	"testing"

	"github.com/vita-dounai/Firework/token"
)

func TestNextToken(t *testing.T) {
	input := `
	five = 5;
	add2
	add_2
	ten = 10;
	
	add = |x, y| {
		x + y;
	};
	
	result = add(five, ten);
	!-/*5
	5 < 10 > 5

	if 5 < 10 {
		return true;
	} else {
		return false;
	}

	5 == 5
	5 != 10

	"foobar"
	"foo bar"
	"foo\nbar"
	"foo\tbar"
	"foo\"bar"

	while x < 2 { break; continue; }
	[1, 2];
	3 % 2;
	{"foo": "bar"}
	macro(x, y) { x + y; };
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "add2"},
		{token.IDENTIFIER, "add_2"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.VERTICAL, "|"},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.VERTICAL, "|"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EXCLAMATION, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.INT, "5"},
		{token.NOT_EQ, "!="},
		{token.INT, "10"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.STRING, "foo\nbar"},
		{token.STRING, "foo\tbar"},
		{token.STRING, "foo\"bar"},
		{token.WHILE, "while"},
		{token.IDENTIFIER, "x"},
		{token.LT, "<"},
		{token.INT, "2"},
		{token.LBRACE, "{"},
		{token.BREAK, "break"},
		{token.SEMICOLON, ";"},
		{token.CONTINUE, "continue"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.INT, "3"},
		{token.PERCENT, "%"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.MACRO, "macro"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expectd=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestTokenPosition(t *testing.T) {
	input := `
    |test| = {
        12 + test;
    }
	`

	tests := []struct {
		expectedType   token.TokenType
		expectedLine   int
		expectedColumn int
	}{
		{token.VERTICAL, 2, 5},
		{token.IDENTIFIER, 2, 6},
		{token.VERTICAL, 2, 10},
		{token.ASSIGN, 2, 12},
		{token.LBRACE, 2, 14},
		{token.INT, 3, 9},
		{token.PLUS, 3, 12},
		{token.IDENTIFIER, 3, 14},
		{token.SEMICOLON, 3, 18},
		{token.RBRACE, 4, 5},
	}

	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expectd=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d", i, tt.expectedColumn, tok.Column)
		}
	}
}
