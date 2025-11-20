package main

import (
	"strings"
)

func newStringBuiltins() map[string]*BuiltinObject {
	return map[string]*BuiltinObject{
		"split": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				s, ok := args[0].(*StringObject)
				if !ok {
					return newError("first argument to `split` must be STRING, got %s", args[0].Type())
				}
				sep, ok := args[1].(*StringObject)
				if !ok {
					return newError("second argument to `split` must be STRING, got %s", args[1].Type())
				}
				parts := strings.Split(s.Value, sep.Value)
				elements := make([]MemoryObject, len(parts))
				for i, p := range parts {
					elements[i] = &StringObject{Value: p}
				}
				return &ListObject{Elements: elements}
			},
		},
		"join": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				list, ok := args[0].(*ListObject)
				if !ok {
					return newError("first argument to `join` must be LIST, got %s", args[0].Type())
				}
				sep, ok := args[1].(*StringObject)
				if !ok {
					return newError("second argument to `join` must be STRING, got %s", args[1].Type())
				}
				var parts []string
				for _, el := range list.Elements {
					s, ok := el.(*StringObject)
					if !ok {
						return newError("all elements in list for `join` must be STRING, got %s", el.Type())
					}
					parts = append(parts, s.Value)
				}
				return &StringObject{Value: strings.Join(parts, sep.Value)}
			},
		},
		"trim": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				s, ok := args[0].(*StringObject)
				if !ok {
					return newError("argument to `trim` must be STRING, got %s", args[0].Type())
				}
				return &StringObject{Value: strings.TrimSpace(s.Value)}
			},
		},
		"upper": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				s, ok := args[0].(*StringObject)
				if !ok {
					return newError("argument to `upper` must be STRING, got %s", args[0].Type())
				}
				return &StringObject{Value: strings.ToUpper(s.Value)}
			},
		},
		"lower": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				s, ok := args[0].(*StringObject)
				if !ok {
					return newError("argument to `lower` must be STRING, got %s", args[0].Type())
				}
				return &StringObject{Value: strings.ToLower(s.Value)}
			},
		},
		"replace": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 3 {
					return newError("wrong number of arguments. got=%d, want=3", len(args))
				}
				s, ok := args[0].(*StringObject)
				if !ok {
					return newError("first argument to `replace` must be STRING, got %s", args[0].Type())
				}
				old, ok := args[1].(*StringObject)
				if !ok {
					return newError("second argument to `replace` must be STRING, got %s", args[1].Type())
				}
				newStr, ok := args[2].(*StringObject)
				if !ok {
					return newError("third argument to `replace` must be STRING, got %s", args[2].Type())
				}
				return &StringObject{Value: strings.ReplaceAll(s.Value, old.Value, newStr.Value)}
			},
		},
		"contains": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				s, ok := args[0].(*StringObject)
				if !ok {
					return newError("first argument to `contains` must be STRING, got %s", args[0].Type())
				}
				sub, ok := args[1].(*StringObject)
				if !ok {
					return newError("second argument to `contains` must be STRING, got %s", args[1].Type())
				}
				if strings.Contains(s.Value, sub.Value) {
					return True
				}
				return False
			},
		},
	}
}
