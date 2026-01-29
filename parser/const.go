package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
)

func (s *syntaxer) parseConstDecl() (ast.Stmt, error) {
	const fn = "parseConstDecl"

	if s.peek().Token() != token.CONST {
		return nil, fmt.Errorf("[%s] expected token.CONST", fn)
	}
	s.consume()

	pos := s.cur.Pos()

	if s.peek().Token() != token.ID {
		return nil, fmt.Errorf("[%s] expected token.ID", fn)
	}
	s.consume()

	name := s.lex.GetValue(s.cur)

	if s.peek().Token() != token.ID {
		return nil, fmt.Errorf("[%s] expected token.ID", fn)
	}
	s.consume()

	typeName := s.lex.GetValue(s.cur)

	if s.peek().Token() != token.ASSIGN {
		return nil, fmt.Errorf("[%s] expected =", fn)
	}
	s.consume()

	if s.peek().Token() == token.EOL {
		s.consume()
	}

	value, err := s.parseExpr(token.LowestPrec)
	if err != nil {
		return nil, fmt.Errorf("[%s] parse err: %w", fn, err)
	}

	return ast.NewConstDecl(pos, name, typeName, value), nil
}
