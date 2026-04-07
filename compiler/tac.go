package compiler

import (
	"fmt"
	"io"
	"loh/ast"
	"loh/token"
	"loh/types"
)

type Instruction interface {
	Operands() []Value

	instr()

	Print(w io.Writer)
}

var (
	_ Instruction = &Assign{}
	_ Instruction = &Call{}

	_ Instruction = &Label{}
	_ Instruction = &Goto{}
	_ Instruction = &IfGoto{}

	_ Instruction = &Alloca{}
	_ Instruction = &Store{}
	_ Instruction = &Load{}

	_ Instruction = &Phi{}

	_ Instruction = &Return{}
)

type loopCtx struct {
	breakLabel    int
	continueLabel int
}

type Tac struct {
	name string

	code []Instruction
	size map[string]int64

	structs map[string]*Struct
	types   map[string]*Struct

	tempCnt  int
	labelCnt int

	loopStack []loopCtx
}

func NewTac(a *ast.CompileUnit) []*Tac {
	tacs := []*Tac{}

	types := make(map[string]*Struct)

	nodes := a.Nodes()
	for i := 0; i < len(nodes); i++ {
		switch d := nodes[i].(type) {
		default:
		case *ast.StructDecl:
			parseStruct(types, d)
		}
	}

	for i := 0; i < len(nodes); i++ {
		t := &Tac{
			size:    make(map[string]int64),
			types:   make(map[string]*Struct),
			structs: types,

			tempCnt:  -1,
			labelCnt: -1,
		}

		switch d := nodes[i].(type) {
		case *ast.FnDecl:
			t.parseFn(d)
			tacs = append(tacs, t)
		}
	}

	return tacs
}

func (t *Tac) pushLoop(breakLabel, continueLabel int) {
	t.loopStack = append(t.loopStack, loopCtx{breakLabel: breakLabel, continueLabel: continueLabel})
}

func (t *Tac) popLoop() {
	t.loopStack = t.loopStack[:len(t.loopStack)-1]
}

func (t *Tac) curLoop() loopCtx {
	if len(t.loopStack) == 0 {
		panic("break/continue outside loop")
	}

	return t.loopStack[len(t.loopStack)-1]
}

func (t *Tac) id() int {
	t.labelCnt++
	return t.labelCnt
}

func (t *Tac) temp() *TempVar {
	t.tempCnt++
	return &TempVar{label: fmt.Sprintf("t%d", t.tempCnt)}
}

type Struct struct {
	info   types.TypeInfo
	fields map[string]*Struct
	offset int
}

func parseStruct(structs map[string]*Struct, s *ast.StructDecl) {
	name := s.Name()
	if _, ok := structs[name]; ok {
		panic("redeclare struct not allowed")
	}

	typ := &Struct{fields: make(map[string]*Struct)}
	structs[name] = typ
	typ.parseStructFields(s.Fields())
}

func (t *Struct) parseStructFields(fields []*ast.StructFieldDecl) {
	offset := 0
	for i := 0; i < len(fields); i++ {
		typ := t.parseType(offset, fields[i].Type())
		t.fields[fields[i].Name()] = typ
		offset += int(typ.info.Size())
	}
}

func (t *Struct) parseType(offset int, typ ast.Type) *Struct {
	switch typ := typ.(type) {
	default:
		panic(fmt.Sprintf("undefined type: %T %+v", typ, typ))
	case *ast.BasicType:
		return &Struct{offset: offset, info: typ.Element()}
	}
}

func (t *Tac) parseFn(fn *ast.FnDecl) {
	t.name = fn.Name()
	regs := 0

	t.code = append(t.code, &Label{id: t.id()})
	args := fn.Args()
	for i := 0; i < len(args); i++ {
		if regs < 8 {
			v := &Var{name: args[i].Name()}
			t.size[v.Name()] = args[i].Size()

			t.code = append(t.code, &Assign{
				target: v,
				op:     token.ASSIGN,
				arg1:   &ArgReg{ID: regs},

				size: int(args[i].Size()),
			})
			regs++

			continue
		}

		panic("not enough registers")
	}

	stmts := fn.Stmts()
	for i := 0; i < len(stmts); i++ {
		t.parseStmt(stmts[i])
	}
}

