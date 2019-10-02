package tokenizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const (
	comment     = `(?m)(^\/\/.*\n)|(?m)(^\/\*.*\*\/)`
	emptyLine   = `(?m)(^\n)`
	keyword     = `(?m)(^class$|^constructor$|^function$|^method$|^field$|^static$|^var$|^int$|^char$|^boolean$|^void$|^true$|^false$|^null$|^this$|^let$|^do$|^if$|^else$|^while$|^return$)`
	symbol      = `(?m)([{}()[\].,;+\-*/&|<>=~])`
	intConst    = `(?m)(\d+)`
	stringConst = `(?m)(\"[^\n]*\")`
	identifier  = `(?m)([A-Za-z_]\w*)`
)

// Tokenizer --
type Tokenizer struct {
	FilePath string
	Tokens   []*Token
	Output   *os.File
	currIdx  int
}

// NewTokenizer -- create Tokenizer struct
func NewTokenizer(filePath string) *Tokenizer {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// remove comments
	bytes = regexp.MustCompile(comment+`|`+emptyLine).ReplaceAll(bytes, []byte{})
	// tokenize
	values := regexp.MustCompile(keyword+`|`+symbol+`|`+intConst+`|`+stringConst+`|`+identifier).FindAllString(string(bytes), -1)

	// push tokens to []Token
	var tokens []*Token
	for _, v := range values {
		tokens = append(tokens, NewToken(v))
	}

	// create file
	filename := regexp.MustCompile(`.jack$`).ReplaceAllString(filepath.Base(filePath), "")
	file, err := os.Create(filename + "T.xml")
	if err != nil {
		panic(err)
	}

	return &Tokenizer{
		FilePath: filePath,
		Tokens:   tokens,
		Output:   file,
		currIdx:  0,
	}
}

// HasMoreToken -- check whether token remains
func (t *Tokenizer) HasMoreToken() bool {
	return t.currIdx < len(t.Tokens)
}

// Advance -- increment CurrIdx
func (t *Tokenizer) Advance() {
	if t.currIdx+1 >= len(t.Tokens) {
		panic("No more token")
	}
	t.currIdx++
}

// GetCurrToken -- get current token
func (t *Tokenizer) GetCurrToken() *Token {
	if t.currIdx >= len(t.Tokens) {
		panic("Out of Index")
	}
	return t.Tokens[t.currIdx]
}

// WriteXML -- write -T.xml file
func (t *Tokenizer) WriteXML() {
	file, err := os.OpenFile(t.Output.Name(), os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// open and close token tag
	fmt.Fprintln(file, "<tokens>")
	defer fmt.Fprintln(file, "</tokens>")

	// write each token
	for _, token := range t.Tokens {
		var value interface{}
		switch token.Type() {
		case Keyword:
			value = token.KeywordType()
		case Symbol:
			value = token.Symbol()
		case Identifier:
			value = token.Identifier()
		case IntConst:
			value = token.IntVal()
		case StringConst:
			value = token.StringVal()
		default:
			value = "undefined token"
		}
		fmt.Fprintf(file, "\t<%v>%v</%v>\n", token.Type(), value, token.Type())
	}
}
