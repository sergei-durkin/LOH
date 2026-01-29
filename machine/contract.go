package machine

import (
	"io"
	"loh/token"
)

type PhysicalRegister int

type IR interface {
	Functions() []Function
}

type Function interface {
	Name() string
	FrameSize() int
	Blocks() []Block
	CalleeSavedRegisters() []PhysicalRegister
}

type Block interface {
	Labeled
	Instructions() []Instruction
}

type Instruction interface {
	Operands() []Value
	Type() InstructionType
}

type Value interface {
	value()
}

type Machine interface {
	Emit(io.Writer, IR)
}

type Labeled interface {
	Label() int
}

type Tokened interface {
	Token() token.Token
}

type Calleeble interface {
	Callee() string
}

type Sizeble interface {
	Size() int
}
