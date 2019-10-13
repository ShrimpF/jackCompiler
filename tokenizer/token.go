package tokenizer

import (
	"fmt"
	"regexp"
	"strconv"
)

// Token --
type Token struct {
	value string
}

// NewToken -- create new token
func NewToken(value string) *Token {
	return &Token{value: value}
}

// Type -- return token type
func (t *Token) Type() TokenType {
	switch {
	case regexp.MustCompile(keyword).MatchString(t.value):
		return Keyword
	case regexp.MustCompile(symbol).MatchString(t.value):
		return Symbol
	case regexp.MustCompile(stringConst).MatchString(t.value):
		return StringConst
	case regexp.MustCompile(identifier).MatchString(t.value):
		return Identifier
	case regexp.MustCompile(intConst).MatchString(t.value):
		return IntConst
	default:
		return UndefinedToken
	}
}

// KeywordType -- return KeywordType
func (t *Token) KeywordType() KeywordType {
	if t.Type() != Keyword {
		return UndefinedKeyword
	}
	switch t.value {
	case "class":
		return Class
	case "constructor":
		return Constructor
	case "function":
		return Function
	case "method":
		return Method
	case "field":
		return Field
	case "static":
		return Static
	case "var":
		return Var
	case "int":
		return Int
	case "char":
		return Char
	case "boolean":
		return Boolean
	case "void":
		return Void
	case "true":
		return True
	case "false":
		return False
	case "null":
		return Null
	case "this":
		return This
	case "let":
		return Let
	case "do":
		return Do
	case "if":
		return If
	case "else":
		return Else
	case "while":
		return While
	case "return":
		return Return
	default:
		return UndefinedKeyword
	}
}

// Symbol return Symbol
func (t *Token) Symbol() string {
	if t.Type() != Symbol {
		return ""
	}

	switch t.value {
	case "<":
		return "&lt;"
	case ">":
		return "&gt;"
	case "&":
		return "&amp;"
	default:
		return t.value
	}
}

// Identifier return identifier
func (t *Token) Identifier() string {
	if t.Type() != Identifier {
		return ""
	}
	return t.value
}

// IntVal return intConst's value
func (t *Token) IntVal() int {
	val, err := strconv.Atoi(t.value)
	if err != nil || t.Type() != IntConst {
		fmt.Println(err)
		return -1
	}
	return val
}

// StringVal return stringConst's value
func (t *Token) StringVal() string {
	if t.Type() != StringConst {
		return ""
	}
	return regexp.MustCompile(`\"`).ReplaceAllString(t.value, "")
}
