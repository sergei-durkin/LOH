package parser

import (
	"fmt"
	"loh/ast"
	"loh/token"
	"strconv"
)

func (s *syntaxer) parseExpr(bp int) (ast.Expr, error) {
	const fn = "parseExpr"

	left, err := s.parseNud(bp)
	if err != nil {
		return nil, WrapErr(fmt.Errorf("[%s] lhs parse err: %w", fn, err), s.peek())
	}

	for bp < s.peek().Token().Priority() {
		left, err = s.parseLed(left)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] led parse err: %w", fn, err), s.peek())
		}
	}

	return left, nil
}

func (s *syntaxer) parseNud(bp int) (ast.Expr, error) {
	const fn = "parseNud"

	switch s.peek().Token() {
	case token.RSQBR:
		return nil, nil
	case token.SCOL:
		return nil, nil
	case token.ID:
		s.consume()

		left := ast.NewIdentExpr(s.cur.Pos(), s.lex.GetValue(s.cur))

		if bp <= token.LowestPrec {
			if s.peek().Token() == token.LBRACE {
				return s.parseStructLitExpr(left)
			}
		}

		if s.peek().Token() == token.DOT {
			return s.parseMemberExpr(left)
		}
		if s.peek().Token() == token.LPAR {
			return s.parseCallExpr(left)
		}
		if s.peek().Token() == token.LSQBR {
			return s.parseArrayAccessExpr(left)
		}

		return left, nil
	case token.NUM:
		s.consume()

		lit, err := strconv.ParseInt(s.lex.GetValue(s.cur), 10, 64)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] err parse num: %w", fn, err), s.peek())
		}

		return ast.NewNumberLitExpr(s.cur.Pos(), lit), nil
	case token.TRUE:
		s.consume()

		return ast.NewBoolLitExpr(s.cur.Pos(), true), nil
	case token.FALSE:
		s.consume()

		return ast.NewBoolLitExpr(s.cur.Pos(), false), nil
	case token.STR:
		s.consume()

		return ast.NewStringLitExpr(s.cur.Pos(), s.lex.GetValue(s.cur)), nil
	case token.MINUS:
		s.consume()

		pos := s.cur.Pos()
		op := s.cur.Token()

		right, err := s.parseExpr(token.UnaryPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse expr err: %w", fn, err), s.peek())
		}

		return ast.NewUnaryOp(pos, op, right), nil
	case token.STAR:
		s.consume()

		pos := s.cur.Pos()
		op := s.cur.Token()

		right, err := s.parseExpr(token.UnaryPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse expr err: %w", fn, err), s.peek())
		}

		return ast.NewUnaryOp(pos, op, right), nil
	case token.AMP:
		s.consume()

		pos := s.cur.Pos()
		op := s.cur.Token()

		right, err := s.parseExpr(token.UnaryPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse expr err: %w", fn, err), s.peek())
		}

		return ast.NewUnaryOp(pos, op, right), nil
	case token.LBRACE:
		s.consume()

		pos := s.cur.Pos()

		list := []ast.Expr{}
		for s.peek().Token() != token.RBRACE {
			elem, err := s.parseExpr(token.LowestPrec)
			if err != nil {
				return nil, WrapErr(fmt.Errorf("[%s] parse array elem led expr err: %w", fn, err), s.peek())
			}
			list = append(list, elem)

			if s.peek().Token() == token.COM {
				s.consume()
			}
		}

		if s.peek().Token() != token.RBRACE {
			return nil, WrapErr(fmt.Errorf("[%s] expected } at pos %d", fn, s.peek().Pos()), s.peek())
		}
		s.consume()

		return ast.NewArrayLitExpr(pos, list), nil
	case token.LPAR:
		s.consume()

		expr, err := s.parseExpr(token.LowestPrec)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse expr err: %w", fn, err), s.peek())
		}

		if s.peek().Token() != token.RPAR {
			return nil, WrapErr(fmt.Errorf("[%s] expected ) at pos %d but given: (%s)", fn, s.peek().Pos(), s.peek().Token().String()), s.peek())
		}
		s.consume()

		return expr, nil
	case token.CHAR:
		s.consume()

		v := s.lex.GetValue(s.cur)
		str, err := strconv.Unquote(v)
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] unquote err: %w", fn, err), s.peek())
		}

		if len(str) > 1 && str[0] != '\\' {
			return nil, WrapErr(fmt.Errorf("[%s] expected CHAR at pos %d but str given: (%s)", fn, s.peek().Pos(), s.peek().Token().String()), s.peek())
		}

		return ast.NewNumberLitExpr(s.cur.Pos(), int64([]byte(str)[0])), nil
	default:
		s.lookahead.Print()
		fmt.Println(s.lex.GetValue(s.lookahead))
	}

	return nil, WrapErr(fmt.Errorf("[%s] unexpected token", fn), s.peek())
}

