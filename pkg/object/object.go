package object

import (
	"bytes"
	"errors"
	"fmt"
	"monkey/pkg/ast"
	"strings"
)

type ObjectType string
type Object interface {
	Type() ObjectType
	Inspect() string
}

func NewEnviromentFromMap(dir map[string]interface{}) (*Environment, error) {
	store := make(map[string]Object, len(dir))
	for k, v := range dir {
		switch v := v.(type) {
		case bool:
			store[k] = Boolean(v)
		case int:
			store[k] = Integer(v)
		case int64:
			store[k] = Integer(v)
		case string:
			store[k] = String(v)
		default:
			return nil, errors.New("invalid type")
		}
	}
	return &Environment{store: store}, nil
}
func ToGoValue(obj Object) (interface{}, error) {
	switch obj := obj.(type) {
	case Boolean:
		return bool(obj), nil
	case Integer:
		return int64(obj), nil
	case String:
		return string(obj), nil
	default:
		return nil, errors.New("invalid type")
	}
}
func NewEnviroment() *Environment {
	return &Environment{store: make(map[string]Object)}
}
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnviroment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// TODO: 调整对象系统 使用golang原生对象 提高求值性能
// TODO: 添加对浮点数据的支持
const (
	// 基础类型

	NULL_OBJ    = "NULL"
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"

	// 指针类型

	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
)

type Null struct{}

var _ Object = (Null)(struct{}{})

func (n Null) Type() ObjectType { return NULL_OBJ }
func (n Null) Inspect() string  { return "null" }

type Integer int64

var _ Object = (Integer)(0)

func (i Integer) Type() ObjectType { return INTEGER_OBJ }
func (i Integer) Inspect() string  { return fmt.Sprintf("%d", i) }

type Boolean bool

var _ Object = (Boolean)(false)

func (b Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b Boolean) Inspect() string  { return fmt.Sprintf("%t", b) }

type String string

var _ Object = (String)("")

func (s String) Type() ObjectType { return STRING_OBJ }
func (s String) Inspect() string  { return string(s) }

// ---

type ReturnValue struct {
	Value Object
}

var _ Object = (*ReturnValue)(nil)

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

var _ Object = (*Error)(nil)

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e Error) Error() string     { return e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

var _ Object = (*Function)(nil)

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

type BuiltinFunction func(args ...Object) Object

var _ Object = (*Builtin)(nil)

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

var _ Object = (*Array)(nil)

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