func (t *Tac) parseStructLit(s *Struct, target *Var, lit *ast.StructLitExpr) {
	list := lit.Fields()
	if len(list) > len(s.fields) {
		panic("too short struct decl")
	}

	for i := 0; i < len(list); i++ {
		typ := s.fields[list[i].Name()]
		offset := typ.offset
		size := typ.info.Size()

		sp := t.temp()
		t.code = append(t.code, &Assign{
			target: sp,
			op:     token.PLUS,
			arg1:   target,
			arg2:   &IntConst{int: int64(offset)},
		})

		t.code = append(t.code, &Store{
			Destination: sp,
			Value:       t.parseExpr(list[i].Value()),
			Size:        int(size),
		})
	}
}

func (t *Tac) parseVarDecl(st *ast.VarDecl) {
	switch typ := st.Type().(type) {
	default:
		panic(fmt.Sprintf("undefined type: %T %+v", typ, typ))
	case *ast.CustomType:
		s := t.structs[typ.Name()]
		size := 0
		for _, v := range s.fields {
			size += int(v.info.Size())
		}
		size = aligned(size, 16)

		target := &Var{name: st.Name()}
		t.types[target.Name()] = s
		t.size[target.Name()] = int64(size)

		tmp := t.temp()
		t.code = append(t.code, &Assign{
			target: tmp,
			op:     token.ASSIGN,
			arg1:   &IntConst{int64(size)},
			// arg2:   t.parseExpr(ln), // TODO:

			size: int(types.Pointer.Size()),
		})

		t.code = append(t.code, &Alloca{
			ptr:  target,
			size: tmp,
		})

		switch v := st.Value().(type) {
		case *ast.IdentExpr:
			t.code = append(t.code, &Assign{
				target: target,
				op:     token.ASSIGN,
				arg1:   t.parseExpr(v),

				size: int(types.Pointer.Size()),

				isTemp: false,
			})
		case *ast.StructLitExpr:
			t.parseStructLit(s, target, v)
		}
		// panic(size)

	case *ast.ArrayType:
		switch ln := typ.Len().(type) {
		case *ast.NumberLitExpr, *ast.IdentExpr:
			target := &Var{name: st.Name()}
			var size int64
			switch s := typ.Element().(type) {
			default:
				panic(fmt.Sprintf("undefined type: %T %+v", s, s))
			case *ast.BasicType:
				size = s.Element().Size()
			}
			t.size[target.Name()] = size

			tmp := t.temp()
			t.code = append(t.code, &Assign{
				target: tmp,
				op:     token.STAR,
				arg1:   &IntConst{size},
				arg2:   t.parseExpr(ln),

				size: int(types.Pointer.Size()),
			})

			t.code = append(t.code, &Alloca{
				ptr:  target,
				size: tmp,
			})

			if st.Value() != nil {
				t.parseArrayList(target, st.Value().(*ast.ArrayLitExpr), ln)
			}
		case nil:
			target := &Var{name: st.Name()}
			var size int64
			switch s := typ.Element().(type) {
			default:
				panic(fmt.Sprintf("undefined type: %T %+v", s, s))
			case *ast.BasicType:
				size = s.Element().Size()
			}
			t.size[target.Name()] = size

			switch v := st.Value().(type) {
			default:
				panic(fmt.Sprintf("undefined value: %T %+v", v, v))
			case *ast.ArrayLitExpr:
				n := t.parseArrayList(target, v, typ.Len())
				if n > 0 {
					size *= int64(n)
				}
				t.code = append(t.code, &Alloca{
					ptr:  target,
					size: &IntConst{int: size},
				})

				return
			case *ast.StringLitExpr:
				if size != int64(types.Uint8) {
					panic("size missmatch")
				}

				s := v.Unquoted()
				if len(s) > 0 {
					size *= int64(len(s))
					size++
				}

				t.code = append(t.code, &Alloca{
					ptr:  target,
					size: &IntConst{int: size},
				})

				for i := 0; i < len(s); i++ {
					t.storeByte(target, i, s[i])
				}
				t.storeByte(target, len(s), 0)
			}
		}
	case *ast.BasicType:
		v := &Var{name: st.Name()}
		t.size[v.Name()] = typ.Element().Size()

		t.code = append(t.code, &Assign{
			target: v,
			op:     token.ASSIGN,
			arg1:   t.parseExpr(st.Value()),

			size: int(typ.Element().Size()),

			isTemp: false,
		})
	}
}

