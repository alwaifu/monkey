package interpreter

import (
	"fmt"

	"github.com/alwaifu/monkey/pkg/ast"
	"github.com/alwaifu/monkey/pkg/object"
)

var (
	TRUE  = object.True
	FALSE = object.False
	NULL  = object.NULL
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	var result object.Object
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	// statment
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if val.Type() == object.ERROR_OBJ {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if val.Type() == object.ERROR_OBJ {
			return val // fail fast
		}
		return &object.ReturnValue{Value: val}
	// expression
	case *ast.IntegerLiteral:
		return object.Integer(node.Value)
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return object.String(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if right.Type() == object.ERROR_OBJ {
			return right // fail fast
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if left.Type() == object.ERROR_OBJ {
			return left // fail fast
		}
		// FIXME: Implement short circuiting. See: https://en.wikipedia.org/wiki/Short-circuit_evaluation
		right := Eval(node.Right, env)
		if right.Type() == object.ERROR_OBJ {
			return right // fail fast
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		condition := Eval(node.Condition, env)
		if condition.Type() == object.ERROR_OBJ {
			return condition // fail fast
		}
		if isTruthy(condition) {
			return Eval(node.Consequence, env)
		} else if node.Alternative != nil {
			return Eval(node.Alternative, env)
		} else {
			return NULL
		}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if function.Type() == object.ERROR_OBJ {
			return function
		}
		args := make([]object.Object, 0, len(node.Arguments))
		for _, e := range node.Arguments {
			evaluated := Eval(e, env)
			if evaluated.Type() == object.ERROR_OBJ {
				return evaluated
			}
			args = append(args, evaluated)
		}
		// args := evalExpressions(node.Arguments, env)
		// if len(args) == 1 && args[0].Type() == object.ERROR_OBJ {
		// 	return args[0]
		// }
		return applyFunction(function, args)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && elements[0].Type() == object.ERROR_OBJ {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if left.Type() == object.ERROR_OBJ {
			return left
		}
		index := Eval(node.Index, env)
		if index.Type() == object.ERROR_OBJ {
			return index
		}
		return evalIndexExpression(left, index)
	}
	return result
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			if rt := result.Type(); rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result // 对比evalProgram函数 此处不解包return_value, 直接传递给外层来中断外层语句块
			}
		}
	}
	return result
}
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(object.Integer)
	return -value
}
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left.(object.Integer), right.(object.Integer))
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left.(object.String), right.(object.String))
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left.(object.Boolean), right.(object.Boolean))
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
func evalIntegerInfixExpression(operator string, left, right object.Integer) object.Object {
	switch operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right
	case "<":
		return nativeBoolToBooleanObject(left < right)
	case "<=":
		return nativeBoolToBooleanObject(left <= right)
	case ">":
		return nativeBoolToBooleanObject(left > right)
	case ">=":
		return nativeBoolToBooleanObject(left >= right)
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
func evalStringInfixExpression(operator string, left, right object.String) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	case "+":
		return left + right
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}
func evalBooleanInfixExpression(operator string, left, right object.Boolean) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	case "and":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case "or":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	} else if buitin := object.GetBuiltinByName(node.Value); buitin != nil {
		return buitin
	} else {
		return newError("identifier not found: " + node.Value)
	}
}
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		arrayObject := left.(*object.Array)
		idx := int64(index.(object.Integer))
		if idx < 0 || idx >= int64(len(arrayObject.Elements)) { // TODO: 数组索引负数&超出数组长度
			return NULL
		}
		return arrayObject.Elements[idx]
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	result := make([]object.Object, 0, len(exps))
	for _, e := range exps {
		evaluated := Eval(e, env)
		if evaluated.Type() == object.ERROR_OBJ {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		env := object.NewEnclosedEnvironment(fn.Env)
		for i, param := range fn.Parameters {
			env.Set(param.Value, args[i])
		}
		evaluated := Eval(fn.Body, env)
		// FIXME: 这里需要解包
		// if evaluated, ok := evaluated.(*object.ReturnValue); ok {
		// 	return evaluated.Value
		// }
		return evaluated
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
func nativeBoolToBooleanObject(input bool) object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
