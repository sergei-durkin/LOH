package compiler

import "fmt"

type Variable interface {
	Label() string
	Name() string
}

type Reg struct {
	ID int
}

func (r *Reg) value() {}

func (r *Reg) Equal(another Value) bool {
	if r == nil || another == nil {
		return r == another
	}

	if v2, ok := another.(*Reg); ok {
		return *r == *v2
	}

	return false
}

type ArgReg struct {
	ID int
}

func (r *ArgReg) value() {}

func (r *ArgReg) Equal(another Value) bool {
	if r == nil || another == nil {
		return r == another
	}

	if v2, ok := another.(*ArgReg); ok {
		return *r == *v2
	}

	return false
}

type FP struct{}

func (f *FP) value() {}
func (f *FP) Equal(another Value) bool {
	if f == nil || another == nil {
		return f == another
	}

	return true
}

type Var struct {
	name string
	ver  int
}

func (v *Var) value() {}

func (v *Var) Label() string {
	return fmt.Sprintf("%s_%d", v.name, v.ver)
}

func (v *Var) Name() string {
	return v.name
}

func (v *Var) Equal(another Value) bool {
	if v == nil || another == nil {
		return v == another
	}

	if v2, ok := another.(*Var); ok {
		return *v == *v2
	}

	return false
}

type TempVar struct {
	label string
}

func (t *TempVar) value() {}
func (t *TempVar) Label() string {
	return t.label
}

func (t *TempVar) Name() string {
	return t.label
}

func (t *TempVar) Equal(another Value) bool {
	if t == nil || another == nil {
		return t == another
	}

	if t2, ok := another.(*TempVar); ok {
		return *t == *t2
	}

	return false
}
