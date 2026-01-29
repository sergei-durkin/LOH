package compiler

import (
	"fmt"
	"loh/token"
)

func initDsu(n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return p
}

func union(a, b int, dsu []int) {
	x, y := find(a, dsu), find(b, dsu)
	if x == y {
		return
	}

	dsu[y] = x
}

func find(a int, dsu []int) int {
	if dsu[a] != a {
		dsu[a] = find(dsu[a], dsu)
	}

	return dsu[a]
}

func (f *Fn) destroy() {
	cnt := make(map[string]int)
	for _, variable := range f.vars {
		switch v := variable.(type) {
		case *Assign:
			if v, ok := v.target.(*Var); ok {
				cnt[v.name]++
			}
		case *Call:
			if v, ok := v.target.(*Var); ok {
				cnt[v.name]++
			}
		case *Store:
			if v, ok := v.Destination.(*Var); ok {
				cnt[v.name]++
			}
		case *Phi:
			cnt[v.Var.name]++
			for range v.Args {
				cnt[v.Var.name]++
			}
		}
	}

	semweb := make(map[string][]int)
	for v, c := range cnt {
		semweb[v] = initDsu(c)
	}

	setVer := func(v Value) {
		if v, ok := v.(*Var); ok && len(semweb[v.name]) > 0 {
			v.ver = semweb[v.name][v.ver]
		}
	}

	for _, bb := range f.po {
		for _, phi := range bb.Phi {
			dsu := semweb[phi.Var.name]
			for _, op := range phi.Args {
				if op, ok := op.(*Var); ok {
					union(dsu[phi.Var.ver], op.ver, dsu)
				}
			}
		}
	}

	for _, bb := range f.po {
		for len(bb.Phi) > 0 {
			phi := bb.Phi[0]
			bb.Phi = bb.Phi[1:]

			dsu := semweb[phi.Var.name]
			for pred, arg := range phi.Args {
				term := f.cfg.mp[pred].Instr
				f.cfg.mp[pred].Instr = []Instruction{}

				switch op := arg.(type) {
				case *Var:
					if find(phi.Var.ver, dsu) != find(op.ver, dsu) {
						f.cfg.mp[pred].Instr = append(f.cfg.mp[pred].Instr, &Assign{
							target: &phi.Var,
							op:     token.ASSIGN,
							arg1:   op,

							size: int(f.cfg.size[phi.Var.name]),
						})
					}
				case *IntConst, *BoolConst, *StringConst:
					f.cfg.mp[pred].Instr = append(f.cfg.mp[pred].Instr, &Assign{
						target: &phi.Var,
						op:     token.ASSIGN,
						arg1:   arg,

						size: int(f.cfg.size[phi.Var.name]),
					})
				}

				f.cfg.mp[pred].Instr = append(f.cfg.mp[pred].Instr, term...)
				delete(phi.Args, pred)
			}

			if len(phi.Args) != 0 {
				panic("noway")
			}
		}

		for _, instr := range bb.Instr {
			switch instr := instr.(type) {
			default:
				panic(fmt.Sprintf("undefined instruction: %T %+v", instr, instr))
			case *Goto:
			case *Alloca:
			case *Store:
				setVer(instr.Destination)
				setVer(instr.Value)
			case *Call:
				setVer(instr.target.(Value))
				for _, arg := range instr.Args() {
					setVer(arg)
				}
			case *Assign:
				setVer(instr.target.(Value))
				for _, op := range instr.Operands() {
					setVer(op)
				}
			case *IfGoto:
				setVer(instr.cond)
			case *Return:
				setVer(instr.value)
			}
		}
	}
}
