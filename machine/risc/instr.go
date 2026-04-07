package risc

import (
	"fmt"
	"loh/machine"
	"loh/token"
)

func (a *RISCV) instruction(instr machine.Instruction) {
	operands := instr.Operands()
	switch instr.Type() {
	default:
		panic(fmt.Sprintf("undefined instruction: %T %+v", instr, instr))
	case machine.AND:
		d, l, r := a.threeOp(operands)
		if r[0] == 'x' {
			a.writeln(asmInstr{s: fmt.Sprintf("AND %s, %s, %s", d, l, r)})
		} else {
			a.writeln(asmInstr{s: fmt.Sprintf("ANDI %s, %s, %s", d, l, r)})
		}
	case machine.OR:
		d, l, r := a.threeOp(operands)

		if r[0] == 'x' {
			a.writeln(asmInstr{s: fmt.Sprintf("OR %s, %s, %s", d, l, r)})
		} else {
			a.writeln(asmInstr{s: fmt.Sprintf("ORI %s, %s, %s", d, l, r)})
		}
	case machine.XOR:
		d, l, r := a.threeOp(operands)

		if r[0] == 'x' {
			a.writeln(asmInstr{s: fmt.Sprintf("XOR %s, %s, %s", d, l, r)})
		} else {
			a.writeln(asmInstr{s: fmt.Sprintf("XORI %s, %s, %s", d, l, r)})
		}
	case machine.SUM:
		d, l, r := a.threeOp(operands)

		if r[0] == 'x' {
			a.writeln(asmInstr{s: fmt.Sprintf("ADD %s, %s, %s", d, l, r)})
		} else {
			a.writeln(asmInstr{s: fmt.Sprintf("ADDI %s, %s, %s", d, l, r)})
		}
	case machine.SUB:
		d, l, r := a.threeOp(operands)

		if r[0] == 'x' {
			a.writeln(asmInstr{s: fmt.Sprintf("SUB %s, %s, %s", d, l, r)})
		} else {
			a.writeln(asmInstr{s: fmt.Sprintf("ADDI %s, %s, -%s", d, l, r)})
		}
	case machine.CALL:
		_ = operands
		target := operands[0]
		operands = operands[1:]
		var reg = 0
		for _, arg := range operands {
			op := a.operand(a.frameSize, arg)
			if op[0] == 'x' {
				a.writeln(asmInstr{s: fmt.Sprintf("ADD %s, x0, %s", ArgRegister(reg), op)})
			} else {
				a.writeln(asmInstr{s: fmt.Sprintf("ADDI %s, x0, %s", ArgRegister(reg), op)})
			}
			reg++
		}

		a.writeln(asmInstr{s: "AUIPC x1, 0x0"})
		a.writeln(asmInstr{base: 1, s: "JALR x1, x1,", jump: fmt.Sprintf("F_%s", instr.(machine.Calleeble).Callee())})

		if target != nil {
			a.writeln(asmInstr{s: fmt.Sprintf("ADD %s, x0, %s", a.operand(a.frameSize, target), regsMap["A0"])})
		}
	case machine.JMP:
		label := instr.(machine.Labeled).Label()
		a.writeln(asmInstr{s: "JAL x0, ", jump: fmt.Sprintf("F_%s_%d", a.context, label)})
	case machine.CMP:
		size := machine.RegisterSize(instr.(machine.Sizeble).Size())
		if size <= 0 {
			panic(size)
		}
		var i string
		_, l, r := a.threeOp(operands)
		i, l, r = cmp(instr.(machine.Tokened).Token(), l, r)
		if l[0] != 'x' {
			t := tmp()
			defer free(t)

			a.writeln(asmInstr{s: fmt.Sprintf("ADDI x%d, x0, %s", t, l)})
			l = fmt.Sprintf("x%d", t)
		}
		if r[0] != 'x' {
			t := tmp()
			defer free(t)

			a.writeln(asmInstr{s: fmt.Sprintf("ADDI x%d, x0, %s", t, r)})
			r = fmt.Sprintf("x%d", t)
		}

		a.writeln(asmInstr{s: fmt.Sprintf("%s %s, %s, .+0x8", i, l, r)})
	case machine.CBZ:
		label := instr.(machine.Labeled).Label()
		a.writeln(asmInstr{s: "JAL x0, ", jump: fmt.Sprintf("F_%s_%d", a.context, label)})
	case machine.MOV:
		l, r := a.operand(a.frameSize, operands[0]), a.operand(a.frameSize, operands[1])

		if r[0] == 'x' {
			a.writeln(asmInstr{s: fmt.Sprintf("ADD %s, x0, %s", l, r)})
		} else {
			a.writeln(asmInstr{s: fmt.Sprintf("ADDI %s, x0, %s", l, r)})
		}
	case machine.STR:
		_ = operands
		size := machine.RegisterSize(instr.(machine.Sizeble).Size())
		l, r := operandWithSize(operands[1], size), a.operand(a.frameSize, operands[0])

		a.writeln(asmInstr{s: fmt.Sprintf("%s %s, 0x0(%s)", instrStoreSize(size), l, r)})
	case machine.LDR:
		_ = operands
		size := machine.RegisterSize(instr.(machine.Sizeble).Size())
		if size <= 0 {
			panic(1)
		}
		l, r := operandWithSize(operands[0], size), a.operand(a.frameSize, operands[1])

		a.writeln(asmInstr{s: fmt.Sprintf("%s %s, 0x0(%s)", instrLoadSize(size), l, r)})
	case machine.RET:
		value := a.operand(a.frameSize, operands[0])

		if value[0] == 'x' {
			a.writeln(asmInstr{s: fmt.Sprintf("ADD %s, x0, %s", regsMap["A0"], value)})
		} else {
			a.writeln(asmInstr{s: fmt.Sprintf("ADDI %s, x0, %s", regsMap["A0"], value)})
		}
		a.writeln(asmInstr{s: "JAL x0, ", jump: fmt.Sprintf(epilog, a.context)})
	case machine.ALLOCA:
		l, r := a.operand(a.frameSize, operands[0]), a.operand(a.frameSize, operands[1])
		a.writeln(asmInstr{s: fmt.Sprintf("ADDI %s, %s, 0x0F", r, r)})
		a.writeln(asmInstr{s: fmt.Sprintf("ANDI %s, %s, -0x10", r, r)})
		a.writeln(asmInstr{s: fmt.Sprintf("SUB %s, %s, %s", regsMap["SP"], regsMap["SP"], r)})
		a.writeln(asmInstr{s: fmt.Sprintf("ADD %s, x0, %s", l, regsMap["SP"])})
	}
}

