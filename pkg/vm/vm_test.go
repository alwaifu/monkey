package vm

import (
	"fmt"
	"testing"

	"github.com/alwaifu/monkey/pkg/ast"
	"github.com/alwaifu/monkey/pkg/lexer"
	"github.com/alwaifu/monkey/pkg/object"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestRunIntegerArithmetic(t *testing.T) {
	testCases := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVmTests(t, testCases)
}
func TestRunBooleanExpressions(t *testing.T) {
	testCases := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
	}
	runVmTests(t, testCases)
}
func TestRunConditionals(t *testing.T) {
	testCases := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", NULL},
		{"if (false) { 10 }", NULL},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}
	runVmTests(t, testCases)
}
func TestRunStringExpressions(t *testing.T) {
	testCases := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}
	runVmTests(t, testCases)
}
func TestRunIndexExpressions(t *testing.T) {
	testCases := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		// {"[][0]", NULL},
		// {"[1, 2, 3][99]", NULL},
		// {"[1][-1]", NULL},

		// {"{1: 1, 2: 2}[1]", 1},
		// {"{1: 1, 2: 2}[2]", 2},
		// {"{1: 1}[0]", NULL},
		// {"{}[0]", NULL},
	}
	runVmTests(t, testCases)
}
func TestRunCallingFunctions(t *testing.T) {
	testCases := []vmTestCase{
		{"let fivePlusTen = fn() { 5 + 10; }; fivePlusTen();", 15},
		{"let one = fn() { 1; }; let two = fn() { 2; }; one() + two();", 3},
		{"let a = fn() { 1; }; let b = fn() { 2; }; let c = fn() { 3; }; a() + b() + c();", 6},
		{"let earlyExit = fn() {return 99; 100;}; earlyExit();", 99},
		{"let earlyExit = fn() {return 99;return 100;}; earlyExit();", 99},
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	runVmTests(t, testCases)
}
func TestRunBuiltinFunctions(t *testing.T) {
	testCases := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
	}
	runVmTests(t, testCases)
}

func runVmTests(t *testing.T, testCases []vmTestCase) {
	t.Helper()
	for _, tt := range testCases {
		defer func() {
			if v := recover(); v != nil {
				t.Fatal(v, "input:", tt.input)
			}
		}()
		l := lexer.NewLexer(tt.input)
		p := ast.NewParser(l)
		program := p.ParseProgram()
		comp := NewCompiler(nil, []object.Object{})
		if err := comp.Compile(program); err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		vm := NewVM(comp, make([]object.Object, GlobalSize))
		if err := vm.Run(); err != nil {
			t.Fatalf("vm error: %s", err)
		}
		result := vm.stack[vm.sp]
		want := fmt.Sprint(tt.expected)
		got := fmt.Sprint(result)
		if got != want {
			t.Fatal("test fatal, input:'", tt.input, "', got", got, "want", want)
		}

	}
}
