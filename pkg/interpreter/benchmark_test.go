package interpreter

import (
	"testing"

	"github.com/alwaifu/monkey/pkg/ast"
	"github.com/alwaifu/monkey/pkg/lexer"
	"github.com/alwaifu/monkey/pkg/object"
)

func BenchmarkEvalIntegerExpression(b *testing.B) {
	tests := []struct {
		env      map[string]interface{}
		input    string
		expected interface{}
	}{
		{
			env:      map[string]interface{}{"a": 1, "b": 0, "c": 1, "d": 0, "e": 1, "f": 0, "g": 1, "h": 0, "i": 1, "j": 0, "k": 1, "l": 0, "m": 1, "n": 0, "o": 1, "p": 0, "q": 1, "r": 0, "s": 1, "t": 0, "u": 1, "v": 0, "w": 1, "x": 0, "y": 1, "z": 0},
			input:    "a!=1 or b<0 or c==-1 or d!=1 or e!=1 or f!=1 or g!=1 or h!=1 or i!=1 or j!=1 or k!=1 or l!=1 or m!=1 or n!=1 or o!=1 or p!=1 or q!=1 or r!=1 or s!=1 or t!=1 or u!=1 or v!=1 or w!=1 or x!=1 or y!=1 or z!=1",
			expected: true,
		},
	}
	for _, tt := range tests {
		program := ast.NewParser(lexer.NewLexer(tt.input)).ParseProgram()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			env, err := object.NewEnviromentFromMap(tt.env)
			if err != nil {
				b.FailNow()
			}
			evaluated, err := object.ToGoValue(Eval(program, env))
			if err != nil {
				b.FailNow()
			}
			if interface{}(evaluated) != interface{}(tt.expected) {
				b.FailNow()
			}
		}
	}
}
