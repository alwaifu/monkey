package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"monkey/pkg/lexer"
)

type Node interface {
	TokenLiteral() string // 输出节点字面量 仅用于调试测试
	String() string       // to debug
}
type Statement interface {
	Node
	statementNode() // 占位方法 用于编译
}
type Expression interface {
	Node
	expressionNode() // 占位方法 用于编译
}

// ---

var _ Node = (*Program)(nil)

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ---

var (
	_ Node      = (*LetStatement)(nil)
	_ Statement = (*LetStatement)(nil)
)

type LetStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ---

var (
	_ Node      = (*ReturnStatement)(nil)
	_ Statement = (*ReturnStatement)(nil)
)

type ReturnStatement struct {
	Token       lexer.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ---

var (
	_ Node      = (*ExpressionStatement)(nil)
	_ Statement = (*ExpressionStatement)(nil)
)

type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// ---

var (
	_ Node       = (*Identifier)(nil)
	_ Expression = (*Identifier)(nil)
)

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// ---

var (
	_ Node       = (*IntegerLiteral)(nil)
	_ Expression = (*IntegerLiteral)(nil)
)

type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// ---

var (
	_ Node       = (*Boolean)(nil)
	_ Expression = (*Boolean)(nil)
)

type Boolean struct {
	Token lexer.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// ---

var (
	_ Node       = (*StringLiteral)(nil)
	_ Expression = (*StringLiteral)(nil)
)

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// ---

var (
	_ Node       = (*ArrayLiteral)(nil)
	_ Expression = (*ArrayLiteral)(nil)
)

type ArrayLiteral struct {
	Token    lexer.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// ---

var (
	_ Node       = (*IndexExpression)(nil)
	_ Expression = (*IndexExpression)(nil)
)

type IndexExpression struct {
	Token lexer.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// ---

var (
	_ Node       = (*PrefixExpression)(nil)
	_ Expression = (*PrefixExpression)(nil)
)

type PrefixExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// ---

var (
	_ Node       = (*InfixExpression)(nil)
	_ Expression = (*InfixExpression)(nil)
)

type InfixExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// ---

var (
	_ Node       = (*IfExpression)(nil)
	_ Expression = (*IfExpression)(nil)
)

type IfExpression struct {
	Token       lexer.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// ---

var (
	_ Node      = (*BlockStatement)(nil)
	_ Statement = (*BlockStatement)(nil)
)

type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ---

var (
	_ Node       = (*FunctionLiteral)(nil)
	_ Expression = (*FunctionLiteral)(nil)
)

type FunctionLiteral struct {
	Token      lexer.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// ---

var (
	_ Node       = (*CallExpression)(nil)
	_ Expression = (*CallExpression)(nil)
)

type CallExpression struct {
	Token     lexer.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// ---

type (
	prefixParseFn func() Expression
	infixParseFns map[lexer.TokenType]func(Expression) Expression
)

// Parser 语法解析器
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]func() Expression
	infixParseFns  map[lexer.TokenType]func(Expression) Expression
}

func NewParser(l *lexer.Lexer) *Parser {
	first, second := l.NextToken(), l.NextToken() // read two tokens, so curToken and peekToken are both set
	p := &Parser{
		l:         l,
		errors:    []string{},
		curToken:  first,
		peekToken: second,
	}
	// 注册解析函数
	p.prefixParseFns = map[lexer.TokenType]func() Expression{
		lexer.IDENT:    p.parseIdentifier,
		lexer.INT:      p.parseIntegerLiteral,
		lexer.TRUE:     p.parseBoolean,
		lexer.FALSE:    p.parseBoolean,
		lexer.STRING:   p.parseStringLiteral,
		lexer.BANG:     p.parsePrefixExpression,
		lexer.MINUS:    p.parsePrefixExpression,
		lexer.LPAREN:   p.parseGroupedExpression,
		lexer.IF:       p.parseIfExpression,
		lexer.FUNCTION: p.parseFunctionLiteral,
		lexer.LBRACKET: p.parseArrayLiteral,
	}
	p.infixParseFns = map[lexer.TokenType]func(Expression) Expression{
		lexer.PLUS:     p.parseInfixExpression,
		lexer.MINUS:    p.parseInfixExpression,
		lexer.SLASH:    p.parseInfixExpression,
		lexer.ASTERISK: p.parseInfixExpression,
		lexer.EQ:       p.parseInfixExpression,
		lexer.NOT_EQ:   p.parseInfixExpression,
		lexer.LT:       p.parseInfixExpression,
		lexer.LE:       p.parseInfixExpression,
		lexer.GT:       p.parseInfixExpression,
		lexer.GE:       p.parseInfixExpression,
		lexer.AND:      p.parseInfixExpression,
		lexer.OR:       p.parseInfixExpression,

		lexer.LPAREN:   p.parseCallExpression,
		lexer.LBRACKET: p.parseIndexExpression,
	}
	return p
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}
	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		stmt := &ExpressionStatement{Token: p.curToken}
		stmt.Expression = p.parseExpression(LOWEST)
		if p.peekToken.Type == (lexer.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	}
}
func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekToken.Type == (lexer.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}
func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekToken.Type == (lexer.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}
func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}
func (p *Parser) parseIntegerLiteral() Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &IntegerLiteral{Token: p.curToken, Value: value}
}
func (p *Parser) parseBoolean() Expression {
	return &Boolean{Token: p.curToken, Value: p.curToken.Type == lexer.TRUE}
}
func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}
func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(lexer.RBRACKET)
	return array
}
func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}
	return exp
}
func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}
func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	curPrecedence := precedence(p.curToken)
	p.nextToken()
	expression.Right = p.parseExpression(curPrecedence)
	return expression
}
func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}
	return exp
}
func (p *Parser) parseIfExpression() Expression {
	expression := &IfExpression{Token: p.curToken}
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	if p.peekToken.Type == lexer.ELSE {
		p.nextToken()
		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}
func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}
	p.nextToken()
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}
func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}
func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}
func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}
	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}
	return identifiers
}
func (p *Parser) parseExpressionList(end lexer.TokenType) []Expression {
	list := []Expression{}
	if p.peekToken.Type == end {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}
func (p *Parser) parseExpression(curPrecedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type))
		return nil
	}
	leftExp := prefix()
	nextPrecedence := precedence(p.peekToken)
	for p.peekToken.Type != lexer.SEMICOLON && curPrecedence < nextPrecedence {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// expectPeek 断言函数 只有符合预期才会继续执行
func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == (t) {
		p.nextToken()
		return true
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
		return false
	}
}
func (p *Parser) Errors() []string {
	return p.errors
}
