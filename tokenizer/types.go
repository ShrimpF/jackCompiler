package tokenizer

// TokenType --
type TokenType int

// Token type enum
const (
	Keyword TokenType = iota + 1
	Symbol
	Identifier
	IntConst
	StringConst
	UndefinedToken
)

func (t TokenType) String() string {
	switch t {
	case Keyword:
		return "keyword"
	case Symbol:
		return "symbol"
	case Identifier:
		return "identifier"
	case IntConst:
		return "integerConstant"
	case StringConst:
		return "stringConstant"
	default:
		return "undefinedToken"
	}
}

// KeywordType enum
type KeywordType int

// keyword type enum
const (
	Class KeywordType = iota + 1
	Constructor
	Function
	Method
	Field
	Static
	Var
	Int
	Char
	Boolean
	Void
	True
	False
	Null
	This
	Let
	Do
	If
	Else
	While
	Return
	UndefinedKeyword
)

func (t KeywordType) String() string {
	switch t {
	case Class:
		return "class"
	case Constructor:
		return "constructor"
	case Function:
		return "function"
	case Method:
		return "method"
	case Field:
		return "field"
	case Static:
		return "static"
	case Var:
		return "var"
	case Int:
		return "int"
	case Char:
		return "char"
	case Boolean:
		return "boolean"
	case Void:
		return "void"
	case True:
		return "true"
	case False:
		return "false"
	case Null:
		return "null"
	case This:
		return "this"
	case Let:
		return "let"
	case Do:
		return "do"
	case If:
		return "if"
	case Else:
		return "else"
	case While:
		return "while"
	case Return:
		return "return"
	default:
		return "undefinedKeyword"
	}
}
