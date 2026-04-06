package parser

import (
	"errors"
	"fmt"
	"loh/ast"
	"loh/lexer"
	"loh/token"
)

type ParserError struct {
	error

	Lexeme *lexer.Lexeme
}

func WrapErr(err error, lexeme *lexer.Lexeme) error {
	return &ParserError{
		error:  err,
		Lexeme: lexeme,
	}
}

type syntaxer struct {
	lex       *lexer.Lexer
	cur       *lexer.Lexeme
	lookahead *lexer.Lexeme
}

func NewSyntaxer(l *lexer.Lexer) *syntaxer {
	s := &syntaxer{
		lex: l,
	}

	s.consume()

	return s
}

func Parse(buf []byte) (*ast.AST, error) {
	if len(buf) == 0 {
		return ast.NewAST(nil), nil
	}

	l := lexer.NewLexer(buf)
	s := NewSyntaxer(l)

	var decls []ast.Node
	for cur, err := s.Next(); cur != nil || err != nil; cur, err = s.Next() {
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[Parse] next token err: %w", err), s.peek())
		}

		decls = append(decls, cur)
	}

	return ast.NewAST(ast.NewCompileUnit("", decls)), nil
}

func (s *syntaxer) Next() (ast.Node, error) {
	for s.peek().Token() == token.EOL ||
		s.peek().Token() == token.COMMENT ||
		s.peek().Token() == token.MLCOMMENT {
		s.consume()
	}

	switch s.peek().Token() {
	case token.CONST:
		return s.parseConstDecl()
	case token.TYPE:
		return s.parseStructDecl()
	case token.FN:
		return s.parseFnDecl()
	case token.EOF, token.EOL, 0:
		return nil, nil
	}

	return nil, WrapErr(errors.New("unknown token"), s.peek())
}

func (s *syntaxer) parseStmt() (ast.Stmt, error) {
	for s.peek().Token() == token.EOL ||
		s.peek().Token() == token.COMMENT ||
		s.peek().Token() == token.MLCOMMENT {
		s.consume()
	}

	switch s.peek().Token() {
	case token.RBRACE:
		return nil, nil
	case token.VAR:
		return s.parseVarDecl()
	case token.CONST:
		return s.parseConstDecl()
	case token.STRUCT:
		return s.parseStructDecl()
	case token.RETURN:
		return s.parseReturnStmt()
	case token.IF:
		return s.parseIfStmt()
	case token.FOR:
		return s.parseForStmt()
	case token.ID:
		return s.parseIDStmt()
	case token.STAR:
		return s.parseStarIDStmt()
	case token.CONTINUE:
		return s.parseContinue()
	case token.BREAK:
		return s.parseBreak()
	}

	return nil, WrapErr(errors.New("unknown stmt"), s.peek())
}

func (s *syntaxer) parseContinue() (ast.Stmt, error) {
	s.consume()
	return ast.NewContinueStmt(s.cur.Pos()), nil
}

func (s *syntaxer) parseBreak() (ast.Stmt, error) {
	s.consume()
	return ast.NewBreakStmt(s.cur.Pos()), nil
}

func (s *syntaxer) parseStarIDStmt() (ast.Stmt, error) {
	const fn = "starID"

	pos := s.cur.Pos()

	if s.peek().Token() != token.STAR {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.STAR", fn), s.peek())
	}
	s.consume()

	var err error

	var left ast.Expr
	if s.peek().Token() == token.LPAR {
		s.consume()
		left, err = s.parseExpr(token.LowestPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse member expr err: %w", fn, err), s.peek())
		}

		if s.peek().Token() != token.RPAR {
			return nil, WrapErr(fmt.Errorf("[%s] expected )", fn), s.peek())
		}
		s.consume()
	} else {
		left, err = s.parseIDExpr()
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse member expr err: %w", fn, err), s.peek())
		}
	}
	left = ast.NewUnaryOp(pos, token.STAR, left)

	switch s.peek().Token() {
	case token.ASSIGN:
		s.consume()

		rhs, err := s.parseExpr(token.LowestPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse rhs err: %w", fn, err), s.peek())
		}

		return ast.NewAssignStmt(left.Pos(), left, rhs), nil
	case token.LPAR:
		call, err := s.parseCallExpr(left)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse call err: %w", fn, err), s.peek())
		}

		return ast.NewExprStmt(call.Pos(), call), nil
	case token.PLUS:
		s.consume()
		pos := s.cur.Pos()

		if s.peek().Token() == token.PLUS {
			s.consume()

			return ast.NewExprStmt(left.Pos(), ast.NewPostIncOp(pos, left)), nil
		}

		return nil, WrapErr(fmt.Errorf("[%s] expected token.PLUS", fn), s.peek())
	case token.MINUS:
		s.consume()
		pos := s.cur.Pos()

		if s.peek().Token() == token.MINUS {
			s.consume()

			return ast.NewExprStmt(left.Pos(), ast.NewPostDecOp(pos, left)), nil
		}

		return nil, WrapErr(fmt.Errorf("[%s] expected token.MINUS", fn), s.peek())
	}

	return nil, WrapErr(fmt.Errorf("[%s] unknown STMT", fn), s.peek())
}

func (s *syntaxer) parseIDStmt() (ast.Stmt, error) {
	const fn = "ID"

	if s.peek().Token() != token.ID {
		return nil, WrapErr(fmt.Errorf("[%s] expected token.ID", fn), s.peek())
	}

	left, err := s.parseIDExpr()
	if err != nil {
		return nil, WrapErr(fmt.Errorf("[%s] parse member expr err: %w", fn, err), s.peek())
	}

	switch s.peek().Token() {
	case token.ASSIGN:
		s.consume()

		rhs, err := s.parseExpr(token.LowestPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse rhs err: %w", fn, err), s.peek())
		}

		return ast.NewAssignStmt(left.Pos(), left, rhs), nil
	case token.LPAR:
		call, err := s.parseCallExpr(left)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse call err: %w", fn, err), s.peek())
		}

		return ast.NewExprStmt(call.Pos(), call), nil
	case token.PLUS:
		s.consume()
		pos := s.cur.Pos()

		if s.peek().Token() == token.PLUS {
			s.consume()

			return ast.NewExprStmt(left.Pos(), ast.NewPostIncOp(pos, left)), nil
		}

		return nil, WrapErr(fmt.Errorf("[%s] expected token.PLUS", fn), s.peek())
	case token.MINUS:
		s.consume()
		pos := s.cur.Pos()

		if s.peek().Token() == token.MINUS {
			s.consume()

			return ast.NewExprStmt(left.Pos(), ast.NewPostDecOp(pos, left)), nil
		}

		return nil, WrapErr(fmt.Errorf("[%s] expected token.MINUS", fn), s.peek())
	}

	return nil, WrapErr(fmt.Errorf("[%s] unknown STMT", fn), s.peek())
}

func (s *syntaxer) peek() *lexer.Lexeme {
	return s.lookahead
}

func (s *syntaxer) consume() {
	n := s.lex.Next()
	s.cur = s.lookahead
	s.lookahead = n
}

func (s *syntaxer) print(t *lexer.Lexeme) {
	fmt.Printf("Token '%d' at '%d': '%s'\n", t.Token(), t.Pos(), s.lex.GetValue(t))
}
