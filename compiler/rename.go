package compiler

func (f *Fn) rename() {
	s := make(map[string][]int)
	cnt := make(map[string]int)

	var processValue func(Value)
	processValue = func(v Value) {
		switch v := v.(type) {
		case *Var:
			if top := top(s, v.name); top != -1 {
				v.ver = top
			}
		case *AddressOf:
			processValue(v.Target)
		}
	}

	var processInstr func(Instruction)
	processInstr = func(instr Instruction) {
		switch instr := instr.(type) {
		case *Store:
			processValue(instr.Destination)
			processValue(instr.Value)
		case *Assign:
			processValue(instr.arg1)
			processValue(instr.arg2)

			if v, ok := instr.target.(*Var); ok {
				idx := cnt[v.name]
				cnt[v.name]++

				s[v.name] = append(s[v.name], idx)
				v.ver = idx
			}
		case *Return:
			processValue(instr.value)
		case *IfGoto:
			processValue(instr.cond)
		case *Call:
			for i := 0; i < len(instr.args); i++ {
				processValue(instr.args[i])
			}
		case *Alloca:
			if v, ok := instr.ptr.(*Var); ok {
				idx := cnt[v.name]
				cnt[v.name]++

				s[v.name] = append(s[v.name], idx)
				v.ver = idx
			} else {
				// panic(instr.ptr)
			}
		}
	}

	var rename func(*BasicBlock)
	rename = func(bb *BasicBlock) {
		for i := 0; i < len(bb.Phi); i++ {
			p := bb.Phi[i]
			idx := cnt[p.Var.name]
			cnt[p.Var.name]++

			s[p.Var.name] = append(s[p.Var.name], idx)
			p.Var.ver = idx
		}

		for i := 0; i < len(bb.Instr); i++ {
			processInstr(bb.Instr[i])
		}

		for i := 0; i < len(bb.Succ); i++ {
			succ := bb.Succ[i]
			for j := 0; j < len(succ.Phi); j++ {
				value := succ.Phi[j].Var.name
				if ver := top(s, value); ver != -1 {
					succ.Phi[j].Args[bb.ID] = &Var{name: value, ver: ver}
				} else {
					delete(succ.Phi[j].Args, bb.ID)
				}
			}
		}

		ch := f.domTreeChildren(bb)
		for i := 0; i < len(ch); i++ {
			rename(ch[i])
		}

		for i := 0; i < len(bb.Phi); i++ {
			_ = pop(s, bb.Phi[i].Var.name)
		}

		for i := 0; i < len(bb.Instr); i++ {
			switch instr := bb.Instr[i].(type) {
			case *Assign:
				if v, ok := instr.target.(*Var); ok {
					_ = pop(s, v.name)
				}
			case *Alloca:
				if v, ok := instr.ptr.(*Var); ok {
					_ = pop(s, v.name)
				}
			}
		}
	}

	rename(f.entry)
}

func (f *Fn) domTreeChildren(x *BasicBlock) []*BasicBlock {
	var res []*BasicBlock
	for c := f.sdom[x.ID].child; c != nil; c = f.sdom[c.ID].sib {
		res = append(res, c)
	}
	return res
}

func top(m map[string][]int, v string) int {
	if m[v] != nil {
		return m[v][len(m[v])-1]
	}

	return -1
}

func pop(m map[string][]int, v string) int {
	if m[v] != nil {
		res := m[v][len(m[v])-1]
		m[v] = m[v][:len(m[v])-1]

		return res
	}

	return -1
}
