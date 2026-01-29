package compiler

import "fmt"

func (f *Fn) aliased() map[string]bool {
	aliased := make(map[string]bool)

	for _, bb := range f.cfg.blocks {
		for _, instr := range bb.Instr {
			extractAliased(instr, aliased)
		}
	}

	return aliased
}

func extractAliased(instr Instruction, aliased map[string]bool) {
	operands := []Value{}

	switch instr := instr.(type) {
	case *Assign:
		operands = append(operands, instr.arg1, instr.arg2)
	case *Alloca:
		operands = append(operands, instr.ptr)
	case *Store:
		operands = append(operands, instr.Destination)
	case *Load:
		operands = append(operands, instr.Source)
	case *Return:
		operands = append(operands, instr.value)
	case *Call:
		operands = append(operands, instr.args...)
	}

	var mark func(Value)
	mark = func(v Value) {
		switch v := v.(type) {
		case *Var:
			aliased[v.name] = true
		case *TempVar:
			aliased[v.label] = true
		}
	}

	for len(operands) > 0 {
		op := operands[0]
		operands = operands[1:]

		switch op := op.(type) {
		default:
			panic(fmt.Sprintf("unknown operation: %T %+v", op, op))

		case *AddressOf:
			mark(op.Target)
		case *Dereference:
			operands = append(operands, op.Addr)

		case *BoolConst:
		case *IntConst:

		case *Var:
		case *TempVar:

		case *Reg:
		case *ArgReg:
		case *FP:

		case nil:
		}
	}
}
