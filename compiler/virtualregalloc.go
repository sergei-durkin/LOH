package compiler

import (
	"fmt"
)

type VirtualRegister int

func (f *Fn) virtualregalloc() {
	cur := VirtualRegister(0)
	allocated := make(map[string]VirtualRegister)

	var getReg func(v Value) Value
	getReg = func(v Value) Value {
		if v, ok := v.(Variable); ok {
			if _, ok := allocated[v.Label()]; !ok {
				allocated[v.Label()] = cur
				cur++
			}
			return &Reg{ID: int(allocated[v.Label()])}
		}

		if arg, ok := v.(*AddressOf); ok {
			arg.Target = getReg(arg.Target)
		}

		if arg, ok := v.(*Dereference); ok {
			arg.Addr = getReg(arg.Addr)
		}

		return v
	}

	for i := len(f.po) - 1; i >= 0; i-- {
		bb := f.po[i]

		for j := 0; j < len(bb.Instr); j++ {
			instr := bb.Instr[j]
			switch in := instr.(type) {
			default:
				panic(fmt.Sprintf("unknown instruction: %T %+v", in, in))
			case *Alloca:
				getReg(in.ptr)
				in.size = getReg(in.size)
			case *IfGoto:
				in.cond = getReg(in.cond)
			case *Goto:
			case *Store:
				in.Destination = getReg(in.Destination)
				in.Value = getReg(in.Value)
			case *Assign:
				getReg(in.target.(Value))
				in.arg1 = getReg(in.arg1)
				in.arg2 = getReg(in.arg2)
			case *Call:
				getReg(in.target.(Value))

				for k := 0; k < len(in.args); k++ {
					in.args[k] = getReg(in.args[k])
				}
			case *Load:
				in.Destination = getReg(in.Destination)
				in.Source = getReg(in.Source)

			case *Return:
				in.value = getReg(in.value)
			}
		}
	}

	f.virtualRegs = allocated
}
