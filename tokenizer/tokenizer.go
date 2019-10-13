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

// Base is tokenizer's base struct
type Base struct {
	FilePath string
	Tokens   []*Token
	Output   *os.File
	CurrIdx  int
}

// NewTokenizer -- create Tokenizer base struct
func NewTokenizer(filePath string) *Base {
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

	return &Base{
		FilePath: filePath,
		Tokens:   tokens,
		CurrIdx:  0,
	}
}

// HasMoreTokens -- check whether token remains
func (base *Base) HasMoreTokens() bool {
	return base.CurrIdx < len(base.Tokens)
}

// Advance -- increment CurrIdx
func (base *Base) Advance() {
	if base.HasMoreTokens() {
		base.CurrIdx++
	}
}

// GetCurrToken -- get current token
func (base *Base) GetCurrToken() *Token {
	if base.CurrIdx >= len(base.Tokens) {
		panic("Out of Index")
	}
	return base.Tokens[base.CurrIdx]
}
