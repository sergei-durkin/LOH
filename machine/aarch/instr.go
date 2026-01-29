package aarch

import (
	"fmt"
	"io"
	"loh/machine"
	"loh/token"
)

func (a *AARCH) instruction(w io.Writer, instr machine.Instruction) {
	operands := instr.Operands()
	switch instr.Type() {
	default:
		panic(fmt.Sprintf("undefined instruction: %T %+v", instr, instr))
	case machine.AND:
		a, b, c := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1]), operand(w, a.frameSize, operands[2])
		fmt.Fprintf(w, "\tAND %s, %s, %s\n", a, b, c)
	case machine.OR:
		a, b, c := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1]), operand(w, a.frameSize, operands[2])
		fmt.Fprintf(w, "\tORR %s, %s, %s\n", a, b, c)
	case machine.XOR:
		a, b, c := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1]), operand(w, a.frameSize, operands[2])
		fmt.Fprintf(w, "\tEOR %s, %s, %s\n", a, b, c)
	case machine.MUL:
		a, b, c := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1]), operand(w, a.frameSize, operands[2])
		fmt.Fprintf(w, "\tMUL %s, %s, %s\n", a, b, c)
	case machine.SUM:
		a, b, c := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1]), operand(w, a.frameSize, operands[2])
		fmt.Fprintf(w, "\tADD %s, %s, %s\n", a, b, c)
	case machine.DIV:
		a, b, c := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1]), operand(w, a.frameSize, operands[2])
		fmt.Fprintf(w, "\tSDIV %s, %s, %s\n", a, b, c)
	case machine.MOD:
		a, b, c := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1]), operand(w, a.frameSize, operands[2])
		fmt.Fprintf(w, "\tSDIV %s, %s, %s\n", a, b, c)
		fmt.Fprintf(w, "\tMSUB %s, %s, %s, %s\n", a, a, c, b)
	case machine.SUB:
		a, b, c := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1]), operand(w, a.frameSize, operands[2])
		fmt.Fprintf(w, "\tSUB %s, %s, %s\n", a, b, c)
	case machine.CALL:
		_ = operands
		target := operands[0]
		operands = operands[1:]
		var reg int = 0
		for _, arg := range operands {
			fmt.Fprintf(w, "\tMOV x%d, %s\n", reg, operand(w, a.frameSize, arg))
			reg++
		}

		fmt.Fprintf(w, "\tBL F_%s\n", instr.(machine.Calleeble).Callee())

		if target != nil {
			fmt.Fprintf(w, "\tMOV %s, x0\n", operand(w, a.frameSize, target))
		}
	case machine.JMP:
		label := instr.(machine.Labeled).Label()
		fmt.Fprintf(w, "\tB F_%s_%d\n", a.context, label)
	case machine.CMP:
		size := machine.RegisterSize(instr.(machine.Sizeble).Size())
		if size <= 0 {
			panic(size)
		}
		a, b, c := operand(w, a.frameSize, operands[0]), operandWithSize(operands[1], size), operandWithSize(operands[2], size)
		fmt.Fprintf(w, "\tCMP %s, %s\n", b, c)
		fmt.Fprintf(w, "\tCSET %s, %s\n", a, cset(instr.(machine.Tokened).Token()))
	case machine.CBZ:
		_ = operands
		label := instr.(machine.Labeled).Label()
		l := operand(w, a.frameSize, operands[0])
		fmt.Fprintf(w, "\tCBZ %s, F_%s_%d\n", l, a.context, label)
	case machine.MOV:
		a, b := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1])

		fmt.Fprintf(w, "\tMOV %s, %s\n", a, b)
	case machine.STR:
		_ = operands
		size := machine.RegisterSize(instr.(machine.Sizeble).Size())
		a, b := operandWithSize(operands[1], size), operand(w, a.frameSize, operands[0])

		fmt.Fprintf(w, "\t%s %s, [%s]\n", instrStoreSize(size), a, b)
	case machine.LDR:
		_ = operands
		size := machine.RegisterSize(instr.(machine.Sizeble).Size())
		if size <= 0 {
			panic(1)
		}
		a, b := operandWithSize(operands[0], size), operand(w, a.frameSize, operands[1])

		fmt.Fprintf(w, "\t%s %s, [%s]\n", instrLoadSize(size), a, b)
	case machine.RET:
		value := operand(w, a.frameSize, operands[0])

		fmt.Fprintf(w, "\tMOV x0, %s\n", value)
		fmt.Fprintf(w, "\tB %s\n", fmt.Sprintf(epilog, a.context))
	case machine.ALLOCA:
		a, b := operand(w, a.frameSize, operands[0]), operand(w, a.frameSize, operands[1])
		fmt.Fprintf(w, "\tADD %s, %s, #15\n", b, b)
		fmt.Fprintf(w, "\tAND %s, %s, #-16\n", b, b)
		fmt.Fprintf(w, "\tSUB sp, sp, %s\n", b)
		fmt.Fprintf(w, "\tMOV %s, sp\n", a)
	}
}

