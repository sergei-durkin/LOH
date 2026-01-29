package compiler

type BasicBlock struct {
	ID    int
	Name  string
	Instr []Instruction

	Pred []*BasicBlock
	Succ []*BasicBlock

	Phi []*Phi
}

func (bb *BasicBlock) removeSucc(id int) {
	for i := 0; i < len(bb.Succ); i++ {
		if bb.Succ[i].ID == id {
			bb.Succ[i], bb.Succ[len(bb.Succ)-1] = bb.Succ[len(bb.Succ)-1], bb.Succ[i]
			bb.Succ = bb.Succ[:len(bb.Succ)-1]
			i--
		}
	}
}

func (bb *BasicBlock) removePred(id int) {
	for i := 0; i < len(bb.Pred); i++ {
		if bb.Pred[i].ID == id {
			bb.Pred[i], bb.Pred[len(bb.Pred)-1] = bb.Pred[len(bb.Pred)-1], bb.Pred[i]
			bb.Pred = bb.Pred[:len(bb.Pred)-1]
			i--
		}
	}
}

type Users map[string][]Instruction

// TODO optimize: можно за константу, пока лень
func (u Users) removeUser(label string, in Instruction) {
	users, ok := u[label]
	if !ok || len(users) == 0 {
		return
	}

	for i := 0; i < len(users); i++ {
		if users[i] != in {
			continue
		}

		users[i], users[len(users)-1] = users[len(users)-1], users[i]
		users = users[:len(users)-1]
		i--
	}

	u[label] = users
}

type Fn struct {
	cfg *Cfg

	defs map[string][]bool

	vars  map[string]Instruction
	users Users

	po       []*BasicBlock // postorder blocks
	idom     []*BasicBlock // immediate dominators
	sdom     SparseTree    // dominators tree
	domfront [][]bool      // dominance frontier

	virtualRegs map[string]VirtualRegister

	frameIndex map[string]int
	frameSize  int

	entry *BasicBlock
}

func NewFn(c *Cfg) *Fn {
	f := &Fn{
		cfg: c,

		defs: make(map[string][]bool),

		entry: c.blocks[0],
	}

	f.defines()

	{ // construct SSA
		f.dominators()
		f.sparseTree()
		f.dominanceFrontier()

		f.placePhi()
		f.rename()
	}

	{ // main passes
		f.variables(false)
		f.copyprop()
		f.constPropagation()
		f.foldUnusedOperations()
	}

	{ // preschedule
		f.allocs()

		f.variables(true)
		f.copyprop()
		f.constPropagation()

		f.variables(true)
		f.unusedVars()

		f.simplify()
		f.critedge()

		f.variables(true)
	}

	{ // schedule
		f.destroy()
		f.variables(true)
		f.virtualregalloc()
	}

	return f
}

func (f *Fn) Name() string {
	return f.cfg.name
}

func (f *Fn) Blocks() []*BasicBlock {
	return f.cfg.blocks
}

func (f *Fn) Postorder() []*BasicBlock {
	return f.po
}

func (f *Fn) FrameIndex() map[string]int {
	return f.frameIndex
}

func (f *Fn) FrameSize() int {
	return f.frameSize
}

func (f *Fn) Regs() map[string]VirtualRegister {
	return f.virtualRegs
}

func (f *Fn) reachable() []bool {
	res := make([]bool, len(f.cfg.blocks))
	res[f.entry.ID] = true

	iwl := []Instruction{}

	bwl := []*BasicBlock{}
	bwl = append(bwl, f.entry)

	for len(iwl) > 0 || len(bwl) > 0 {
		if len(iwl) > 0 {
			in := iwl[len(iwl)-1]
			iwl = iwl[:len(iwl)-1]

			switch in := in.(type) {
			case *IfGoto:
				a := in.label
				if !res[a] {
					res[a] = true
					bwl = append(bwl, f.cfg.mp[a])
				}
				a = in.fall
				if !res[a] {
					res[a] = true
					bwl = append(bwl, f.cfg.mp[a])
				}
			case *Goto:
				a := in.label
				if !res[a] {
					res[a] = true
					bwl = append(bwl, f.cfg.mp[a])
				}
			}
		} else {
			bb := bwl[len(bwl)-1]
			bwl = bwl[:len(bwl)-1]

			for i := 0; i < len(bb.Instr); i++ {
				if _, ok := bb.Instr[i].(*Return); ok {
					break
				}

				iwl = append(iwl, bb.Instr[i])
			}
		}
	}

	return res
}

type bi struct {
	bb  *BasicBlock
	ind int
}

func (f *Fn) postorder() {
	if f.po != nil {
		return
	}

	seen := make([]bool, len(f.cfg.blocks))
	f.po = make([]*BasicBlock, 0, len(f.cfg.blocks))

	s := []bi{{bb: f.cfg.mp[0]}}
	for len(s) > 0 {
		ind := len(s) - 1
		curBi := s[ind]
		cur := curBi.bb
		if i := curBi.ind; i < len(cur.Succ) {
			s[ind].ind++
			bb := cur.Succ[i]
			if !seen[bb.ID] {
				seen[bb.ID] = true
				s = append(s, bi{bb: bb})
			}
			continue
		}

		s = s[:ind]
		f.po = append(f.po, cur)
	}
}
