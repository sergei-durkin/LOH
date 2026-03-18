package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
)

func (s *syntaxer) parseCallExpr(left ast.Expr) (ast.Expr, error) {
	const fn = "parseCallArgs"

	if s.peek().Token() != token.LPAR {
		return nil, WrapErr(fmt.Errorf("[%s] expected (", fn), s.peek())
	}
	s.consume()

	args := []ast.Expr{}
	for s.peek().Token() != token.RPAR {
		expr, err := s.parseExpr(token.LowestPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse expr err: %w", fn, err), s.peek())
		}

		args = append(args, expr)

		if s.peek().Token() == token.COM {
			s.consume()
		}
	}

	if s.peek().Token() != token.RPAR {
		return nil, WrapErr(fmt.Errorf("[%s] expected )", fn), s.peek())
	}
	s.consume()

	return ast.NewCallExpr(left.Pos(), left, args), nil
}
