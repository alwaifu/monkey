package monkey

import (
	"fmt"
	"monkey/pkg/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case object.String:
				return object.Integer(len(arg))
			case *object.Array:
				return object.Integer(len(arg.Elements))
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
		// TODO: 添加字符串操作函数(字符串包含, 正则匹配 ...)
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return NULL
		},
	},
}
