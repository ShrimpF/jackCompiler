package compiler

import (
	"fmt"
	"os"

	"github.com/ShrimpF/jackCompiler/tokenizer"
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

// GetToken -- get current token. shorten Tokenizer.GetCurrToken()
func (b *Base) GetToken() *tokenizer.Token {
	return b.Tokenizer.GetCurrToken()
}

// write -- general use for write string
func (b *Base) write(value string) {
	fmt.Fprintln(b.Output, value)
}

// WriteOpenTag -- write <XXX>
func (b *Base) WriteOpenTag(value interface{}) {
	b.write(fmt.Sprintf("<%v>", value))
}

// WriteCloseTag -- write </XXX>
func (b *Base) WriteCloseTag(value interface{}) {
	b.write(fmt.Sprintf("</%v>", value))
}

// WriteTerminal -- write <XXX>value</XXX>,then advace next token
func (b *Base) WriteTerminal() {
	token := b.GetToken()
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

	b.write(fmt.Sprintf("\t<%v>%v</%v>", token.Type(), value, token.Type()))
	b.Tokenizer.Advance()
}

// CompileClass -- write class xml
func (b *Base) CompileClass() {
	b.WriteOpenTag("class")
	b.WriteTerminal()

	b.WriteTerminal() // write class-name
	b.WriteTerminal() // write "{"

	b.CompileClassVarDec()
	b.CompileSubroutineDec()

	b.WriteTerminal() // write "}"
	b.WriteCloseTag("class")
}

// CompileClassVarDec -- write class variables declaration
func (b *Base) CompileClassVarDec() {
	if !isClassVarDec(b.GetToken()) {
		return
	}
	b.WriteOpenTag("classVarDec")
	b.WriteTerminal() // "field" or "static"
	b.WriteTerminal() // type
	b.WriteTerminal() // variable name
	for b.GetToken().Symbol() == "," {
		b.WriteTerminal()
		b.WriteTerminal()
	}
	b.WriteTerminal() // ";"
	b.WriteCloseTag("classVarDec")
	b.CompileClassVarDec() // call next classVarDec
}

// CompileSubroutineDec -- write func method constructot ...etc
func (b *Base) CompileSubroutineDec() {
	if !isSubroutineDec(b.GetToken()) {
		return
	}
	b.WriteOpenTag("subroutineDec")
	b.WriteTerminal() // constructor, function, method
	b.WriteTerminal() // void , type
	b.WriteTerminal() // subroutine name
	b.WriteTerminal() // (
	b.CompileParameterList()
	b.WriteTerminal() // )

	b.CompileSubroutineBody()

	b.WriteCloseTag("subroutineDec")
}

// CompileParameterList -- write parameter list like (int x,int y)
func (b *Base) CompileParameterList() {
	if b.GetToken().Symbol() == ")" {
		return
	}
	b.WriteOpenTag("parameterList")
	b.WriteTerminal() // type
	b.WriteTerminal() // varName
	for b.GetToken().Symbol() == "," {
		b.WriteTerminal() // ,
		b.WriteTerminal() // type
		b.WriteTerminal() // varName
	}
	b.WriteCloseTag("parameterList")
}

// CompileSubroutineBody -- write subroutine body
func (b *Base) CompileSubroutineBody() {
	b.WriteOpenTag("subroutineBody")
	b.WriteTerminal() // {

	b.CompileVarDec()
	b.CompileStatements()

	b.WriteTerminal() // }
	b.WriteCloseTag("subroutineBody")
}

// CompileVarDec -- write variable declearation
func (b *Base) CompileVarDec() {
	if b.GetToken().KeywordType() != tokenizer.Var {
		return
	}
	b.WriteOpenTag("varDec")
	b.WriteTerminal() // var
	b.WriteTerminal() // type
	b.WriteTerminal() // var name
	for b.GetToken().Symbol() == "," {
		b.WriteTerminal() // ,
		b.WriteTerminal() // var name
	}
	b.WriteCloseTag("varDec")
	b.CompileVarDec()
}

// CompileStatements -- write statements like let if while do return
func (b *Base) CompileStatements() {
	switch b.GetToken().KeywordType() {
	case tokenizer.Let:
		//TODO compile let statements
	case tokenizer.If:
		//TODO compile if statements
	case tokenizer.While:
		//TODO compile while statements
	case tokenizer.Do:
		//TODO compile do statements
	case tokenizer.Return:
		//TODO compile return statements
	default:
		return
	}
	b.CompileStatements()
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