func (t *Tac) parseAssignStmt(stmt *ast.AssignStmt) {
	switch target := stmt.Target().(type) {
	case *ast.UnaryOp: // *a = ast.expr
		if target.Op() != token.STAR {
			panic("unknown op")
		}

		tmp2 := t.temp()
		t.code = append(t.code, &Assign{
			target: tmp2,
			op:     token.ASSIGN,
			arg1:   t.parseExpr(stmt.Value()),

			size: int(types.Pointer.Size()),

			isTemp: true,
		})

		tmp := t.temp()
		t.code = append(t.code, &Assign{
			target: tmp,
			op:     token.ASSIGN,
			arg1:   t.parseExpr(stmt.Target()),

			isTemp: true,
		})

		t.code = append(t.code, &Store{
			Destination: tmp,
			Value:       tmp2,
			Size:        8,
		})
	case *ast.ArrayAccessExpr:
		ttmp := t.temp()
		arrayTarget, size := t.parseArrayAccess(target)
		t.code = append(t.code, &Assign{
			target: ttmp,
			op:     token.ASSIGN,
			arg1:   arrayTarget,

			isTemp: true,
		})

		vtmp := t.temp()
		t.code = append(t.code, &Assign{
			target: vtmp,
			op:     token.ASSIGN,
			arg1:   t.parseExpr(stmt.Value()),

			isTemp: true,
		})

		t.code = append(t.code, &Store{
			Destination: ttmp,
			Value:       vtmp,
			Size:        int(size),
		})
	case *ast.MemberExpr:
		panic("not implemented")
	default:
		t.code = append(t.code, &Assign{
			target: &Var{name: t.resolveIdentity(stmt.Target())},
			op:     token.ASSIGN,
			arg1:   t.parseExpr(stmt.Value()),

			size: -1,

			isTemp: false,
		})
	}
}

func (t *Tac) parseIfStmt(stmt *ast.IfStmt) {
	cond := t.parseExpr(stmt.Condition())

	els := t.id()
	end := t.id()

	next := t.id()

	t.code = append(t.code, &IfGoto{
		cond:  cond,
		label: els,
		fall:  next,
	})

	t.code = append(t.code, &Label{id: next})

	b := stmt.Body()
	t.parseStmt(&b)

	t.code = append(t.code, &Goto{label: end})
	t.code = append(t.code, &Label{id: els})

	t.parseStmt(stmt.Else())

	t.code = append(t.code, &Goto{label: end})
	t.code = append(t.code, &Label{id: end})
}

func (t *Tac) parseForStmt(stmt *ast.ForStmt) {
	t.parseStmt(stmt.Init())

	loop := t.id()
	t.code = append(t.code, &Label{id: loop})

	cond := t.parseExpr(stmt.Condition())

	loopBody := t.id()
	postLabel := t.id()
	end := t.id()

	t.code = append(t.code, &IfGoto{
		cond:  cond,
		label: end,
		fall:  loopBody,
	})
	t.code = append(t.code, &Label{id: loopBody})

	t.pushLoop(end, postLabel)

	b := stmt.Body()
	t.parseStmt(&b)

	t.popLoop()

	t.code = append(t.code, &Label{id: postLabel})
	t.parseStmt(stmt.Post())
	t.code = append(t.code, &Goto{label: loop})

	t.code = append(t.code, &Label{id: end})
}

func (t *Tac) parseStmt(stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		stmts := stmt.Stmts()
		for i := 0; i < len(stmts); i++ {
			t.parseStmt(stmts[i])
		}
	case *ast.VarDecl:
		t.parseVarDecl(stmt)
	case *ast.AssignStmt:
		t.parseAssignStmt(stmt)
	case *ast.IfStmt:
		t.parseIfStmt(stmt)
	case *ast.BreakStmt:
		cur := t.curLoop()
		t.code = append(t.code, &Goto{cur.breakLabel})
	case *ast.ContinueStmt:
		cur := t.curLoop()
		t.code = append(t.code, &Goto{cur.continueLabel})
	case *ast.ForStmt:
		t.parseForStmt(stmt)
	case *ast.ReturnStmt:
		ret := t.parseExpr(stmt.Value())
		t.code = append(t.code, &Return{value: ret})
	case *ast.ExprStmt:
		switch expr := stmt.Expr().(type) {
		default:
			panic(fmt.Sprintf("undefined expr: %T %+v", expr, expr))
		case *ast.CallExpr: // test(1)
			t.parseCall(expr)
		case *ast.PostIncOp: // a++
			v := t.resolveIdentity(expr.Value())
			t.code = append(t.code, &Assign{
				target: &Var{name: v},
				op:     token.PLUS,
				arg1:   &Var{name: v},
				arg2:   &IntConst{int: 1},

				isTemp: false,
			})
		case *ast.PostDecOp: // a--
			v := t.resolveIdentity(expr.Value())
			t.code = append(t.code, &Assign{
				target: &Var{name: v},
				op:     token.MINUS,
				arg1:   &Var{name: v},
				arg2:   &IntConst{int: 1},

				isTemp: false,
			})
		}
	}
}

