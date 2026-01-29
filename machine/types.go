package machine

type InstructionType int
type CMPType int

const (
	ALLOCA InstructionType = iota
	MUL
	SUM
	DIV
	MOD
	OR
	AND
	XOR
	SUB
	CALL
	JMP
	CMP
	CBZ
	MOV
	STR
	LDR
	RET
)
