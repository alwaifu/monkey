package vm

import (
	"bufio"
	"fmt"
	"io"

	"github.com/alwaifu/monkey/pkg/ast"
	"github.com/alwaifu/monkey/pkg/lexer"
	"github.com/alwaifu/monkey/pkg/object"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	constants := []object.Object{}
	globals := make([]object.Object, GlobalSize)
	symbolTable := NewSymbolTable(nil)
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Fprint(out, "> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := ast.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			continue
		}
		compiler := NewCompiler(symbolTable, constants)
		if err := compiler.Compile(program); err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}
		machine := NewVM(compiler, globals)
		if err := machine.Run(); err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}
		result := machine.stack[machine.sp]
		_, _ = io.WriteString(out, result.Inspect())
		_, _ = io.WriteString(out, "\n")
	}
}
func printErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\t")
	}
}