func operandWithSize(v machine.Value, size machine.RegisterSize) string {
	switch v := v.(type) {
	default:
		panic(fmt.Sprintf("undefined value: %T %+v", v, v))
	case *machine.ArgReg:
		return ArgRegister(v.ID)
	case *machine.Reg:
		return fmt.Sprintf("%s%s", regsize(size), CalleeSavedRegister(v.ID)[1:])
	case *machine.IntConst:
		return fmt.Sprintf("#%d", v.Int)
	case *machine.BoolConst:
		r := 0
		if v.Bool {
			r = 1
		}

		return fmt.Sprintf("#%d", r)
	}
}

var temp = []int{5, 6, 7, 28, 29, 30, 31}

func tmp() int {
	t := temp[0]
	temp = temp[1:]
	return t
}
func free(t int) {
	temp = append(temp, t)
}

func (a *RISCV) operand(frameSize int, v machine.Value) string {
	switch v := v.(type) {
	default:
		panic(fmt.Sprintf("undefined value: %T %+v", v, v))
	case *machine.ArgReg:
		return ArgRegister(v.ID)
	case *machine.FP:
		t := tmp()
		defer free(t)

		a.writeln(asmInstr{s: fmt.Sprintf("ADDI x%d, %s, -%#x", t, regsMap["FP"], frameSize)})
		return fmt.Sprintf("x%d", t)
	case *machine.Reg:
		return CalleeSavedRegister(v.ID)
	case *machine.IntConst:
		if v.Int >= 4096 {
			t := tmp()
			defer free(t)

			a.emitMOV(fmt.Sprintf("x%d", t), v.Int)
			return fmt.Sprintf("x%d", t)
		}
		return fmt.Sprintf("%#x", v.Int)
	case *machine.BoolConst:
		r := 0
		if v.Bool {
			r = 1
		}

		return fmt.Sprintf("%#x", r)
	}
}

func (a *RISCV) threeOp(operands []machine.Value) (string, string, string) {
	return a.operand(a.frameSize, operands[0]),
		a.operand(a.frameSize, operands[1]),
		a.operand(a.frameSize, operands[2])
}

func regsize(s machine.RegisterSize) string {
	return "x"
}

func instrStoreSize(s machine.RegisterSize) string {
	switch s {
	default:
		panic(fmt.Sprintf("undefined size: %T %+v", s, s))
	case machine.INT8:
		return "SB"
	case machine.INT16:
		return "SH"
	case machine.INT32:
		return "SW"
	}
}

func instrLoadSize(s machine.RegisterSize) string {
	switch s {
	default:
		return "LW"
		panic(fmt.Sprintf("undefined size: %T %+v", s, s))
	case machine.INT8:
		return "LB"
	case machine.INT16:
		return "LH"
	case machine.INT32:
		return "LW"
	}
}

func cmp(t token.Token, l, r string) (instr, a, b string) {
	switch t {
	default:
		panic(t)
	case token.GT:
		return "BLT", r, l
	case token.GTE:
		return "BLE", r, l
	case token.LT:
		return "BGT", r, l
	case token.LTE:
		return "BGE", r, l
	case token.EQ:
		return "BEQ", l, r
	case token.NE:
		return "BNE", l, r
	}
}

func (a *RISCV) emitMOV(reg string, val int64) {
	if val >= 0 && val <= 4095 {
		a.writeln(asmInstr{s: fmt.Sprintf("ADDI %s, x0, %#x", reg, val)})
		return
	}
	panic("not implemented")
}
