package compiler

import (
	"fmt"
	"os"

	"github.com/ShrimpF/jackCompiler/symboltable"
	"github.com/ShrimpF/jackCompiler/vmwriter"
	"github.com/ShrimpF/jackcompiler/tokenizer"
)

// Compiler -- main struct
type Compiler struct {
	Tokenizer   *tokenizer.Tokenizer
	symbolTable *symboltable.SymbolTable
	writer      *vmwriter.VMWriter
}

// map keywordtype => kind
var kind = map[string]string{
	"static": "static",
	"field":  "this",
	"arg":    "argumnet",
	"var":    "local",
}

// New -- create compiler's base struct
func New(path string, output *os.File) *Compiler {
	return &Compiler{
		Tokenizer:   tokenizer.New(path),
		symbolTable: symboltable.New(),
		writer:      vmwriter.New(output),
	}
}

// Start -- start compiling
func (c *Compiler) Start() {
	c.compileClass()
}

// getToken -- get current token. shorten Tokenizer.GetCurrToken()
func (c *Compiler) getToken() *tokenizer.Token {
	return c.Tokenizer.GetCurrToken()
}

// getTokenAndAdvance -- get current token and count up tokenizer.currIdx
func (c *Compiler) getTokenAndAdvance() *tokenizer.Token {
	defer c.Tokenizer.Advance()
	return c.Tokenizer.GetCurrToken()
}

// write -- general use for write string
func (c *Compiler) write(value string) {
	// fmt.Fprintln(c.Output, value)
}

// writeOpenTag -- write <XXX>
func (c *Compiler) writeOpenTag(value interface{}) {
	c.write(fmt.Sprintf("<%v>", value))
}

// writeCloseTag -- write </XXX>
func (c *Compiler) writeCloseTag(value interface{}) {
	c.write(fmt.Sprintf("</%v>", value))
}

// writeTerminal -- write <XXX>value</XXX>,then advace next token
func (c *Compiler) writeTerminal() {
	token := c.getToken()
	var value interface{}

	switch token.Type() {
	case tokenizer.Keyword:
		value = token.KeywordType()
	case tokenizer.Symbol:
		value = token.Symbol()
	case tokenizer.Identifier:
		value = token.Identifier()
	case tokenizer.IntConst:
		value = token.IntVal()
	case tokenizer.StringConst:
		value = token.StringVal()
	default:
		value = "undefined token"
	}

	c.write(fmt.Sprintf("<%v> %v </%v>", token.Type(), value, token.Type()))
	c.Tokenizer.Advance()
}

// compileClass -- write class xml
func (c *Compiler) compileClass() {
	c.Tokenizer.Advance() // class
	c.Tokenizer.Advance() // class-name
	c.Tokenizer.Advance() // {

	c.compileClassVarDec()
	// c.compileSubroutineDec()
	fmt.Println(c.symbolTable)
	c.Tokenizer.Advance() // write "}"
}

// compileClassVarDec -- write class variables declaration
func (c *Compiler) compileClassVarDec() {
	if !isClassVarDec(c.getToken()) {
		return
	}
	kind := c.getTokenAndAdvance().KeywordType().String()     // "field" or "static"
	typeName := c.getTokenAndAdvance().KeywordType().String() // type name
	name := c.getTokenAndAdvance().Identifier()               // name
	c.symbolTable.Define(name, typeName, kind)
	for c.getToken().Symbol() == "," {
		c.Tokenizer.Advance()                      // ,
		name = c.getTokenAndAdvance().Identifier() // name
		c.symbolTable.Define(name, typeName, kind)
	}
	c.Tokenizer.Advance()  // ";"
	c.compileClassVarDec() // call next classVarDec
}

// compileSubroutineDec -- write func method constructot ...etc
func (c *Compiler) compileSubroutineDec() {
	if !isSubroutineDec(c.getToken()) {
		return
	}
	c.writeOpenTag("subroutineDec")
	c.writeTerminal() // constructor, function, method
	c.writeTerminal() // void , type
	c.writeTerminal() // subroutine name
	c.writeTerminal() // (
	c.compileParameterList()
	c.writeTerminal() // )
	c.compileSubroutineBody()
	c.writeCloseTag("subroutineDec")
	c.compileSubroutineDec()
}

// compileParameterList -- write parameter list like (int x,int y)
func (c *Compiler) compileParameterList() {
	c.writeOpenTag("parameterList")
	if c.getToken().Symbol() != ")" {
		c.writeTerminal() // type
		c.writeTerminal() // varName
		for c.getToken().Symbol() == "," {
			c.writeTerminal() // ,
			c.writeTerminal() // type
			c.writeTerminal() // varName
		}
	}
	c.writeCloseTag("parameterList")
}

// compileSubroutineBody -- write subroutine body
func (c *Compiler) compileSubroutineBody() {
	c.writeOpenTag("subroutineBody")
	c.writeTerminal() // {

	c.compileVarDec()
	c.compileStatements()

	c.writeTerminal() // }
	c.writeCloseTag("subroutineBody")
}

// compileVarDec -- write variable declearation
func (c *Compiler) compileVarDec() {
	if c.getToken().KeywordType() != tokenizer.Var {
		return
	}
	c.writeOpenTag("varDec")
	c.writeTerminal() // var
	c.writeTerminal() // type
	c.writeTerminal() // var name
	for c.getToken().Symbol() == "," {
		c.writeTerminal() // ,
		c.writeTerminal() // var name
	}
	c.writeTerminal() // ;
	c.writeCloseTag("varDec")
	c.compileVarDec()
}

