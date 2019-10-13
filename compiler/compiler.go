package compiler

import (
	"fmt"
	"os"

	"github.com/ShrimpF/jackcompiler/tokenizer"
)

// Base -- main struct
type Base struct {
	Tokenizer *tokenizer.Base
	Output    *os.File
}

// New -- create compiler's base struct
func New(t *tokenizer.Base, output *os.File) *Base {
	return &Base{Tokenizer: t, Output: output}
}

// Start -- start compiling
func (b *Base) Start() {
	b.compileClass()
}

// getToken -- get current token. shorten Tokenizer.GetCurrToken()
func (b *Base) getToken() *tokenizer.Token {
	return b.Tokenizer.GetCurrToken()
}

// write -- general use for write string
func (b *Base) write(value string) {
	fmt.Fprintln(b.Output, value)
}

// writeOpenTag -- write <XXX>
func (b *Base) writeOpenTag(value interface{}) {
	b.write(fmt.Sprintf("<%v>", value))
}

// writeCloseTag -- write </XXX>
func (b *Base) writeCloseTag(value interface{}) {
	b.write(fmt.Sprintf("</%v>", value))
}

// writeTerminal -- write <XXX>value</XXX>,then advace next token
func (b *Base) writeTerminal() {
	token := b.getToken()
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

	b.write(fmt.Sprintf("<%v> %v </%v>", token.Type(), value, token.Type()))
	b.Tokenizer.Advance()
}

// compileClass -- write class xml
func (b *Base) compileClass() {
	b.writeOpenTag("class")
	b.writeTerminal()

	b.writeTerminal() // write class-name
	b.writeTerminal() // write "{"

	b.compileClassVarDec()
	b.compileSubroutineDec()

	b.writeTerminal() // write "}"
	b.writeCloseTag("class")
}

// compileClassVarDec -- write class variables declaration
func (b *Base) compileClassVarDec() {
	if !isClassVarDec(b.getToken()) {
		return
	}
	b.writeOpenTag("classVarDec")
	b.writeTerminal() // "field" or "static"
	b.writeTerminal() // type
	b.writeTerminal() // variable name
	for b.getToken().Symbol() == "," {
		b.writeTerminal()
		b.writeTerminal()
	}
	b.writeTerminal() // ";"
	b.writeCloseTag("classVarDec")
	b.compileClassVarDec() // call next classVarDec
}

// compileSubroutineDec -- write func method constructot ...etc
func (b *Base) compileSubroutineDec() {
	if !isSubroutineDec(b.getToken()) {
		return
	}
	b.writeOpenTag("subroutineDec")
	b.writeTerminal() // constructor, function, method
	b.writeTerminal() // void , type
	b.writeTerminal() // subroutine name
	b.writeTerminal() // (
	b.compileParameterList()
	b.writeTerminal() // )
	b.compileSubroutineBody()
	b.writeCloseTag("subroutineDec")
	b.compileSubroutineDec()
}

// compileParameterList -- write parameter list like (int x,int y)
func (b *Base) compileParameterList() {
	b.writeOpenTag("parameterList")
	if b.getToken().Symbol() != ")" {
		b.writeTerminal() // type
		b.writeTerminal() // varName
		for b.getToken().Symbol() == "," {
			b.writeTerminal() // ,
			b.writeTerminal() // type
			b.writeTerminal() // varName
		}
	}
	b.writeCloseTag("parameterList")
}

// compileSubroutineBody -- write subroutine body
func (b *Base) compileSubroutineBody() {
	b.writeOpenTag("subroutineBody")
	b.writeTerminal() // {

	b.compileVarDec()
	b.compileStatements()

	b.writeTerminal() // }
	b.writeCloseTag("subroutineBody")
}

// compileVarDec -- write variable declearation
func (b *Base) compileVarDec() {
	if b.getToken().KeywordType() != tokenizer.Var {
		return
	}
	b.writeOpenTag("varDec")
	b.writeTerminal() // var
	b.writeTerminal() // type
	b.writeTerminal() // var name
	for b.getToken().Symbol() == "," {
		b.writeTerminal() // ,
		b.writeTerminal() // var name
	}
	b.writeTerminal() // ;
	b.writeCloseTag("varDec")
	b.compileVarDec()
}

