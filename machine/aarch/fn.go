package aarch

import (
	"fmt"
	"io"
	"loh/machine"
)

func (a *AARCH) saveRegisters(w io.Writer, fn machine.Function) {
	regs := fn.CalleeSavedRegisters()
	for len(regs) > 0 {
		if len(regs) == 1 {
			fmt.Fprintf(w, "\tSTR x%d, [SP, -16]!\n", regs[0])
			break
		}

		r1, r2 := regs[0], regs[1]
		regs = regs[2:]

		fmt.Fprintf(w, "\tSTP x%d, x%d, [SP, -16]!\n", r1, r2)
	}
}

func (a *AARCH) restoreRegisters(w io.Writer, fn machine.Function) {
	regs := append([]machine.PhysicalRegister{}, fn.CalleeSavedRegisters()...)

	res := []string{}
	for len(regs) > 0 {
		if len(regs) == 1 {
			res = append(res, fmt.Sprintf("\tLDR x%d, [SP], 16\n", regs[0]))
			break
		}

		r1, r2 := regs[0], regs[1]
		regs = regs[2:]

		res = append(res, fmt.Sprintf("\tLDP x%d, x%d, [SP], 16\n", r1, r2))
	}

	for i := len(res) - 1; i >= 0; i-- {
		fmt.Fprint(w, res[i])
	}

	fmt.Fprintln(w)
}

func (a *AARCH) prolog(w io.Writer, fn machine.Function) {
	fmt.Fprintf(w, "F_%s:\n", fn.Name())

	a.saveRegisters(w, fn)

	fmt.Fprint(w, "\tSTP X29, x30, [SP, -16]!\n")
	fmt.Fprint(w, "\tMOV X29, SP\n")
	fmt.Fprintln(w)

	if fn.FrameSize() > 0 {
		fmt.Fprintf(w, "\tSUB SP, SP, #%d\n", fn.FrameSize())
	}

	fmt.Fprintln(w)
}

func (a *AARCH) function(w io.Writer, fn machine.Function) {
	a.prolog(w, fn)

	blocks := fn.Blocks()
	for i := 0; i < len(blocks); i++ {
		a.block(w, blocks[i])
	}
	fmt.Fprintln(w)

	a.epilog(w, fn)
}

func (a *AARCH) epilog(w io.Writer, fn machine.Function) {
	fmt.Fprintf(w, "%s:\n", fmt.Sprintf(epilog, a.context))

	fmt.Fprint(w, "\tMOV SP, X29\n")

	fmt.Fprint(w, "\tLDP X29, x30, [SP], 16\n")

	a.restoreRegisters(w, fn)

	fmt.Fprint(w, "\tRET\n")
	fmt.Fprintln(w)
}
