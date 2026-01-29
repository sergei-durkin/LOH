package token

import (
	"fmt"
)

type Token int

const (
	UND Token = iota

	ID  // hellO_
	NUM // 1234

	LT  // <
	LTE // <=
	GT  // >
	GTE // >=
	EQ  // ==
	NE  // !=

	XOR // ^

	PLUS  // +
	MINUS // -

	AMP // &
	AND // &&

	PIPE // |
	OR   // ||

	STAR    // *
	SLASH   // /
	PERCENT // %

	ASSIGN // =
	EXCL   // !

	LPAR   // (
	RPAR   // )
	LSQBR  // [
	RSQBR  // ]
	LBRACE // {
	RBRACE // }

	DOT  // .
	COM  // ,
	COL  // :
	SCOL // ;

	STR  // "any"
	CHAR // '0'

	MLCOMMENT // /* ... */
	COMMENT   // //

	// keywords
	TRUE
	FALSE

	CONST
	FOR
	IF
	ELSE
	RETURN
	FN
	TYPE
	STRUCT
	VAR
	BREAK
	CONTINUE

	EOL
	EOF
)

var tokens = map[Token]string{
	ID:  "id",
	NUM: "num",

	ASSIGN: "=",
	EXCL:   "!",
	LT:     "<",
	GT:     ">",
	LTE:    "<=",
	GTE:    ">=",
	EQ:     "==",
	NE:     "!=",

	XOR: "^",

	PLUS:  "+",
	MINUS: "-",

	AMP: "&",
	AND: "&&",

	PIPE: "|",
	OR:   "||",

	STAR:    "*",
	SLASH:   "/",
	PERCENT: "%",

	LPAR:   "(",
	RPAR:   ")",
	LSQBR:  "[",
	RSQBR:  "]",
	LBRACE: "{",
	RBRACE: "}",

	DOT:  ".",
	COM:  ",",
	COL:  ":",
	SCOL: ";",

	MLCOMMENT: "/* ... */",
	COMMENT:   "//",

	STR:  "string",
	CHAR: "'0'",

	// keywords
	TRUE:  "true",
	FALSE: "false",

	CONST:    "const",
	FOR:      "for",
	IF:       "if",
	ELSE:     "else",
	RETURN:   "return",
	FN:       "fn",
	TYPE:     "type",
	STRUCT:   "struct",
	VAR:      "var",
	BREAK:    "break",
	CONTINUE: "continue",

	EOL: "eol",
	EOF: "eof",
}

var keywords = map[string]Token{
	"true":     TRUE,
	"false":    FALSE,
	"const":    CONST,
	"for":      FOR,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"fn":       FN,
	"type":     TYPE,
	"struct":   STRUCT,
	"var":      VAR,
	"break":    BREAK,
	"continue": CONTINUE,
}

const (
	LowestPrec  = 0
	UnaryPrec   = 90
	HighestPrec = 100
)

func (t Token) Priority() int {
	switch t {
	case IF:
		return 1
	case ASSIGN, EXCL:
		return 10
	case OR:
		return 20
	case AND:
		return 30
	case EQ, NE, LT, LTE, GT, GTE:
		return 50
	case PLUS, MINUS, PIPE, XOR:
		return 60
	case STAR, SLASH, AMP, PERCENT:
		return 70
	}

	return LowestPrec
}

func (t Token) String() string {
	if s, ok := tokens[t]; ok {
		return s
	}

	return fmt.Sprintf("tok(%d)", t)
}

func IsKeyword(s string) (Token, bool) {
	t, ok := keywords[s]
	return t, ok
}

func (t Token) IsTerm() bool {
	return t == SCOL || t == EOL
}
