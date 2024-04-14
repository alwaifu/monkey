package lexer

type Token struct {
	Type    TokenType
	Literal string
}

type TokenType = string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"  // 标识符
	INT    = "INT"    // int字面量
	STRING = "STRING" // string字面量

	ASSIGN   = "ASSIGN"   // =
	PLUS     = "PLUS"     // +
	MINUS    = "MINUS"    // -
	BANG     = "BANG"     // !
	ASTERISK = "ASTERISK" // *
	SLASH    = "SLASH"    // /
	LT       = "LT"       // <
	LE       = "LE"       // <=
	GT       = "GT"       // >
	GE       = "GE"       // >=
	EQ       = "EQ"       // ==
	NOT_EQ   = "NOT_EQ"   // !=
	AND      = "AND"      // and
	OR       = "OR"       // or

	COMMA     = "," // ,
	SEMICOLON = ";" // ;

	LPAREN   = "(" // (
	RPAREN   = ")" // )
	LBRACE   = "{" // {
	RBRACE   = "}" // }
	LBRACKET = "[" // [
	RBRACKET = "]" // ]

	FUNCTION = "FUNCTION" // function
	LET      = "LET"      // let
	TRUE     = "TRUE"     // true
	FALSE    = "FALSE"    // false
	IF       = "IF"       // if
	ELSE     = "ELSE"     // else
	RETURN   = "RETURN"   // return

)

func LookupIdent(ident string) TokenType {
	switch ident {
	case "fn":
		return FUNCTION
	case "let":
		return LET
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "if":
		return IF
	case "else":
		return ELSE
	case "return":
		return RETURN
	case "and":
		return AND
	case "or":
		return OR
	default:
		return IDENT
	}
}