// compileStatements -- write statements like let if while do return
func (c *Compiler) compileStatements() {
	c.writeOpenTag("statements")
STATEMENTS_LOOP:
	for {
		switch c.getToken().KeywordType() {
		case tokenizer.Let:
			c.compileLet()
		case tokenizer.If:
			c.compileIf()
		case tokenizer.While:
			c.compileWhile()
		case tokenizer.Do:
			c.compileDo()
		case tokenizer.Return:
			c.compileReturn()
		default:
			break STATEMENTS_LOOP
		}
	}
	c.writeCloseTag("statements")
}

// compileLet -- write let statements
func (c *Compiler) compileLet() {
	c.writeOpenTag("letStatement")
	c.writeTerminal() // let
	c.writeTerminal() // varName
	if c.getToken().Symbol() == "[" {
		c.writeTerminal() // [
		c.compileExpression()
		c.writeTerminal() // ]
	}
	c.writeTerminal() // =
	c.compileExpression()
	c.writeTerminal() // ;
	c.writeCloseTag("letStatement")
}

// compileIf -- write if statements
func (c *Compiler) compileIf() {
	c.writeOpenTag("ifStatement")
	c.writeTerminal() // if
	c.writeTerminal() // (
	c.compileExpression()
	c.writeTerminal() // )
	c.writeTerminal() // {
	c.compileStatements()
	c.writeTerminal() // }
	if c.getToken().KeywordType() == tokenizer.Else {
		c.writeTerminal() // else
		c.writeTerminal() // {
		c.compileStatements()
		c.writeTerminal() // }
	}
	c.writeCloseTag("ifStatement")
}

// compileWhile -- write while statements
func (c *Compiler) compileWhile() {
	c.writeOpenTag("whileStatement")
	c.writeTerminal() // while
	c.writeTerminal() // (
	c.compileExpression()
	c.writeTerminal() // )
	c.writeTerminal() // {
	c.compileStatements()
	c.writeTerminal() // }
	c.writeCloseTag("whileStatement")
}

// compileDo -- write do statementes
func (c *Compiler) compileDo() {
	c.writeOpenTag("doStatement")
	c.writeTerminal() // do
	c.writeTerminal() // subroutine name, class name , var name
	switch c.getToken().Symbol() {
	case "(":
		c.writeTerminal() // (
		c.compileExpressionList()
		c.writeTerminal() // )
	case ".":
		c.writeTerminal() // .
		c.writeTerminal() // subroutine name
		c.writeTerminal() // (
		c.compileExpressionList()
		c.writeTerminal() // )
	}
	c.writeTerminal() // ;
	c.writeCloseTag("doStatement")
}

// compileReturn -- write return statements
func (c *Compiler) compileReturn() {
	c.writeOpenTag("returnStatement")
	c.writeTerminal() // return
	if c.getToken().Symbol() != ";" {
		c.compileExpression()
	}
	c.writeTerminal() // ;
	c.writeCloseTag("returnStatement")
}

// compileExpression -- write expressioon
func (c *Compiler) compileExpression() {
	c.writeOpenTag("expression")
	c.compileTerm()
	for isOperand(c.getToken()) {
		c.writeTerminal() // write operand
		c.compileTerm()
	}
	c.writeCloseTag("expression")
}

// compileTerm -- write term
func (c *Compiler) compileTerm() {
	c.writeOpenTag("term")
	defer c.writeCloseTag("term")

	switch c.getToken().Symbol() {
	case "-", "~":
		c.writeTerminal() // - or ~
		c.compileTerm()
		return
	case "(":
		c.writeTerminal() // (
		c.compileExpression()
		c.writeTerminal() // )
		return
	default:
		c.writeTerminal()
	}

	switch c.getToken().Symbol() {
	case "[":
		c.writeTerminal() // [
		c.compileExpression()
		c.writeTerminal() // ]
		return
	// subroutine call
	case "(":
		c.writeTerminal() // (
		c.compileExpressionList()
		c.writeTerminal() // )
	case ".":
		c.writeTerminal() // .
		c.writeTerminal() // subrourine name
		c.writeTerminal() // (
		c.compileExpressionList()
		c.writeTerminal() // )
	}

}

// compileExpressionList -- compile a bunch of expressions
func (c *Compiler) compileExpressionList() {
	c.writeOpenTag("expressionList")
	if c.getToken().Symbol() != ")" {
		c.compileExpression()
		for c.getToken().Symbol() == "," {
			c.writeTerminal()
			c.compileExpression()
		}
	}
	c.writeCloseTag("expressionList")
}

func isClassVarDec(t *tokenizer.Token) bool {
	switch t.KeywordType() {
	case tokenizer.Field, tokenizer.Static:
		return true
	default:
		return false
	}
}

func isSubroutineDec(t *tokenizer.Token) bool {
	switch t.KeywordType() {
	case tokenizer.Constructor, tokenizer.Function, tokenizer.Method:
		return true
	default:
		return false
	}
}

func isOperand(t *tokenizer.Token) bool {
	switch t.Symbol() {
	case "+", "-", "*", "/", "&amp;", "|", "&lt;", "&gt;", "=":
		return true
	default:
		return false
	}
}
