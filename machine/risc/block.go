package risc

import (
	"fmt"
	"loh/machine"
)

func (a *RISCV) block(b machine.Block) {
	a.writeln(asmInstr{label: fmt.Sprintf("F_%s_%d", a.context, b.Label())})

	instructions := b.Instructions()
	for i := 0; i < len(instructions); i++ {
		a.instruction(instructions[i])
	}
}
