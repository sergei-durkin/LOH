package lir

import (
	"fmt"
	"loh/compiler"
	"loh/machine"
	"maps"
	"slices"
	"sort"
)

type liveRange struct {
	reg      machine.PhysicalRegister
	from, to int
}

type blockLiveness struct {
	use map[compiler.VirtualRegister]bool
	def map[compiler.VirtualRegister]bool
	in  map[compiler.VirtualRegister]bool
	out map[compiler.VirtualRegister]bool
}

func (f *Fn) liveness() map[compiler.VirtualRegister]*liveRange {
	blocks := make([]*blockLiveness, len(f.blocks))
	for i := 0; i < len(blocks); i++ {
		blocks[i] = &blockLiveness{
			use: make(map[compiler.VirtualRegister]bool),
			def: make(map[compiler.VirtualRegister]bool),
			in:  make(map[compiler.VirtualRegister]bool),
			out: make(map[compiler.VirtualRegister]bool),
		}
	}

	for i := 0; i < len(f.po); i++ {
		bb := f.po[i]
		bl := blocks[bb.ID]

		for _, instr := range bb.Instr {
			var target machine.Value
			ops := instr.Operands()
			if len(ops) > 1 {
				target = ops[0]
				ops = ops[1:]
			}

			for j := 0; j < len(ops); j++ {
				if v, ok := ops[j].(*machine.Reg); ok {
					if !bl.def[compiler.VirtualRegister(v.ID)] {
						bl.use[compiler.VirtualRegister(v.ID)] = true
					}
				}
			}

			if target, ok := target.(*machine.Reg); ok {
				bl.def[compiler.VirtualRegister(target.ID)] = true
			}
		}
	}

	wl := make([]*BasicBlock, 0)
	wl = append(wl, f.po...)

	union := func(a, b map[compiler.VirtualRegister]bool) map[compiler.VirtualRegister]bool {
		res := make(map[compiler.VirtualRegister]bool)

		for k := range a {
			res[k] = true
		}
		for k := range b {
			res[k] = true
		}

		return res
	}

	setMinus := func(a, b map[compiler.VirtualRegister]bool) map[compiler.VirtualRegister]bool {
		result := make(map[compiler.VirtualRegister]bool)
		for k := range a {
			if !b[k] {
				result[k] = true
			}
		}
		return result
	}

	for len(wl) > 0 {
		bb := wl[len(wl)-1]
		wl = wl[:len(wl)-1]
		bl := blocks[bb.ID]

		newOut := make(map[compiler.VirtualRegister]bool)
		for i := 0; i < len(bb.Succ); i++ {
			newOut = union(newOut, blocks[bb.Succ[i].ID].in)
		}

		newIn := union(bl.use, setMinus(newOut, bl.def))

		if !maps.Equal(bl.in, newIn) || !maps.Equal(bl.out, newOut) {
			bl.in = newIn
			bl.out = newOut

			wl = append(wl, bb.Pred...)
		}
	}

	blockStart := make([]int, len(f.po))
	blockEnd := make([]int, len(f.po))

	instrPos := make(map[Instruction]int)
	pos := 0
	for i := len(f.po) - 1; i >= 0; i-- {
		bb := f.po[i]
		blockStart[bb.ID] = pos

		for j := 0; j < len(bb.Instr); j++ {
			instrPos[bb.Instr[j]] = pos
			pos += 2
		}

		blockEnd[bb.ID] = pos
		if len(bb.Instr) > 0 {
			blockEnd[bb.ID] -= 2
		}
	}

	ranges := make(map[compiler.VirtualRegister]*liveRange)

	ensure := func(id compiler.VirtualRegister) *liveRange {
		lr, ok := ranges[id]
		if !ok {
			lr = &liveRange{
				from: 1 << 30,
				to:   -1,
			}
			ranges[id] = lr
		}

		return lr
	}

	for _, bb := range f.po {
		bl := blocks[bb.ID]

		for name := range bl.in {
			lr := ensure(name)
			lr.from = min(lr.from, blockStart[bb.ID])
			lr.to = max(lr.to, blockEnd[bb.ID])
		}
		for name := range bl.out {
			lr := ensure(name)
			lr.from = min(lr.from, blockStart[bb.ID])
			lr.to = max(lr.to, blockEnd[bb.ID])
		}

		for _, instr := range bb.Instr {
			var target machine.Value
			ops := instr.Operands()
			if _, ok := instr.(*STR); !ok && len(ops) > 1 {
				target = ops[0]
				ops = ops[1:]
			}

			if target, ok := target.(*machine.Reg); ok {
				lr := ensure(compiler.VirtualRegister(target.ID))
				lr.from = min(lr.from, instrPos[instr])
			}

			for _, op := range ops {
				switch v := op.(type) {
				default:
					panic(fmt.Sprintf("unknown operation: %T %+v", op, op))
				case *machine.IntConst:
				case *machine.BoolConst:
				case *machine.FP:
				case *machine.ArgReg:
				case *machine.Reg:
					lr := ensure(compiler.VirtualRegister(v.ID))
					lr.to = max(lr.to, instrPos[instr])
				}
			}
		}
	}

	// for id, bl := range blocks {
	// 	fmt.Printf("BB%d\n", id)
	//
	// 	fmt.Printf("\tDEF: ")
	// 	for def := range bl.def {
	// 		fmt.Printf("%d, ", def)
	// 	}
	// 	fmt.Print("\n")
	//
	// 	fmt.Printf("\tUSE: ")
	// 	for def := range bl.use {
	// 		fmt.Printf("%d, ", def)
	// 	}
	// 	fmt.Print("\n")
	//
	// 	fmt.Printf("\tIN: ")
	// 	for def := range bl.in {
	// 		fmt.Printf("%d, ", def)
	// 	}
	// 	fmt.Print("\n")
	//
	// 	fmt.Printf("\tOUT: ")
	// 	for def := range bl.out {
	// 		fmt.Printf("%d, ", def)
	// 	}
	// 	fmt.Print("\n")
	// }

	return ranges
}

