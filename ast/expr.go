package ast

import (
	"fmt"
	"strings"
)

var (
	_ Expr = &StructLitExpr{}
	_ Expr = &NumberLitExpr{}

	_ Expr = &IdentExpr{}
	_ Expr = &MemberExpr{}
	_ Expr = &CallExpr{}

	_ Expr = &UnaryOp{}
	_ Expr = &BinaryOp{}
	_ Expr = &PostDecOp{}
	_ Expr = &PostIncOp{}
)

type Expr interface {
	Node
	expr()
}

type StructLitExpr struct {
	pos int

	name   Expr
	fields []*StructLitExprField
}

func NewStructLitExpr(pos int, name Expr, fields []*StructLitExprField) *StructLitExpr {
	return &StructLitExpr{
		pos:    pos,
		name:   name,
		fields: fields,
	}
}

func (e *StructLitExpr) expr() {}

func (e *StructLitExpr) Name() Expr {
	return e.name
}

func (e *StructLitExpr) Fields() []*StructLitExprField {
	return e.fields
}

func (e *StructLitExpr) Pos() int {
	return e.pos
}

func (e *StructLitExpr) Print() {
	inc()

	fmt.Printf("%T {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)

	fmt.Printf("%[1]*s name: ", level(), " ")
	e.name.Print()

	fmt.Printf("%[1]*s fields: ", level(), " ")
	if len(e.fields) > 0 {
		fmt.Print("[\n")

		inc()
		for i := 0; i < len(e.fields); i++ {
			fmt.Print(strings.Repeat(" ", level()))
			e.fields[i].Print()
		}
		dec()

		fmt.Printf("%[1]*s],\n", level()+1, " ")
	} else {
		fmt.Print("[]\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type StructLitExprField struct {
	pos int

	name  string
	value Expr
}

func NewStructLitExprField(pos int, name string, value Expr) *StructLitExprField {
	return &StructLitExprField{
		pos:   pos,
		name:  name,
		value: value,
	}
}

func (s *StructLitExprField) Pos() int {
	return s.pos
}

func (s *StructLitExprField) Name() string {
	return s.name
}

func (s *StructLitExprField) Value() Expr {
	return s.value
}

func (s *StructLitExprField) Print() {
	inc()

	fmt.Printf("%T {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)
	fmt.Printf("%[1]*s name: %q\n", level(), " ", s.name)
	fmt.Printf("%[1]*s value: ", level(), " ")

	if s.value != nil {
		s.value.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type ArrayLitExpr struct {
	pos int

	list []Expr
}

func NewArrayLitExpr(pos int, list []Expr) *ArrayLitExpr {
	return &ArrayLitExpr{
		pos:  pos,
		list: list,
	}
}

func (e *ArrayLitExpr) expr() {}

func (e *ArrayLitExpr) Pos() int {
	return e.pos
}
func (e *ArrayLitExpr) List() []Expr {
	return e.list
}

func (e *ArrayLitExpr) Print() {
	inc()

	fmt.Printf("%T {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)

	fmt.Printf("%[1]*s list: ", level(), " ")
	if len(e.list) > 0 {
		fmt.Print("[\n")

		inc()
		for i := 0; i < len(e.list); i++ {
			fmt.Print(strings.Repeat(" ", level()))
			e.list[i].Print()
		}
		dec()

		fmt.Printf("%[1]*s],\n", level()+1, " ")
	} else {
		fmt.Print("[]\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type NumberLitExpr struct {
	pos int

	value int64
}

func NewNumberLitExpr(pos int, value int64) *NumberLitExpr {
	return &NumberLitExpr{
		pos:   pos,
		value: value,
	}
}

func (e *NumberLitExpr) expr() {}

func (e *NumberLitExpr) Value() int64 {
	return e.value
}

func (e *NumberLitExpr) Pos() int {
	return e.pos
}

func (e *NumberLitExpr) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)
	fmt.Printf("%[1]*s value: %d\n", level(), " ", e.value)

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type BoolLitExpr struct {
	pos int

	value bool
}

func NewBoolLitExpr(pos int, value bool) *BoolLitExpr {
	return &BoolLitExpr{
		pos:   pos,
		value: value,
	}
}

func (e *BoolLitExpr) expr() {}

func (e *BoolLitExpr) Value() bool {
	return e.value
}

func (e *BoolLitExpr) Pos() int {
	return e.pos
}

func (e *BoolLitExpr) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)
	fmt.Printf("%[1]*s value: %t\n", level(), " ", e.value)

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type MemberExpr struct {
	pos int

	obj  Expr
	prop Expr
}

func NewMemberExpr(pos int, obj, prop Expr) *MemberExpr {
	return &MemberExpr{
		pos:  pos,
		obj:  obj,
		prop: prop,
	}
}

func (e *MemberExpr) expr() {}

func (e *MemberExpr) Obj() Expr {
	return e.obj
}

func (e *MemberExpr) Prop() Expr {
	return e.prop
}

func (e *MemberExpr) Pos() int {
	return e.pos
}

func (e *MemberExpr) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)

	fmt.Printf("%[1]*s obj: ", level(), " ")
	if e.obj != nil {
		e.obj.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	fmt.Printf("%[1]*s prop: ", level(), " ")
	e.prop.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type ArrayAccessExpr struct {
	pos    int
	target Expr
	addr   Expr
}

func NewArrayAccessExpr(pos int, target, addr Expr) *ArrayAccessExpr {
	return &ArrayAccessExpr{
		pos:    pos,
		target: target,
		addr:   addr,
	}
}

func (e *ArrayAccessExpr) expr() {}

func (e *ArrayAccessExpr) Target() Expr {
	return e.target
}

func (e *ArrayAccessExpr) Address() Expr {
	return e.addr
}

func (e *ArrayAccessExpr) Pos() int {
	return e.pos
}

func (e *ArrayAccessExpr) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)

	fmt.Printf("%[1]*s target: ", level(), " ")
	e.target.Print()

	fmt.Printf("%[1]*s address: ", level(), " ")
	e.addr.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type CallExpr struct {
	pos int

	callee Expr
	args   []Expr
}

func NewCallExpr(pos int, callee Expr, args []Expr) *CallExpr {
	return &CallExpr{
		pos:    pos,
		callee: callee,
		args:   args,
	}
}

func (e *CallExpr) expr() {}

func (e *CallExpr) Callee() Expr {
	return e.callee
}

func (e *CallExpr) Args() []Expr {
	return e.args
}

func (e *CallExpr) Pos() int {
	return e.pos
}

func (e *CallExpr) Print() {
	inc()

	fmt.Printf("%T {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)

	fmt.Printf("%[1]*s callee: ", level(), " ")
	e.callee.Print()

	fmt.Printf("%[1]*s args: ", level(), " ")
	if len(e.args) > 0 {
		fmt.Printf("[\n")

		inc()
		for i := 0; i < len(e.args); i++ {
			fmt.Print(strings.Repeat(" ", level()))
			e.args[i].Print()
		}
		dec()

		fmt.Printf("%[1]*s],\n", level()+1, " ")
	} else {
		fmt.Printf("[]\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type CallArgExpr struct {
	pos int

	expr Expr
}

func NewCallArgExpr(pos int, expr Expr) *CallArgExpr {
	return &CallArgExpr{
		pos:  pos,
		expr: expr,
	}
}

func (e *CallArgExpr) Pos() int {
	return e.pos
}

func (e *CallArgExpr) Print() {
	inc()
	fmt.Printf("%T: %+v\n", e, e)

	inc()
	fmt.Printf("%[1]*s", level(), "\t")
	e.expr.Print()
	dec()

	dec()
}

type IdentExpr struct {
	pos int

	value string
}

func NewIdentExpr(pos int, value string) *IdentExpr {
	return &IdentExpr{
		pos:   pos,
		value: value,
	}
}

func (e *IdentExpr) expr() {}

func (e *IdentExpr) Value() string {
	return e.value
}

func (e *IdentExpr) Pos() int {
	return e.pos
}

func (e *IdentExpr) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)
	fmt.Printf("%[1]*s value: %s\n", level(), " ", e.value)

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type StringLitExpr struct {
	pos int

	value string
}

func NewStringLitExpr(pos int, value string) *StringLitExpr {
	return &StringLitExpr{
		pos:   pos,
		value: value,
	}
}

func (e *StringLitExpr) expr() {}

func (e *StringLitExpr) Value() string {
	return e.value
}

func (e *StringLitExpr) Unquoted() string {
	return e.value[1 : len(e.value)-1]
}

func (e *StringLitExpr) Pos() int {
	return e.pos
}

func (e *StringLitExpr) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)

	fmt.Printf("%[1]*s value: %s\n", level(), " ", e.value)

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}
