package analyzer

import (
	"strings"

	. "interval/pkg/domain"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type Cfg struct {
	fid    string
	sig    *types.FuncType
	params []*ir.Param
	blocks []Node
	succ   map[Node][]Node
	pred   map[Node][]Node
}

func NewCfg(f *ir.Func) Cfg {
	var cfg Cfg
	cfg.fid = f.Ident()
	cfg.sig = f.Sig
	cfg.params = f.Params
	for _, b := range f.Blocks {
		cfg.blocks = append(cfg.blocks, b)
	}
	cfg.succ = make(map[Node][]Node)
	cfg.pred = make(map[Node][]Node)
	for _, block := range f.Blocks {
		succs := block.Term.Succs()
		for _, s := range succs {
			cfg.succ[block] = append(cfg.succ[block], s)
		}
		for _, s := range succs {
			b := false
			for _, p := range cfg.pred[s] {
				if p == block {
					b = true
				}
			}
			// add pred only if it is unique
			if b == false {
				cfg.pred[s] = append(cfg.pred[s], block)
			}
		}
	}
	return cfg
}

func (c *Cfg) GetFid() string {
	return c.fid
}

func (c *Cfg) EntryBlock() Node {
	entry := c.blocks[0]
	// sanity check
	if len(c.pred[entry]) > 0 {
		panic("The entry block mismatch")
	}
	for _, to := range c.succ {
		for _, b := range to {
			if b == entry {
				panic("The first block is not the entry block")
			}
		}
	}
	return entry
}

func (c *Cfg) IsEntry(b Node) bool {
	if b == c.EntryBlock() {
		return true
	} else {
		return false
	}
}

func (c *Cfg) Succ(n Node) []Node {
	return c.succ[n]
}

func (c *Cfg) Pred(n Node) []Node {
	return c.pred[n]
}

func (c *Cfg) Dot() string {
	buf := &strings.Builder{}
	buf.WriteString("digraph " + c.fid + " {\n")

	// for _, block := range c.blocks {
	// 	fmt.Fprintf(buf, "\t%q -> %q\n", block, c.succ[block])
	// }

	buf.WriteString("}")
	return buf.String()
}

type Module struct {
	typeDefs []types.Type
	globals  []*ir.Global
	Cfgs     []Cfg
}

func NewModule(m *ir.Module) Module {
	var md Module
	md.typeDefs = m.TypeDefs
	md.globals = m.Globals
	for _, f := range m.Funcs {
		if f.Blocks != nil {
			md.Cfgs = append(md.Cfgs, NewCfg(f))
		}
	}
	return md
}
