package ast

import "monkey/pkg/lexer"

const (
	_ int = iota
	LOWEST
	OR          // ||
	AND         // &&
	EQUALS      // == or !=
	LESSGREATER // > or >= or < or <=
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

func precedence(token lexer.Token) int {
	switch token.Type {
	case lexer.OR:
		return OR
	case lexer.AND:
		return AND
	case lexer.EQ, lexer.NOT_EQ:
		return EQUALS
	case lexer.LT, lexer.LE, lexer.GT, lexer.GE:
		return LESSGREATER
	case lexer.PLUS, lexer.MINUS:
		return SUM
	case lexer.SLASH, lexer.ASTERISK:
		return PRODUCT
	case lexer.LPAREN:
		return CALL
	case lexer.LBRACKET:
		return INDEX
	default:
		return LOWEST
	}
}
