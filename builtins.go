package main

func newBuiltins() map[string]*BuiltinObject {
	builtins := map[string]*BuiltinObject{
		"type": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				} else {
					return &StringObject{Value: string(args[0].Type())}
				}
			},
		},
	}

	for name, builtin := range newIOBuiltins() {
		builtins[name] = builtin
	}

	for name, builtin := range newListBuiltins() {
		builtins[name] = builtin
	}

	for name, builtin := range newStringBuiltins() {
		builtins[name] = builtin
	}

	for name, builtin := range newTypeBuiltins() {
		builtins[name] = builtin
	}

	for name, builtin := range newMathBuiltins() {
		builtins[name] = builtin
	}

	return builtins
}

var builtins = newBuiltins()
