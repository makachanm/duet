package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var builtins = map[string]*BuiltinObject{
	"type": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			} else {
				return &StringObject{Value: string(args[0].Type())}
			}
		},
	},
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
	"read": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*StringObject)
			if !ok {
				return newError("argument to `read` must be STRING, got %s", args[0].Type())
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("could not read file: %s", err)
			}
			return &StringObject{Value: string(data)}
		},
	},
	"write": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok := args[0].(*StringObject)
			if !ok {
				return newError("first argument to `write` must be STRING, got %s", args[0].Type())
			}
			content, ok := args[1].(*StringObject)
			if !ok {
				return newError("second argument to `write` must be STRING, got %s", args[1].Type())
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("could not write file: %s", err)
			}
			return True
		},
	},
	"lines": {
		Fn: func(args ...MemoryObject) MemoryObject {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*StringObject)
			if !ok {
				return newError("argument to `lines` must be STRING, got %s", args[0].Type())
			}
			file, err := os.Open(path.Value)
			if err != nil {
				return newError("could not open file: %s", err)
			}
			defer file.Close()
			var lines []MemoryObject
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, &StringObject{Value: scanner.Text()})
			}
			return &ListObject{Elements: lines}
		},
	},
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
