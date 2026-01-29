package compiler

import (
	"fmt"
	"loh/token"
)

func (f *Fn) copyprop() {
	reachable := f.reachable()
	done := make(map[Instruction]bool)

	iwl := []Instruction{}
	bwl := []*BasicBlock{f.entry}

	var replace func(Instruction, Variable, Variable) bool
	replace = func(use Instruction, old Variable, new Variable) bool {
		changed := false

		switch use := use.(type) {
		default:
			panic(fmt.Sprintf("unexpected instruction %T %+v", use, use))
		case nil, *Alloca, *Phi, *Store:
			return false

		case *Call:
			for i := 0; i < len(use.args); i++ {
				if arg, ok := use.args[i].(Variable); ok && arg.Label() == old.Label() {
					use.args[i] = new.(Value)
					changed = true
				}
			}

		case *Assign:
			if arg, ok := use.arg1.(Variable); ok && arg.Label() == old.Label() {
				use.arg1 = new.(Value)
				changed = true
			}
			if arg, ok := use.arg2.(Variable); ok && arg.Label() == old.Label() {
				use.arg2 = new.(Value)
				changed = true
			}

		case *Return:
			if arg, ok := use.value.(Variable); ok && arg.Label() == old.Label() {
				use.value = new.(Value)
				changed = true
			}

		case *IfGoto:
			if arg, ok := use.cond.(Variable); ok && arg.Label() == old.Label() {
				use.cond = new.(Value)
				changed = true
			}
		}

		if changed {
			f.users[new.Label()] = append(f.users[new.Label()], use)
		}

		return changed
	}

	var processVariable func(Variable, Value)
	processVariable = func(old Variable, new Value) {
		switch new := new.(type) {
		default:
			panic(fmt.Sprintf("unexpected instruction %T: %+v", new, new))
		case nil:
		case *Dereference:
		case *FP:
		case *Reg:
		case *ArgReg:
		case *IntConst:
		case *BoolConst:
		case *AddressOf:
		case Variable:
			removed := []int{}
			for i := 0; i < len(f.users[old.Label()]); i++ {
				if replace(f.users[old.Label()][i], old, new) {
					removed = append(removed, i)
				}
			}

			users := f.users[old.Label()]
			label := old.Label()
			if len(removed) < len(users) {
				for i := 0; i < len(removed); i++ {
					users[removed[i]] = nil
				}
				for i := 0; i < len(users); i++ {
					ln := len(users)
					users[i], users[ln-1] = users[ln-1], users[i]
					users = users[:ln-1]
					i--
				}
			} else {
				f.users[label] = nil
				delete(f.vars, label)
			}

			if len(removed) > 0 {
				iwl = append(iwl, f.users[new.Label()]...)
			}
		}
	}

	var simplify func(*Assign)
	simplify = func(in *Assign) {
		if arg1, ok := in.arg1.(Variable); ok {
			if arg2, ok := in.arg2.(Variable); ok && arg1.Label() == arg2.Label() {
				switch in.op {
				case token.MINUS:
					in.op = token.EQ
					in.arg1, in.arg2 = &IntConst{int: 0}, nil
				case token.EQ, token.LTE, token.GTE:
					in.op = token.EQ
					in.arg1, in.arg2 = &BoolConst{bool: true}, nil
				case token.NE, token.LT, token.GT:
					in.op = token.EQ
					in.arg1, in.arg2 = &BoolConst{bool: false}, nil
				default:
					return
				}

				f.users.removeUser(arg1.Label(), in)
			}
		}
	}

	var processInstruction func(instr Instruction)
	processInstruction = func(instr Instruction) {
		switch instr := instr.(type) {
		default:
			panic(fmt.Sprintf("unexpected instruction %T: %+v", instr, instr))
		case *Alloca:
		case *Store:
		case *Return:
		case *Phi:
		case *Call:
		case *Assign:
			if _, ok := f.vars[instr.target.Label()]; !ok {
				return
			}
			if len(f.users[instr.target.Label()]) == 0 {
				return
			}
			if instr.arg1 != nil && instr.arg2 != nil {
				simplify(instr)
				return
			}

			processVariable(instr.target, instr.arg1)
			processVariable(instr.target, instr.arg2)
		case *Goto:
			if done[instr] {
				return
			}

			bwl = append(bwl, f.cfg.mp[instr.label])
			reachable[instr.label] = true
		case *IfGoto:
			if done[instr] {
				return
			}

			if i, ok := instr.cond.(*BoolConst); ok {
				if i.bool {
					bwl = append(bwl, f.cfg.mp[instr.fall])
					reachable[instr.fall] = true
				} else {
					bwl = append(bwl, f.cfg.mp[instr.label])
					reachable[instr.label] = true
				}
			} else {
				bwl = append(bwl, f.cfg.mp[instr.fall])
				bwl = append(bwl, f.cfg.mp[instr.label])

				reachable[instr.label] = true
				reachable[instr.fall] = true
			}

			done[instr] = true
		}
	}

	for len(iwl) > 0 || len(bwl) > 0 {
		if len(iwl) > 0 {
			instr := iwl[0]
			iwl = iwl[1:]

			processInstruction(instr)
		} else {
			bb := bwl[0]
			bwl = bwl[1:]
			if !reachable[bb.ID] {
				continue
			}

			iwl = append(iwl, bb.Instr...)
			for i := 0; i < len(bb.Phi); i++ {
				iwl = append(iwl, bb.Phi[i])
			}
		}
	}
}
