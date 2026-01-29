package lexer

import (
	"fmt"
	"loh/token"
)

type Lexeme struct {
	token token.Token
	pos   int
	len   int
}

func (l *Lexeme) Pos() int {
	return l.pos
}

func (l *Lexeme) Token() token.Token {
	return l.token
}

func (l *Lexeme) String(buf []byte) string {
	// TODO: use unsafeptr instead of alloc new str
	return string(buf[l.pos : l.pos+l.len])
}

func (l *Lexeme) Print() {
	fmt.Println(l.token.String(), l.pos, l.pos+l.len)
}

type Lexer struct {
	cur int
	pos int
	buf []byte
	ch  byte
}

func NewLexer(buf []byte) *Lexer {
	if len(buf) == 0 {
		return nil
	}

	return &Lexer{
		buf: buf,
	}
}

func (l *Lexer) Next() *Lexeme {
	l.consume()
	l.skipWhitespace()

	switch l.ch {
	case '\n', '\r':
		return &Lexeme{token: token.EOL, pos: l.cur, len: 1}
	case '^':
		return &Lexeme{token: token.XOR, pos: l.cur, len: 1}
	case '&':
		return &Lexeme{token: token.AMP, pos: l.cur, len: 1}
	case '|':
		return &Lexeme{token: token.PIPE, pos: l.cur, len: 1}
	case '+':
		return &Lexeme{token: token.PLUS, pos: l.cur, len: 1}
	case '-':
		return &Lexeme{token: token.MINUS, pos: l.cur, len: 1}
	case '*':
		return &Lexeme{token: token.STAR, pos: l.cur, len: 1}
	case '/':
		if l.peek() == '*' {
			return l.readMLComment()
		}

		if l.peek() == '/' {
			return l.readComment()
		}

		return &Lexeme{token: token.SLASH, pos: l.cur, len: 1}
	case '%':
		return &Lexeme{token: token.PERCENT, pos: l.cur, len: 1}
	case '"':
		return l.readString()
	case '\'':
		return l.readByte()

	case '.':
		return &Lexeme{token: token.DOT, pos: l.cur, len: 1}
	case ',':
		return &Lexeme{token: token.COM, pos: l.cur, len: 1}
	case ':':
		return &Lexeme{token: token.COL, pos: l.cur, len: 1}
	case ';':
		return &Lexeme{token: token.SCOL, pos: l.cur, len: 1}

	case '(':
		return &Lexeme{token: token.LPAR, pos: l.cur, len: 1}
	case ')':
		return &Lexeme{token: token.RPAR, pos: l.cur, len: 1}

	case '[':
		return &Lexeme{token: token.LSQBR, pos: l.cur, len: 1}
	case ']':
		return &Lexeme{token: token.RSQBR, pos: l.cur, len: 1}

	case '{':
		return &Lexeme{token: token.LBRACE, pos: l.cur, len: 1}
	case '}':
		return &Lexeme{token: token.RBRACE, pos: l.cur, len: 1}

	case '>':
		return &Lexeme{token: token.GT, pos: l.cur, len: 1}
	case '<':
		return &Lexeme{token: token.LT, pos: l.cur, len: 1}
	case '=':
		return &Lexeme{token: token.ASSIGN, pos: l.cur, len: 1}
	case '!':
		return &Lexeme{token: token.EXCL, pos: l.cur}
	case 0:
		return &Lexeme{token: token.EOF, pos: l.cur}
	}

	switch {
	case isDigit(l.ch):
		return l.readNum()
	case isIDStart(l.ch):
		return l.readIdentify()
	}

	return &Lexeme{token: token.UND, pos: l.pos}
}

func (l *Lexer) GetValue(t *Lexeme) string {
	return t.String(l.buf)
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.buf) {
		return 0
	}

	return l.buf[l.pos]
}

func (l *Lexer) consume() {
	if l.pos >= len(l.buf) {
		l.ch = 0

		return
	}

	l.ch = l.buf[l.pos]
	l.cur = l.pos
	l.pos++
}

func (l *Lexer) skipWhitespace() {
	for ; l.ch == ' ' || l.ch == '\t'; l.consume() {
	}
}

func (l *Lexer) readNum() *Lexeme {
	ptr := l.cur

	for ; isDigit(l.peek()); l.consume() {
	}

	return &Lexeme{
		token: token.NUM,
		pos:   ptr,
		len:   l.cur - ptr + 1,
	}
}

func (l *Lexer) readByte() *Lexeme {
	ptr := l.cur

	l.consume()
	for l.ch != '\'' {
		l.consume()
	}

	if l.ch == 0 {
		return &Lexeme{token: token.UND}
	}

	return &Lexeme{
		token: token.CHAR,
		pos:   ptr,
		len:   l.cur - ptr + 1,
	}
}

func (l *Lexer) readComment() *Lexeme {
	ptr := l.cur

	l.consume()
	l.consume()

	for l.ch != 0 {
		l.consume()

		if l.ch == 0 {
			return &Lexeme{token: token.UND}
		}

		if l.ch == '\n' {
			break
		}
	}

	return &Lexeme{
		token: token.COMMENT,
		pos:   ptr,
		len:   l.cur - ptr + 1,
	}
}

func (l *Lexer) readMLComment() *Lexeme {
	ptr := l.cur

	l.consume()
	l.consume()

	for l.ch != 0 {
		cur := l.ch
		l.consume()

		if l.ch == 0 {
			return &Lexeme{token: token.UND}
		}

		if cur == '*' && l.ch == '/' {
			break
		}
	}

	return &Lexeme{
		token: token.MLCOMMENT,
		pos:   ptr,
		len:   l.cur - ptr + 1,
	}
}

func (l *Lexer) readString() *Lexeme {
	ptr := l.cur

	l.consume()
	for l.ch != '"' {
		l.consume()
	}

	if l.ch == 0 {
		return &Lexeme{token: token.UND}
	}

	return &Lexeme{
		token: token.STR,
		pos:   ptr,
		len:   l.cur - ptr + 1,
	}
}

func (l *Lexer) readIdentify() *Lexeme {
	ptr := l.cur

	for ; isIDSymbol(l.peek()); l.consume() {
	}

	lex := Lexeme{
		token: token.ID,
		pos:   ptr,
		len:   l.cur - ptr + 1,
	}

	if t, ok := token.IsKeyword(lex.String(l.buf)); ok {
		lex.token = t
	}

	return &lex
}

func isAlpha(ch byte) bool    { return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') }
func isDigit(ch byte) bool    { return '0' <= ch && ch <= '9' }
func isIDStart(ch byte) bool  { return isAlpha(ch) || ch == '_' }
func isIDSymbol(ch byte) bool { return isAlpha(ch) || isDigit(ch) || ch == '_' }
