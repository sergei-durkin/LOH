package compiler

func (f *Fn) unusedVars() {
	reachable := f.reachable()
	live := make(map[Instruction]bool)
	aliased := f.aliased()

	bwl := []*BasicBlock{f.entry}
	iwl := []Instruction{}

	for len(iwl) > 0 || len(bwl) > 0 {
		if len(iwl) > 0 {
			instr := iwl[len(iwl)-1]
			iwl = iwl[:len(iwl)-1]

			switch in := instr.(type) {
			case *IfGoto:
				if live[in] {
					continue
				}

				bwl = append(bwl, f.cfg.mp[in.label])
				bwl = append(bwl, f.cfg.mp[in.fall])
			case *Goto:
				if live[in] {
					continue
				}

				bwl = append(bwl, f.cfg.mp[in.label])
			case *Phi:
				t := in.Var.Label()
				if len(f.users[t]) == 0 {
					live[instr] = false

					for _, op := range in.Operands() {
						if v, ok := op.(Variable); ok {
							f.users.removeUser(v.Label(), in)
							iwl = append(iwl, f.vars[v.Label()])
						}
					}
					continue
				}
			case *Assign:
				t := in.target.Label()
				if len(f.users[t]) == 0 {
					if _, ok := in.arg1.(*AddressOf); ok {
						continue
					}
					if _, ok := in.arg2.(*AddressOf); ok {
						continue
					}
					live[instr] = false

					for _, op := range in.Operands() {
						if v, ok := op.(Variable); ok {
							f.users.removeUser(v.Label(), in)
							iwl = append(iwl, f.vars[v.Label()])
						}
					}
					continue
				}
			}

			live[instr] = true
		} else {
			bb := bwl[len(bwl)-1]
			bwl = bwl[:len(bwl)-1]

			if !reachable[bb.ID] {
				continue
			}

			for i := 0; i < len(bb.Phi); i++ {
				iwl = append(iwl, bb.Phi[i])
			}

			for i := 0; i < len(bb.Instr); i++ {
				iwl = append(iwl, bb.Instr[i])
				if _, ok := bb.Instr[i].(*Return); ok {
					break
				}
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
			if !live[bb.Phi[j]] {
				continue
			}

			phi = append(phi, bb.Phi[j])
		}
		bb.Phi = phi

		for j := 0; j < len(bb.Instr); j++ {
			if instr, ok := bb.Instr[j].(*Assign); ok {
				if t, ok := instr.target.(*Var); ok {
					if aliased[t.name] {
						instrs = append(instrs, bb.Instr[j])
						continue
					}
				}
				if arg, ok := instr.arg1.(*AddressOf); ok {
					if t, ok := arg.Target.(*Var); ok {
						if aliased[t.name] {
							instrs = append(instrs, bb.Instr[j])
							continue
						}
					}
				}
				if arg, ok := instr.arg2.(*AddressOf); ok {
					if t, ok := arg.Target.(*Var); ok {
						if aliased[t.name] {
							instrs = append(instrs, bb.Instr[j])
							continue
						}
					}
				}
			}
			if !live[bb.Instr[j]] {
				continue
			}

			instrs = append(instrs, bb.Instr[j])
		}
		bb.Instr = instrs
	}
}
