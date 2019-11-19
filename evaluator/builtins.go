package evaluator

import (
	"fmt"
	"strings"

	"github.com/vita-dounai/Firework/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of arguments, got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("Argument to `len` not supported, got %s",
					arg.Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			if len(array.Elements) > 0 {
				return array.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `last` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length > 0 {
				return array.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `rest` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length > 0 {
				rest := make([]object.Object, length-1, length-1)
				copy(rest, array.Elements[1:length])
				return &object.Array{Elements: rest}
			}

			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)

			newArray := make([]object.Object, length+1, length+1)
			copy(newArray, array.Elements)
			newArray[length] = args[1]
			return &object.Array{Elements: newArray}
		},
	},
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			length := len(args)
			for i, arg := range args {
				var inspect string

				if _, ok := arg.(*object.String); ok {
					inspect = strings.Trim(arg.Inspect(), "\"")
				} else {
					inspect = arg.Inspect()
				}

				fmt.Print(inspect)
				if i < length {
					fmt.Print(" ")
				}
			}
			fmt.Print("\n")
			return nil
		},
	},
}
