package parser

import (
	"strconv"

	"github.com/vita-dounai/Firework/ast"
	"github.com/vita-dounai/Firework/lexer"
	"github.com/vita-dounai/Firework/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	EXP         // **
	PREFIX      // - or !
	CALL        // funcion call
	INDEX       // array[index]
)

var (
	BREAK_STATEMENT    = &ast.BreakStatement{}
	CONTINUE_STATEMENT = &ast.ContinueStatement{}

	UNEXPECTED_EOF   = &UnexpectedEOF{}
	ILLEGAL_BREAK    = &IllegalBreak{}
	ILLEGAL_CONTINUE = &IllegalContinue{}
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.PERCENT:  PRODUCT,
	token.EXP:      EXP,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []ParseError

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	ident  int
	inLoop int
}

func (p *Parser) Init(l *lexer.Lexer) {
	p.l = l

	// Set parser.curToken to first token in lexer
	p.nextToken()
	p.nextToken()

	// Remove all previous parsing errors
	p.errors = p.errors[0:0]
	p.ident = 0
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	p.peekError(tokenType)
	return false
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) Ident() int {
	return p.ident
}

func (p *Parser) checkUnexpectedEOF() bool {
	length := len(p.errors)
	if length >= 1 {
		if p.errors[length-1] == UNEXPECTED_EOF {
			return true
		}
	}

	return false
}

func (p *Parser) peekError(tokenType token.TokenType) {
	if p.peekTokenIs(token.EOF) {
		if !p.checkUnexpectedEOF() {
			p.errors = append(p.errors, UNEXPECTED_EOF)
		}
	} else {
		p.errors = append(p.errors, &IllegalSyntax{Expected: tokenType, Got: p.peekToken})
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseOptionalSemicolon() ast.Statement {
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return nil
}

func (p *Parser) parseIdentifierCommon() ast.Statement {
	if p.peekTokenIs(token.ASSIGN) {
		// Assign statement
		return p.parseAssignStatement()
	}

	// Expression statement
	return p.parseExpressionStatement()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.SEMICOLON:
		return p.parseOptionalSemicolon()
	case token.IDENTIFIER:
		return p.parseIdentifierCommon()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.LBRACE:
		return p.parseBlockCommon()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBlockCommon() ast.Statement {
	switch p.peekToken.Type {
	case "return":
		fallthrough
	case "while":
		fallthrough
	case ";":
		fallthrough
	case "break":
		fallthrough
	case "continue":
		fallthrough
	case "}":
		return p.parseBlockStatement()
	case "{":
		p.nextToken()
		p.ident++
		nestedBlock := p.parseBlockCommon()
		p.ident--

		if _, ok := nestedBlock.(*ast.BlockStatement); ok {
			p.nextToken()
			return p.parseBlockStatement2(nestedBlock)
		}

		expressionStatement := nestedBlock.(*ast.ExpressionStatement)

		if p.peekTokenIs(token.COLON) {
			// Skip colon
			p.nextToken()
			p.nextToken()

			firstValue := p.parseExpression(LOWEST)

			p.nextToken()
			mapLiteral := p.parseMapLiteral2(expressionStatement.Expression, firstValue)

			expressionStatement := &ast.ExpressionStatement{}
			expressionStatement.Expression = p.parseExpression2(LOWEST, mapLiteral)
			return expressionStatement
		}

		p.nextToken()
		return p.parseBlockStatement2(nestedBlock)
	}

	p.nextToken()
	piece := p.parseExpression(LOWEST)

	if identifier, ok := piece.(*ast.Identifier); ok {
		if p.peekTokenIs(token.ASSIGN) {
			p.ident++
			assignStatement := p.parseAssignStatement2(identifier)
			p.ident--

			p.nextToken()
			blockStatement := p.parseBlockStatement2(assignStatement)
			return blockStatement
		}

		// Skip optional semicolon
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		expressionStatement := &ast.ExpressionStatement{}
		expressionStatement.Expression = identifier

		p.nextToken()
		return p.parseBlockStatement2(expressionStatement)
	}

	if p.peekTokenIs(token.COLON) {
		// Skip colon
		p.nextToken()
		p.nextToken()

		firstValue := p.parseExpression(LOWEST)

		p.nextToken()
		mapLiteral := p.parseMapLiteral2(piece, firstValue)

		expressionStatement := &ast.ExpressionStatement{}
		expressionStatement.Expression = p.parseExpression2(LOWEST, mapLiteral)
		return expressionStatement
	}

	expressionStatement := &ast.ExpressionStatement{}
	expressionStatement.Expression = piece

	p.nextToken()
	return p.parseBlockStatement2(expressionStatement)
}

func (p *Parser) parseMapLiteralCommon(mapLiteral *ast.MapLiteral) ast.Expression {
	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)
		mapLiteral.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return mapLiteral
}

func (p *Parser) parseMapLiteral() ast.Expression {
	mapLiteral := &ast.MapLiteral{}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	return p.parseMapLiteralCommon(mapLiteral)
}

func (p *Parser) parseMapLiteral2(firstKey ast.Expression, firstValue ast.Expression) ast.Expression {
	mapLiteral := &ast.MapLiteral{}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)
	mapLiteral.Pairs[firstKey] = firstValue

	if p.curTokenIs(token.RBRACE) {
		return mapLiteral
	}

	return p.parseMapLiteralCommon(mapLiteral)
}

func (p *Parser) parseAssignStatementCommon(statement *ast.AssignStatement) *ast.AssignStatement {
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	statement := &ast.AssignStatement{}
	statement.Name = &ast.Identifier{Value: p.curToken.Literal}

	return p.parseAssignStatementCommon(statement)
}

func (p *Parser) parseAssignStatement2(identifier *ast.Identifier) *ast.AssignStatement {
	statement := &ast.AssignStatement{}
	statement.Name = identifier

	return p.parseAssignStatementCommon(statement)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	p.inLoop++
	statement := &ast.WhileStatement{}

	p.nextToken()

	statement.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	statement.Body = p.parseBlockStatement()

	p.inLoop--
	return statement
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	// Swallow optional semicolon first to avoid triggering extra no prefix function error
	// when break statement is not in a loop statement
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.inLoop == 0 {
		p.errors = append(p.errors, ILLEGAL_BREAK)
		return nil
	}

	return BREAK_STATEMENT
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	// Swallow optional semicolon first to avoid triggering extra no prefix function error
	// when continue statement is not in a loop statement
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.inLoop == 0 {
		p.errors = append(p.errors, ILLEGAL_CONTINUE)
		return nil
	}

	return CONTINUE_STATEMENT
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{}

	statement.Expression = p.parseExpression(LOWEST)

	// Skip optional semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	switch p.curToken.Type {
	case token.EOF:
		if !p.checkUnexpectedEOF() {
			p.errors = append(p.errors, UNEXPECTED_EOF)
		}
	case token.ILLEGAL:
		p.errors = append(p.errors, &IllegalSymbol{p.curToken.Literal})
	default:
		p.errors = append(p.errors, &NoPrefixFunction{Token: t})
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}

	leftExp := prefix()

	return p.parseExpression2(precedence, leftExp)
}