func (t *Tac) parseCall(expr *ast.CallExpr) Value {
	eargs := expr.Args()

	args := make([]Value, 0, len(eargs))
	for i := 0; i < len(eargs); i++ {
		args = append(args, t.parseExpr(eargs[i]))
	}

	tmp := t.temp()
	t.code = append(t.code, &Call{
		target: tmp,
		callee: t.resolveIdentity(expr.Callee()),
		args:   args,
	})
	return tmp
}

func (t *Tac) resolveIdentity(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.UnaryOp:
		if e.Op() != token.STAR {
			panic(fmt.Sprintf("unknown unary id: %+v", e.Op()))
		}
		return t.resolveIdentity(e.Value())
	case *ast.IdentExpr:
		return e.Value()
	case *ast.MemberExpr:
		return fmt.Sprintf("%s.%s", t.resolveIdentity(e.Obj()), t.resolveIdentity(e.Prop()))
	}

	return ""
}

func (t *Tac) parseArrayList(v *Var, expr *ast.ArrayLitExpr, ln ast.Expr) int {
	list := expr.List()
	if ln, ok := ln.(*ast.NumberLitExpr); ok && len(list) > int(ln.Value()) {
		panic("too short decl")
	}

	size, ok := t.size[v.Name()]
	if !ok {
		panic(fmt.Sprintf("unknown var: %s", v.Label()))
	}

	for i := 0; i < len(list); i++ {
		sp := t.temp()
		t.code = append(t.code, &Assign{
			target: sp,
			op:     token.PLUS,
			arg1:   v,
			arg2:   &IntConst{int: int64(i) * size},
		})

		t.code = append(t.code, &Store{
			Destination: sp,
			Value:       t.parseExpr(list[i]),
			Size:        int(size),
		})
	}

	return len(list)
}

func (t *Tac) parseArrayAccess(expr *ast.ArrayAccessExpr) (Value, int64) {
	target := t.parseExpr(expr.Target())
	size := types.Pointer.Size()
	if v, ok := target.(Variable); ok {
		size = t.size[v.Name()]
	}

	offset := t.temp()
	t.code = append(t.code, &Assign{
		target: offset,
		op:     token.STAR,
		arg1:   t.parseExpr(expr.Address()),
		arg2:   &IntConst{int: size},
	})

	tmp := t.temp()
	t.code = append(t.code, &Assign{
		target: tmp,
		op:     token.PLUS,
		arg1:   target,
		arg2:   offset,

		isTemp: true,
	})

	return tmp, size
}

func (t *Tac) parseMemberExpr(expr *ast.MemberExpr) (Value, int) {
	target, ok := expr.Obj().(*ast.IdentExpr)
	if !ok {
		panic("not implemented")
	}

	property, ok := expr.Prop().(*ast.IdentExpr)
	if !ok {
		panic("not implemented")
	}

	s := t.types[target.Value()]
	typ := s.fields[property.Value()]

	size := typ.info.Size()
	offset := typ.offset

	tmp := t.temp()
	t.code = append(t.code, &Assign{
		target: tmp,
		op:     token.PLUS,
		arg1:   &Var{name: target.Value()},
		arg2:   &IntConst{int: int64(offset)}, // TODO: byte array

		isTemp: true,
	})

	return tmp, int(size)
}

