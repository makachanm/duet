package main

import (
	"fmt"
)

// ExcutionEngine은 AST와 실행 환경(메모리)을 가집니다.
type ExcutionEngine struct {
	Program *Program
	Memory  *Memory
}

// NewExcutionEngine은 새로운 실행 엔진을 생성합니다.
func NewExcutionEngine(program *Program, memory *Memory) *ExcutionEngine {
	if memory == nil {
		memory = NewMemory()
	}
	return &ExcutionEngine{Program: program, Memory: memory}
}

// Run은 프로그램 실행의 진입점입니다.
func (e *ExcutionEngine) Run() MemoryObject {
	return Eval(e.Program, e.Memory)
}

// Eval은 AST 노드를 받아 평가하고 MemoryObject를 반환하는 핵심 함수입니다.
func Eval(node Node, mem *Memory) MemoryObject {
	switch node := node.(type) {
	// 문 (Statements)
	case *Program:
		return evalProgram(node, mem)
	case *ExpressionStatement:
		return Eval(node.Expression, mem)
	case *FunctionStatement:
		fn := &FunctionObject{
			Name:       node.Name,
			Token:      node.Token,
			Parameters: node.Parameters,
			ReturnType: node.ReturnType,
			Body:       node.Body,
			Mem:        mem,
		}
		mem.Set(string(node.Name.Value), fn)
		return nil // 함수 정의는 값을 반환하지 않습니다.

	case *FailExpression:
		return &FailObject{Message: node.Message}

	// 표현식 (Expressions)
	case *Identifier:
		return evalIdentifier(node, mem)
	case *IntegerLiteral:
		return &IntegerObject{Value: node.Value}
	case *FloatLiteral:
		return &FloatObject{Value: node.Value}
	case *StringLiteral:
		return &StringObject{Value: node.Value}
	case *BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *NilLiteral:
		return Nil
	case *PrefixExpression:
		right := Eval(node.Right, mem)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		if node.Operator == "|>" {
			return evalPipelineExpression(node, mem)
		}
		left := Eval(node.Left, mem)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, mem)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *IfExpression:
		return evalIfExpression(node, mem)
	case *ForExpression:
		return evalForExpression(node, mem)
	case *CallExpression:
		function := Eval(node.Function, mem)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, mem)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, false)
	case *ListLiteral:
		elements := evalExpressions(node.Elements, mem)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &ListObject{Elements: elements}
	case *MapLiteral:
		return evalMapLiteral(node, mem)
	case *IndexExpression:
		left := Eval(node.Left, mem)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, mem)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
}

