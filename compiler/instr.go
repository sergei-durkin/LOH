package compiler

import "loh/token"

type Assign struct {
	target Variable
	op     token.Token
	arg1   Value
	arg2   Value

	size int

	isTemp bool
}

func (a *Assign) instr() {}

func (a *Assign) Target() Variable {
	return a.target
}

func (a *Assign) Arg1() Value {
	return a.arg1
}
func (a *Assign) Arg2() Value {
	return a.arg2
}
func (a *Assign) Size() int {
	return a.size
}

func (a *Assign) Operation() token.Token {
	return a.op
}

func (a *Assign) Operands() []Value {
	return []Value{a.arg1, a.arg2}
}

type Call struct {
	target Variable
	callee string
	args   []Value
}

func (c *Call) instr() {}

func (c *Call) Target() Variable {
	return c.target
}

func (c *Call) Callee() string {
	return c.callee
}

func (c *Call) Args() []Value {
	return c.args
}

func (c *Call) Operands() []Value {
	return c.args
}

type Alloca struct {
	ptr  Value
	size Value
}

func (a *Alloca) instr() {}

func (a *Alloca) Pointer() Value {
	return a.ptr
}

func (a *Alloca) Size() Value {
	return a.size
}

func (a *Alloca) Operands() []Value {
	return []Value{a.size}
}

type Store struct {
	Destination Value
	Value       Value
	Size        int
}

func (i *Store) instr() {}
func (i *Store) Operands() []Value {
	return []Value{i.Destination, i.Value}
}

type Load struct {
	Destination Value
	Source      Value
	Size        int
}

func (i *Load) instr() {}
func (i *Load) Operands() []Value {
	return []Value{i.Destination, i.Source}
}

type Label struct {
	id int
}

func (l *Label) instr() {}
func (l *Label) Operands() []Value {
	return nil
}

type Goto struct {
	label int
}

func (g *Goto) instr() {}

func (g *Goto) Label() int {
	return g.label
}

func (g *Goto) Operands() []Value {
	return nil
}

type IfGoto struct {
	cond  Value
	label int
	fall  int
}

func (i *IfGoto) instr() {}

func (i *IfGoto) Condition() Value {
	return i.cond
}

func (i *IfGoto) Label() int {
	return i.label
}
func (i *IfGoto) Fall() int {
	return i.fall
}

func (i *IfGoto) Operands() []Value {
	return []Value{i.cond}
}

type Return struct {
	value Value
}

func (r *Return) instr() {}
func (r *Return) Value() Value {
	return r.value
}
func (r *Return) Operands() []Value {
	return []Value{r.value}
}

type Phi struct {
	BbID int
	Var  Var
	Args map[int]Value
}

func (p *Phi) instr() {}
func (p *Phi) Operands() []Value {
	res := make([]Value, 0, len(p.Args))
	for _, v := range p.Args {
		res = append(res, v)
	}
	return res
}
