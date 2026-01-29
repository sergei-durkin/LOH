package machine

type RegisterSize int

const (
	NULL  RegisterSize = 0
	INT8  RegisterSize = 1
	INT16 RegisterSize = 2
	INT32 RegisterSize = 4
	INT64 RegisterSize = 8
)

type FP struct{}

func (s *FP) value() {}

type ArgReg struct {
	ID int
}

func (s *ArgReg) value() {}

type Reg struct {
	ID   int
	Size RegisterSize
}

func (r *Reg) value() {}

type IntConst struct {
	Int int64
}

func (i *IntConst) value() {}

type BoolConst struct {
	Bool bool
}

func (b *BoolConst) value() {}
