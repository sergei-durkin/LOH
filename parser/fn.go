package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
	"loh/types"
)

func (s *syntaxer) parseFnDecl() (ast.Stmt, error) {
	const fn = "parseFnDecl"

	if s.peek().Token() != token.FN {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.FN", fn), s.peek())
	}
	s.consume()

	pos := s.cur.Pos()

	if s.peek().Token() != token.ID {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.ID", fn), s.peek())
	}
	s.consume()

	name := s.lex.GetValue(s.cur)

	if s.peek().Token() != token.LPAR {
		return nil, WrapErr(fmt.Errorf("[%s] expected (", fn), s.peek())
	}
	s.consume()

	args, err := s.parseFnArgsDecl()
	if err != nil {
		return nil, WrapErr(fmt.Errorf("[%s] error parse args: %w", fn, err), s.peek())
	}

	if s.peek().Token() != token.RPAR {
		return nil, WrapErr(fmt.Errorf("[%s] expected ) at pos %d", fn, s.peek().Pos()), s.peek())
	}
	s.consume()

	var rt ast.Expr
	if s.peek().Token() == token.ID {
		s.consume()

		rt = ast.NewIdentExpr(s.cur.Pos(), s.lex.GetValue(s.cur))
	} else if s.peek().Token() == token.STAR {
		s.consume()
		pos := s.cur.Pos()
		op := s.cur.Token()

		if s.peek().Token() == token.ID {
			s.consume()

			id := ast.NewIdentExpr(s.cur.Pos(), s.lex.GetValue(s.cur))
			rt = ast.NewUnaryOp(pos, op, id)
		}
	}

	block, err := s.parseBlockStmts()
	if err != nil {
		return nil, WrapErr(fmt.Errorf("[%s] parseStmt err: %w", fn, err), s.peek())
	}

	return ast.NewFnDecl(pos, name, args, rt, *block), nil
}

func (s *syntaxer) parseBlockStmts() (*ast.BlockStmt, error) {
	const fn = "parseStmts"

	pos := s.cur.Pos()

	if s.peek().Token() != token.LBRACE {
		return nil, WrapErr(fmt.Errorf("[%s] expected { but %s given", fn, s.peek().Token().String()), s.peek())
	}
	s.consume()

	for s.peek().Token() == token.EOL {
		s.consume()
	}

	var stmts []ast.Stmt
	for s.peek().Token() != token.RBRACE {
		stmt, err := s.parseStmt()
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parseStmt err: %w", fn, err), s.peek())
		}

		stmts = append(stmts, stmt)

		for s.peek().Token() == token.EOL {
			s.consume()
		}
	}

	if s.peek().Token() != token.RBRACE {
		return nil, WrapErr(fmt.Errorf("[%s] expected }", fn), s.peek())
	}
	s.consume()

	return ast.NewBlockStmt(pos, stmts), nil
}

func (s *syntaxer) parseFnArgsDecl() ([]*ast.FnArgDecl, error) {
	const fn = "parseFnArgsDecl"

	args := []*ast.FnArgDecl{}
	for s.peek().Token() != token.RPAR {
		if s.peek().Token() != token.ID {
			return nil, WrapErr(fmt.Errorf("[%s] expected token.ID", fn), s.peek())
		}
		s.consume()

		apos := s.cur.Pos()
		aname := s.lex.GetValue(s.cur)

		switch s.peek().Token() {
		default:
			panic(fmt.Sprintf("unexpected token %d %s", s.peek().Token(), s.peek().Token().String()))
		case token.STAR:
			s.consume()

			pos := s.cur.Pos()
			op := s.cur.Token()

			if s.peek().Token() != token.ID {
				return nil, WrapErr(fmt.Errorf("[%s] expected token.ID", fn), s.peek())
			}
			s.consume()

			id := ast.NewIdentExpr(s.cur.Pos(), s.lex.GetValue(s.cur))
			typeExpr := ast.NewUnaryOp(pos, op, id)
			size := types.Pointer.Size()

			args = append(args, ast.NewFnArgDecl(apos, aname, typeExpr, size))
		case token.ID:
			s.consume()

			typeName := s.lex.GetValue(s.cur)
			typeExpr := ast.NewIdentExpr(s.cur.Pos(), typeName)
			size := types.Info(typeName).Size()

			args = append(args, ast.NewFnArgDecl(apos, aname, typeExpr, size))
		}

		if s.peek().Token() == token.COM {
			s.consume()
		}
	}

	return args, nil
}
