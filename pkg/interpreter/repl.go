package interpreter

import (
	"bufio"
	"fmt"
	"io"

	"github.com/alwaifu/monkey/pkg/ast"
	"github.com/alwaifu/monkey/pkg/lexer"
	"github.com/alwaifu/monkey/pkg/object"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnviroment()
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.NewLexer(line)
		// for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
		// 	fmt.Fprintf(out, "%v\n", tok)
		// }
		p := ast.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			continue
		}
		if result := Eval(program, env); result != nil {
			_, _ = io.WriteString(out, result.Inspect())
			_, _ = io.WriteString(out, "\n")
		}

	}
}

func printErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\t")
	}
}
