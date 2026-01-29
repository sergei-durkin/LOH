package ast

import (
	"fmt"
	"loh/token"
	"strings"
)

var (
	_ Stmt = &FnDecl{}
	_ Stmt = &StructDecl{}
	_ Stmt = &ConstDecl{}
	_ Stmt = &VarDecl{}
	_ Stmt = &FnArgDecl{}

	_ Stmt = &BlockStmt{}

	_ Stmt = &IfStmt{}
	_ Stmt = &AssignStmt{}
	_ Stmt = &ForStmt{}
	_ Stmt = &ContinueStmt{}
	_ Stmt = &BreakStmt{}
	_ Stmt = &ReturnStmt{}

	_ Stmt = &ExprStmt{}
)

type Stmt interface {
	Node
	stmt()
}

type CompileUnit struct {
	path  string
	nodes []Node
}

func NewCompileUnit(path string, nodes []Node) *CompileUnit {
	return &CompileUnit{
		path:  path,
		nodes: nodes,
	}
}

func (s *CompileUnit) Nodes() []Node {
	return s.nodes
}

func (s *CompileUnit) Print() {
	inc()

	fmt.Printf("%T {\n", s)

	fmt.Printf("%[1]*s Nodes: ", level(), " ")
	if len(s.nodes) > 0 {
		fmt.Print("[\n")

		inc()
		for i := 0; i < len(s.nodes); i++ {
			fmt.Print(strings.Repeat(" ", level()))
			s.nodes[i].Print()
		}
		dec()

		fmt.Printf("%[1]*s],\n", level()+1, " ")
	} else {
		fmt.Print("[]\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type FnDecl struct {
	pos int

	name       string
	args       []*FnArgDecl
	returnType Expr
	body       BlockStmt
}

func NewFnDecl(pos int, name string, args []*FnArgDecl, returnType Expr, body BlockStmt) *FnDecl {
	return &FnDecl{
		pos:        pos,
		name:       name,
		args:       args,
		returnType: returnType,
		body:       body,
	}
}

func (s *FnDecl) stmt() {}

func (s *FnDecl) Name() string {
	return s.name
}

func (s *FnDecl) Stmts() []Stmt {
	return s.body.stmts
}

func (s *FnDecl) Args() []*FnArgDecl {
	return s.args
}

func (s *FnDecl) Pos() int {
	return s.pos
}

func (s *FnDecl) Print() {
	inc()

	fmt.Printf("%T {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)
	fmt.Printf("%[1]*s name: %q\n", level(), " ", s.name)

	fmt.Printf("%[1]*s args: ", level(), " ")
	if len(s.args) > 0 {
		fmt.Print("[\n")

		inc()
		for i := 0; i < len(s.args); i++ {
			fmt.Print(strings.Repeat(" ", level()))
			s.args[i].Print()
		}
		dec()

		fmt.Printf("%[1]*s],\n", level()+1, " ")
	} else {
		fmt.Print("[]\n")
	}

	fmt.Printf("%[1]*s returnType: ", level(), " ")
	if s.returnType != nil {
		s.returnType.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	fmt.Printf("%[1]*s body: ", level(), " ")
	s.body.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type StructDecl struct {
	pos int

	name   string
	fields []*StructFieldDecl
}

func NewStructDecl(pos int, name string, fields []*StructFieldDecl) *StructDecl {
	return &StructDecl{
		pos:    pos,
		name:   name,
		fields: fields,
	}
}

func (s *StructDecl) stmt() {}

func (s *StructDecl) Pos() int {
	return s.pos
}
func (s *StructDecl) Name() string {
	return s.name
}
func (s *StructDecl) Fields() []*StructFieldDecl {
	return s.fields
}

func (s *StructDecl) Print() {
	inc()

	fmt.Printf("%T {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)
	fmt.Printf("%[1]*s name: %q\n", level(), " ", s.name)

	fmt.Printf("%[1]*s fields: ", level(), " ")
	if len(s.fields) > 0 {
		fmt.Print("[\n")

		inc()
		for i := 0; i < len(s.fields); i++ {
			fmt.Print(strings.Repeat(" ", level()))
			s.fields[i].Print()
		}
		dec()

		fmt.Printf("%[1]*s],\n", level()+1, " ")
	} else {
		fmt.Print("[]\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type StructFieldDecl struct {
	pos int

	name string
	typ  Type
}

func NewStructFieldDecl(pos int, name string, typ Type) *StructFieldDecl {
	return &StructFieldDecl{
		pos:  pos,
		name: name,
		typ:  typ,
	}
}

func (s *StructFieldDecl) Pos() int {
	return s.pos
}
func (s *StructFieldDecl) Name() string {
	return s.name
}
func (s *StructFieldDecl) Type() Type {
	return s.typ
}

func (s *StructFieldDecl) Print() {
	inc()

	fmt.Printf("%T {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)
	fmt.Printf("%[1]*s name: %q\n", level(), " ", s.name)
	fmt.Printf("%[1]*s type: \n", level(), " ")
	s.typ.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type ConstDecl struct {
	pos int

	name     string
	typeName string
	value    Expr
}

func NewConstDecl(pos int, name string, typeName string, value Expr) *ConstDecl {
	return &ConstDecl{
		pos:      pos,
		name:     name,
		typeName: typeName,
		value:    value,
	}
}

func (s *ConstDecl) stmt() {}

func (s *ConstDecl) Pos() int {
	return s.pos
}

func (s *ConstDecl) Print() {
	inc()

	fmt.Printf("%T {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)
	fmt.Printf("%[1]*s name: %q\n", level(), " ", s.name)
	fmt.Printf("%[1]*s type: %q\n", level(), " ", s.typeName)
	fmt.Printf("%[1]*s value: ", level(), " ")

	if s.value != nil {
		s.value.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type VarDecl struct {
	pos int

	name string
	typ  Type

	value Expr
}

func NewVarDecl(pos int, name string, typ Type, value Expr) *VarDecl {
	return &VarDecl{
		pos:   pos,
		name:  name,
		typ:   typ,
		value: value,
	}
}

func (s *VarDecl) stmt() {}

func (s *VarDecl) Name() string {
	return s.name
}
func (s *VarDecl) Type() Type {
	return s.typ
}

func (s *VarDecl) Value() Expr {
	return s.value
}

func (s *VarDecl) Pos() int {
	return s.pos
}

func (s *VarDecl) Print() {
	inc()

	fmt.Printf("%T {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)
	fmt.Printf("%[1]*s name: %q\n", level(), " ", s.name)

	fmt.Printf("%[1]*s type: \n", level(), " ")
	s.typ.Print()

	fmt.Printf("%[1]*s value: ", level(), " ")

	if s.value != nil {
		s.value.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type FnArgDecl struct {
	pos int

	name     string
	typeName Expr
	size     int64
}

func NewFnArgDecl(pos int, name string, typeName Expr, size int64) *FnArgDecl {
	return &FnArgDecl{
		pos:      pos,
		name:     name,
		typeName: typeName,
		size:     size,
	}
}

func (s *FnArgDecl) stmt() {}

func (s *FnArgDecl) Name() string {
	return s.name
}

func (s *FnArgDecl) Size() int64 {
	return s.size
}

func (s *FnArgDecl) Pos() int {
	return s.pos
}

func (s *FnArgDecl) Print() {
	inc()

	fmt.Printf("%T: {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)
	fmt.Printf("%[1]*s name: %s\n", level(), " ", s.name)

	fmt.Printf("%[1]*s type: ", level(), " ")
	if s.typeName != nil {
		s.typeName.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type BlockStmt struct {
	pos int

	stmts []Stmt
}

func NewBlockStmt(pos int, stmts []Stmt) *BlockStmt {
	return &BlockStmt{
		pos:   pos,
		stmts: stmts,
	}
}

func (s *BlockStmt) stmt() {}

func (s *BlockStmt) Stmts() []Stmt {
	return s.stmts
}

func (s *BlockStmt) Pos() int {
	return s.pos
}

func (s *BlockStmt) Print() {
	inc()

	fmt.Printf("%T {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)

	fmt.Printf("%[1]*s stmts: ", level(), " ")
	if len(s.stmts) > 0 {
		fmt.Printf("[\n")

		inc()
		for i := 0; i < len(s.stmts); i++ {
			fmt.Print(strings.Repeat(" ", level()))
			s.stmts[i].Print()
		}
		dec()

		fmt.Printf("%[1]*s],\n", level()+1, " ")
	} else {
		fmt.Printf("[]\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type ExprStmt struct {
	pos int

	expr Expr
}

func NewExprStmt(pos int, expr Expr) *ExprStmt {
	return &ExprStmt{
		pos:  pos,
		expr: expr,
	}
}

func (s *ExprStmt) stmt() {}

func (s *ExprStmt) Expr() Expr {
	return s.expr
}

func (s *ExprStmt) Pos() int {
	return s.pos
}

func (s *ExprStmt) Print() {
	inc()

	fmt.Printf("%T: {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)

	fmt.Printf("%[1]*s expr: ", level(), " ")
	s.expr.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type IfStmt struct {
	pos int

	condition Expr
	body      BlockStmt
	els       Stmt
}

func NewIfStmt(pos int, cond Expr, body BlockStmt, els Stmt) *IfStmt {
	return &IfStmt{
		pos:       pos,
		condition: cond,
		body:      body,
		els:       els,
	}
}

func (s *IfStmt) stmt() {}

func (s *IfStmt) Condition() Expr {
	return s.condition
}

func (s *IfStmt) Body() BlockStmt {
	return s.body
}

func (s *IfStmt) Else() Stmt {
	return s.els
}

func (s *IfStmt) Pos() int {
	return s.pos
}

func (s *IfStmt) Print() {
	inc()

	fmt.Printf("%T: {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)

	fmt.Printf("%[1]*s cond: ", level(), " ")
	s.condition.Print()

	fmt.Printf("%[1]*s body: ", level(), " ")
	s.body.Print()

	fmt.Printf("%[1]*s else: ", level(), " ")
	if s.els != nil {
		s.els.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type AssignStmt struct {
	pos int

	target Expr
	value  Expr
}

func NewAssignStmt(pos int, target, value Expr) *AssignStmt {
	return &AssignStmt{
		pos:    pos,
		target: target,
		value:  value,
	}
}

func (s *AssignStmt) stmt() {}

func (s *AssignStmt) Target() Expr {
	return s.target
}

func (s *AssignStmt) Value() Expr {
	return s.value
}

func (s *AssignStmt) Pos() int {
	return s.pos
}

func (s *AssignStmt) Print() {
	inc()

	fmt.Printf("%T: {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)

	fmt.Printf("%[1]*s target: ", level(), " ")
	s.target.Print()

	fmt.Printf("%[1]*s value: ", level(), " ")
	s.value.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type ForStmt struct {
	pos int

	init      Stmt
	condition Expr
	post      Stmt
	body      BlockStmt
}

func NewForStmt(pos int, init Stmt, condition Expr, post Stmt, body BlockStmt) *ForStmt {
	return &ForStmt{
		pos:       pos,
		init:      init,
		condition: condition,
		post:      post,
		body:      body,
	}
}

func (s *ForStmt) stmt() {}

func (s *ForStmt) Init() Stmt {
	return s.init
}

func (s *ForStmt) Condition() Expr {
	return s.condition
}

func (s *ForStmt) Post() Stmt {
	return s.post
}

func (s *ForStmt) Body() BlockStmt {
	return s.body
}

func (s *ForStmt) Pos() int {
	return s.pos
}

func (s *ForStmt) Print() {
	inc()

	fmt.Printf("%T: {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)

	fmt.Printf("%[1]*s init: ", level(), " ")
	if s.init != nil {
		s.init.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	fmt.Printf("%[1]*s cond: ", level(), " ")
	if s.condition != nil {
		s.condition.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	fmt.Printf("%[1]*s post: ", level(), " ")
	if s.post != nil {
		s.post.Print()
	} else {
		fmt.Print("<nil>\n")
	}

	fmt.Printf("%[1]*s body: ", level(), " ")
	s.body.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type ContinueStmt struct {
	pos int
}

func NewContinueStmt(pos int) *ContinueStmt {
	return &ContinueStmt{
		pos: pos,
	}
}

func (s *ContinueStmt) stmt() {}

func (s *ContinueStmt) Pos() int {
	return s.pos
}

func (s *ContinueStmt) Print() {
	inc()

	fmt.Printf("%T: {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type BreakStmt struct {
	pos int
}

func NewBreakStmt(pos int) *BreakStmt {
	return &BreakStmt{
		pos: pos,
	}
}

func (s *BreakStmt) stmt() {}

func (s *BreakStmt) Pos() int {
	return s.pos
}

func (s *BreakStmt) Print() {
	inc()

	fmt.Printf("%T: {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type ReturnStmt struct {
	pos int

	value Expr
}

func NewReturnStmt(pos int, value Expr) *ReturnStmt {
	return &ReturnStmt{
		pos:   pos,
		value: value,
	}
}

func (s *ReturnStmt) stmt() {}

func (s *ReturnStmt) Value() Expr {
	return s.value
}

func (s *ReturnStmt) Pos() int {
	return s.pos
}

func (s *ReturnStmt) Print() {
	inc()

	fmt.Printf("%T: {\n", s)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", s.pos)

	fmt.Printf("%[1]*s value: ", level(), " ")
	s.value.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type UnaryOp struct {
	pos int

	op    token.Token
	value Expr
}

func NewUnaryOp(pos int, op token.Token, value Expr) *UnaryOp {
	return &UnaryOp{
		pos:   pos,
		op:    op,
		value: value,
	}
}

func (e *UnaryOp) Op() token.Token {
	return e.op
}

func (e *UnaryOp) Value() Expr {
	return e.value
}

func (e *UnaryOp) expr() {}

func (e *UnaryOp) Pos() int {
	return e.pos
}

func (e *UnaryOp) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)
	fmt.Printf("%[1]*s op: {int: %d, str: %s}\n", level(), " ", e.op, e.op.String())

	fmt.Printf("%[1]*s value: ", level(), " ")
	e.value.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type BinaryOp struct {
	pos int

	op    token.Token
	left  Expr
	right Expr
}

func NewBinaryOp(pos int, op token.Token, left, right Expr) *BinaryOp {
	return &BinaryOp{
		pos:   pos,
		op:    op,
		left:  left,
		right: right,
	}
}

func (e *BinaryOp) expr() {}

func (e *BinaryOp) Op() token.Token {
	return e.op
}

func (e *BinaryOp) Left() Expr {
	return e.left
}

func (e *BinaryOp) Right() Expr {
	return e.right
}

func (e *BinaryOp) Pos() int {
	return e.pos
}

func (e *BinaryOp) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)
	fmt.Printf("%[1]*s op: {int: %d, str: %s}\n", level(), " ", e.op, e.op.String())

	fmt.Printf("%[1]*s left: ", level(), " ")
	e.left.Print()

	fmt.Printf("%[1]*s right: ", level(), " ")
	e.right.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type PostDecOp struct {
	pos int

	value Expr
}

func NewPostDecOp(pos int, value Expr) *PostDecOp {
	return &PostDecOp{
		pos:   pos,
		value: value,
	}
}
func (e *PostDecOp) expr() {}

func (e *PostDecOp) Value() Expr {
	return e.value
}

func (e *PostDecOp) Pos() int {
	return e.pos
}

func (e *PostDecOp) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)

	fmt.Printf("%[1]*s value: ", level(), " ")
	e.value.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}

type PostIncOp struct {
	pos int

	value Expr
}

func NewPostIncOp(pos int, value Expr) *PostIncOp {
	return &PostIncOp{
		pos:   pos,
		value: value,
	}
}

func (e *PostIncOp) expr() {}

func (e *PostIncOp) Value() Expr {
	return e.value
}

func (e *PostIncOp) Pos() int {
	return e.pos
}

func (e *PostIncOp) Print() {
	inc()

	fmt.Printf("%T: {\n", e)
	fmt.Printf("%[1]*s pos: %d\n", level(), " ", e.pos)

	fmt.Printf("%[1]*s value: ", level(), " ")
	e.value.Print()

	dec()
	fmt.Printf("%[1]*s},\n", level()+1, " ")
}
