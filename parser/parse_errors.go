package parser

import (
	"fmt"

	"github.com/vita-dounai/Firework/token"
)

const (
	UNEXPECTED_EOF_ERROR    = "UNEXPECTED_EOF"
	ILLEGAL_SYNTAX_ERROR    = "ILLEGAL_SYNTAX"
	ILLEGAL_SYMBOL_ERROR    = "ILLEGAL_SYMBOL"
	NOPREFIX_FUNCTION_ERROR = "NOPREFIX_FUNCTION"
	ILLEGAL_INTEGER_ERROR   = "ILLEGAL_INTEGER"
	ILLEGAL_BREAK_ERROR     = "ILLEGAL_BREAK"
	ILLEGAL_CONTINUE_ERROR  = "ILLEGAL_CONTINUE"
)

type ParseError interface {
	Type() string
	Info() string
}

type IllegalSyntax struct {
	Expected token.TokenType
	Got      token.Token
}

func (is *IllegalSyntax) Type() string {
	return ILLEGAL_SYNTAX_ERROR
}

func (is *IllegalSyntax) Info() string {
	msg := fmt.Sprintf("expected next token to be `%s`, got `%s` instead, at line: %d, column: %d", is.Expected, is.Got.Literal, is.Got.Line, is.Got.Column)
	return msg
}

type IllegalSymbol struct {
	Symbol string
}

func (is *IllegalSymbol) Type() string {
	return ILLEGAL_SYMBOL_ERROR
}

func (is *IllegalSymbol) Info() string {
	msg := fmt.Sprintf("symbol not recognized `%s`", is.Symbol)
	return msg
}

type UnexpectedEOF struct {
}

func (ue *UnexpectedEOF) Type() string {
	return UNEXPECTED_EOF_ERROR
}

func (ue *UnexpectedEOF) Info() string {
	return "Unexpected EOF"
}

type NoPrefixFunction struct {
	Token token.Token
}

func (npf *NoPrefixFunction) Type() string {
	return NOPREFIX_FUNCTION_ERROR
}

func (npf *NoPrefixFunction) Info() string {
	msg := fmt.Sprintf("no prefix parse function for `%s` found, line: %d, column: %d",
		npf.Token.Literal, npf.Token.Line, npf.Token.Column)
	return msg
}

type IllegalInteger struct {
	Literal string
}

func (ii *IllegalInteger) Type() string {
	return ILLEGAL_INTEGER_ERROR
}

func (ii *IllegalInteger) Info() string {
	msg := fmt.Sprintf("counld not parse `%s` as integer", ii.Literal)
	return msg
}

type IllegalBreak struct{}

func (ib *IllegalBreak) Type() string {
	return ILLEGAL_BREAK_ERROR
}

func (ib *IllegalBreak) Info() string {
	msg := fmt.Sprintf("break should be used in loop statement")
	return msg
}

type IllegalContinue struct{}

func (ib *IllegalContinue) Type() string {
	return ILLEGAL_CONTINUE_ERROR
}

func (ib *IllegalContinue) Info() string {
	msg := fmt.Sprintf("continue should be used in loop statement")
	return msg
}
