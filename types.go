package main

import (
	"fmt"
	"strconv"
)

func newTypeBuiltins() map[string]*BuiltinObject {
	return map[string]*BuiltinObject{
		"int": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newFail("wrong number of arguments. got=%d, want=1", len(args))
				}
				switch arg := args[0].(type) {
				case *StringObject:
					i, err := strconv.ParseInt(arg.Value, 10, 64)
					if err != nil {
						return newFail("could not parse string to int: %s", arg.Value)
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
					return newFail("argument to `int` not supported, got %s", args[0].Type())
				}
			},
		},
		"string": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newFail("wrong number of arguments. got=%d, want=1", len(args))
				}
				return &StringObject{Value: args[0].Inspect()}
			},
		},
		"bool": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newFail("wrong number of arguments. got=%d, want=1", len(args))
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
					return newFail("could not parse string to bool: %s", arg.Value)
				case *IntegerObject:
					if arg.Value != 0 {
						return True
					}
					return False
				default:
					return newFail("argument to `bool` not supported, got %s", args[0].Type())
				}
			},
		},
		"type": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newFail("wrong number of arguments. got=%d, want=1", len(args))
				}
				return &StringObject{Value: fmt.Sprintf("%s", args[0].Type())}
			},
		},
		"is_fail": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newFail("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() == FAIL_OBJ {
					return True
				}
				return False
			},
		},
	}
}
