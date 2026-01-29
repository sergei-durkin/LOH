package lir

import (
	"fmt"
	"io"
	"loh/machine"
	"loh/token"
	"slices"
	"strings"
)

type Instruction interface {
	Operands() []machine.Value
	instr()
	Print(w io.Writer)
	Type() machine.InstructionType
}

type Replaceable interface {
	ReplaceAll(operands ...machine.Value)
}

type MUL struct {
	destination, left, right machine.Value
}

func (i *MUL) Type() machine.InstructionType {
	return machine.MUL
}

func (i *MUL) instr() {}
func (i *MUL) Print(w io.Writer) {
	fmt.Fprintf(w, "\tMUL %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *MUL) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}
func (i *MUL) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type SUM struct {
	destination, left, right machine.Value
}

func (i *SUM) Type() machine.InstructionType {
	return machine.SUM
}

func (i *SUM) instr() {}
func (i *SUM) Print(w io.Writer) {
	fmt.Fprintf(w, "\tSUM %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *SUM) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}
func (i *SUM) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type DIV struct {
	destination, left, right machine.Value
}

func (i *DIV) Type() machine.InstructionType {
	return machine.DIV
}

func (i *DIV) instr() {}
func (i *DIV) Print(w io.Writer) {
	fmt.Fprintf(w, "\tDIV %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *DIV) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}

func (i *DIV) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type XOR struct {
	destination, left, right machine.Value
}

func (i *XOR) Type() machine.InstructionType {
	return machine.XOR
}

func (i *XOR) instr() {}
func (i *XOR) Print(w io.Writer) {
	fmt.Fprintf(w, "\tXOR %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *XOR) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}
func (i *XOR) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type OR struct {
	destination, left, right machine.Value
}

func (i *OR) Type() machine.InstructionType {
	return machine.OR
}

func (i *OR) instr() {}
func (i *OR) Print(w io.Writer) {
	fmt.Fprintf(w, "\tOR %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *OR) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}
func (i *OR) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type AND struct {
	destination, left, right machine.Value
}

func (i *AND) Type() machine.InstructionType {
	return machine.AND
}

func (i *AND) instr() {}
func (i *AND) Print(w io.Writer) {
	fmt.Fprintf(w, "\tAND %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *AND) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}
func (i *AND) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type MOD struct {
	destination, left, right machine.Value
}

func (i *MOD) Type() machine.InstructionType {
	return machine.MOD
}

func (i *MOD) instr() {}
func (i *MOD) Print(w io.Writer) {
	fmt.Fprintf(w, "\tMOD %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *MOD) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}
func (i *MOD) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type SUB struct {
	destination, left, right machine.Value
}

func (i *SUB) Type() machine.InstructionType {
	return machine.SUB
}

func (i *SUB) instr() {}
func (i *SUB) Print(w io.Writer) {
	fmt.Fprintf(w, "\tSUB %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *SUB) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}
func (i *SUB) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type CALL struct {
	callee string

	destination machine.Value
	args        []machine.Value
}

func (i *CALL) Callee() string {
	return i.callee
}

func (i *CALL) Type() machine.InstructionType {
	return machine.CALL
}

func (i *CALL) instr() {}
func (i *CALL) Print(w io.Writer) {
	args := []string{}
	for j := 0; j < len(i.args); j++ {
		args = append(args, toStr(i.args[j]))
	}
	fmt.Fprintf(w, "\tCALL %s = %s (%s);\n", toStr(i.destination), i.callee, strings.Join(args, ","))
}

func (i *CALL) Operands() []machine.Value {
	return append([]machine.Value{i.destination}, i.args...)
}
func (i *CALL) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.args = slices.Clone(operands[1:])
}

type JMP struct {
	label int
	fn    string
}

func (i *JMP) Label() int {
	return i.label
}
func (i *JMP) Type() machine.InstructionType {
	return machine.JMP
}

func (i *JMP) instr() {}
func (i *JMP) Print(w io.Writer) {
	fmt.Fprintf(w, "\tJMP L_%s_%d;\n", i.fn, i.label)
}

func (i *JMP) Operands() []machine.Value {
	return nil
}

type CMP struct {
	destination, left, right machine.Value
	token                    token.Token
	size                     int
}

func (i *CMP) Type() machine.InstructionType {
	return machine.CMP
}

func (i *CMP) Token() token.Token {
	return i.token
}

func (i *CMP) Size() int {
	return i.size
}

func (i *CMP) instr() {}
func (i *CMP) Print(w io.Writer) {
	fmt.Fprintf(w, "\tCMP %s = %s, %s;\n", toStr(i.destination), toStr(i.left), toStr(i.right))
}

func (i *CMP) Operands() []machine.Value {
	return []machine.Value{i.destination, i.left, i.right}
}
func (i *CMP) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.left = operands[1]
	i.right = operands[2]
}

type CBZ struct {
	label int
	fn    string
	left  machine.Value
}

func (i *CBZ) Label() int {
	return i.label
}
func (i *CBZ) Type() machine.InstructionType {
	return machine.CBZ
}

func (i *CBZ) instr() {}
func (i *CBZ) Print(w io.Writer) {
	fmt.Fprintf(w, "\tCBZ %s, L_%s_%d;\n", toStr(i.left), i.fn, i.label)
}

func (i *CBZ) Operands() []machine.Value {
	return []machine.Value{i.left}
}
func (i *CBZ) ReplaceAll(operands ...machine.Value) {
	i.left = operands[0]
}

type MOV struct {
	destination, source machine.Value
}

func (i *MOV) Type() machine.InstructionType {
	return machine.MOV
}

func (i *MOV) instr() {}
func (i *MOV) Print(w io.Writer) {
	fmt.Fprintf(w, "\tMOV %s, %s;\n", toStr(i.destination), toStr(i.source))
}

func (i *MOV) Operands() []machine.Value {
	return []machine.Value{i.destination, i.source}
}
func (i *MOV) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.source = operands[1]
}

type ALLOCA struct {
	destination, size machine.Value
}

func (i *ALLOCA) Type() machine.InstructionType {
	return machine.ALLOCA
}

func (i *ALLOCA) instr() {}
func (i *ALLOCA) Print(w io.Writer) {
	fmt.Fprintf(w, "\tALLOCA %s, %s;\n", toStr(i.destination), toStr(i.size))
}

func (i *ALLOCA) Operands() []machine.Value {
	return []machine.Value{i.destination, i.size}
}
func (i *ALLOCA) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.size = operands[1]
}

type STR struct {
	destination, source machine.Value
	size                int
}

func (i *STR) Size() int {
	return i.size
}

func (i *STR) Type() machine.InstructionType {
	return machine.STR
}

func (i *STR) instr() {}
func (i *STR) Print(w io.Writer) {
	fmt.Fprintf(w, "\tSTR %s, %s;\n", toStr(i.source), toStr(i.destination))
}

func (i *STR) Operands() []machine.Value {
	return []machine.Value{i.destination, i.source}
}
func (i *STR) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.source = operands[1]
}

type LDR struct {
	destination, source machine.Value
	size                int
}

func (i *LDR) Size() int {
	return i.size
}

func (i *LDR) Type() machine.InstructionType {
	return machine.LDR
}

func (i *LDR) instr() {}
func (i *LDR) Print(w io.Writer) {
	fmt.Fprintf(w, "\tLDR %s, %s;\n", toStr(i.destination), toStr(i.source))
}

func (i *LDR) Operands() []machine.Value {
	return []machine.Value{i.destination, i.source}
}
func (i *LDR) ReplaceAll(operands ...machine.Value) {
	i.destination = operands[0]
	i.source = operands[1]
}

type RET struct {
	Value machine.Value
}

func (i *RET) Type() machine.InstructionType {
	return machine.RET
}

func (i *RET) instr() {}
func (i *RET) Print(w io.Writer) {
	fmt.Fprintf(w, "\tRET %s;\n", toStr(i.Value))
}

func (i *RET) Operands() []machine.Value {
	return []machine.Value{i.Value}
}
func (i *RET) ReplaceAll(operands ...machine.Value) {
	i.Value = operands[0]
}

func toStr(arg any) string {
	switch a := arg.(type) {
	default:
		return fmt.Sprintf("Unknown %T: %+v", a, a)
	case nil:
		return "<nil>"
	case *machine.ArgReg:
		return fmt.Sprintf("X%d", a.ID)
	case *machine.Reg:
		return fmt.Sprintf("x%d", a.ID)
	case *machine.BoolConst:
		return fmt.Sprintf("%v", a.Bool)
	case *machine.IntConst:
		return fmt.Sprintf("%d", a.Int)
	}
}
