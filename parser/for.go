package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
)

func (s *syntaxer) parseForStmt() (ast.Stmt, error) {
	const fn = "parseForStmt"

	if s.peek().Token() != token.FOR {
		return nil, fmt.Errorf("[%s] expected token.IF", fn)
	}
	s.consume()

	pos := s.cur.Pos()

	if s.peek().Token() == token.LBRACE {
		body, err := s.parseBlockStmts()
		if err != nil {
			return nil, fmt.Errorf("[%s] parse block stmts err: %w", fn, err)
		}
		return ast.NewForStmt(pos, nil, nil, nil, *body), nil
	}

	init, condition, post, err := s.parseForHeader()
	if err != nil {
		return nil, fmt.Errorf("[%s] parse for header err: %w", fn, err)
	}

	body, err := s.parseBlockStmts()
	if err != nil {
		return nil, fmt.Errorf("[%s] parse block stmts err: %w", fn, err)
	}

	return ast.NewForStmt(pos, init, condition, post, *body), nil
}

func (s *syntaxer) parseForHeader() (ast.Stmt, ast.Expr, ast.Stmt, error) {
	const fn = "parseForHeader"

	init, err := s.parseForInit()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("[%s] parse init stmt err: %w", fn, err)
	}

	if s.peek().Token() != token.SCOL {
		return nil, nil, nil, fmt.Errorf("[%s] expected ;", fn)
	}
	s.consume()

	condition, err := s.parseExpr(token.LowestPrec)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("[%s] parse expr err: %w", fn, err)
	}

	if s.peek().Token() != token.SCOL {
		return nil, nil, nil, fmt.Errorf("[%s] expected ;", fn)
	}
	s.consume()

	var post ast.Stmt
	if s.peek().Token() != token.LBRACE {
		post, err = s.parseIDStmt()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("[%s] parse post stmt err: %w", fn, err)
		}
	}

	return init, condition, post, nil
}

func (s *syntaxer) parseForInit() (ast.Stmt, error) {
	switch s.peek().Token() {
	case token.SCOL:
		return nil, nil
	case token.VAR:
		return s.parseVarDecl()
	case token.ID:
		return s.parseIDStmt()
	default:
		return nil, fmt.Errorf("unexpected token")
	}
}
