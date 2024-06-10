《用Go语言自制解释器》 读书笔记

基于《用Go语言自制解释器》实现解释器 直接解释 AST Tree 执行

# 词法分析
解释源代码，需要两次转换
![](https://res.weread.qq.com/wrepub/CB_3300018889_2.jpg)
第一次为词法分析
比如源码
```
let x = 5 + 5
```
词法生成结果，词法单元列表
```
[
    LET,
    IDENTIFIER("x"),
    EQUAL_SIGN,
    INTERGER(5),
    PLUS_SIGN,
    INTERGER(5),
    SEMICOLON,
]
```

定义词法单元数据结构, TODO: 将文件名和行号附加到词法单元中，便于追踪编译错误
```go
type Token struct {
	Type    TokenType
	Literal string
}
```
定义词法单元类型
```go
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
```
遍历源码字符串，挨个字节解析
* 判断字符为`+-*/<,;()[]{}`直接构造成对应token,
* 判断字符为`=!<>`额外读取下一个字符判断是否为`==`,`!=`,`<=`,`>=`
* 判断字符为`"`构造为 `STRING` TOKEN
* 判断字符为数字构造为 `INT` TOKEN
* 其他情况下构造为 `IDENT` 或者 内置标识符

# 语法分析
语法分析器将输入的内容转换成对应的数据结构(AST)。
类似JSON解析器将输入的文本构建成一个能表示这个输入的数据结构。
```python
> var input = '{"name": "Thorsten", "age": 28}';
> var output = JSON.parse(input);
> output
{ name: 'Thorsten', age: 28 }
> output.name
'Thorsten'
> output.age
28
>
```
JavaScript使用MagicLexer和MagicParser生成AST示例
```python
> var input = 'if (3 * 5 > 10) { return "hello"; } else { return "goodbye"; }';
> var tokens = MagicLexer.parse(input);
> MagicParser.parse(tokens);
{
  type: "if-statement",
  condition: {
    type: "operator-expression",
    operator: ">",
    left: {
      type: "operator-expression",
      operator: "*",
      left: { type: "integer-literal", value: 3 },
      right: { type: "integer-literal", value: 5 }
    },
    right: { type: "integer-literal", value: 10 }
  },
  consequence: {
    type: "return-statement",
    returnValue: { type: "string-literal", value: "hello" }
  },
  alternative: {
    type: "return-statement",
    returnValue: { type: "string-literal", value: "goodbye" }
  }
}
```
各个AST的实现都非常相似，概念上相同，但是细节有区别，没有一个通用的AST格式供所有语法分析器使用。这里会定义适用于Monkey语言的AST，并递归解析词法单元来构建这个AST

此处语法分析器是递归下降语法分析器。基于自上而下的运算符优先级分析法的语法分析器。因为发明人是沃恩·普拉特（Vaughan Pratt），所以有时它也称为普拉特语法分析器。

## 语句和表达式
表达式会产生值而语句不会。
### 解析语句
在Monkey语言中只有两种语句：let语句、return语句。其他都是表达式
```
let <标识符> = <表达式>;
return <表达式>;
```
![](https://res.weread.qq.com/wrepub/CB_3300018889_3.jpg)
移动词法单元指针并检查当前词法单元，根据当前词法单元来构造AST节点
### 解析表达式
难点：
* 运算符优先级
* 相同类型的词法单元可能出现在多个位置, 在表达式开头作为前缀运算符，在表达式中间作为中缀运算符

给每个词法单元绑定一个前缀解析函数和一个中缀解析函数
```go
type (
    prefixParseFn func() ast.Expression
    infixParseFn  func(ast.Expression) ast.Expression
)
```
定义前缀表达式节点
```go
type PrefixExpression struct {
    Token    token.Token // 前缀词法单元，如!
    Operator string
    Right    Expression
}
```
定义中缀表达式节点
```go
type InfixExpression struct {
    Token    token.Token // 运算符词法单元，如+
    Left     Expression
    Operator string
    Right    Expression
}
```
例如解析表达式 `1 + 2 + 3;` 最大的挑战不是在最终的AST中表示每个运算符和操作数，而是如何正确嵌套AST的节点。最后得到的应该是一个AST如下所示：
![](https://res.weread.qq.com/wrepub/CB_3300018889_4.jpg)

1. 首先检查是否有一个与当前`p.curToken`关联的 `prefixParseFn` 函数。
此处 token.INT 关联的前缀解析函数得到 `*ast.IntegerLiteral`
![](https://res.weread.qq.com/wrepub/CB_3300018889_5.jpg)
2. 循环判断`p.peekToken`优先级是否更高，如果`p.peekToken`优先级更高，调用`p.peekToken`的`infixParseFn`函数。
![](https://res.weread.qq.com/wrepub/CB_3300018889_6.jpg)
3. `infixParseFn`递归调用表达式解析出`*ast.IntegerLiteral`节点并生成`*ast.InfixExpression`节点返回
![](https://res.weread.qq.com/wrepub/CB_3300018889_8.jpg)
4. 前移词法单元
![](https://res.weread.qq.com/wrepub/CB_3300018889_7.jpg)
5. 并使用前面生成的`*ast.InfixExpression`作为做节点再次生成中缀表达式
![](https://res.weread.qq.com/wrepub/CB_3300018889_9.jpg)

定义优先级
```go
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
```
绑定前缀表达式生成函数
```go
    {
		lexer.IDENT:    p.parseIdentifier,
		lexer.INT:      p.parseIntegerLiteral,
		lexer.TRUE:     p.parseBooleanLiteral,
		lexer.FALSE:    p.parseBooleanLiteral,
		lexer.STRING:   p.parseStringLiteral,
		lexer.BANG:     p.parsePrefixExpression, //!
		lexer.MINUS:    p.parsePrefixExpression, //-
		lexer.LPAREN:   p.parseGroupedExpression,
		lexer.IF:       p.parseIfExpression,
		lexer.FUNCTION: p.parseFunctionLiteral,
		lexer.LBRACKET: p.parseArrayLiteral,
	}
```
绑定中缀表达式生成函数
```go
    {
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
```
# 解释 AST Tree 执行
如果不进行求值，那么类似1 + 2的表达式转换后也只是一组字符、一组词法单元或一个树结构，并没有含义。经过求值，1 + 2会得到3；5 > 1得到true；5 < 1得到false；而puts("Hello World!")则能输出一条众所周知的问候语。
```go
function eval(astNode) {
  if (astNode is integerliteral) {
    return astNode.integerValue
  } else if (astNode is booleanLiteral) {
    return astNode.booleanValue
  } else if (astNode is infixExpression) {
    leftEvaluated = eval(astNode.Left)
    rightEvaluated = eval(astNode.Right)
    if astNode.Operator == "+" {
      return leftEvaluated + rightEvaluated
    } else if ast.Operator == "-" {
      return leftEvaluated - rightEvaluated
    }
  }
}
```