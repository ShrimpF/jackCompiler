package symboltable

// enum for kind
const (
	static   = "STATIC"
	field    = "FIELD"
	argument = "ARG"
	variable = "VAR"
)

type symbol struct {
	typeName string
	kind     string
	number   int
}

// SymbolTable has class-scope,subroutine scope map and count map.
type SymbolTable struct {
	classTable map[string]*symbol
	subTable   map[string]*symbol
	count      map[string]int
}

// New -- create a new symboltable
func New() *SymbolTable {
	count := map[string]int{
		static:   0,
		field:    0,
		argument: 0,
		variable: 0,
	}
	return &SymbolTable{count: count}
}

// StartSubroutine -- reset subtable map and arg/var count
func (st *SymbolTable) StartSubroutine() {
	st.subTable = make(map[string]*symbol)
	st.count[argument], st.count[variable] = 0, 0
}

// Define -- define a new identifier of given name,type and kind
func (st *SymbolTable) Define(name, typeName, kind string) {
	number, ok := st.count[kind]
	if !ok {
		return
	}
	st.count[kind]++

	switch kind {
	case static, field:
		st.classTable[name] = &symbol{
			typeName: typeName,
			kind:     kind,
			number:   number,
		}
	case argument, variable:
		st.subTable[name] = &symbol{
			typeName: typeName,
			kind:     kind,
			number:   number,
		}
	}
}

// VarCount -- return count
func (st *SymbolTable) VarCount(kind string) int {
	if c, ok := st.count[kind]; ok {
		return c
	}
	return -1
}

// KindOf -- return kind
func (st *SymbolTable) KindOf(name string) string {
	if class, ok := st.classTable[name]; ok {
		return class.kind
	}
	if sub, ok := st.subTable[name]; ok {
		return sub.kind
	}
	return ""
}

// TypeOf --return typename
func (st *SymbolTable) TypeOf(name string) string {
	if class, ok := st.classTable[name]; ok {
		return class.typeName
	}
	if sub, ok := st.subTable[name]; ok {
		return sub.typeName
	}
	return ""
}

// IndexOf --return typename
func (st *SymbolTable) IndexOf(name string) int {
	if class, ok := st.classTable[name]; ok {
		return class.number
	}
	if sub, ok := st.subTable[name]; ok {
		return sub.number
	}
	return -1
}
