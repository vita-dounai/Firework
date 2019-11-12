package ast

import "testing"
import "github.com/vita-dounai/Firework/token"

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&AssignStatement{
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%q", program.String())
	}
}
