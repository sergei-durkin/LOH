package compiler

type Ssa struct {
	funcs []*Fn
}

func NewSSA(c []*Cfg) *Ssa {
	s := &Ssa{}

	for i := 0; i < len(c); i++ {
		s.funcs = append(s.funcs, NewFn(c[i]))
	}

	return s
}

func (s *Ssa) Funcs() []*Fn {
	return s.funcs
}
