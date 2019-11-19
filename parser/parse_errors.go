package parser

import (
	"fmt"

	"github.com/vita-dounai/Firework/token"
)

var UNEXPECTED_EOF = &UnexpectedEOF{}

type ParseError interface {
	Type() string
	Info() string
}

type IllegalSyntax struct {
	Expected token.TokenType
	Got      token.TokenType
}

func (is *IllegalSyntax) Type() string {
	return "IllegalSyntax"
}

func (is *IllegalSyntax) Info() string {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", is.Expected, is.Got)
	return msg
}

type IllegalSymbol struct {
	Symbol string
}

func (is *IllegalSymbol) Type() string {
	return "IllegalSymbol"
}

func (is *IllegalSymbol) Info() string {
	msg := fmt.Sprintf("symbol not recognized %q", is.Symbol)
	return msg
}

type UnexpectedEOF struct {
}

func (ue *UnexpectedEOF) Type() string {
	// A special type of IllegalSyntax
	return "IllegalSyntax"
}

func (ue *UnexpectedEOF) Info() string {
	return "Unexpected EOF"
}

type NoPrefixFunction struct {
	Token token.TokenType
}

func (npf *NoPrefixFunction) Type() string {
	return "NoPrefixFunction"
}

func (npf *NoPrefixFunction) Info() string {
	msg := fmt.Sprintf("no prefix parse function for %s found", npf.Token)
	return msg
}

type IllegalInteger struct {
	Literal string
}

func (ii *IllegalInteger) Type() string {
	return "IllegalInteger"
}

func (ii *IllegalInteger) Info() string {
	msg := fmt.Sprintf("counld not parse %q as integer", ii.Literal)
	return msg
}