func (f *Fn) regalloc() {
	liveness := f.liveness()

	stp := []machine.PhysicalRegister{}
	stpdedup := make(map[machine.PhysicalRegister]struct{})

	free := []machine.PhysicalRegister{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9,
	}

	ranges := make([]*liveRange, 0, len(liveness))
	for _, lr := range liveness {
		ranges = append(ranges, lr)
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].from < ranges[j].from
	})

	active := []*liveRange{}

	for _, lr := range ranges {
		active = slices.DeleteFunc(active, func(a *liveRange) bool {
			if a.to < lr.from {
				free = append(free, a.reg)
				return true
			}
			return false
		})

		if len(free) == 0 {
			panic("todo spills")
		} else {
			reg := free[len(free)-1]
			free = free[:len(free)-1]
			lr.reg = reg

			if _, ok := stpdedup[reg]; !ok {
				stpdedup[reg] = struct{}{}
				stp = append(stp, reg)
			}

			active = append(active, lr)
		}
	}

	for i := len(f.po) - 1; i >= 0; i-- {
		bb := f.po[i]

		for j := 0; j < len(bb.Instr); j++ {
			instr := bb.Instr[j]
			operands := instr.Operands()
			for k := 0; k < len(operands); k++ {
				switch op := operands[k].(type) {
				default:
					panic(fmt.Sprintf("undefined operand: %T %+v", op, op))
				case *machine.IntConst:
				case *machine.BoolConst:
				case *machine.FP:
				case *machine.ArgReg:
				case *machine.Reg:
					r := liveness[compiler.VirtualRegister(op.ID)].reg
					operands[k] = &machine.Reg{ID: int(r)}
				}
			}
			if repl, ok := instr.(Replaceable); ok {
				repl.ReplaceAll(operands...)
			}
		}
	}

	f.regs = stp
}
