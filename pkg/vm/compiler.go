package vm

import (
	"fmt"

	"github.com/alwaifu/monkey/pkg/ast"
	"github.com/alwaifu/monkey/pkg/object"
)

type CompilationScope struct {
	instructions        Instructions
	lastInsPosition     int // position of last instruction
	previousInsPosition int // position of previous instruction
}

type Compiler struct {
	Constants []object.Object

	symbolTable *SymbolTable

	scopes     []*CompilationScope
	scopeIndex int
}

func NewCompiler(s *SymbolTable, constants []object.Object) *Compiler {
	mainScope := &CompilationScope{
		instructions:        make(Instructions, 0, 256),
		lastInsPosition:     0,
		previousInsPosition: 0,
	}
	c := &Compiler{
		symbolTable: NewSymbolTable(nil),
		scopes:      []*CompilationScope{mainScope},
		scopeIndex:  0,
	}
	if s == nil {
		for i, v := range object.Builtins {
			c.symbolTable.DefineBuiltin(i, v.Name)
		}
	} else {
		c.symbolTable = s
	}
	if constants != nil {
		c.Constants = constants
	}
	return c
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(OpPop)
	case *ast.InfixExpression:
		if node.Operator == "<" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			if err := c.Compile(node.Left); err != nil {
				return err
			}
			c.emit(OpGt)
			break
		}
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(OpAdd)
		case "-":
			c.emit(OpSub)
		case "*":
			c.emit(OpMul)
		case "/":
			c.emit(OpDiv)
		case ">":
			c.emit(OpGt)
		case "==":
			c.emit(OpEqual)
		case "!=":
			c.emit(OpNotEqual)
		case "and":
			c.emit(OpAnd)
		case "or":
			c.emit(OpOr)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.PrefixExpression:
		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(OpBang)
		case "-":
			c.emit(OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IfExpression:
		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		// FIXME: may be need to remove last pop for condition

		jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)
		if err := c.Compile(node.Consequence); err != nil {
			return err
		}
		c.removeLastPop()
		jumpPos := c.emit(OpJump, 9999)
		afterConsequencePos := len(c.scopes[c.scopeIndex].instructions)
		c.replaceInstruction(jumpNotTruthyPos, MakeInstruction(OpJumpNotTruthy, afterConsequencePos))

		if node.Alternative == nil {
			c.emit(OpNull)
		} else {
			if err := c.Compile(node.Alternative); err != nil {
				return err
			}
			c.removeLastPop()
		}
		afterAlternativePos := len(c.scopes[c.scopeIndex].instructions)
		c.replaceInstruction(jumpPos, MakeInstruction(OpJump, afterAlternativePos))
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.IntegerLiteral:
		constIdx := c.addConstant(object.Integer(node.Value))
		c.emit(OpConstant, constIdx)
	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(OpTrue)
		} else {
			c.emit(OpFalse)
		}
	case *ast.StringLiteral:
		c.emit(OpConstant, c.addConstant(object.String(node.Value)))
	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			if err := c.Compile(el); err != nil {
				return err
			}
		}
		c.emit(OpArray, len(node.Elements))
	case *ast.FunctionLiteral:
		c.enterScope()
		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}
		if err := c.Compile(node.Body); err != nil {
			return err
		}
		c.replaceFunctionLastPopWithReturn()
		numLocals := len(c.symbolTable.store) // number of local var
		instructions := c.leaveScope()
		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}
		c.emit(OpConstant, c.addConstant(compiledFn))
	case *ast.ReturnStatement:
		if err := c.Compile(node.ReturnValue); err != nil {
			return err
		}
		c.emit(OpReturnValue)
	case *ast.CallExpression:
		// node.Function can be FunctionLiteral or Identifier, so there is no way to verify arguments at compiler
		if err := c.Compile(node.Function); err != nil {
			return err
		}
		for _, a := range node.Arguments {
			if err := c.Compile(a); err != nil {
				return err
			}
		}
		c.emit(OpCall, len(node.Arguments))
	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		if symbol.Scope == GlobalScope {
			c.emit(OpSetGlobal, symbol.Index)
		} else {
			c.emit(OpSetLocal, symbol.Index)
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		switch symbol.Scope {
		case GlobalScope:
			c.emit(OpGetGlobal, symbol.Index)
		case LocalScope:
			c.emit(OpGetLocal, symbol.Index)
		case BuiltinScope:
			c.emit(OpGetBuiltin, symbol.Index)
		default:
			return fmt.Errorf("undefined variable %s", node.Value)
		}
	case *ast.IndexExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(OpIndex)
	}
	return nil
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.Constants = append(c.Constants, obj)
	return len(c.Constants) - 1
}
func (c *Compiler) emit(op Opcode, operands ...int) int {
	scope := c.scopes[c.scopeIndex]
	ins := MakeInstruction(op, operands...)
	pos := len(scope.instructions)
	scope.instructions = append(scope.instructions, ins...)
	scope.previousInsPosition = scope.lastInsPosition
	scope.lastInsPosition = pos
	return pos // position of this instruction
}
func (c *Compiler) removeLastPop() {
	scope := c.scopes[c.scopeIndex]
	if scope.instructions[scope.lastInsPosition] == byte(OpPop) {
		scope.instructions = scope.instructions[:scope.lastInsPosition]
		scope.lastInsPosition = scope.previousInsPosition
	}
}
func (c *Compiler) replaceFunctionLastPopWithReturn() {
	scope := c.scopes[c.scopeIndex]
	if len(scope.instructions) > 0 {
		c.replaceInstruction(scope.lastInsPosition, MakeInstruction(OpReturnValue))
	} else {
		c.emit(OpReturn)
	}
}
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	scope := c.scopes[c.scopeIndex]
	for i := 0; i < len(newInstruction); i++ {
		scope.instructions[pos+i] = newInstruction[i]
	}
}
func (c *Compiler) enterScope() {
	c.scopeIndex++
	c.scopes = append(c.scopes, &CompilationScope{
		instructions:        make(Instructions, 0, 256),
		lastInsPosition:     0,
		previousInsPosition: 0,
	})
	c.symbolTable = NewSymbolTable(c.symbolTable)
}
func (c *Compiler) leaveScope() Instructions {
	ins := c.scopes[c.scopeIndex].instructions
	c.scopes = c.scopes[:c.scopeIndex]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.outer
	return ins
}
