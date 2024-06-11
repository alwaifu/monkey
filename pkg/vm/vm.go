package vm

import (
	"fmt"

	"github.com/alwaifu/monkey/pkg/object"
)

const StackSize = 2048
const GlobalSize = 65536
const FramesSize = 1024

var (
	True  = object.True
	False = object.False
	NULL  = object.NULL
)

type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int

	globals []object.Object

	frames []*Frame
}

func NewVM(c *Compiler, globals []object.Object) *VM {
	mainFrame := NewFrame(&object.CompiledFunction{Instructions: c.scopes[c.scopeIndex].instructions}, 0)
	frames := make([]*Frame, 0, FramesSize)
	frames = append(frames, mainFrame)
	return &VM{
		constants: c.Constants,
		stack:     make([]object.Object, StackSize),
		sp:        0,
		globals:   globals,
		frames:    frames,
	}
}
func (vm *VM) Run() error {
	for caller := vm.frames[0]; caller.pc < len(caller.fn.Instructions); caller = vm.frames[len(vm.frames)-1] {
		ins := caller.fn.Instructions[caller.pc]
		caller.pc++

		op := Opcode(ins)
		switch op {
		case OpConstant:
			constIdx := caller.readInsOprandUint16()
			vm.push(vm.constants[constIdx])
		case OpPop:
			vm.pop()
		case OpAdd:
			right := vm.pop()
			left := vm.pop()
			if r, err := add(left, right); err != nil {
				return err
			} else {
				vm.push(r)
			}
		case OpSub:
			right := vm.pop().(object.Integer)
			left := vm.pop().(object.Integer)
			vm.push(object.Integer(left - right))
		case OpMul:
			right := vm.pop().(object.Integer)
			left := vm.pop().(object.Integer)
			vm.push(object.Integer(left * right))
		case OpDiv:
			right := vm.pop().(object.Integer)
			left := vm.pop().(object.Integer)
			vm.push(object.Integer(left / right))
		case OpTrue:
			vm.push(True)
		case OpFalse:
			vm.push(False)
		case OpAnd:
			right := isTruthy(vm.pop())
			left := isTruthy(vm.pop())
			vm.push(object.Boolean(left && right))
		case OpOr:
			right := isTruthy(vm.pop())
			left := isTruthy(vm.pop())
			vm.push(object.Boolean(left || right))
		case OpEqual:
			right := vm.pop()
			left := vm.pop()
			vm.push(object.Boolean(left == right))
		case OpNotEqual:
			right := vm.pop()
			left := vm.pop()
			vm.push(object.Boolean(left != right))
		case OpGt:
			right := vm.pop().(object.Integer)
			left := vm.pop().(object.Integer)
			vm.push(object.Boolean(left > right))
		case OpBang:
			operand := vm.pop()
			if isTruthy(operand) {
				vm.push(False)
			} else {
				vm.push(True)
			}
		case OpMinus:
			vm.push(object.Integer(-vm.pop().(object.Integer)))
		case OpJump:
			pos := int(caller.readInsOprandUint16())
			caller.pc = pos //jump to pos
		case OpJumpNotTruthy:
			pos := int(caller.readInsOprandUint16())
			if condition := vm.pop(); !isTruthy(condition) {
				caller.pc = pos //jump to pos
			}
		case OpNull:
			vm.push(NULL)
		case OpSetGlobal:
			idx := int(caller.readInsOprandUint16())
			vm.globals[idx] = vm.pop()
		case OpGetGlobal:
			idx := int(caller.readInsOprandUint16())
			vm.push(vm.globals[idx])
		case OpSetLocal:
			idx := int(caller.readInsOprandUint8())
			vm.stack[caller.basePointer+idx] = vm.pop()
		case OpGetLocal:
			idx := int(caller.readInsOprandUint8())
			vm.push(vm.stack[caller.basePointer+idx])
		case OpArray:
			numElements := int(caller.readInsOprandUint16())
			arr := make([]object.Object, 0, numElements)
			arr = append(arr, vm.stack[vm.sp-numElements:vm.sp]...)
			vm.sp = vm.sp - numElements
			vm.push(&object.Array{Elements: arr})
		case OpIndex:
			index := vm.pop()
			left := vm.pop()
			switch {
			case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
				arr := left.(*object.Array)
				index := index.(object.Integer)
				if index < 0 || int(index) >= len(arr.Elements) {
					return fmt.Errorf("array index out of bounds: %d", index)
				}
				vm.push(arr.Elements[index])
			// TODO: implement OpIndex for map
			default:
				return fmt.Errorf("index operator not supported: %s", left.Type())
			}
		case OpCall:
			numArgs := int(caller.readInsOprandUint8())
			fn := vm.stack[vm.sp-numArgs-1]
			switch fn := fn.(type) {
			case *object.CompiledFunction:
				if numArgs != fn.NumParameters {
					return fmt.Errorf("wrong number of arguments: want=%d, got=%d", fn.NumParameters, numArgs)
				}
				callee := NewFrame(fn, vm.sp-numArgs)
				vm.frames = append(vm.frames, callee)
				vm.sp = callee.basePointer + fn.NumLocals
			case *object.Builtin:
				args := vm.stack[vm.sp-numArgs : vm.sp]
				result := fn.Fn(args...)
				vm.sp = vm.sp - numArgs - 1
				vm.push(result)
			default:
				return fmt.Errorf("calling non-function: %T", vm.stack[vm.sp-1])
			}
		case OpReturnValue:
			returnValue := vm.pop()
			vm.sp = vm.frames[len(vm.frames)-1].basePointer - 1
			vm.frames = vm.frames[:len(vm.frames)-1]
			vm.push(returnValue)
		case OpReturn:
			vm.sp = vm.frames[len(vm.frames)-1].basePointer - 1
			vm.frames = vm.frames[:len(vm.frames)-1]
		case OpGetBuiltin:
			idx := int(caller.readInsOprandUint8())
			vm.push(object.Builtins[idx].Builtin)
		default:
			return fmt.Errorf("unsupported opcode: %d", op)
		}
	}
	return nil
}
func (vm *VM) push(o object.Object) {
	if vm.sp >= len(vm.stack) {
		vm.stack = append(vm.stack, make([]object.Object, StackSize)...)
	}
	vm.stack[vm.sp] = o
	vm.sp++
}
func (vm *VM) pop() object.Object {
	vm.sp--
	return vm.stack[vm.sp]
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case object.Boolean:
		return bool(obj)
	case object.Null:
		return false
	default:
		return true
	}
}
func add(left, right object.Object) (object.Object, error) {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return left.(object.Integer) + right.(object.Integer), nil
	} else if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return left.(object.String) + right.(object.String), nil
	} else {
		return NULL, fmt.Errorf("type mismatch: %s + %s", left.Type(), right.Type())
	}
}