func (t *Tac) parseExpr(expr ast.Expr) Value {
	switch expr := expr.(type) {
	default:
		panic(fmt.Sprintf("not implemented: %T %+v", expr, expr))
	case *ast.MemberExpr:
		tmp := t.temp()
		target, size := t.parseMemberExpr(expr)
		t.code = append(t.code, &Assign{
			target: tmp,
			op:     token.ASSIGN,
			arg1:   &Dereference{Addr: target},

			size: int(size),

			isTemp: true,
		})
		return tmp
	case *ast.ArrayAccessExpr:
		tmp := t.temp()
		arrayTarget, size := t.parseArrayAccess(expr)
		t.code = append(t.code, &Assign{
			target: tmp,
			op:     token.ASSIGN,
			arg1:   &Dereference{Addr: arrayTarget},

			size: int(size),

			isTemp: true,
		})
		return tmp
	case *ast.CallExpr:
		return t.parseCall(expr)
	case *ast.IdentExpr:
		v := &Var{name: expr.Value()}
		size, ok := t.size[v.Name()]
		if !ok {
			panic(fmt.Sprintf("unknown size of var %s", v.Label()))
		}

		tmp := t.temp()
		t.code = append(t.code, &Assign{
			target: tmp,
			op:     token.ASSIGN,
			arg1:   v,

			size: int(size),

			isTemp: true,
		})
		t.size[tmp.Name()] = size

		return tmp
	case *ast.BoolLitExpr:
		return &BoolConst{bool: expr.Value()}
	case *ast.NumberLitExpr:
		return &IntConst{int: expr.Value()}
	case *ast.StringLitExpr:
		s := expr.Unquoted()
		size := int64(len(s)) + 1

		tmp := t.temp()
		t.code = append(t.code, &Alloca{
			ptr:  tmp,
			size: &IntConst{int: size},
		})

		for i := 0; i < len(s); i++ {
			t.storeByte(tmp, i, s[i])
		}
		t.storeByte(tmp, len(s), 0)

		return tmp
	case *ast.BinaryOp:
		l, r := t.parseExpr(expr.Left()), t.parseExpr(expr.Right())

		if f := eval(expr.Op(), l, r); f != nil {
			tmp := t.temp()
			t.code = append(t.code, &Assign{
				target: tmp,
				op:     token.ASSIGN,
				arg1:   f,

				size: -1,

				isTemp: true,
			})
			return tmp
		} else {
			tmp := t.temp()
			t.code = append(t.code, &Assign{
				target: tmp,
				// TODO FIX THIS ЩИТ
				op:   expr.Op(),
				arg1: l,
				arg2: r,

				size: max(t.sizeof(l), t.sizeof(r)),

				isTemp: true,
			})
			return tmp
		}

	case *ast.UnaryOp:
		return t.parseUnaryOpExpr(expr)
	case nil:
		return nil
	}
}

func (t *Tac) parseUnaryOpExpr(expr *ast.UnaryOp) Value {
	switch expr.Op() {
	default:
		panic(fmt.Sprintf("not implemented: %s", expr.Op().String()))
	case token.MINUS:
		tmp := t.temp()

		r := t.parseExpr(expr.Value())
		if c, ok := r.(*IntConst); ok {
			c.int = -c.int
			t.code = append(t.code, &Assign{
				target: tmp,
				op:     token.ASSIGN,
				arg1:   c,

				size: -1,

				isTemp: true,
			})
		} else {
			t.code = append(t.code, &Assign{
				target: tmp,
				arg1:   r,
				op:     token.XOR,
				arg2:   &IntConst{int: -1},

				size: -1,

				isTemp: true,
			})
		}

		return tmp
	case token.AMP:
		switch r := expr.Value().(type) {
		default:
			panic(fmt.Sprintf("not implemented: %+v", r))
		case *ast.IdentExpr:
			addr := t.temp()
			t.code = append(t.code, &Assign{
				target: addr,
				op:     token.ASSIGN,
				arg1:   &AddressOf{Target: &Var{name: r.Value()}},

				size: int(types.Pointer.Size()),
			})
			return addr
		case *ast.ArrayAccessExpr:
			tmp := t.temp()
			arrayTarget, _ := t.parseArrayAccess(r)
			t.code = append(t.code, &Assign{
				target: tmp,
				op:     token.ASSIGN,
				arg1:   arrayTarget,

				size: int(types.Pointer.Size()),

				isTemp: true,
			})
			return tmp
		}
	case token.STAR:
		tmp := t.temp()
		t.code = append(t.code, &Assign{
			target: tmp,
			op:     token.ASSIGN,
			arg1:   &Dereference{Addr: t.parseExpr(expr.Value())},

			size: -1,

			isTemp: true,
		})
		return tmp
	}
}

func (t *Tac) storeByte(target Value, at int, b byte) {
	tmp := t.temp()
	t.code = append(t.code, &Assign{
		target: tmp,
		op:     token.PLUS,
		arg1:   target,
		arg2:   &IntConst{int: int64(at)},

		size: int(types.Pointer.Size()),
	})

	t.code = append(t.code, &Store{
		Destination: tmp,
		Value:       &IntConst{int: int64(b)},
		Size:        1,
	})
}

func (t *Tac) sizeof(v Value) int {
	switch v := v.(type) {
	default:
		panic(fmt.Sprintf("undefined type: %T %+v", v, v))
	case Variable:
		return int(t.size[v.Name()])
	case *IntConst:
		return 1
	case *BoolConst:
		return 1
	}
}
