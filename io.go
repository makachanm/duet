package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func newIOBuiltins() map[string]*BuiltinObject {
	return map[string]*BuiltinObject{
		"print": {
			Fn: func(args ...MemoryObject) MemoryObject {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return Nil
			},
		},
		"readln": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 0 {
					return newFail("wrong number of arguments. got=%d, want=0", len(args))
				}
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				return &StringObject{Value: strings.TrimSpace(text)}
			},
		},
		"read": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newFail("wrong number of arguments. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*StringObject)
				if !ok {
					return newFail("argument to `read` must be STRING, got %s", args[0].Type())
				}
				data, err := os.ReadFile(path.Value)
				if err != nil {
					return newFail("could not read file: %s", err)
				}
				return &StringObject{Value: string(data)}
			},
		},
		"write": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 2 {
					return newFail("wrong number of arguments. got=%d, want=2", len(args))
				}
				path, ok := args[0].(*StringObject)
				if !ok {
					return newFail("first argument to `write` must be STRING, got %s", args[0].Type())
				}
				content, ok := args[1].(*StringObject)
				if !ok {
					return newFail("second argument to `write` must be STRING, got %s", args[1].Type())
				}
				err := os.WriteFile(path.Value, []byte(content.Value), 0644)
				if err != nil {
					return newFail("could not write file: %s", err)
				}
				return True
			},
		},
		"lines": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newFail("wrong number of arguments. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*StringObject)
				if !ok {
					return newFail("argument to `lines` must be STRING, got %s", args[0].Type())
				}
				file, err := os.Open(path.Value)
				if err != nil {
					return newFail("could not open file: %s", err)
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
	}
}
