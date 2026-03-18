package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
)

func (s *syntaxer) parseIfStmt() (ast.Stmt, error) {
	const fn = "parseIfStmt"

	if s.peek().Token() != token.IF {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.IF", fn), s.peek())
	}
	s.consume()

	pos := s.cur.Pos()

	expr, err := s.parseExpr(s.cur.Token().Priority())
	if err != nil {
		return nil, err
	}

	body, err := s.parseBlockStmts()
	if err != nil {
		return nil, err
	}

	if s.peek().Token() != token.ELSE {
		return ast.NewIfStmt(pos, expr, *body, nil), nil
	}
	s.consume()

	if s.peek().Token() == token.LBRACE {
		block, err := s.parseBlockStmts()
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parseBlockStmts err: %w", fn, err), s.peek())
		}

		return ast.NewIfStmt(pos, expr, *body, block), nil
	}

	if s.peek().Token() != token.IF {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.IF", fn), s.peek())
	}

	elseIf, err := s.parseIfStmt()
	if err != nil {
		return nil, WrapErr(fmt.Errorf("[%s] parseIfStmt err: %w", fn, err), s.peek())
	}

	return ast.NewIfStmt(pos, expr, *body, elseIf), nil
}
