package main

func newListBuiltins() map[string]*BuiltinObject {
	return map[string]*BuiltinObject{
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
	}
}
