package vm

import (
	"fmt"
	"testing"

	"github.com/alwaifu/monkey/pkg/ast"
	"github.com/alwaifu/monkey/pkg/lexer"
	"github.com/alwaifu/monkey/pkg/object"
)

func TestCompileIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpAdd),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpPop),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpSub),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpMul),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpDiv),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpMinus),
				MakeInstruction(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
func TestComplileBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				MakeInstruction(OpTrue),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				MakeInstruction(OpFalse),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpGt),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpGt),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpEqual),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpNotEqual),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				MakeInstruction(OpTrue),
				MakeInstruction(OpFalse),
				MakeInstruction(OpEqual),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				MakeInstruction(OpTrue),
				MakeInstruction(OpFalse),
				MakeInstruction(OpNotEqual),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				MakeInstruction(OpTrue),
				MakeInstruction(OpBang),
				MakeInstruction(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
func TestCompileConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if (true) { 10 }; 3333;
			`,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []Instructions{
				// 0000
				MakeInstruction(OpTrue),
				// 0001
				MakeInstruction(OpJumpNotTruthy, 10),
				// 0004
				MakeInstruction(OpConstant, 0),
				// 0007
				MakeInstruction(OpJump, 11),
				// 0010
				MakeInstruction(OpNull),
				// 0011
				MakeInstruction(OpPop),
				// 0012
				MakeInstruction(OpConstant, 1),
				// 0015
				MakeInstruction(OpPop),
			},
		},
		{
			input: `
			if (true) { 10 } else { 20 }; 3333;
			`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []Instructions{
				// 0000
				MakeInstruction(OpTrue),
				// 0001
				MakeInstruction(OpJumpNotTruthy, 10),
				// 0004
				MakeInstruction(OpConstant, 0),
				// 0007
				MakeInstruction(OpJump, 13),
				// 0010
				MakeInstruction(OpConstant, 1),
				// 0013
				MakeInstruction(OpPop),
				// 0014
				MakeInstruction(OpConstant, 2),
				// 0017
				MakeInstruction(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
func TestCompileGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let one = 1;
			let two = 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpSetGlobal, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpSetGlobal, 1),
			},
		},
		{
			input: `
			let one = 1;
			one;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpSetGlobal, 0),
				MakeInstruction(OpGetGlobal, 0),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `
			let one = 1;
			let two = one;
			two;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpSetGlobal, 0),
				MakeInstruction(OpGetGlobal, 0),
				MakeInstruction(OpSetGlobal, 1),
				MakeInstruction(OpGetGlobal, 1),
				MakeInstruction(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
func TestCompileArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				MakeInstruction(OpArray, 0),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpConstant, 2),
				MakeInstruction(OpArray, 3),
				MakeInstruction(OpPop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpAdd),
				MakeInstruction(OpConstant, 2),
				MakeInstruction(OpConstant, 3),
				MakeInstruction(OpSub),
				MakeInstruction(OpConstant, 4),
				MakeInstruction(OpConstant, 5),
				MakeInstruction(OpMul),
				MakeInstruction(OpArray, 3),
				MakeInstruction(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
func TestCompileIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []interface{}{1, 2, 3, 1, 1},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpConstant, 2),
				MakeInstruction(OpArray, 3),
				MakeInstruction(OpConstant, 3),
				MakeInstruction(OpConstant, 4),
				MakeInstruction(OpAdd),
				MakeInstruction(OpIndex),
				MakeInstruction(OpPop),
			},
		},
		// {
		// 	input:             "{1: 2}[2 - 1]",
		// 	expectedConstants: []interface{}{1, 2, 2, 1},
		// 	expectedInstructions: []Instructions{
		// 		MakeInstruction(OpConstant, 0),
		// 		MakeInstruction(OpConstant, 1),
		// 		MakeInstruction(OpHash, 2),
		// 		MakeInstruction(OpConstant, 2),
		// 		MakeInstruction(OpConstant, 3),
		// 		MakeInstruction(OpSub),
		// 		MakeInstruction(OpIndex),
		// 		MakeInstruction(OpPop),
		// 	},
		// },
	}
	runCompilerTests(t, tests)
}
func TestCompileFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() {return 5+10}`,
			expectedConstants: []interface{}{
				5,
				10,
				[]Instructions{
					MakeInstruction(OpConstant, 0),
					MakeInstruction(OpConstant, 1),
					MakeInstruction(OpAdd),
					MakeInstruction(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 2),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `fn() {5+10}`,
			expectedConstants: []interface{}{
				5,
				10,
				[]Instructions{
					MakeInstruction(OpConstant, 0),
					MakeInstruction(OpConstant, 1),
					MakeInstruction(OpAdd),
					MakeInstruction(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 2),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `fn() {1; 2}`,
			expectedConstants: []interface{}{
				1,
				2,
				[]Instructions{
					MakeInstruction(OpConstant, 0),
					MakeInstruction(OpPop),
					MakeInstruction(OpConstant, 1),
					MakeInstruction(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 2),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `fn() {}`,
			expectedConstants: []interface{}{
				[]Instructions{
					MakeInstruction(OpReturn),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
func TestCompileFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { 24 }();`,
			expectedConstants: []interface{}{
				24,
				[]Instructions{
					MakeInstruction(OpConstant, 0),
					MakeInstruction(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpCall, 0),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `
			let noArg = fn() { 24 };
			noArg();`,
			expectedConstants: []interface{}{
				24,
				[]Instructions{
					MakeInstruction(OpConstant, 0),
					MakeInstruction(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpSetGlobal, 0),
				MakeInstruction(OpGetGlobal, 0),
				MakeInstruction(OpCall, 0),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `
				let oneArg = fn(a) { a };
				oneArg(24);`,
			expectedConstants: []interface{}{
				[]Instructions{
					MakeInstruction(OpGetLocal, 0),
					MakeInstruction(OpReturnValue),
				},
				24,
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpSetGlobal, 0),
				MakeInstruction(OpGetGlobal, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpCall, 1),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `
				let manyArg = fn(a, b, c) { a; b; c };
				manyArg(24, 25, 26);`,
			expectedConstants: []interface{}{
				[]Instructions{
					MakeInstruction(OpGetLocal, 0),
					MakeInstruction(OpPop),
					MakeInstruction(OpGetLocal, 1),
					MakeInstruction(OpPop),
					MakeInstruction(OpGetLocal, 2),
					MakeInstruction(OpReturnValue),
				},
				24,
				25,
				26,
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpSetGlobal, 0),
				MakeInstruction(OpGetGlobal, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpConstant, 2),
				MakeInstruction(OpConstant, 3),
				MakeInstruction(OpCall, 3),
				MakeInstruction(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
func TestCompileLetStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let num = 55;
			fn() { num }
			`,
			expectedConstants: []interface{}{
				55,
				[]Instructions{
					MakeInstruction(OpGetGlobal, 0),
					MakeInstruction(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 0),
				MakeInstruction(OpSetGlobal, 0),
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `
			fn() {
				let num = 55;
				num
			}
			`,
			expectedConstants: []interface{}{
				55,
				[]Instructions{
					MakeInstruction(OpConstant, 0),
					MakeInstruction(OpSetLocal, 0),
					MakeInstruction(OpGetLocal, 0),
					MakeInstruction(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 1),
				MakeInstruction(OpPop),
			},
		},
		{
			input: `
			fn() {
				let a = 55;
				let b = 77;
				a + b
			}
			`,
			expectedConstants: []interface{}{
				55,
				77,
				[]Instructions{
					MakeInstruction(OpConstant, 0),
					MakeInstruction(OpSetLocal, 0),
					MakeInstruction(OpConstant, 1),
					MakeInstruction(OpSetLocal, 1),
					MakeInstruction(OpGetLocal, 0),
					MakeInstruction(OpGetLocal, 1),
					MakeInstruction(OpAdd),
					MakeInstruction(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				MakeInstruction(OpConstant, 2),
				MakeInstruction(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []Instructions
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := ast.NewParser(l)
		program := p.ParseProgram()
		comp := NewCompiler(NewSymbolTable(nil), []object.Object{})
		if err := comp.Compile(program); err != nil {
			t.Fatalf("compiler error: %s \ninput: %s", err, tt.input)
		}

		if err := testInstructions(tt.expectedInstructions, comp.scopes[comp.scopeIndex].instructions); err != nil {
			t.Fatalf("testInstructions failed: %s \ninput: %s", err, tt.input)
		}

		if err := testConstants(tt.expectedConstants, comp.Constants); err != nil {
			t.Fatalf("testConstants failed: %s \ninput: %s", err, tt.input)
		}
	}
}

func testInstructions(expected []Instructions, actual []byte) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%x\ngot =%x", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%x\ngot =%x", i, concatted, actual)
		}
	}

	return nil
}

func concatInstructions(s []Instructions) Instructions {
	out := make([]byte, 0)
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testConstants(expectedList []interface{}, actualList []object.Object) error {
	if len(expectedList) != len(actualList) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d", len(actualList), len(expectedList))
	}
	for i, expected := range expectedList {
		want := fmt.Sprint(expected)
		if expected, ok := expected.([]Instructions); ok {
			want = fmt.Sprint(concatInstructions(expected))
		}
		got := fmt.Sprint(actualList[i])
		if f, ok := actualList[i].(*object.CompiledFunction); ok {
			got = fmt.Sprint(Instructions(f.Instructions))
		}
		if got != want {
			return fmt.Errorf("test constant %d failed- got: %s, want: %s", i, got, want)
		}
	}
	return nil
}
