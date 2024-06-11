package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alwaifu/monkey/pkg/ast"
)

type ObjectType string
type Object interface {
	Type() ObjectType
	Inspect() string
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

	RETURN_VALUE_OBJ      = "RETURN_VALUE"
	ERROR_OBJ             = "ERROR"
	FUNCTION_OBJ          = "FUNCTION"
	BUILTIN_OBJ           = "BUILTIN"
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
	ARRAY_OBJ             = "ARRAY"
)

var (
	True  = Boolean(true)
	False = Boolean(false)
	NULL  = Null{}
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

// ReturnValue
type ReturnValue struct {
	Value Object
}

var _ Object = (*ReturnValue)(nil)

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error
type Error struct {
	Message string
}

var _ Object = (*Error)(nil)

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e Error) Error() string     { return e.Message }

// Function for interpreter
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

// BuiltinFunction
type BuiltinFunction func(args ...Object) Object

var _ Object = (*Builtin)(nil)

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// CompiledFunction for vm
type CompiledFunction struct {
	Instructions  []byte
	NumLocals     int
	NumParameters int
}

var _ Object = (*CompiledFunction)(nil)

func (cf *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OBJ }
func (cf *CompiledFunction) Inspect() string  { return fmt.Sprintf("CompiledFunction[%p]", cf) }

// Array
type Array struct {
	Elements []Object
}

var _ Object = (*Array)(nil)

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
