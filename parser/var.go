package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
	"loh/types"
	"strconv"
)

func (s *syntaxer) parseVarDecl() (ast.Stmt, error) {
	const fn = "parseVarDecl"

	if s.peek().Token() == token.SCOL {
		return nil, nil
	}

	if s.peek().Token() != token.VAR {
		return nil, fmt.Errorf("[%s] expected lexer.VAR", fn)
	}
	s.consume()

	pos := s.cur.Pos()

	if s.peek().Token() != token.ID {
		return nil, fmt.Errorf("[%s] expected lexer.ID", fn)
	}
	s.consume()

	name := s.lex.GetValue(s.cur)

	typ, err := s.parseVarDeclType()
	if err != nil {
		return nil, fmt.Errorf("[%s] parse typeDecl err: %w", fn, err)
	}

	if s.peek().Token().IsTerm() {
		return ast.NewVarDecl(pos, name, typ, nil), nil
	}

	if s.peek().Token() != token.ASSIGN {
		return nil, fmt.Errorf("[%s] expected =", fn)
	}
	s.consume()

	value, err := s.parseExpr(token.LowestPrec)
	if err != nil {
		return nil, fmt.Errorf("[%s] parse var decl err: %w", fn, err)
	}

	return ast.NewVarDecl(pos, name, typ, value), nil
}

func (s *syntaxer) parseVarDeclType() (ast.Type, error) {
	const fn = "parseVarDeclType"

	pos := s.peek().Pos()

	arrayDecl, ok, err := s.parseArrayDecl()
	if err != nil {
		return nil, fmt.Errorf("[%s] parse arrayDecl decl err: %w", fn, err)
	}

	if s.peek().Token() != token.ID {
		return nil, fmt.Errorf("[%s] expected lexer.ID", fn)
	}
	s.consume()

	tname := s.lex.GetValue(s.cur)
	typeInfo := types.Info(tname)

	var varType ast.Type
	if typeInfo == nil {
		varType = ast.NewCustomType(pos, tname)
	} else {
		varType = ast.NewBasicType(pos, typeInfo)
	}

	if ok {
		return ast.NewArrayType(pos, varType, arrayDecl), nil
	}

	return varType, nil
}

func (s *syntaxer) parseArrayDecl() (arrLen ast.Expr, ok bool, err error) {
	const fn = "parseArrayDecl"

	if s.peek().Token() != token.LSQBR {
		return nil, false, nil
	}

	s.consume()

	pos := s.cur.Pos()

	switch s.peek().Token() {
	default:
		return nil, false, fmt.Errorf("[%s] parse var decl expected num or id", fn)
	case token.NUM:
		s.consume()
		arr, err := strconv.ParseInt(s.lex.GetValue(s.cur), 10, 64)
		if err != nil {
			return nil, false, fmt.Errorf("[%s] parse var size decl err: %w", fn, err)
		}
		arrLen = ast.NewNumberLitExpr(pos, arr)
	case token.ID:
		arrLen, err = s.parseExpr(token.LowestPrec)
		if err != nil {
			return nil, false, fmt.Errorf("[%s] parse var decl expr err: %w", fn, err)
		}
	case token.RSQBR:
	}

	if s.peek().Token() != token.RSQBR {
		return nil, false, fmt.Errorf("[%s] expected ]", fn)
	}
	s.consume()

	return arrLen, true, nil
}
