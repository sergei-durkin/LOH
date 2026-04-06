package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
	"loh/types"
)

func (s *syntaxer) parseStructDecl() (ast.Stmt, error) {
	const fn = "parseStructDecl"

	if s.peek().Token() != token.TYPE {
		return nil, WrapErr(fmt.Errorf("[%s] expected lexer.STRUCT", fn), s.peek())
	}
	s.consume()

	pos := s.cur.Pos()

	if s.peek().Token() != token.ID {
		return nil, WrapErr(fmt.Errorf("[%s] expected lexer.ID", fn), s.peek())
	}
	s.consume()

	name := s.lex.GetValue(s.cur)

	if s.peek().Token() != token.STRUCT {
		return nil, WrapErr(fmt.Errorf("[%s] expected `struct`", fn), s.peek())
	}
	s.consume()

	if s.peek().Token() != token.LBRACE {
		return nil, WrapErr(fmt.Errorf("[%s] expected {", fn), s.peek())
	}
	s.consume()

	if s.peek().Token() == token.EOL {
		s.consume()
	}

	var fields []*ast.StructFieldDecl
	for s.peek().Token() != token.RBRACE {
		if s.peek().Token() != token.ID {
			return nil, WrapErr(fmt.Errorf("[%s] expected lexer.ID", fn), s.peek())
		}
		s.consume()

		fpos := s.cur.Pos()
		fname := s.lex.GetValue(s.cur)

		var isPtr bool
		if s.peek().Token() == token.LSQBR {
			s.consume()

			if s.peek().Token() != token.RSQBR {
				return nil, WrapErr(fmt.Errorf("[%s] expected ]", fn), s.peek())
			}
			s.consume()

			isPtr = true
		}

		if s.peek().Token() != token.ID {
			return nil, WrapErr(fmt.Errorf("[%s] expected lexer.ID", fn), s.peek())
		}
		s.consume()

		ftyp := types.Info(s.lex.GetValue(s.cur))
		if isPtr {
			ftyp = types.Pointer
		}

		fields = append(fields, ast.NewStructFieldDecl(fpos, fname, ast.NewBasicType(s.cur.Pos(), ftyp)))

		if s.peek().Token() == token.EOL {
			s.consume()
		}
	}

	if s.peek().Token() != token.RBRACE {
		return nil, WrapErr(fmt.Errorf("[%s] expected } at pos %d", fn, s.peek().Pos()), s.peek())
	}
	s.consume()

	if s.peek().Token() == token.EOL {
		s.consume()
	}

	return ast.NewStructDecl(pos, name, fields), nil
}

func (s *syntaxer) parseStructLitExpr(left ast.Expr) (ast.Expr, error) {
	const fn = "parseStructLiteral"

	if s.peek().Token() != token.LBRACE {
		return nil, WrapErr(fmt.Errorf("[%s] expected { at pos %d", fn, s.peek().Pos()), s.peek())
	}
	s.consume()

	if s.peek().Token() == token.RBRACE {
		return ast.NewStructLitExpr(left.Pos(), left, nil), nil
	}

	if s.peek().Token() == token.EOL {
		s.consume()
	}

	var fields []*ast.StructLitExprField
	for s.peek().Token() == token.ID {
		s.consume()
		pos := s.cur.Pos()

		name := s.lex.GetValue(s.cur)

		if s.peek().Token() != token.COL {
			return nil, WrapErr(fmt.Errorf("[%s] expected : at pos %d", fn, s.peek().Pos()), s.peek())
		}
		s.consume()

		expr, err := s.parseExpr(token.LowestPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse expr err: %w", fn, err), s.peek())
		}

		if s.peek().Token() == token.COM {
			s.consume()
		}

		for s.peek().Token() == token.EOL {
			s.consume()
		}

		fields = append(fields, ast.NewStructLitExprField(pos, name, expr))

		if s.peek().Token() == token.RBRACE {
			break
		}
	}

	if s.peek().Token() == token.EOL {
		s.consume()
	}

	if s.peek().Token() != token.RBRACE {
		return nil, WrapErr(fmt.Errorf("[%s] expected } at pos %d", fn, s.peek().Pos()), s.peek())
	}
	s.consume()

	return ast.NewStructLitExpr(left.Pos(), left, fields), nil
}
