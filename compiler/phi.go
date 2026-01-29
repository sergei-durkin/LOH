package compiler

func (f *Fn) placePhi() {
	has := make([]map[string]bool, len(f.cfg.blocks))
	for i := 0; i < len(f.cfg.blocks); i++ {
		has[i] = make(map[string]bool)
	}

	for v, ids := range f.defs {
		q := []int{}
		for i := 0; i < len(ids); i++ {
			if ids[i] {
				q = append(q, i)
			}
		}

		for len(q) > 0 {
			cur := f.cfg.mp[q[0]]
			q = q[1:]

			for i := 0; i < len(f.cfg.blocks); i++ {
				if !f.domfront[cur.ID][i] {
					continue
				}

				bb := f.cfg.mp[i]

				if has[bb.ID][v] {
					continue
				}

				r := Var{name: v}
				p := &Phi{
					BbID: bb.ID,
					Var:  r,
					Args: make(map[int]Value, len(bb.Pred)),
				}
				for j := 0; j < len(bb.Pred); j++ {
					p.Args[bb.Pred[j].ID] = &r
				}

				bb.Phi = append(bb.Phi, p)
				has[bb.ID][v] = true

				if !f.defs[v][bb.ID] {
					q = append(q, bb.ID)
				}
			}
		}
	}
}
