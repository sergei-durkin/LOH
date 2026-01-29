package compiler

import (
	"loh/token"
)

// тупой, но пока пойдет
func eval(op token.Token, a, b Value) Value {
	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	var (
		ai, bi   *IntConst
		aok, bok bool
	)

	ai, aok = a.(*IntConst)
	bi, bok = b.(*IntConst)
	if aok != bok {
		return nil
	}

	if !aok && !bok {
		return nil
	}

	res := &IntConst{}
	switch op {
	default:
		return nil
	case token.PLUS:
		res.int = ai.int + bi.int
	case token.MINUS:
		res.int = ai.int - bi.int
	case token.STAR:
		res.int = ai.int * bi.int
	case token.SLASH:
		if bi.int == 0 {
			panic("div by zero")
		}
		res.int = ai.int / bi.int
	case token.XOR:
		res.int = ai.int ^ bi.int
	case token.PIPE:
		res.int = ai.int | bi.int
	case token.AMP:
		res.int = ai.int & bi.int
	}

	return res
}
