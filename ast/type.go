package ast

import (
	"fmt"
	"loh/types"
)

var (
	_ Type = &ArrayType{}
	_ Type = &BasicType{}
)

type Type interface {
	Print()
	typ()
}

type ArrayType struct {
	pos int

	info Type
	len  Expr
}

func (e *ArrayType) typ() {}

func (e *ArrayType) Element() Type {
	return e.info
}

func (e *ArrayType) Len() Expr {
	return e.len
}

func (e *ArrayType) Print() {
	inc()

	fmt.Printf("%T {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)
	fmt.Printf("%[1]*s info: ", level(), " ")
	e.info.Print()

	fmt.Printf("%[1]*s len: ", level(), " ")
	if e.len != nil {
		e.len.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

func NewArrayType(pos int, info Type, len Expr) *ArrayType {
	return &ArrayType{
		pos:  pos,
		info: info,
		len:  len,
	}
}

type BasicType struct {
	pos int

	info types.TypeInfo
}

func (e *BasicType) typ() {}

func (e *BasicType) Element() types.TypeInfo {
	return e.info
}
func (e *BasicType) Print() {
	inc()

	fmt.Printf("%T {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)
	fmt.Printf("%[1]*s info: %T %d\n", level(), " ", e.info, e.info.Size())

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

func NewBasicType(pos int, info types.TypeInfo) *BasicType {
	return &BasicType{
		pos:  pos,
		info: info,
	}
}

type CustomType struct {
	pos int

	name string
}

func (e *CustomType) typ() {}

func (e *CustomType) Name() string {
	return e.name
}

func (e *CustomType) Print() {
	inc()

	fmt.Printf("%T {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)
	fmt.Printf("%[1]*s name: %s\n", level(), " ", e.name)

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

func NewCustomType(pos int, name string) *CustomType {
	return &CustomType{
		pos:  pos,
		name: name,
	}
}
