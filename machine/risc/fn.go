package risc

import (
	"fmt"
	"loh/machine"
)

func (a *RISCV) saveRegisters(fn machine.Function) {
	regs := fn.CalleeSavedRegisters()
	for len(regs) > 0 {
		r1 := CalleeSavedRegister(int(regs[0]))
		regs = regs[1:]

		a.writeln(asmInstr{s: "ADDI x2, x2, -0x08"})
		a.writeln(asmInstr{s: fmt.Sprintf("SW %s, 0x0(x2)", r1)})
	}
}

func (a *RISCV) restoreRegisters(fn machine.Function) {
	regs := append([]machine.PhysicalRegister{}, fn.CalleeSavedRegisters()...)

	for len(regs) > 0 {
		r1 := CalleeSavedRegister(int(regs[len(regs)-1]))
		regs = regs[:len(regs)-1]

		a.writeln(asmInstr{s: fmt.Sprintf("LW %s, 0x0(x2)", r1)})
		a.writeln(asmInstr{s: "ADDI x2, x2, 0x08"})
	}
}

func (a *RISCV) prolog(fn machine.Function) {
	a.writeln(asmInstr{label: fmt.Sprintf("F_%s", fn.Name())})

	a.saveRegisters(fn)

	a.writeRows(`
	ADDI x2, x2, -0x10
	SW x1, 0x0C(x2)
	SW x8, 0x08(x2)
	ADDI x8, x2, 0x0`)

	if fn.FrameSize() > 0 {
		a.writeln(asmInstr{s: fmt.Sprintf("ADDI x2, x2, %#x", -fn.FrameSize())})
	}
}

func (a *RISCV) function(fn machine.Function) {
	a.prolog(fn)

	blocks := fn.Blocks()
	for i := 0; i < len(blocks); i++ {
		a.block(blocks[i])
	}

	a.epilog(fn)
}

func (a *RISCV) epilog(fn machine.Function) {
	a.writeln(asmInstr{label: fmt.Sprintf(epilog, a.context)})

	a.writeRows(`
	ADD x2, x0, x8
	LW x1, 0x0C(x2)
	LW x8, 0x08(x2)
	ADDI x2, x2, 0x10`)

	a.restoreRegisters(fn)

	a.writeln(asmInstr{s: "RET"})
}
