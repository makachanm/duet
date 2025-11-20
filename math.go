package main

import (
	"math"
)

func getFloat(obj MemoryObject) (float64, bool) {
	switch obj := obj.(type) {
	case *IntegerObject:
		return float64(obj.Value), true
	case *FloatObject:
		return obj.Value, true
	default:
		return 0, false
	}
}

func newMathBuiltins() map[string]*BuiltinObject {
	return map[string]*BuiltinObject{
		"abs": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				val, ok := getFloat(args[0])
				if !ok {
					return newError("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
				}
				return &FloatObject{Value: math.Abs(val)}
			},
		},
		"sqrt": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				val, ok := getFloat(args[0])
				if !ok {
					return newError("argument to `sqrt` must be INTEGER or FLOAT, got %s", args[0].Type())
				}
				return &FloatObject{Value: math.Sqrt(val)}
			},
		},
		"pow": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				base, ok := getFloat(args[0])
				if !ok {
					return newError("base for `pow` must be INTEGER or FLOAT, got %s", args[0].Type())
				}
				exp, ok := getFloat(args[1])
				if !ok {
					return newError("exponent for `pow` must be INTEGER or FLOAT, got %s", args[1].Type())
				}
				return &FloatObject{Value: math.Pow(base, exp)}
			},
		},
		"sin": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				val, ok := getFloat(args[0])
				if !ok {
					return newError("argument to `sin` must be INTEGER or FLOAT, got %s", args[0].Type())
				}
				return &FloatObject{Value: math.Sin(val)}
			},
		},
		"cos": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				val, ok := getFloat(args[0])
				if !ok {
					return newError("argument to `cos` must be INTEGER or FLOAT, got %s", args[0].Type())
				}
				return &FloatObject{Value: math.Cos(val)}
			},
		},
		"tan": {
			Fn: func(args ...MemoryObject) MemoryObject {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				val, ok := getFloat(args[0])
				if !ok {
					return newError("argument to `tan` must be INTEGER or FLOAT, got %s", args[0].Type())
				}
				return &FloatObject{Value: math.Tan(val)}
			},
		},
	}
}
