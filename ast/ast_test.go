package ast

import (
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&AssignStatement{
				Name: &Identifier{
					Value: "myVar",
				},
				Value: &Identifier{
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%q", program.String())
	}
}
