package ast

import (
	"fmt"
	"loh/types"
)

var (
	step int = 4
	l    int = 0
)

func level() int {
	return l
}

func inc() {
	l += step
}

func dec() {
	l -= step
}

type Node interface {
	Pos() int
	Print()
}

type AST struct {
	Vars map[string]types.TypeInfo
	Unit *CompileUnit
}

func NewAST(unit *CompileUnit) *AST {
	return &AST{
		Unit: unit,
	}
}

func (a *AST) stmt() {}

func (a *AST) Print() {
	inc()

	fmt.Printf("%T {\n", a)

	fmt.Printf("%[1]*s Unit: ", level(), " ")
	a.Unit.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}
