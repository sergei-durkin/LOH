package compiler

func (f *Fn) critedge() {
	for i := 0; i < len(f.po); i++ {
		bb := f.po[i]
		if len(bb.Pred) <= 1 {
			continue
		}

		for j := 0; j < len(bb.Pred); j++ {
			p := bb.Pred[j]
			if len(p.Succ) <= 1 {
				continue
			}
			newBB := f.cfg.NewBasicBlock(len(f.cfg.blocks))

			bb.removePred(p.ID)
			p.removeSucc(bb.ID)

			p.Succ = append(p.Succ, newBB)
			bb.Pred = append(bb.Pred, newBB)

			newBB.Succ = []*BasicBlock{bb}
			newBB.Pred = []*BasicBlock{p}

			newBB.Instr = append(newBB.Instr, &Goto{label: bb.ID})
			f.cfg.mp[newBB.ID] = newBB

			for k := 0; k < len(p.Instr); k++ {
				switch in := p.Instr[k].(type) {
				case *IfGoto:
					if in.label == bb.ID {
						in.label = newBB.ID
					}
					if in.fall == bb.ID {
						in.fall = newBB.ID
					}
				case *Goto:
					if in.label == bb.ID {
						in.label = newBB.ID
					}
				}
			}

			for k := 0; k < len(bb.Phi); k++ {
				phi := bb.Phi[k]
				for id, arg := range phi.Args {
					if id == p.ID {
						delete(phi.Args, id)
						phi.Args[newBB.ID] = arg
					}
				}
			}

			f.cfg.blocks = append(f.cfg.blocks, newBB)
		}
	}

	f.po = nil
	f.idom = nil
	f.sdom = nil
	f.domfront = nil

	f.dominators()
	f.sparseTree()
	f.dominanceFrontier()
}