func operandWithSize(v machine.Value, size machine.RegisterSize) string {
	switch v := v.(type) {
	default:
		panic(fmt.Sprintf("undefined value: %T %+v", v, v))
	case *machine.ArgReg:
		return fmt.Sprintf("X%d", v.ID)
	case *machine.Reg:
		return fmt.Sprintf("%s%d", regsize(size), v.ID)
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

var temp = []int{11, 12, 13, 14, 15}

func tmp() int {
	t := temp[0]
	temp = temp[1:]
	return t
}
func free(t int) {
	temp = append(temp, t)
}

func operand(w io.Writer, frameSize int, v machine.Value) string {
	switch v := v.(type) {
	default:
		panic(fmt.Sprintf("undefined value: %T %+v", v, v))
	case *machine.ArgReg:
		return fmt.Sprintf("X%d", v.ID)
	case *machine.FP:
		t := tmp()
		defer free(t)

		fmt.Fprintf(w, "\tSUB X%d, X29, #%d\n", t, frameSize)
		return fmt.Sprintf("X%d", t)
	case *machine.Reg:
		return fmt.Sprintf("%s%d", regsize(v.Size), v.ID)
	case *machine.IntConst:
		if v.Int >= 4096 {
			t := tmp()
			defer free(t)

			emitMOV(w, fmt.Sprintf("X%d", t), v.Int)
			return fmt.Sprintf("X%d", t)
		}
		return fmt.Sprintf("#%d", v.Int)
	case *machine.BoolConst:
		r := 0
		if v.Bool {
			r = 1
		}

		return fmt.Sprintf("#%d", r)
	}
}

func regsize(s machine.RegisterSize) string {
	switch s {
	default:
		return "X"
		panic(fmt.Sprintf("undefined size: %T %+v", s, s))
	case machine.INT8:
		return "W"
	case machine.INT16:
		return "W"
	case machine.INT32:
		return "W"
	case machine.INT64:
		return "X"
	}
}

func instrStoreSize(s machine.RegisterSize) string {
	switch s {
	default:
		panic(fmt.Sprintf("undefined size: %T %+v", s, s))
		return "STR"
	case machine.INT8:
		return "STRB"
	case machine.INT16:
		return "STRH"
	case machine.INT32:
		return "STR"
	case machine.INT64:
		return "STR"
	}
}

func instrLoadSize(s machine.RegisterSize) string {
	switch s {
	default:
		return "LDR"
		panic(fmt.Sprintf("undefined size: %T %+v", s, s))
	case machine.INT8:
		return "LDRB"
	case machine.INT16:
		return "LDRH"
	case machine.INT32:
		return "LDR"
	case machine.INT64:
		return "LDR"
	}
}

func cset(t token.Token) string {
	switch t {
	default:
		panic(t)
	case token.GT:
		return "GT"
	case token.GTE:
		return "GE"
	case token.LT:
		return "LT"
	case token.LTE:
		return "LE"
	case token.EQ:
		return "EQ"
	case token.NE:
		return "NE"
	}
}

func emitMOV(w io.Writer, reg string, val int64) {
	if val >= 0 && val <= 65535 {
		fmt.Fprintf(w, "\tMOV %s, #%d\n", reg, val)
		return
	}
	first := true
	for shift := 0; shift < 64; shift += 16 {
		chunk := (val >> shift) & 0xFFFF
		if chunk == 0 {
			continue
		}
		if first {
			fmt.Fprintf(w, "\tMOVZ %s, #0x%X, LSL #%d\n", reg, chunk, shift)
			first = false
		} else {
			fmt.Fprintf(w, "\tMOVK %s, #0x%X, LSL #%d\n", reg, chunk, shift)
		}
	}
	if first { // val == 0
		fmt.Fprintf(w, "\tMOV %s, #0\n", reg)
	}
}
