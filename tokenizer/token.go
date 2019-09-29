package tokenizer

import (
	"fmt"
	"regexp"
	"strconv"
)

// Token --
type Token struct {
	Value string
}

// NewToken -- create new token
func NewToken(word string) *Token {
	return &Token{Value: word}
}

// Type -- return token type
func (t *Token) Type() TokenType {
	switch {
	case regexp.MustCompile(keyword).MatchString(t.Value):
		return Keyword
	case regexp.MustCompile(symbol).MatchString(t.Value):
		return Symbol
	case regexp.MustCompile(stringConst).MatchString(t.Value):
		return StringConst
	case regexp.MustCompile(identifier).MatchString(t.Value):
		return Identifier
	case regexp.MustCompile(intConst).MatchString(t.Value):
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
	switch t.Value {
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

	switch t.Value {
	case "<":
		return "&lt;"
	case ">":
		return "&gt;"
	case "&":
		return "&amp;"
	default:
		return t.Value
	}
}

// Identifier return identifier
func (t *Token) Identifier() string {
	if t.Type() != Identifier {
		return ""
	}
	return t.Value
}

// IntVal return intConst's value
func (t *Token) IntVal() int {
	val, err := strconv.Atoi(t.Value)
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
	return regexp.MustCompile(`\"`).ReplaceAllString(t.Value, "")
}
