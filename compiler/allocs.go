package compiler

import (
	"fmt"
	"loh/token"
	"loh/types"
)

func (f *Fn) allocs() {
	aliased := f.aliased()

	tmp := 42
	offset := 0
	frameIndex := make(map[string]int)

	for _, bb := range f.cfg.blocks {
		instrs := make([]Instruction, 0, len(bb.Instr))
		for _, instr := range bb.Instr {
			tmp++

			instrs = append(instrs, instr)
			switch instr := instr.(type) {
			case *Assign:
				target := instr.target
				if !aliased[target.Name()] {
					continue
				}

				if ptr, ok := frameIndex[target.Name()]; ok {
					addr := &TempVar{label: fmt.Sprintf("%s_%d", target.Label(), tmp)}
					instrs[len(instrs)-1] = &Assign{
						target: addr,
						op:     token.PLUS,
						arg1:   &FP{},
						arg2:   &IntConst{int: int64(ptr)},

						size: int(types.Pointer.Size()),
					}
					instrs = append(instrs, &Store{
						Value:       instr.arg1,
						Destination: addr,
						Size:        int(f.cfg.size[target.Name()]),
					})
					offset += aligned(0, 8)
				} else {
					frameIndex[target.Name()] = offset
					addr := &TempVar{label: fmt.Sprintf("%s_%d", target.Label(), tmp)}
					instrs = append(instrs, &Assign{
						target: addr,
						op:     token.PLUS,
						arg1:   &FP{},
						arg2:   &IntConst{int: int64(offset)},

						size: int(types.Pointer.Size()),
					})
					instrs = append(instrs, &Store{
						Value:       target.(Value),
						Destination: addr,
						Size:        int(f.cfg.size[target.Name()]),
					})
					offset += aligned(0, 8)
				}

			case *Alloca:
				// switch ln := instr.len.(type) {
				// default:
				// 	panic("not implemented")
				// 	// case *IntConst:
				// 	// 	frameIndex[instr.ptr.(Variable).Name()] = offset
				// 	// 	offset += aligned(int(instr.elSize*ln.Int()), 16)
				// }
			}
		}

		for _, instr := range instrs {
			switch instr := instr.(type) {
			case *Assign:
				if arg, ok := instr.arg1.(*AddressOf); ok {
					if arg, ok := arg.Target.(Variable); ok {
						if f, ok := frameIndex[arg.Name()]; ok {
							instr.op = token.PLUS
							instr.arg1 = &FP{}
							instr.arg2 = &IntConst{int: int64(f)}
						} else {
							instr.arg2 = arg.(Value)
							//panic(arg)
						}
					}
				}
				if _, ok := instr.arg2.(*AddressOf); ok {
					panic(1)
				}
			}
		}

		bb.Instr = instrs
	}

	f.frameIndex = frameIndex
	f.frameSize = offset
}

func aligned(size int, align int) int {
	size = max(size, align)
	if size%align != 0 {
		size += align - size%align
	}

	return size
}