func (s *syntaxer) parseLed(left ast.Expr) (ast.Expr, error) {
	const fn = "parseLed"

	switch s.peek().Token() {
	case token.ASSIGN, token.LT, token.GT, token.EXCL:
		pos := s.cur.Pos()
		var op token.Token

		switch s.peek().Token() {
		case token.ASSIGN:
			s.consume()
			op = token.EQ

			if s.peek().Token() != token.ASSIGN {
				return nil, WrapErr(fmt.Errorf("[%s] expected = but given %s at pos %d", fn, s.peek().Token().String(), s.peek().Pos()), s.peek())
			}
			s.consume()
		case token.LT:
			s.consume()
			op = token.LT

			if s.peek().Token() == token.ASSIGN {
				s.consume()
				op = token.LTE
			}
		case token.GT:
			s.consume()
			op = token.GT

			if s.peek().Token() == token.ASSIGN {
				s.consume()
				op = token.GTE
			}
		case token.EXCL:
			s.consume()
			op = token.NE

			if s.peek().Token() != token.ASSIGN {
				return nil, WrapErr(fmt.Errorf("[%s] expected = at pos %d", fn, s.peek().Pos()), s.peek())
			}
			s.consume()
		default:
			panic("noway")
		}

		right, err := s.parseExpr(op.Priority())
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse bin led expr err: %w", fn, err), s.peek())
		}

		return ast.NewBinaryOp(pos, op, left, right), nil

	case token.PLUS:
		s.consume()

		pos := s.cur.Pos()

		if s.peek().Token() == token.PLUS {
			s.consume()

			return ast.NewPostIncOp(pos, left), nil
		}

		op := s.cur.Token()

		right, err := s.parseExpr(op.Priority())
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse bin led expr err: %w", fn, err), s.peek())
		}

		return ast.NewBinaryOp(pos, op, left, right), nil

	case token.PERCENT:
		s.consume()

		pos := s.cur.Pos()

		op := s.cur.Token()

		right, err := s.parseExpr(op.Priority())
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse bin led expr err: %w", fn, err), s.peek())
		}

		return ast.NewBinaryOp(pos, op, left, right), nil

	case token.MINUS:
		s.consume()

		pos := s.cur.Pos()

		if s.peek().Token() == token.MINUS {
			s.consume()

			return ast.NewPostDecOp(pos, left), nil
		}

		op := s.cur.Token()

		right, err := s.parseExpr(op.Priority())
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] prse bin led expr err: %w", fn, err), s.peek())
		}

		return ast.NewBinaryOp(pos, op, left, right), nil
	case token.STAR, token.SLASH, token.XOR:
		s.consume()

		pos := s.cur.Pos()
		op := s.cur.Token()

		right, err := s.parseExpr(op.Priority())
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse bin led expr err: %w", fn, err), s.peek())
		}

		return ast.NewBinaryOp(pos, op, left, right), nil
	case token.AMP:
		s.consume()

		pos := s.cur.Pos()
		op := s.cur.Token()

		if s.peek().Token() == token.AMP {
			s.consume()
			op = token.AND
		}

		right, err := s.parseExpr(op.Priority())
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse bin led expr err: %w", fn, err), s.peek())
		}

		return ast.NewBinaryOp(pos, op, left, right), nil
	case token.PIPE:
		s.consume()

		pos := s.cur.Pos()
		op := s.cur.Token()

		if s.peek().Token() == token.PIPE {
			s.consume()
			op = token.OR
		}

		right, err := s.parseExpr(op.Priority())
		if err != nil {
			return nil, WrapErr(fmt.Errorf("[%s] parse bin led expr err: %w", fn, err), s.peek())
		}

		return ast.NewBinaryOp(pos, op, left, right), nil
	case token.LPAR:
		s.consume()

		var args []ast.Expr
		for s.peek().Token() != token.RPAR {
			arg, err := s.parseExpr(token.LowestPrec)
			if err != nil {
				return nil, WrapErr(fmt.Errorf("[%s] parse arg led expr err: %w", fn, err), s.peek())
			}
			args = append(args, arg)

			if s.peek().Token() == token.COM {
				s.consume()
			}
		}

		if s.peek().Token() != token.RPAR {
			return nil, WrapErr(fmt.Errorf("[%s] expected ) at pos %d", fn, s.peek().Pos()), s.peek())
		}
		s.consume()

		return ast.NewCallExpr(left.Pos(), left, args), nil
	case token.DOT:
		s.consume()

		if s.peek().Token() != token.ID {
			return nil, WrapErr(fmt.Errorf("[%s] expected ID at pos %d", fn, s.peek().Pos()), s.peek())
		}
		s.consume()

		cur := ast.NewIdentExpr(s.cur.Pos(), s.lex.GetValue(s.cur))

		return ast.NewMemberExpr(left.Pos(), left, cur), nil
	default:
		s.print(s.lookahead)
	}

	return nil, WrapErr(fmt.Errorf("[%s] unexpected token", fn), s.peek())
}
