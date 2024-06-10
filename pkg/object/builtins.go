package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case String:
				return Integer(len(arg))
			case *Array:
				return Integer(len(arg.Elements))
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		}},
	},
	// TODO: 添加字符串操作函数(字符串包含, 正则匹配 ...)
	{
		"print",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return NULL
		}},
	},
}

func GetBuiltinByName(name string) *Builtin {
	for _, item := range Builtins {
		if item.Name == name {
			return item.Builtin
		}
	}
	return nil
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
