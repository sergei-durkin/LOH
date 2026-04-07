package risc

import (
	"fmt"
	"io"
	"loh/machine"
	"strings"
)

const epilog = "F_%s_END"

type RISCV struct {
	ir machine.IR

	frameSize int
	context   string

	labels map[string]int
	output []string
	asm    []asmInstr
	cnt    int
}

type asmInstr struct {
	s string

	label string
	jump  string

	base int
}

func (a asmInstr) String() string {
	return string(a.s)
}

var debug = false

func (a *RISCV) Emit(w io.Writer, ir machine.IR) {
	a.output = []string{""}
	a.labels = make(map[string]int)
	a.cnt = 0

	a.ir = ir
	a._start()
	funcs := ir.Functions()
	for i := 0; i < len(funcs); i++ {
		a.context = funcs[i].Name()
		a.frameSize = funcs[i].FrameSize()
		a.function(funcs[i])
	}
	a.print()

	rows := make([]string, 0, len(a.asm))
	for i, instr := range a.asm {
		if instr.jump != "" {
			offset := offset(i-instr.base, a.labels[instr.jump])
			sign := ""
			if instr.base == 0 {
				sign = ".+"
				if offset < 0 {
					sign = ".-"
					offset = -offset
				}
			}
			rows = append(rows, fmt.Sprintf("%s%s%#x", instr.s, sign, offset))
			continue
		}
		if len(instr.s) == 0 {
			continue
		}

		rows = append(rows, instr.s)
	}

	w.Write([]byte(strings.Join(rows, "\n")))
}

func offset(cur int, dst int) int {
	return (dst - cur) * 4
}

func (a *RISCV) writeln(i asmInstr) {
	if i.label != "" {
		a.labels[i.label] = a.cnt
		return
	}
	if i.label == "" {
		a.cnt++
	}

	a.asm = append(a.asm, i)
}
func (a *RISCV) writeRows(s string) {
	for _, str := range strings.Split(s, "\n") {
		str = strings.Trim(str, "\t")
		if str != "" {
			a.asm = append(a.asm, asmInstr{s: str})
			a.cnt++
		}
	}
}

func (a *RISCV) _start() {
	a.writeRows(`
.section .text
.globl _start
_start:
	LUI  x2, 0x0
	ADDI x2, x2, 0x1FC
	AUIPC x1, 0x0`)

	a.writeln(asmInstr{base: 1, s: "JALR x1, x1,", jump: "F_main"})

	a.writeRows(`JAL x0, .+0x0`)
}

func (a *RISCV) print() {
	a.writeln(asmInstr{label: "F_print"})

	a.writeRows(`
	ADDI x2, x2, -0x08
	SW x19, 0x0(x2)
	ADDI x2, x2, -0x08
	SW x20, 0x0(x2)
	ADDI x2, x2, -0x08
	SW x21, 0x0(x2)
	ADDI x2, x2, -0x08
	SW x22, 0x0(x2)
	ADDI x2, x2, -0x08
	SW x23, 0x0(x2)
	ADDI x2, x2, -0x08
	SW x24, 0x0(x2)
	ADDI x2, x2, -0x08
	SW x25, 0x0(x2)
	ADDI x2, x2, -0x08
	SW x26, 0x0(x2)
	ADDI x2, x2, -0x10
	SW x1, 0x0C(x2)
	SW x8, 0x08(x2)
	ADDI x8, x2, 0x0
	ADDI x2, x2, -0x10
	ADDI x20, x0, 0x0
	ADD x21, x0, x10
	ADD x26, x0, x11
	ADD x19, x0, x12
	ADDI x23, x19, 0x1
	ADD x24, x0, x23
	ADDI x24, x24, 0x0F
	ANDI x24, x24, -0x10
	SUB x2, x2, x24
	LUI x22, 0x07000
	ADD x22, x22, x26`)

	a.writeln(asmInstr{
		s:    "JAL x0, ",
		jump: "F_print_1",
	})

	a.writeln(asmInstr{label: "F_print_1"})

	a.writeln(asmInstr{
		s:    "BGE x20, x19,",
		jump: "F_print_4",
	})

	a.writeln(asmInstr{label: "F_print_2"})

	a.writeRows(`
	ADD x24, x22, x20
	ADD x23, x0, x24
	ADD x24, x21, x20
	LB x25, 0x0(x24)
	ADD x24, x0, x25
	SB x24, 0x0(x23)`)

	a.writeln(asmInstr{label: "F_print_3"})
	a.writeRows("ADDI x20, x20, 0x1")

	a.writeln(asmInstr{
		s:    "JAL x0, ",
		jump: "F_print_1",
	})

	a.writeln(asmInstr{label: "F_print_4"})

	a.writeRows(`
	ADDI x10, x0, 0x0
	ADD x2, x0, x8
	LW x1, 0x0C(x2)
	LW x8, 0x08(x2)
	ADDI x2, x2, 0x10
	LW x26, 0x0(x2)
	ADDI x2, x2, 0x08
	LW x25, 0x0(x2)
	ADDI x2, x2, 0x08
	LW x24, 0x0(x2)
	ADDI x2, x2, 0x08
	LW x23, 0x0(x2)
	ADDI x2, x2, 0x08
	LW x22, 0x0(x2)
	ADDI x2, x2, 0x08
	LW x21, 0x0(x2)
	ADDI x2, x2, 0x08
	LW x20, 0x0(x2)
	ADDI x2, x2, 0x08
	LW x19, 0x0(x2)
	ADDI x2, x2, 0x08
	RET`)
}
