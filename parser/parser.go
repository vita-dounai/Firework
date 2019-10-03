package parser

import (
	"github.com/vita-dounai/Firework/ast"
	"github.com/vita-dounai/Firework/lexer"
	"github.com/vita-dounai/Firework/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func NewParser(l *lexer.Lexer) *Parser {
	parser := &Parser{l: l}

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
