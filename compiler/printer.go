package compiler

import (
	"fmt"
	"io"
	"strings"
)

func (t *Tac) Print(w io.Writer) {
	for i := 0; i < len(t.code); i++ {
		t.code[i].Print(w)
	}
}

func (c *Cfg) Print(w io.Writer) {
	for i := 0; i < len(c.blocks); i++ {
		fmt.Fprintf(w, "BB%d:\n", c.blocks[i].ID)
		for j := 0; j < len(c.blocks[i].Phi); j++ {
			c.blocks[i].Phi[j].Print(w)
		}
		for j := 0; j < len(c.blocks[i].Instr); j++ {
			c.blocks[i].Instr[j].Print(w)
		}
	}
}

func (s *Ssa) Print(w io.Writer) {
	for i := 0; i < len(s.funcs); i++ {
		s.funcs[i].Print(w)
	}
}

func (s SparseTree) Print(w io.Writer) {
	for i := 0; i < len(s); i++ {
		fmt.Fprintf(w, "ID: %d\n", i)

		fmt.Fprintf(w, "\tParent: ")
		if s[i].par != nil {
			fmt.Fprintf(w, "%d\n", s[i].par.ID)
		} else {
			fmt.Fprintf(w, "<nil>\n")
		}

		fmt.Fprintf(w, "\tChild: ")
		if s[i].child != nil {
			fmt.Fprintf(w, "%d\n", s[i].child.ID)
		} else {
			fmt.Fprintf(w, "<nil>\n")
		}

		fmt.Fprintf(w, "\tSibling: ")
		if s[i].sib != nil {
			fmt.Fprintf(w, "%d\n", s[i].sib.ID)
		} else {
			fmt.Fprintf(w, "<nil>\n")
		}

		fmt.Fprintf(w, "\tEntry: %d\n", s[i].entry)
		fmt.Fprintf(w, "\tExit: %d\n", s[i].exit)
	}
}

func (bb *BasicBlock) Print(w io.Writer) {
	fmt.Fprintf(w, "BB%d:\n", bb.ID)
	fmt.Fprintf(w, "\tLabel: %s\n", bb.Name)

	fmt.Fprintf(w, "\tPred: ")
	buf := make([]int, 0, len(bb.Pred))
	for i := 0; i < len(bb.Pred); i++ {
		buf = append(buf, bb.Pred[i].ID)
	}
	fmt.Fprintln(w, buf)

	fmt.Fprintf(w, "\tSucc: ")
	buf = make([]int, 0, len(bb.Succ))
	for i := 0; i < len(bb.Succ); i++ {
		buf = append(buf, bb.Succ[i].ID)
	}
	fmt.Fprintln(w, buf)
}

func (f *Fn) Print(w io.Writer) {
	reachable := f.reachable()

	for i := 0; i < len(f.po); i++ {
		f.po[i].Print(w)
		if f.idom[f.po[i].ID] != nil {
			fmt.Fprintf(w, "\tDom: %d\n", f.idom[f.po[i].ID].ID)
		}

		fmt.Fprintf(w, "\tDF: ")
		buf := make([]int, 0, len(f.cfg.blocks))
		for j := 0; j < len(f.domfront[f.po[i].ID]); j++ {
			if f.domfront[f.po[i].ID][j] {
				buf = append(buf, j)
			}
		}
		fmt.Fprintln(w, buf)
	}

	for i := len(f.po) - 1; i >= 0; i-- {
		bb := f.po[i]
		if !reachable[bb.ID] {
			continue
		}

		fmt.Fprintf(w, "BB%d:\n", bb.ID)
		for j := 0; j < len(bb.Phi); j++ {
			if reg, ok := f.virtualRegs[bb.Phi[j].Var.Label()]; ok {
				fmt.Fprintf(w, "\tREG_%d:", reg)
			}
			bb.Phi[j].Print(w)
		}
		for j := 0; j < len(bb.Instr); j++ {
			if in, ok := bb.Instr[j].(*Alloca); ok {
				if reg, ok := f.virtualRegs[in.ptr.(Variable).Label()]; ok {
					fmt.Fprintf(w, "\tREG_%d:", reg)
				}
			}
			if in, ok := bb.Instr[j].(*Assign); ok {
				if reg, ok := f.virtualRegs[in.target.Label()]; ok {
					fmt.Fprintf(w, "\tREG_%d:", reg)
				}
			}
			if in, ok := bb.Instr[j].(*Call); ok {
				if reg, ok := f.virtualRegs[in.target.Label()]; ok {
					fmt.Fprintf(w, "\tREG_%d:", reg)
				}
			}
			bb.Instr[j].Print(w)
		}
	}
}