func evalIndexExpression(left, index MemoryObject) MemoryObject {
	switch {
	case left.Type() == LIST_OBJ && index.Type() == INTEGER_OBJ:
		return evalListIndexExpression(left, index)
	case left.Type() == MAP_OBJ:
		return evalMapIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalListIndexExpression(list, index MemoryObject) MemoryObject {
	listObject := list.(*ListObject)
	idx := index.(*IntegerObject).Value
	max := int64(len(listObject.Elements) - 1)
	if idx < 0 || idx > max {
		return Nil
	}
	return listObject.Elements[idx]
}
func evalPipelineExpression(node *InfixExpression, mem *Memory) MemoryObject {
	left := Eval(node.Left, mem)
	if isError(left) {
		return left
	}

	// If the left side is a zero-argument function (a supplier), invoke it
	// so the pipeline forwards the produced value instead of the function object.
	switch lf := left.(type) {
	case *FunctionObject:
		if len(lf.Parameters) == 0 {
			produced := applyFunction(lf, []MemoryObject{}, true)
			if isError(produced) {
				return produced
			}
			left = produced
		}
	case *BuiltinObject:
		// If left is a builtin and takes no args, call it to get its value.
		// Most builtins expect args, so this is a best-effort behavior.
		produced := lf.Fn()
		if isError(produced) {
			return produced
		}
		left = produced
	}

	// Case 1: The right side is a call expression, e.g., `data |> process(1, 2)`
	if call, ok := node.Right.(*CallExpression); ok {
		function := Eval(call.Function, mem)
		if isError(function) {
			return function
		}

		args := evalExpressions(call.Arguments, mem)
		if len(args) > 0 && isError(args[0]) {
			return args[0]
		}

		allArgs := append([]MemoryObject{left}, args...)
		return applyFunction(function, allArgs, true)
	}

	// Case 2: The right side is an identifier or other expression that yields a function, e.g., `data |> process`
	right := Eval(node.Right, mem)
	if isError(right) {
		return right
	}

	return applyFunction(right, []MemoryObject{left}, true)
}

func evalProgram(program *Program, mem *Memory) MemoryObject {
	var result MemoryObject
	for _, statement := range program.Statements {
		result = Eval(statement, mem)

		switch result := result.(type) {
		case *ReturnValueObject:
			return result.Value
		case *ErrorObject:
			return result
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *BooleanObject {
	if input {
		return True
	}
	return False
}

func evalPrefixExpression(operator string, right MemoryObject) MemoryObject {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right MemoryObject) MemoryObject {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		if left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ {
			return nativeBoolToBooleanObject(left.(*BooleanObject).Value == right.(*BooleanObject).Value)
		}
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		if left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ {
			return nativeBoolToBooleanObject(left.(*BooleanObject).Value != right.(*BooleanObject).Value)
		}
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right MemoryObject) MemoryObject {
	switch right {
	case True:
		return False
	case False:
		return True
	case Nil:
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperatorExpression(right MemoryObject) MemoryObject {
	if right.Type() == INTEGER_OBJ {
		value := right.(*IntegerObject).Value
		return &IntegerObject{Value: -value}
	}
	if right.Type() == FLOAT_OBJ {
		value := right.(*FloatObject).Value
		return &FloatObject{Value: -value}
	}
	return newError("unknown operator: -%s", right.Type())
}

func evalIntegerInfixExpression(operator string, left, right MemoryObject) MemoryObject {
	leftVal := left.(*IntegerObject).Value
	rightVal := right.(*IntegerObject).Value
	switch operator {
	case "+":
		return &IntegerObject{Value: leftVal + rightVal}
	case "-":
		return &IntegerObject{Value: leftVal - rightVal}
	case "*":
		return &IntegerObject{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &IntegerObject{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &IntegerObject{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right MemoryObject) MemoryObject {
	leftVal := left.(*FloatObject).Value
	rightVal := right.(*FloatObject).Value
	switch operator {
	case "+":
		return &FloatObject{Value: leftVal + rightVal}
	case "-":
		return &FloatObject{Value: leftVal - rightVal}
	case "*":
		return &FloatObject{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0.0 {
			return newError("division by zero")
		}
		return &FloatObject{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right MemoryObject) MemoryObject {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*StringObject).Value
	rightVal := right.(*StringObject).Value
	return &StringObject{Value: leftVal + rightVal}
}

func evalIfExpression(ie *IfExpression, mem *Memory) MemoryObject {
	condition := Eval(ie.Condition, mem)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, mem)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, mem)
	} else {
		return Nil
	}
}

func evalForExpression(fe *ForExpression, mem *Memory) MemoryObject {
	collection := Eval(fe.Collection, mem)
	if isError(collection) {
		return collection
	}

	list, ok := collection.(*ListObject)
	if !ok {
		return newError("for loop must iterate over a list, got %s", collection.Type())
	}

	results := []MemoryObject{}
	for _, el := range list.Elements {
		loopMem := NewEnclosedMemory(mem)
		loopMem.Set(fe.Variable.Value, el)
		result := Eval(fe.Body, loopMem)
		if isError(result) {
			return result
		}
		results = append(results, result)
	}

	return &ListObject{Elements: results}
}

func evalMapLiteral(node *MapLiteral, mem *Memory) MemoryObject {
	pairs := make(map[string]MapPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, mem)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, mem)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = MapPair{Key: key, Value: value}
	}

	return &MapObject{Pairs: pairs}
}

func evalMapIndexExpression(m, index MemoryObject) MemoryObject {
	mapObject := m.(*MapObject)
	key, ok := index.(Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := mapObject.Pairs[key.HashKey()]
	if !ok {
		return Nil
	}
	return pair.Value
}

func evalIdentifier(node *Identifier, mem *Memory) MemoryObject {
	if val, ok := mem.Get(string(node.Value)); ok {
		// If the identifier refers to a zero-argument supplier (supp/esupp),
		// invoke it and return the produced value instead of the function object.
		if fn, ok := val.(*FunctionObject); ok {
			if (fn.Token.Type == SUPP || fn.Token.Type == ESUPP) && len(fn.Parameters) == 0 {
				produced := applyFunction(fn, []MemoryObject{}, false)
				return produced
			}
		}
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + string(node.Value))
}

func evalExpressions(exps []Expression, mem *Memory) []MemoryObject {
	var result []MemoryObject
	for _, e := range exps {
		evaluated := Eval(e, mem)
		if _, ok := e.(*FailExpression); ok {
			return []MemoryObject{evaluated}
		}

		if isError(evaluated) {
			return []MemoryObject{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn MemoryObject, args []MemoryObject, isPipeline bool) MemoryObject {
	switch fn := fn.(type) {
	case *FunctionObject:
		extendedMem := extendFunctionMem(fn, args)
		evaluated := Eval(fn.Body, extendedMem)

		// Unwrap return value if it's wrapped in a ReturnValueObject
		if returnValue, ok := evaluated.(*ReturnValueObject); ok {
			evaluated = returnValue.Value
		}

		// Check if the return type matches the function's signature
		if fn.ReturnType != nil {
			expectedType := fn.ReturnType.Value
			actualType := evaluated.Type()

			// A simple type mapping check
			typeMatch := false
			switch expectedType {
			case "int":
				typeMatch = (actualType == INTEGER_OBJ)
			case "float":
				typeMatch = (actualType == FLOAT_OBJ)
			case "string":
				typeMatch = (actualType == STRING_OBJ)
			case "bool":
				typeMatch = (actualType == BOOLEAN_OBJ)
			case "list":
				typeMatch = (actualType == LIST_OBJ)
			case "map":
				typeMatch = (actualType == MAP_OBJ)
			}

			if (fn.Token.Type == ESUPP || fn.Token.Type == EPROC) && !typeMatch {
				typeMatch = (actualType == FAIL_OBJ)
			}

			if !typeMatch {
				return newError("type error: function %s returned %s, but expected %s", fn.Name.Value, actualType, expectedType)
			}
		}
		return evaluated

	case *BuiltinObject:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionMem(fn *FunctionObject, args []MemoryObject) *Memory {
	mem := NewEnclosedMemory(fn.Mem)
	for i, param := range fn.Parameters {
		mem.Set(string(param.Name.Value), args[i])
	}
	return mem
}

func isTruthy(obj MemoryObject) bool {
	switch obj {
	case Nil:
		return false
	case True:
		return true
	case False:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *ErrorObject {
	return &ErrorObject{Message: fmt.Sprintf(format, a...)}
}

func isError(obj MemoryObject) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
