package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
)

func (s *syntaxer) parseReturnStmt() (ast.Stmt, error) {
	const fn = "parseReturnStmt"

	if s.peek().Token() != token.RETURN {
		return nil, fmt.Errorf("[%s] expected lexer.RETURN", fn)
	}
	s.consume()

	pos := s.cur.Pos()

	if s.peek().Token() == token.EOL {
		return ast.NewReturnStmt(pos, nil), nil
	}

	expr, err := s.parseExpr(token.LowestPrec)
	if err != nil {
		return nil, err
	}

	return ast.NewReturnStmt(pos, expr), nil
}
