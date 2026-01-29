package aarch

import (
	"fmt"
	"io"
	"loh/machine"
)

const epilog = "F_%s_END"

type AARCH struct {
	ir machine.IR

	frameSize int
	context   string
}

func (a *AARCH) Emit(w io.Writer, ir machine.IR) {
	a.ir = ir
	a._start(w)
	funcs := ir.Functions()
	for i := 0; i < len(funcs); i++ {
		a.context = funcs[i].Name()
		a.frameSize = funcs[i].FrameSize()
		a.function(w, funcs[i])
	}
	a.syscall(w)
}

func (a *AARCH) _start(w io.Writer) {
	fmt.Fprint(w, `
	.text
	.global main
	.align 2

main:
	B F_main
`)
}
func (a *AARCH) syscall(w io.Writer) {
	fmt.Fprint(w, `
F_syscall:
	STP X29, X30, [SP, #-16]!
	MOV X29, SP

	MOV X16, X0
	MOV X0, X1
	MOV X1, X2
	MOV X2, X3
	MOV X3, X4
	MOV X4, X5
	MOV X5, X6
	MOV X6, X7

	SVC #0x80

	LDP X29, X30, [SP], #16
	RET
`)
}
