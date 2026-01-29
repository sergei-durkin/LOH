package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
)

func (s *syntaxer) parseCallExpr(left ast.Expr) (ast.Expr, error) {
	const fn = "parseCallArgs"

	if s.peek().Token() != token.LPAR {
		return nil, fmt.Errorf("[%s] expected (", fn)
	}
	s.consume()

	args := []ast.Expr{}
	for s.peek().Token() != token.RPAR {
		expr, err := s.parseExpr(token.LowestPrec)
		if err != nil {
			return nil, fmt.Errorf("[%s] parse expr err: %w", fn, err)
		}

		args = append(args, expr)

		if s.peek().Token() == token.COM {
			s.consume()
		}
	}

	if s.peek().Token() != token.RPAR {
		return nil, fmt.Errorf("[%s] expected )", fn)
	}
	s.consume()

	return ast.NewCallExpr(left.Pos(), left, args), nil
}
