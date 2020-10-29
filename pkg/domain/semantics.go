package domain

import (
	//  "fmt"
	"github.com/llir/llvm/ir"
	//	"github.com/llir/llvm/ir/value"
	//	"github.com/llir/llvm/ir/constant"
	//	"github.com/llir/llvm/ir/enum"
)

func (s *State) transferInst(inst ir.Instruction) {
	// TODO
	panic("Not implemented")
}

func (s *State) TransferBlock(insts []ir.Instruction) {
	for _, inst := range insts {
		s.transferInst(inst)
	}
}
