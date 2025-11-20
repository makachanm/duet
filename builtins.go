package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var builtins = map[string]*BuiltinObject{
	"print": {
		Fn: func(args ...MemoryObject) MemoryObject {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return Nil
		},
	},
	"len": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *StringObject:
				return &IntegerObject{Value: int64(len(arg.Value))}
			case *ListObject:
				return &IntegerObject{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != LIST_OBJ {
				return newError("argument to `first` must be LIST, got %s", args[0].Type())
			}
			list := args[0].(*ListObject)
			if len(list.Elements) > 0 {
				return list.Elements[0]
			}
			return Nil
		},
	},
	"last": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != LIST_OBJ {
				return newError("argument to `last` must be LIST, got %s", args[0].Type())
			}
			list := args[0].(*ListObject)
			length := len(list.Elements)
			if length > 0 {
				return list.Elements[length-1]
			}
			return Nil
		},
	},
	"rest": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != LIST_OBJ {
				return newError("argument to `rest` must be LIST, got %s", args[0].Type())
			}
			list := args[0].(*ListObject)
			length := len(list.Elements)
			if length > 0 {
				newElements := make([]MemoryObject, length-1)
				copy(newElements, list.Elements[1:length])
				return &ListObject{Elements: newElements}
			}
			return Nil
		},
	},
	"push": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != LIST_OBJ {
				return newError("argument to `push` must be LIST, got %s", args[0].Type())
			}
			list := args[0].(*ListObject)
			length := len(list.Elements)
			newElements := make([]MemoryObject, length+1)
			copy(newElements, list.Elements)
			newElements[length] = args[1]
			return &ListObject{Elements: newElements}
		},
	},
	"readln": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 0 {
				return newError("wrong number of arguments. got=%d, want=0", len(args))
			}
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			return &StringObject{Value: strings.TrimSpace(text)}
		},
	},
	"int": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *StringObject:
				i, err := strconv.ParseInt(arg.Value, 10, 64)
				if err != nil {
					return newError("could not parse string to int: %s", arg.Value)
				}
				return &IntegerObject{Value: i}
			case *IntegerObject:
				return arg
			case *BooleanObject:
				if arg.Value {
					return &IntegerObject{Value: 1}
				}
				return &IntegerObject{Value: 0}
			default:
				return newError("argument to `int` not supported, got %s", args[0].Type())
			}
		},
	},
	"string": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &StringObject{Value: args[0].Inspect()}
		},
	},
	"bool": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *BooleanObject:
				return arg
			case *StringObject:
				if arg.Value == "true" {
					return True
				}
				if arg.Value == "false" {
					return False
				}
				return newError("could not parse string to bool: %s", arg.Value)
			case *IntegerObject:
				if arg.Value != 0 {
					return True
				}
				return False
			default:
				return newError("argument to `bool` not supported, got %s", args[0].Type())
			}
		},
	},
}
