package compiler

func (f *Fn) defines() {
	aliased := f.aliased()

	resolveDef := func(v any) (string, bool) {
		switch vv := v.(type) {
		default:
			return "", false
		case *Var:
			return vv.name, true
		}
	}

	for i := 0; i < len(f.cfg.blocks); i++ {
		bb := f.cfg.blocks[i]
		for j := 0; j < len(bb.Instr); j++ {
			switch instr := bb.Instr[j].(type) {
			case *Call:
				if t, ok := resolveDef(instr.target); ok && !aliased[t] {
					if len(f.defs[t]) == 0 {
						f.defs[t] = make([]bool, len(f.cfg.blocks))
					}
					f.defs[t][bb.ID] = true
				}
			case *Assign:
				if t, ok := resolveDef(instr.target); ok && !instr.isTemp && !aliased[t] {
					if len(f.defs[t]) == 0 {
						f.defs[t] = make([]bool, len(f.cfg.blocks))
					}
					f.defs[t][bb.ID] = true
				}
			case *Alloca:
				if t, ok := resolveDef(instr.ptr); ok {
					if len(f.defs[t]) == 0 {
						f.defs[t] = make([]bool, len(f.cfg.blocks))
					}
					f.defs[t][bb.ID] = true
				}
			}
		}
	}
}
