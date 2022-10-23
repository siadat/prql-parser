package token

import (
	"fmt"
	"strconv"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WHITESPACE
	NEWLINE
	PIPE
	COMMENT

	literal_beg
	IDENTIFIER // abc | _abc | $abd | `abc`
	INTEGER    // 12345
	FLOAT      // 123.45
	BOOLEAN    // true | false
	STRING     // "abc" | 'abc' | f"abc{efg}" | s"version()"
	DATE       // @2022-12-31
	TIME       // @00:00:00
	TIMESTAMP  // @2022-12-31T00:00:00
	INTERVAL   // 123microseconds | 123milliseconds | 123seconds | 123minutes | 123hours | 123days | 123weeks | 123months | 123years
	literal_end

	operator_beg
	LPAREN   // (
	RPAREN   // )
	LBRACK   // [
	RBRACK   // ]
	LBRACE   // {
	RBRACE   // }
	COLON    // :
	COMMA    // ,
	PERIOD   // .
	ADD      // +
	SUB      // -
	NOT      // !
	MUL      // *
	QUO      // /
	REM      // %
	EQL      // ==
	LSS      // <
	GTR      // >
	ASSIGN   // =
	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	AND      // and
	OR       // or
	COALESCE // ??
	ARROW    // ->
	operator_end

	keyword_beg
	FUNC  // func
	TABLE // table
	PRQL  // prql
	// TRUE  // true we already have BOOLEAN
	// FALSE // false we already have BOOLEAN
	NULL // null
	keyword_end
)

const (
	AnyTyp Token  = -1
	AnyLit string = ""
)

var tokens = [...]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	WHITESPACE: "WHITESPACE",
	NEWLINE:    "NEWLINE",
	PIPE:       "PIPE",
	COMMENT:    "COMMENT",

	IDENTIFIER: "IDENTIFIER",
	INTEGER:    "INTEGER",
	FLOAT:      "FLOAT",
	BOOLEAN:    "BOOLEAN",
	STRING:     "STRING",
	DATE:       "DATE",
	TIME:       "TIME",
	TIMESTAMP:  "TIMESTAMP",
	INTERVAL:   "INTERVAL",

	LPAREN:   "LPAREN",
	RPAREN:   "RPAREN",
	LBRACK:   "LBRACK",
	RBRACK:   "RBRACK",
	LBRACE:   "LBRACE",
	RBRACE:   "RBRACE",
	COLON:    "COLON",
	COMMA:    "COMMA",
	PERIOD:   "PERIOD",
	ADD:      "ADD",
	SUB:      "SUB",
	NOT:      "NOT",
	MUL:      "MUL",
	QUO:      "QUO",
	REM:      "REM",
	EQL:      "EQL",
	LSS:      "LSS",
	GTR:      "GTR",
	ASSIGN:   "ASSIGN",
	NEQ:      "NEQ",
	LEQ:      "LEQ",
	GEQ:      "GEQ",
	AND:      "AND",
	OR:       "OR",
	COALESCE: "COALESCE",
	ARROW:    "ARROW",

	FUNC:  "FUNC",
	TABLE: "TABLE",
	PRQL:  "PRQL",
	NULL:  "NULL",
}

var Units = [...]string{
	"microseconds",
	"milliseconds",
	"seconds",
	"minutes",
	"hours",
	"days",
	"weeks",
	"months",
	"years",
}

type Precedence int

var LowestPrecedence Precedence = 0
var Precedences = map[Token]Precedence{
	ADD: 1,
	SUB: 1,
	MUL: 2,
	QUO: 2,
}

func (tok Token) String() string {
	if tok == AnyTyp {
		return ":AnyTyp:"
	}
	var s = ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

type LiteralStringer string

func (lit LiteralStringer) String() string {
	if string(lit) == AnyLit {
		return ":AnyLit:"
	}
	return fmt.Sprintf("%q", string(lit))
}
