package vm

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Index int
	Scope SymbolScope
}
type SymbolTable struct {
	outer *SymbolTable
	store map[string]Symbol
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		outer: outer,
		store: make(map[string]Symbol),
	}
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: len(s.store)}
	if s.outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	s.store[name] = symbol
	return symbol
}
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	if symbol, ok := s.store[name]; ok || s.outer == nil {
		return symbol, ok
	} else {
		symbol, ok := s.outer.Resolve(name)
		return symbol, ok
	}
}
func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
	return symbol
}
