package main

import (
	"bytes"
	"fmt"
	"strings"
)

// MemoryObjectType은 메모리에 저장될 객체의 타입을 나타냅니다.
type MemoryObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NIL_OBJ          = "NIL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FAIL_OBJ         = "FAIL"
	FUNCTION_OBJ     = "FUNCTION"
	LIST_OBJ         = "LIST"
	BUILTIN_OBJ      = "BUILTIN"
	MAP_OBJ          = "MAP"
)

// MemoryObject는 인터프리터에서 다루는 모든 값(객체)이 구현해야 하는 인터페이스입니다.
type MemoryObject interface {
	Type() MemoryObjectType
	Inspect() string
}

// 각 데이터 타입을 위한 구조체 정의

type IntegerObject struct {
	Value int64
}

func (i *IntegerObject) Type() MemoryObjectType { return INTEGER_OBJ }
func (i *IntegerObject) Inspect() string        { return fmt.Sprintf("%d", i.Value) }

type FloatObject struct {
	Value float64
}

func (f *FloatObject) Type() MemoryObjectType { return FLOAT_OBJ }
func (f *FloatObject) Inspect() string        { return fmt.Sprintf("%f", f.Value) }

type StringObject struct {
	Value string
}

func (s *StringObject) Type() MemoryObjectType { return STRING_OBJ }
func (s *StringObject) Inspect() string        { return s.Value }

type BooleanObject struct {
	Value bool
}

func (b *BooleanObject) Type() MemoryObjectType { return BOOLEAN_OBJ }
func (b *BooleanObject) Inspect() string        { return fmt.Sprintf("%t", b.Value) }

type NilObject struct{}

func (n *NilObject) Type() MemoryObjectType { return NIL_OBJ }
func (n *NilObject) Inspect() string        { return "nil" }

type ReturnValueObject struct {
	Value MemoryObject
}

func (rv *ReturnValueObject) Type() MemoryObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValueObject) Inspect() string        { return rv.Value.Inspect() }

type ErrorObject struct {
	Message string
}

func (e *ErrorObject) Type() MemoryObjectType { return ERROR_OBJ }
func (e *ErrorObject) Inspect() string        { return "ERROR: " + e.Message }

type FailObject struct {
	Message string
}

func (e *FailObject) Type() MemoryObjectType { return FAIL_OBJ }
func (e *FailObject) Inspect() string        { return e.Message }

type FunctionObject struct {
	Name       *Identifier
	Token      Token // The function type token (e.g., PROC, CONS, SUPP)
	Parameters []*Parameter
	ReturnType *Identifier
	Body       Expression
	Mem        *Memory
}

func (f *FunctionObject) Type() MemoryObjectType { return FUNCTION_OBJ }
func (f *FunctionObject) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.Name.String()+":"+p.Type.String())
	}

	out.WriteString(f.Token.Literal)
	out.WriteString(" ")
	out.WriteString(f.Name.String())
	if len(f.Parameters) > 0 {
		out.WriteString("(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(")")
	}
	if f.ReturnType != nil {
		out.WriteString(":")
		out.WriteString(f.ReturnType.String())
	}
	out.WriteString(" -> ")
	out.WriteString(f.Body.String())

	return out.String()
}

type ListObject struct {
	Elements []MemoryObject
}

func (l *ListObject) Type() MemoryObjectType { return LIST_OBJ }
func (l *ListObject) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range l.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type MapPair struct {
	Key   MemoryObject
	Value MemoryObject
}

type MapObject struct {
	Pairs map[string]MapPair
}

func (m *MapObject) Type() MemoryObjectType { return MAP_OBJ }
func (m *MapObject) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Hashable interface {
	HashKey() string
}

func (i *IntegerObject) HashKey() string { return fmt.Sprintf("i:%d", i.Value) }
func (s *StringObject) HashKey() string  { return fmt.Sprintf("s:%s", s.Value) }
func (b *BooleanObject) HashKey() string { return fmt.Sprintf("b:%t", b.Value) }

type BuiltinFunction func(args ...MemoryObject) MemoryObject

type BuiltinObject struct {
	Fn BuiltinFunction
}

func (b *BuiltinObject) Type() MemoryObjectType { return BUILTIN_OBJ }
func (b *BuiltinObject) Inspect() string        { return "builtin function" }

var (
	True  = &BooleanObject{Value: true}
	False = &BooleanObject{Value: false}
	Nil   = &NilObject{}
)

// Memory는 변수와 함수를 저장하는 환경(Environment)입니다.
// outer 필드를 통해 중첩된 스코프(lexical scope)를 구현합니다.
type Memory struct {
	store map[string]MemoryObject
	outer *Memory
}

// NewMemory는 새로운 최상위 메모리(전역 스코프)를 생성합니다.
func NewMemory() *Memory {
	s := make(map[string]MemoryObject)
	return &Memory{store: s, outer: nil}
}

// NewEnclosedMemory는 외부 스코프를 감싸는 새로운 내부 스코프를 생성합니다.
// 함수 호출 시 지역 변수를 관리하기 위해 사용됩니다.
func NewEnclosedMemory(outer *Memory) *Memory {
	mem := NewMemory()
	mem.outer = outer
	return mem
}

// Get은 현재 스코프 또는 외부 스코프에서 변수 값을 찾습니다.
func (m *Memory) Get(name string) (MemoryObject, bool) {
	obj, ok := m.store[name]
	if !ok && m.outer != nil {
		obj, ok = m.outer.Get(name)
	}
	return obj, ok
}

// Set은 현재 스코프에 변수 값을 설정(또는 생성)합니다.
func (m *Memory) Set(name string, val MemoryObject) MemoryObject {
	m.store[name] = val
	return val
}
