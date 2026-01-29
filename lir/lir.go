package lir

import (
	"loh/compiler"
	"loh/machine"
)

type Block interface {
	Label() int
	Instructions() []Instruction
}

type Lir struct {
	funcs []*Fn
}

func NewLir(ssa *compiler.Ssa) *Lir {
	ssaFuncs := ssa.Funcs()
	funcs := make([]*Fn, len(ssaFuncs))

	for i := 0; i < len(funcs); i++ {
		f := ssaFuncs[i]
		funcs[i] = NewFn(f.Name(), f.Blocks(), f.Postorder(), f.FrameIndex(), f.FrameSize(), f.Regs())
	}

	return &Lir{
		funcs: funcs,
	}
}

func (l *Lir) Functions() []machine.Function {
	funcs := []machine.Function{}
	for i := 0; i < len(l.funcs); i++ {
		funcs = append(funcs, l.funcs[i])
	}
	return funcs
}