func (p *Parser) parseExpression2(precedence int, leftExp ast.Expression) ast.Expression {
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{}
	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	p.expectPeek(token.LBRACE)

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		p.expectPeek(token.LBRACE)
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatementCommon(block *ast.BlockStatement) *ast.BlockStatement {
	p.ident++
	block.Ident = p.ident

	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.EOF) {
			if !p.checkUnexpectedEOF() {
				p.errors = append(p.errors, UNEXPECTED_EOF)
			}
			return nil
		}

		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}

	p.ident--
	return block
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	p.nextToken()
	block = p.parseBlockStatementCommon(block)
	return block
}

func (p *Parser) parseBlockStatement2(firstStatement ast.Statement) *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{firstStatement}

	block = p.parseBlockStatementCommon(block)
	return block
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	expression := &ast.IntegerLiteral{}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, &IllegalInteger{Literal: p.curToken.Literal})
		return nil
	}

	expression.Value = value
	return expression
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Operator: p.curToken.Literal}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	function := &ast.FunctionLiteral{}

	function.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	function.Body = p.parseBlockStatement()

	return function
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	p.nextToken()

	if p.curTokenIs(token.VERTICAL) {
		return identifiers
	}

	identifier := &ast.Identifier{Value: p.curToken.Literal}
	identifiers = append(identifiers, identifier)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		identifier := &ast.Identifier{Value: p.curToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.VERTICAL) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Function: function}
	expression.Arguments = p.parseExpressionList(token.RPAREN)
	return expression
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func NewParser() *Parser {
	parser := &Parser{l: nil, errors: []ParseError{}, ident: 0}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.EXCLAMATION, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.VERTICAL, parser.parseFunctionLiteral)
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)
	parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral)
	parser.registerPrefix(token.LBRACE, parser.parseMapLiteral)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.LTE, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.GTE, parser.parseInfixExpression)
	parser.registerInfix(token.EXP, parser.parseInfixExpression)
	parser.registerInfix(token.PERCENT, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)

	return parser
}
