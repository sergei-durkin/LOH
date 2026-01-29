package lir

import (
	"loh/machine"
)

func legal(b []*BasicBlock, n int) {
	tmp := func(instr Instruction, instrs *[]Instruction) {
		ops := instr.Operands()
		newOps := []machine.Value{}
		for i := 0; i < len(ops); i++ {
			switch op := ops[i].(type) {
			default:
				newOps = append(newOps, op)
			case *machine.IntConst:
				n++
				tmpReg := &machine.Reg{ID: n}
				*instrs = append(*instrs, &MOV{destination: tmpReg, source: op})
				newOps = append(newOps, tmpReg)
			}
		}
		instr.(Replaceable).ReplaceAll(newOps...)
		*instrs = append(*instrs, instr)
	}

	swap2 := func(instr Instruction, instrs *[]Instruction) bool {
		ops := instr.Operands()
		switch ops[0].(type) {
		case *machine.FP:
		case *machine.Reg:
		default:
			return false
		}

		switch ops[1].(type) {
		case *machine.FP:
		case *machine.Reg:
		default:
			return false
		}

		*instrs = append(*instrs, instr)
		return true
	}
	swap3 := func(instr Instruction, instrs *[]Instruction) bool {
		ops := instr.Operands()
		switch ops[1].(type) {
		case *machine.FP, *machine.Reg:
			*instrs = append(*instrs, instr)
			return true
		}

		switch ops[2].(type) {
		case *machine.FP, *machine.Reg:
			ops[1], ops[2] = ops[2], ops[1]
			instr.(Replaceable).ReplaceAll(ops...)
			*instrs = append(*instrs, instr)

			return true
		}

		return false
	}

	for i := 0; i < len(b); i++ {
		bb := b[i]

		instrs := []Instruction{}
		for j := 0; j < len(bb.Instr); j++ {
			switch instr := bb.Instr[j].(type) {
			default:
				instrs = append(instrs, instr)
			case *AND, *OR, *XOR, *MUL, *DIV, *MOD, *SUB, *STR, *LDR:
				tmp(instr, &instrs)
			case *ALLOCA:
				if !swap2(instr, &instrs) {
					tmp(instr, &instrs)
				}
			case *SUM:
				if !swap3(instr, &instrs) {
					tmp(instr, &instrs)
				}
			}
		}
		bb.Instr = instrs
	}
}
