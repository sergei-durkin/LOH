package compiler

import (
	"fmt"
)

func (f *Fn) variables(withReach bool) {
	reachable := f.reachable()

	vars := make(map[string]Instruction)
	users := make(map[string][]Instruction)

	var addVar func(v Value, instr Instruction)
	addVar = func(v Value, instr Instruction) {
		switch v := v.(type) {
		case Variable:
			vars[v.Label()] = instr
		case *AddressOf:
			addVar(v.Target, instr)
		case *Dereference:
			addVar(v.Addr, instr)
		}
	}

	var addVarUser func(v Value, instr Instruction)
	addVarUser = func(v Value, instr Instruction) {
		switch v := v.(type) {
		default:
			panic(fmt.Sprintf("undefined value: %T %+v", v, v))
		case nil:
		case *IntConst:
		case *BoolConst:
		case *Reg:
		case *ArgReg:
		case *FP:
		case Variable:
			users[v.Label()] = append(users[v.Label()], instr)
		case *AddressOf:
			addVarUser(v.Target, instr)
		case *Dereference:
			addVarUser(v.Addr, instr)
		}
	}

	for i := 0; i < len(f.cfg.blocks); i++ {
		bb := f.cfg.blocks[i]
		if withReach && !reachable[bb.ID] {
			continue
		}

		for j := 0; j < len(bb.Phi); j++ {
			addVar(&bb.Phi[j].Var, bb.Phi[j])

			for _, arg := range bb.Phi[j].Args {
				addVarUser(arg, bb.Phi[j])
			}
		}

		for j := 0; j < len(bb.Instr); j++ {
			switch instr := bb.Instr[j].(type) {
			default:
				panic(fmt.Sprintf("unexpected instruction: %T %+v", instr, instr))
			case *Goto:
			case *IfGoto:
				addVarUser(instr.cond, instr)
			case *Assign:
				addVar(instr.target.(Value), instr)

				addVarUser(instr.arg1, instr)
				addVarUser(instr.arg2, instr)
			case *Alloca:
				addVar(instr.ptr, instr)
				addVarUser(instr.size, instr)
			case *Store:
				addVarUser(instr.Destination, instr)
				addVarUser(instr.Value, instr)
			case *Load:
				addVarUser(instr.Source, instr)
				addVarUser(instr.Destination, instr)
			case *Return:
				addVarUser(instr.value, instr)
			case *Call:
				addVar(instr.target.(Value), instr)

				for _, arg := range instr.args {
					addVarUser(arg, instr)
				}
			}
		}

		for j := 0; j < len(bb.Phi); j++ {
			t := bb.Phi[j].Var.Label()
			vars[t] = bb.Phi[j]

			for id, arg := range bb.Phi[j].Args {
				if withReach && !reachable[id] {
					continue
				}

				addVarUser(arg, bb.Phi[j])
			}
		}
	}

	f.users = users
	f.vars = vars
}
