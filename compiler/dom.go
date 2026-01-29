package compiler

func (f *Fn) dominators() {
	if f.idom != nil {
		return
	}

	f.idom = make([]*BasicBlock, len(f.cfg.blocks))
	f.postorder()

	pnum := make([]int, len(f.po))
	for i := 0; i < len(f.po); i++ {
		pnum[f.po[i].ID] = i
	}

	entry := f.cfg.mp[0]
	f.idom[entry.ID] = entry

	if pnum[entry.ID] != len(f.po)-1 {
		panic("entry should be last in postorder")
	}

	ok := true
	for ok {
		ok = false

		for i := len(f.po) - 2; i >= 0; i-- {
			bb := f.po[i]
			var d *BasicBlock
			for j := 0; j < len(bb.Pred); j++ {
				pb := bb.Pred[j]
				if f.idom[pb.ID] == nil {
					continue
				}
				if d == nil {
					d = pb
					continue
				}

				d = intersect(d, pb, pnum, f.idom)
			}
			if d != f.idom[bb.ID] {
				f.idom[bb.ID] = d
				ok = true
			}
		}
	}
	f.idom[entry.ID] = nil
}

func intersect(block, candidate *BasicBlock, pnum []int, idom []*BasicBlock) *BasicBlock {
	for block != candidate {
		if pnum[block.ID] < pnum[candidate.ID] {
			block = idom[block.ID]
		} else {
			candidate = idom[candidate.ID]
		}
	}

	return block
}

func (f *Fn) dominanceFrontier() {
	if f.domfront != nil {
		return
	}

	f.domfront = make([][]bool, len(f.cfg.blocks))
	for i := range len(f.cfg.blocks) {
		f.domfront[i] = make([]bool, len(f.cfg.blocks))
	}

	for i := 0; i < len(f.po); i++ {
		if len(f.po[i].Pred) < 2 {
			continue
		}
		for j := 0; j < len(f.po[i].Pred); j++ {
			r := f.po[i].Pred[j]
			for r.ID != f.idom[f.po[i].ID].ID {
				f.domfront[r.ID][f.po[i].ID] = true
				r = f.idom[r.ID]
			}
		}
	}
}
