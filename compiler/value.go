package compiler

type Value interface {
	value()
	Equal(another Value) bool
}

type AddressOf struct {
	Target Value
}

func (r *AddressOf) value() {}
func (r *AddressOf) Equal(another Value) bool {
	if r == nil || another == nil {
		return r == another
	}

	if v2, ok := another.(*AddressOf); ok {
		return r.Target.Equal(v2.Target)
	}

	return false
}

type Dereference struct {
	Addr Value
}

func (d *Dereference) value() {}
func (d *Dereference) Equal(another Value) bool {
	if d == nil || another == nil {
		return d == another
	}

	if v2, ok := another.(*Dereference); ok {
		if d.Addr == nil || v2.Addr == nil {
			return d.Addr == v2.Addr
		}
		return d.Addr.Equal(v2.Addr)
	}

	return false
}

type BoolConst struct {
	bool bool
}

func (b *BoolConst) value() {}
func (b *BoolConst) Bool() bool {
	return b.bool
}

func (b *BoolConst) Equal(another Value) bool {
	if b == nil || another == nil {
		return b == another
	}

	if b2, ok := another.(*BoolConst); ok {
		return *b == *b2
	}

	return false
}

type IntConst struct {
	int int64
}

func (i *IntConst) value() {}
func (i *IntConst) Int() int64 {
	return i.int
}

func (i *IntConst) Equal(another Value) bool {
	if i == nil || another == nil {
		return i == another
	}

	if i2, ok := another.(*IntConst); ok {
		return *i == *i2
	}

	return false
}

type StringConst struct {
	string string
}

func (s *StringConst) value() {}

func (s *StringConst) Equal(another Value) bool {
	if s == nil || another == nil {
		return s == another
	}

	if s2, ok := another.(*StringConst); ok {
		return *s == *s2
	}

	return false
}
