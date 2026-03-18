package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
)

func (s *syntaxer) parseIDExpr() (ast.Expr, error) {
	const fn = "parseIDExpr"

	if s.peek().Token() != token.ID {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.ID", fn), s.peek())
	}
	s.consume()

	pos := s.cur.Pos()
	left := ast.NewIdentExpr(pos, s.lex.GetValue(s.cur))

	if s.peek().Token() == token.DOT {
		return s.parseMemberExpr(left)
	}

	if s.peek().Token() == token.LSQBR {
		return s.parseArrayAccessExpr(left)
	}

	return left, nil
}

func (s *syntaxer) parseArrayAccessExpr(left ast.Expr) (ast.Expr, error) {
	const fn = "parseArrayAccessExpr"

	if s.peek().Token() != token.LSQBR {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.LBR", fn), s.peek())
	}
	s.consume()

	expr, err := s.parseExpr(token.LowestPrec)
	if err != nil {
		return nil, WrapErr(fmt.Errorf("[%s] parse array access expr err: %w", fn, err), s.peek())
	}
	s.consume()

	return ast.NewArrayAccessExpr(left.Pos(), left, expr), nil
}

func (s *syntaxer) parseMemberExpr(left ast.Expr) (ast.Expr, error) {
	const fn = "parseMemberExpr"

	if s.peek().Token() != token.DOT {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.DOT", fn), s.peek())
	}
	s.consume()

	if s.peek().Token() != token.ID {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.ID", fn), s.peek())
	}
	s.consume()

	pos := s.cur.Pos()
	left = ast.NewMemberExpr(left.Pos(), left, ast.NewIdentExpr(pos, s.lex.GetValue(s.cur)))

	if s.peek().Token() == token.DOT {
		return s.parseMemberExpr(left)
	}

	return left, nil
}
