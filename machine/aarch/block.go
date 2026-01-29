package aarch

import (
	"fmt"
	"io"
	"loh/machine"
)

func (a *AARCH) block(w io.Writer, b machine.Block) {
	fmt.Fprintf(w, "F_%s_%d:\n", a.context, b.Label())

	instructions := b.Instructions()
	for i := 0; i < len(instructions); i++ {
		a.instruction(w, instructions[i])
	}
	fmt.Fprintln(w)
}
