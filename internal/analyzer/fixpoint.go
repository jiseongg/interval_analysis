package analyzer

import (
	. "interval/pkg/domain"
)

func InputOf(here Node, cfg Cfg, tbl Table) State {
	if cfg.IsEntry(here) {
		res := EmptyState()
		for _, p := range cfg.params {
			res.Bind(p.Ident(), InterTop())
		}
		return res
	} else {
		res := EmptyState()
		for _, p := range cfg.Pred(here) {
			res = StateJoin(res, tbl.Find(p))
		}
		return res
	}
}

func AnalyzeStep(cfg Cfg, tbl *Table, widen bool) {
	worklist := NewWorklist()
	worklist.AddSet(cfg.blocks)
	for !worklist.IsEmpty() {
		here := worklist.Choose()
		state := InputOf(here, cfg, *tbl)
		state.TransferBlock(here.Insts)
		old_state := tbl.Find(here)

		if widen {
			tbl.Bind(here, StateWiden(old_state, state))
			worklist.AddSet(cfg.Succ(here))
		} else {
			tbl.Bind(here, StateNarrow(old_state, state))
			worklist.AddSet(cfg.Succ(here))
		}
	}
}

func Analyze(cfg Cfg) Table {
	tbl := NewTable()

	AnalyzeStep(cfg, &tbl, true)  // widening
	AnalyzeStep(cfg, &tbl, false) // narrowing

	return tbl
}
