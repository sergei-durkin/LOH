package compiler

import (
	"fmt"
	"loh/token"
)

const (
	undefined = iota
	defined
	overdefined
)

type lat struct {
	kind int
	val  Value
}

func (l *lat) eq(another *lat) bool {
	if l == nil || another == nil {
		return l == another
	}

	if l.kind != another.kind {
		return false
	}

	if l.kind != defined {
		return true
	}

	if l.val == nil || another.val == nil {
		return l.val == another.val
	}

	return l.val.Equal(another.val)
}

func (f *Fn) constPropagation() {
	reachable := f.reachable()
	done := make([]bool, len(f.cfg.blocks))

	lats := make(map[string]*lat)

	iwl := []Instruction{}
	bwl := append([]*BasicBlock{}, f.entry)

	aliased := f.aliased()

	var lattice func(v Value) *lat
	lattice = func(v Value) *lat {
		switch val := v.(type) {
		default:
			panic(fmt.Sprintf("undefined value %T %+v", v, v))
		case *FP, *AddressOf, *Reg, *ArgReg:
			return &lat{kind: overdefined}
		case nil, *IntConst, *BoolConst:
			return &lat{kind: defined, val: val}
		case Variable:
			l := val.Label()
			if aliased[l] || aliased[val.Name()] {
				lats[l] = &lat{kind: overdefined}

				return lats[l]
			}

			l1 := lats[l]
			if l1 != nil {
				return l1
			}

			lats[l] = &lat{kind: undefined}

			return lats[l]
		case *Dereference:
			switch addr := val.Addr.(type) {
			case Variable, *Dereference:
				l := lattice(val.Addr)
				if l != nil && l.kind == defined {
					return &lat{kind: defined, val: &Dereference{Addr: l.val}}
				}

				return &lat{kind: l.kind}
			case *IntConst:
				return &lat{kind: defined, val: &Dereference{Addr: addr}}
			default:
				panic(fmt.Sprintf("unexpected addr %T %+v", addr, addr))
			}
		}
	}

	evaluate := func(op token.Token, l1, l2 *lat) *lat {
		if l1 == nil {
			panic("noway")
		}

		if l2 == nil {
			return &lat{kind: l1.kind, val: l1.val}
		}

		if l1.kind == overdefined || l2.kind == overdefined {
			return &lat{kind: overdefined}
		}

		if l1.kind != defined || l2.kind != defined {
			return &lat{kind: undefined}
		}

		ai, aok := l1.val.(*IntConst)
		bi, bok := l2.val.(*IntConst)
		if aok && bok {
			var res Value
			switch op {
			default:
				panic(fmt.Sprintf("op %s not defined on int", op.String()))
			case token.PLUS:
				res = &IntConst{int: ai.int + bi.int}
			case token.MINUS:
				res = &IntConst{int: ai.int - bi.int}
			case token.STAR:
				res = &IntConst{int: ai.int * bi.int}
			case token.SLASH:
				if bi.int == 0 {
					panic("div by zero")
				}
				res = &IntConst{int: ai.int / bi.int}
			case token.XOR:
				res = &IntConst{int: ai.int ^ bi.int}
			case token.PIPE:
				res = &IntConst{int: ai.int | bi.int}
			case token.AMP:
				res = &IntConst{int: ai.int & bi.int}
			case token.LT:
				res = &BoolConst{bool: ai.int < bi.int}
			case token.GT:
				res = &BoolConst{bool: ai.int > bi.int}
			case token.LTE:
				res = &BoolConst{bool: ai.int <= bi.int}
			case token.GTE:
				res = &BoolConst{bool: ai.int >= bi.int}
			case token.NE:
				res = &BoolConst{bool: ai.int != bi.int}
			case token.EQ:
				res = &BoolConst{bool: ai.int == bi.int}
			}

			return &lat{kind: defined, val: res}
		}

		ab, aok := l1.val.(*BoolConst)
		bb, bok := l2.val.(*BoolConst)
		if aok && bok {
			var res Value

			switch op {
			default:
				panic(fmt.Sprintf("op %s not defined on bool", op.String()))
			case token.OR:
				res = &BoolConst{bool: ab.bool || bb.bool}
			case token.AND:
				res = &BoolConst{bool: ab.bool && bb.bool}
			case token.NE:
				res = &BoolConst{bool: ab.bool != bb.bool}
			case token.EQ:
				res = &BoolConst{bool: ab.bool == bb.bool}
			}

			return &lat{kind: defined, val: res}
		}

		return &lat{kind: overdefined}
	}

	merge := func(l1, l2 *lat) *lat {
		if l1 == nil || l2 == nil {
			panic("noway")
		}

		if l1.kind == overdefined || l2.kind == overdefined {
			return &lat{kind: overdefined}
		}

		if l1.kind == undefined {
			return &lat{kind: l2.kind, val: l2.val}
		}

		if l2.kind == undefined {
			return &lat{kind: l1.kind, val: l1.val}
		}

		if l1.val.Equal(l2.val) {
			return &lat{kind: defined, val: l1.val}
		}

		return &lat{kind: overdefined}
	}

	isTrue := func(v Value) (bool, bool) {
		switch val := v.(type) {
		case *BoolConst:
			return val.bool, true
		case *IntConst:
			return val.int == 1, val.int == 0 || val.int == 1
		default:
			return false, false
		}
	}

	var processInstruction func(instr Instruction)
	processInstruction = func(instr Instruction) {
		switch in := instr.(type) {
		default:
			panic(fmt.Sprintf("undefined instruction: %T %+v", instr, instr))
		case nil:
		case *Store:
		case *Phi:
			if !reachable[in.BbID] {
				return
			}

			label := in.Var.Label()
			oldVal := lats[label]
			if oldVal == nil {
				oldVal = &lat{kind: undefined}
				lats[label] = oldVal
			}

			if oldVal.kind == overdefined {
				return
			}

			var newVal *lat
			for pred, v := range in.Args {
				if !reachable[pred] {
					continue
				}

				arg := lattice(v)
				if newVal == nil {
					newVal = arg
				} else {
					newVal = merge(newVal, arg)
				}

				if newVal.kind == overdefined {
					break
				}
			}

			if newVal == nil {
				newVal = &lat{kind: undefined}
			}

			if !oldVal.eq(newVal) {
				lats[label] = newVal
				iwl = append(iwl, f.users[label]...)
			}
		case *Alloca:
			label := in.ptr.(Variable).Label()
			offset, ok := f.frameIndex[in.ptr.(Variable).Name()]
			if !ok {
				lats[label] = &lat{kind: overdefined}
				return
			}
			lats[label] = &lat{kind: defined, val: &IntConst{int: int64(offset)}}
		case *Call:
			lats[in.target.Label()] = &lat{kind: overdefined}
		case *Assign:
			label := in.target.Label()
			oldVal := lats[label]
			if oldVal == nil {
				oldVal = &lat{kind: undefined}
				lats[label] = oldVal
			}

			if oldVal.kind == overdefined {
				return
			}

			l1 := lattice(in.arg1)
			var l2 *lat
			if in.arg2 != nil {
				l2 = lattice(in.arg2)
			}

			newVal := evaluate(in.op, l1, l2)

			if !oldVal.eq(newVal) {
				lats[label] = newVal
				iwl = append(iwl, f.users[label]...)
			}
		case *Goto:
			if !done[in.label] {
				done[in.label] = true
				bwl = append(bwl, f.cfg.mp[in.label])
			}
		case *Return:
			l := lattice(in.value)
			switch l.kind {
			case undefined:
				t, ok := in.value.(Variable)
				if ok {
					iwl = append(iwl, f.users[t.Label()]...)
				}
			case defined:
			case overdefined:
			}
		case *IfGoto:
			oldVal := lattice(in.cond)
			switch oldVal.kind {
			case undefined:
			case overdefined:
				if !done[in.label] {
					done[in.label] = true
					bwl = append(bwl, f.cfg.mp[in.label])
				}
				if !done[in.fall] {
					done[in.fall] = true
					bwl = append(bwl, f.cfg.mp[in.fall])
				}
			case defined:
				val := lattice(in.cond)
				take := in.label
				if v, ok := isTrue(val.val); ok && v {
					take = in.fall
				}

				if !done[take] {
					done[take] = true
					bwl = append(bwl, f.cfg.mp[take])
				}
			}
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

	for i := 0; i < len(f.po); i++ {
		bb := f.po[i]
		if !reachable[bb.ID] {
			continue
		}

		phi := make([]*Phi, 0, len(bb.Phi))
		instrs := make([]Instruction, 0, len(bb.Instr))

		for j := 0; j < len(bb.Phi); j++ {
			l := lattice(&bb.Phi[j].Var)
			if l != nil && l.kind == defined {
				bb.Phi[j].Args = make(map[int]Value)
				bb.Phi[j].Args[0] = l.val
				instrs = append(instrs, &Assign{size: int(f.cfg.size[bb.Phi[j].Var.name]), target: &bb.Phi[j].Var, op: token.ASSIGN, arg1: l.val})
			} else {
				for k, v := range bb.Phi[j].Args {
					l := lattice(v)
					if l != nil && l.kind == defined {
						bb.Phi[j].Args[k] = l.val
					}
				}
				phi = append(phi, bb.Phi[j])
			}
		}
		bb.Phi = phi

	loop:
		for j := 0; j < len(bb.Instr); j++ {
			instr := bb.Instr[j]
			switch in := instr.(type) {
			default:
				panic(fmt.Sprintf("undefined instruction %T %+v", instr, instr))
			case *Call:
			case *Goto:
			case *Alloca:
			case *Store:
			case *Assign:
				l := lats[in.target.Label()]
				if l != nil && l.kind == defined {
					in.op = token.ASSIGN
					in.arg1 = l.val
					in.arg2 = nil
				}
			case *IfGoto:
				l := lattice(in.cond)
				if l != nil && l.kind == defined {
					take := in.label
					if v, ok := isTrue(l.val); ok && v {
						take = in.fall
					}
					bb.Instr[j] = &Goto{label: take}
				}
			case *Return:
				l := lattice(in.value)
				if l != nil && l.kind == defined {
					in.value = l.val
				}
				instrs = append(instrs, bb.Instr[j])
				break loop
			}
			instrs = append(instrs, bb.Instr[j])
		}
		bb.Instr = instrs
	}
}
