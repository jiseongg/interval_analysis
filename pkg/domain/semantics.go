package domain

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func evalArgument(v value.Value, s *State) Interval {
	switch v := v.(type) {
	case *constant.Int:
		i := int(v.X.Int64())
		return InterRange(Endpoint{i}, Endpoint{i})
	default:
		loc := v.Ident()
		itv := s.Find(loc)
		return itv
	}
}

func (s *State) transferInstAdd(inst *ir.InstAdd) {
	loc := inst.LocalIdent.Ident()
	vx := evalArgument(inst.X, s)
	vy := evalArgument(inst.Y, s)
	s.Bind(loc, InterPlus(vx, vy))
}

func (s *State) transferInstSub(inst *ir.InstSub) {
	loc := inst.LocalIdent.Ident()
	vx := evalArgument(inst.X, s)
	vy := evalArgument(inst.Y, s)
	s.Bind(loc, InterMinus(vx, vy))
}

func (s *State) transferInstMul(inst *ir.InstMul) {
	loc := inst.LocalIdent.Ident()
	vx := evalArgument(inst.X, s)
	vy := evalArgument(inst.Y, s)
	s.Bind(loc, InterMult(vx, vy))
}

func (s *State) transferInstICmp(inst *ir.InstICmp) {
	loc := inst.LocalIdent.Ident()
	vx := evalArgument(inst.X, s)
	vy := evalArgument(inst.Y, s)
	itv := InterTop()
	switch inst.Pred {
	case enum.IPredSLT:
		itv = InterSLT(vx, vy)
	}
	s.Bind(loc, itv)
}

func (s *State) transferInstPhi(inst *ir.InstPhi) {
	loc := inst.LocalIdent.Ident()
	itv := InterBot()
	for _, inc := range inst.Incs {
		itv = InterJoin(itv, evalArgument(inc.X, s))
	}
	s.Bind(loc, itv)
}

func (s *State) transferInst(inst ir.Instruction) {
	switch inst := inst.(type) {
	case *ir.InstAdd:
		s.transferInstAdd(inst)
	case *ir.InstSub:
		s.transferInstSub(inst)
	case *ir.InstMul:
		s.transferInstMul(inst)
	case *ir.InstPhi:
		s.transferInstPhi(inst)
	case *ir.InstICmp:
		s.transferInstICmp(inst)
	default:
		fmt.Printf("Unsupported instructions: %T\n", inst)
	}
}

func (s *State) TransferBlock(insts []ir.Instruction) {
	for _, inst := range insts {
		s.transferInst(inst)
	}
}