func (a *Assign) Print(w io.Writer) {
	arg1 := toStr(a.arg1)
	if a.arg2 != nil {
		arg2 := toStr(a.arg2)
		fmt.Fprintf(w, "\t%s = %s %s %s\n", a.target.Label(), arg1, a.op, arg2)
		return
	}

	fmt.Fprintf(w, "\t%s = %s\n", a.target.Label(), arg1)
}

func (c *Call) Print(w io.Writer) {
	values := []string{}
	for i := 0; i < len(c.args); i++ {
		switch arg := c.args[i].(type) {
		case *TempVar:
			values = append(values, arg.Label())
		case *Var:
			values = append(values, arg.Label())
		case *IntConst:
			values = append(values, fmt.Sprint(arg.int))
		case *BoolConst:
			values = append(values, fmt.Sprint(arg.bool))
		case *AddressOf:
			values = append(values, fmt.Sprintf("&%s", toStr(arg.Target)))
		case *Dereference:
			values = append(values, toStr(arg))
		}
	}

	fmt.Fprintf(w, "\t%s = call %s (%s)\n", c.target.Label(), c.callee, strings.Join(values, ", "))
}

func (a *Alloca) Print(w io.Writer) {
	fmt.Fprintf(w, "\t%s = alloca (%s)\n", toStr(a.ptr), toStr(a.size))
}

func (l *Label) Print(w io.Writer) {
	fmt.Fprintf(w, "%d:\n", l.id)
}

func (g *Goto) Print(w io.Writer) {
	fmt.Fprintf(w, "\tgoto %d\n", g.label)
}

func (i *IfGoto) Print(w io.Writer) {
	fmt.Fprintf(w, "\tif not %s goto %d else goto %d\n", toStr(i.cond), i.label, i.fall)
}

func (r *Return) Print(w io.Writer) {
	fmt.Fprintf(w, "\treturn %s\n", toStr(r.value))
}

func (i *Store) Print(w io.Writer) {
	fmt.Fprintf(w, "\tstore at %s valueOf(%s)\n", toStr(i.Destination), toStr(i.Value))
}

func (i *Load) Print(w io.Writer) {
	fmt.Fprintf(w, "\tload %s into %s\n", toStr(i.Source), toStr(i.Destination))
}

func (p *Phi) Print(w io.Writer) {
	args := make([]string, 0, len(p.Args))
	for p, a := range p.Args {
		switch arg := a.(type) {
		case *Var:
			args = append(args, fmt.Sprintf("%%BB%d:%s_%d", p, arg.name, arg.ver))
		case *IntConst:
			args = append(args, fmt.Sprintf("%%BB%d:%d", p, arg.int))
		case *BoolConst:
			args = append(args, fmt.Sprintf("%%BB%d:%v", p, arg.bool))
		}
	}

	fmt.Fprintf(w, "\t%s_%d = Phi: (%s)\n", p.Var.name, p.Var.ver, strings.Join(args, ", "))
}

func toStr(arg any) string {
	switch a := arg.(type) {
	default:
		return fmt.Sprintf("Unknown %T: %+v", a, a)
	case nil:
		return "<nil>"
	case *ArgReg:
		return fmt.Sprintf("X%d", a.ID)
	case *Label:
		return fmt.Sprintf("%d", a.id)
	case *AddressOf:
		return fmt.Sprintf("&%s", toStr(a.Target))
	case *Dereference:
		return fmt.Sprintf("*(%s)", toStr(a.Addr))
	case *Reg:
		return fmt.Sprintf("x%d", a.ID)
	case *TempVar:
		return a.label
	case *Var:
		return fmt.Sprintf("%s_%d", a.name, a.ver)
	case *BoolConst:
		return fmt.Sprintf("%v", a.bool)
	case *IntConst:
		return fmt.Sprintf("%d", a.int)
	case *StringConst:
		return a.string
	}
}
