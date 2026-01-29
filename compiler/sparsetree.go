package compiler

type SparseTreeNode struct {
	child *BasicBlock
	sib   *BasicBlock
	par   *BasicBlock

	entry, exit int
}

type SparseTree []SparseTreeNode

func (f *Fn) sparseTree() {
	if f.sdom != nil {
		return
	}

	res := make(SparseTree, len(f.po))

	for i := 0; i < len(f.po); i++ {
		cur := &res[f.po[i].ID]
		par := f.idom[f.po[i].ID]
		if par == nil {
			continue
		}

		cur.par = par
		cur.sib = res[par.ID].child
		res[par.ID].child = f.po[i]
	}

	res.numberBlock(f.entry, 1)

	f.sdom = res
}

func (s SparseTree) numberBlock(b *BasicBlock, n int) int {
	n++
	s[b.ID].entry = n

	n += 2
	for c := s[b.ID].child; c != nil; c = s[c.ID].sib {
		n = s.numberBlock(c, n)
	}

	n++
	s[b.ID].exit = n

	return n + 2
}
