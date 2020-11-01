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

func AnalyzeWiden(cfg Cfg, tbl *Table) {
	worklist := NewWorklist()
	worklist.AddSet(cfg.blocks)
	for !worklist.IsEmpty() {
		here := worklist.Choose()
		state := InputOf(here, cfg, *tbl)
		state.TransferBlock(here.Insts)
		old_state := tbl.Find(here)

		if !StateOrder(state, old_state) {
			if StateOrder(state, old_state) {
				continue
			}
			tbl.Bind(here, StateWiden(old_state, state))
			worklist.AddSet(cfg.Succ(here))
		}
	}
}

func AnalyzeNarrow(cfg Cfg, tbl *Table) {
	worklist := NewWorklist()
	worklist.AddSet(cfg.blocks)
	for !worklist.IsEmpty() {
		here := worklist.Choose()
		state := InputOf(here, cfg, *tbl)
		state.TransferBlock(here.Insts)
		old_state := tbl.Find(here)
		if StateOrder(state, old_state) {
			if StateOrder(state, old_state) {
				continue
			}
			tbl.Bind(here, StateNarrow(old_state, state))
			worklist.AddSet(cfg.Succ(here))
		}
	}
}

func Analyze(cfg Cfg) Table {
	tbl := NewTable()

	AnalyzeWiden(cfg, &tbl)
	AnalyzeNarrow(cfg, &tbl)

	return tbl
}
