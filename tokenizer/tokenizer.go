package tokenizer

import (
	"io/ioutil"
	"os"
	"regexp"
)

const (
	comment     = `(?m)(\/\/.*)|(?s)(\/\*.*?\*\/)`
	emptyLine   = `(?m)(^\n)`
	keyword     = `(?m)(^class$|^constructor$|^function$|^method$|^field$|^static$|^var$|^int$|^char$|^boolean$|^void$|^true$|^false$|^null$|^this$|^let$|^do$|^if$|^else$|^while$|^return$)`
	symbol      = `(?m)([{}()[\].,;+\-*/&|<>=~])`
	intConst    = `(?m)(\d+)`
	stringConst = `(?m)(\"[^\n]*\")`
	identifier  = `(?m)([A-Za-z_]\w*)`
)

// Tokenizer is tokenizer's base struct
type Tokenizer struct {
	FilePath string
	Tokens   []*Token
	Output   *os.File
	CurrIdx  int
}

// New -- create Tokenizer base struct
func New(filePath string) *Tokenizer {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// remove comments
	bytes = regexp.MustCompile(comment+`|`+emptyLine).ReplaceAll(bytes, []byte{})
	bytes = regexp.MustCompile(emptyLine).ReplaceAll(bytes, []byte{})

	// tokenize
	values := regexp.MustCompile(keyword+`|`+symbol+`|`+intConst+`|`+stringConst+`|`+identifier).FindAllString(string(bytes), -1)

	// push tokens to []Token
	var tokens []*Token
	for _, v := range values {
		tokens = append(tokens, NewToken(v))
	}

	return &Tokenizer{
		FilePath: filePath,
		Tokens:   tokens,
		CurrIdx:  0,
	}
}

// HasMoreTokens -- check whether token remains
func (t *Tokenizer) HasMoreTokens() bool {
	return t.CurrIdx < len(t.Tokens)
}

// Advance -- increment CurrIdx
func (t *Tokenizer) Advance() {
	if t.HasMoreTokens() {
		t.CurrIdx++
	}
}

// GetCurrToken -- get current token
func (t *Tokenizer) GetCurrToken() *Token {
	if t.CurrIdx >= len(t.Tokens) {
		panic("Out of Index")
	}
	return t.Tokens[t.CurrIdx]
}
