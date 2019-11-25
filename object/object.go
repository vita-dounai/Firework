package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vita-dounai/Firework/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	BREAK            = "BREAK"
	CONTINUE         = "CONTINUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	parameters := []string{}
	for _, p := range f.Parameters {
		parameters = append(parameters, p.String())
	}

	out.WriteString("|")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString("| ")
	out.WriteString(f.Body.String())

	return out.String()
}
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

type String struct {
	Value string
}

func (s *String) Inspect() string  { return "\"" + s.Value + "\"" }
func (s *String) Type() ObjectType { return STRING_OBJ }

type BuiltFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltFunction
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (a *Array) Type() ObjectType { return ARRAY_OBJ }

type LoopControl struct {
	ControlType ObjectType
}

func (lc *LoopControl) Inspect() string  { return "" }
func (lc *LoopControl) Type() ObjectType { return lc.ControlType }
