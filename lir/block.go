package lir

import "loh/machine"

type BasicBlock struct {
	ID int

	Instr []Instruction

	Succ []*BasicBlock
	Pred []*BasicBlock
}

func (b *BasicBlock) Label() int {
	return b.ID
}

func (b *BasicBlock) Instructions() []machine.Instruction {
	instr := []machine.Instruction{}
	for i := 0; i < len(b.Instr); i++ {
		instr = append(instr, b.Instr[i])
	}
	return instr
}
