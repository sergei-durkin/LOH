package compiler

import (
	"fmt"
	"loh/token"
)

func (f *Fn) foldUnusedOperations() {
	reachable := f.reachable()
	done := make(map[Instruction]bool)

	iwl := []Instruction{}
	bwl := []*BasicBlock{f.entry}

	isConst := func(v Value) bool { _, ok := v.(*IntConst); _, ok2 := v.(*BoolConst); return ok || ok2 }

	isZero := func(v Value) bool { c, ok := v.(*IntConst); return ok && c.int == 0 }
	isOne := func(v Value) bool { c, ok := v.(*IntConst); return ok && c.int == 1 }

	toInt := func(v Value) int64 {
		switch v := v.(type) {
		default:
			panic(fmt.Sprintf("unknown value: %T, %+v", v, v))
		case *IntConst:
			return v.int
		case *BoolConst:
			if v.bool {
				return 1
			}
			return 0
		}
	}

	eqTypes := func(v1, v2 Value) bool {
		if _, ok := v1.(*IntConst); ok {
			_, ok := v2.(*IntConst)
			return ok
		}

		if _, ok := v1.(*BoolConst); ok {
			_, ok := v2.(*BoolConst)
			return ok
		}

		return false
	}

	cmp := func(op token.Token, v1, v2 Value) bool {
		if !eqTypes(v1, v2) {
			panic("dynamic cast int to bool not allowed")
		}

		l, r := toInt(v1), toInt(v2)
		switch op {
		default:
			panic(fmt.Sprintf("unknown operation: %+v", op))
		case token.LT:
			return l < r
		case token.LTE:
			return l <= r
		case token.GT:
			return l > r
		case token.GTE:
			return l >= r
		case token.NE:
			return l != r
		case token.EQ:
			return l == r
		}
	}

	processAssign := func(a *Assign) *Assign {
		fold := func(result Value) *Assign {
			return &Assign{size: a.size, target: a.target, op: token.ASSIGN, arg1: result, isTemp: a.isTemp}
		}

		switch a.op {
		default:
			panic(fmt.Sprintf("unknown operation: %+v", a.op))
		case token.LT, token.LTE, token.GT, token.GTE, token.NE, token.EQ:
			if !isConst(a.arg1) || !isConst(a.arg2) {
				return a
			}

			return fold(&BoolConst{bool: cmp(a.op, a.arg1, a.arg2)})
		case token.PLUS:
			if isZero(a.arg1) {
				return fold(a.arg2)
			}
			if isZero(a.arg2) {
				return fold(a.arg1)
			}
		case token.MINUS:
			if isZero(a.arg2) {
				return fold(a.arg1)
			}
		case token.STAR:
			if isOne(a.arg1) {
				return fold(a.arg2)
			}
			if isZero(a.arg1) {
				return fold(&IntConst{0})
			}

			if isOne(a.arg2) {
				return fold(a.arg1)
			}
			if isZero(a.arg2) {
				return fold(&IntConst{0})
			}
		case token.SLASH:
			if isZero(a.arg1) {
				return fold(&IntConst{0})
			}
			if isOne(a.arg2) {
				return fold(a.arg1)
			}
		case token.PERCENT:
			if isOne(a.arg2) {
				return fold(&IntConst{0})
			}
		case token.AMP:
			if isZero(a.arg1) {
				return fold(&IntConst{0})
			}
			if isZero(a.arg2) {
				return fold(&IntConst{0})
			}
		case token.PIPE:
			if isZero(a.arg1) {
				return fold(a.arg2)
			}
			if isZero(a.arg2) {
				return fold(a.arg1)
			}
		case token.XOR:
			if isZero(a.arg1) {
				return fold(a.arg2)
			}
			if isZero(a.arg2) {
				return fold(a.arg1)
			}
		}

		return a
	}

	var processInstruction func(Instruction)
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
			if instr.arg1 == nil || instr.arg2 == nil {
				return
			}

			newInstr := processAssign(instr)

			instr.target = newInstr.target
			instr.op = newInstr.op
			instr.arg1 = newInstr.arg1
			instr.arg2 = newInstr.arg2
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