// compileStatements -- write statements like let if while do return
func (b *Base) compileStatements() {
	b.writeOpenTag("statements")
STATEMENTS_LOOP:
	for {
		switch b.getToken().KeywordType() {
		case tokenizer.Let:
			b.compileLet()
		case tokenizer.If:
			b.compileIf()
		case tokenizer.While:
			b.compileWhile()
		case tokenizer.Do:
			b.compileDo()
		case tokenizer.Return:
			b.compileReturn()
		default:
			break STATEMENTS_LOOP
		}
	}
	b.writeCloseTag("statements")
}

// compileLet -- write let statements
func (b *Base) compileLet() {
	b.writeOpenTag("letStatement")
	b.writeTerminal() // let
	b.writeTerminal() // varName
	if b.getToken().Symbol() == "[" {
		b.writeTerminal() // [
		b.compileExpression()
		b.writeTerminal() // ]
	}
	b.writeTerminal() // =
	b.compileExpression()
	b.writeTerminal() // ;
	b.writeCloseTag("letStatement")
}

// compileIf -- write if statements
func (b *Base) compileIf() {
	b.writeOpenTag("ifStatement")
	b.writeTerminal() // if
	b.writeTerminal() // (
	b.compileExpression()
	b.writeTerminal() // )
	b.writeTerminal() // {
	b.compileStatements()
	b.writeTerminal() // }
	if b.getToken().KeywordType() == tokenizer.Else {
		b.writeTerminal() // else
		b.writeTerminal() // {
		b.compileStatements()
		b.writeTerminal() // }
	}
	b.writeCloseTag("ifStatement")
}

// compileWhile -- write while statements
func (b *Base) compileWhile() {
	b.writeOpenTag("whileStatement")
	b.writeTerminal() // while
	b.writeTerminal() // (
	b.compileExpression()
	b.writeTerminal() // )
	b.writeTerminal() // {
	b.compileStatements()
	b.writeTerminal() // }
	b.writeCloseTag("whileStatement")
}

// compileDo -- write do statementes
func (b *Base) compileDo() {
	b.writeOpenTag("doStatement")
	b.writeTerminal() // do
	b.writeTerminal() // subroutine name, class name , var name
	switch b.getToken().Symbol() {
	case "(":
		b.writeTerminal() // (
		b.compileExpressionList()
		b.writeTerminal() // )
	case ".":
		b.writeTerminal() // .
		b.writeTerminal() // subroutine name
		b.writeTerminal() // (
		b.compileExpressionList()
		b.writeTerminal() // )
	}
	b.writeTerminal() // ;
	b.writeCloseTag("doStatement")
}

// compileReturn -- write return statements
func (b *Base) compileReturn() {
	b.writeOpenTag("returnStatement")
	b.writeTerminal() // return
	if b.getToken().Symbol() != ";" {
		b.compileExpression()
	}
	b.writeTerminal() // ;
	b.writeCloseTag("returnStatement")
}

// compileExpression -- write expressioon
func (b *Base) compileExpression() {
	b.writeOpenTag("expression")
	b.compileTerm()
	for isOperand(b.getToken()) {
		b.writeTerminal() // write operand
		b.compileTerm()
	}
	b.writeCloseTag("expression")
}

// compileTerm -- write term
func (b *Base) compileTerm() {
	b.writeOpenTag("term")
	defer b.writeCloseTag("term")

	switch b.getToken().Symbol() {
	case "-", "~":
		b.writeTerminal() // - or ~
		b.compileTerm()
		return
	case "(":
		b.writeTerminal() // (
		b.compileExpression()
		b.writeTerminal() // )
		return
	default:
		b.writeTerminal()
	}

	switch b.getToken().Symbol() {
	case "[":
		b.writeTerminal() // [
		b.compileExpression()
		b.writeTerminal() // ]
		return
	// subroutine call
	case "(":
		b.writeTerminal() // (
		b.compileExpressionList()
		b.writeTerminal() // )
	case ".":
		b.writeTerminal() // .
		b.writeTerminal() // subrourine name
		b.writeTerminal() // (
		b.compileExpressionList()
		b.writeTerminal() // )
	}

}

// compileExpressionList -- compile a bunch of expressions
func (b *Base) compileExpressionList() {
	b.writeOpenTag("expressionList")
	if b.getToken().Symbol() != ")" {
		b.compileExpression()
		for b.getToken().Symbol() == "," {
			b.writeTerminal()
			b.compileExpression()
		}
	}
	b.writeCloseTag("expressionList")
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
