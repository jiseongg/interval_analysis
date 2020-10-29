package domain

import (
	"github.com/llir/llvm/ir"
)

/////////////////////
// Abstract Domain //
/////////////////////

type Node *ir.Block
type Table map[Node]State

func NewTable() Table {
	return make(Table)
}

func (t *Table) String() string {
	var res string
	for n, s := range *t {
		res += "   " + n.LocalIdent.Ident() + "\n"
		res += s.String() + "\n"
	}
	return res
}

func (t *Table) Bind(n Node, s State) { (*t)[n] = s }
func (t *Table) Find(n Node) State {
	s, ok := (*t)[n]
	if !ok {
		return EmptyState()
	}
	return s
}
