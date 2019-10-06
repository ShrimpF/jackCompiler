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

// write -- general use for write string
func (base *Base) write(value string) {
	fmt.Fprintln(base.Output, value)
}

// WriteOpenTag -- write <XXX>
func (base *Base) WriteOpenTag(value interface{}) {
	base.write(fmt.Sprintf("<%v>", value))
}

// WriteCloseTag -- write </XXX>
func (base *Base) WriteCloseTag(value interface{}) {
	base.write(fmt.Sprintf("</%v>", value))
}

// WriteTerminal -- write <XXX>value</XXX>
func (base *Base) WriteTerminal() {
	token := base.Tokenizer.GetCurrToken()
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

	base.write(fmt.Sprintf("\t<%v>%v</%v>", token.Type(), value, token.Type()))
	base.Tokenizer.Advance()
}

// CompileClass -- write class xml
func (base *Base) CompileClass() {
	base.WriteOpenTag(tokenizer.Class.String())
	defer base.WriteCloseTag(tokenizer.Class.String())
	base.Tokenizer.Advance()

	base.WriteTerminal() // write class-name
	base.WriteTerminal() // write "{"

	//TODO call CompileClassVarDec
	//TODO call CompileSubroutineDex

	base.WriteTerminal() // write "}"
}

// CompileClassVarDec -- write class variables declaration
func (base *Base) CompileClassVarDec() {
	keywordType := base.Tokenizer.GetCurrToken().KeywordType()
	if keywordType != tokenizer.Static && keywordType != tokenizer.Field {
		return
	}
}
