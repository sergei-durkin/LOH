package lir

import (
	"fmt"
	"io"
)

func (l *Lir) Print(w io.Writer) {
	for i := 0; i < len(l.funcs); i++ {
		l.funcs[i].Print(w)
	}
}

func (f *Fn) Print(w io.Writer) {
	for i := len(f.po) - 1; i >= 0; i-- {
		bb := f.po[i]
		fmt.Fprintf(w, "L_%s_%d:\n", f.name, bb.ID)
		for j := 0; j < len(bb.Instr); j++ {
			bb.Instr[j].Print(w)
		}
	}
}
