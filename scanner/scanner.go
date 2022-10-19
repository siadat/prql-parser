package scanner

import (
	"fmt"
	"io"
	"strings"

	"github.com/siadat/prql-parser/token"
)

type Scanner struct {
	src []rune

	currToken    Token
	currRune     rune
	position     int
	readPosition int

	nextRune       rune
	skipWhitespace bool
}

type Pos int

type Token struct {
	Typ token.Token
	Lit string
	Pos Pos
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %d", t.Typ, t.Lit, t.Pos)
}

func NewScanner(src io.Reader) *Scanner {
	var s, err = io.ReadAll(src)
	if err != nil {
		panic(err)
	}
	var scanner = &Scanner{
		src: []rune(string(s)),
	}
	scanner.readRune()
	return scanner
}

const EndOfInput = 0

func (s *Scanner) SetSkipWhitespace(v bool) {
	s.skipWhitespace = true
}

func (s *Scanner) readRune() {
	if s.readPosition >= len(s.src) {
		s.currRune = EndOfInput
	} else {
		s.currRune = s.src[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition += 1

	// nextRune
	if s.readPosition >= len(s.src) {
		s.nextRune = EndOfInput
	} else {
		s.nextRune = s.src[s.readPosition]
	}
}

func (s *Scanner) NextToken() (Token, error) {
	var t, err = s.nextToken()
	s.currToken = t
	if err != nil {
		s.readRune() // skip
	}
	return t, err
}

func (s *Scanner) isEndOfExpression() bool {
	var ch = s.currRune
	if ch == ' ' || ch == ',' || ch == ']' || ch == '\n' || ch == EndOfInput {
		return true
	}
	// ".."
	if ch == '.' && s.nextRune == '.' {
		return true
	}
	return false
}

func (s *Scanner) Eof() bool {
	return s.currToken.Typ == token.EOF
}

func (s *Scanner) CurrPosition() int {
	return s.position
}

func (s *Scanner) CurrToken() Token {
	return s.currToken
}

func (s *Scanner) nextToken() (Token, error) {
	// s.PrintCursor("debug")
	var start = s.position

	switch s.currRune {
	case '\n', '\r':
		var tok = Token{token.NEWLINE, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case '|':
		var tok = Token{token.PIPE, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case '[':
		var tok = Token{token.LBRACK, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case ']':
		var tok = Token{token.RBRACK, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case '(':
		var tok = Token{token.LPAREN, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case ')':
		var tok = Token{token.RPAREN, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case ',':
		var tok = Token{token.COMMA, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case '%':
		var tok = Token{token.REM, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case '*':
		var tok = Token{token.MUL, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case '/':
		var tok = Token{token.QUO, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case '.':
		var tok = Token{token.PERIOD, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case ':':
		var tok = Token{token.COLON, fmt.Sprintf("%c", s.currRune), Pos(start)}
		s.readRune()
		return tok, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.readNumber()
	case '-':
		// this can be 'a - b' or '-12' or '->'
		if isNumerical(s.nextRune) {
			return s.readNumber()
		}

		if s.nextRune == '>' {
			var tok = Token{token.ARROW, fmt.Sprintf("%c%c", s.currRune, s.nextRune), Pos(start)}
			s.readRune()
			s.readRune() // 2nd time for '=='
			return tok, nil
		} else {
			var tok = Token{token.SUB, fmt.Sprintf("%c", s.currRune), Pos(start)}
			s.readRune()
			return tok, nil
		}
	case '+':
		// this can be 'a + b' or '+12' ('+' followed immediately by a number)
		if isNumerical(s.nextRune) {
			return s.readNumber()
		} else {
			var tok = Token{token.ADD, fmt.Sprintf("%c", s.currRune), Pos(start)}
			s.readRune()
			return tok, nil
		}
	case '=':
		// this can be '=' or '=='
		if s.nextRune == '=' {
			var tok = Token{token.EQL, fmt.Sprintf("%c%c", s.currRune, s.nextRune), Pos(start)}
			s.readRune()
			s.readRune() // 2nd time for '=='
			return tok, nil
		} else {
			var tok = Token{token.ASSIGN, fmt.Sprintf("%c", s.currRune), Pos(start)}
			s.readRune()
			return tok, nil
		}
	case '>':
		// this can be '>' or '>='
		if s.nextRune == '=' {
			var tok = Token{token.GEQ, fmt.Sprintf("%c%c", s.currRune, s.nextRune), Pos(start)}
			s.readRune()
			s.readRune() // 2nd time for '=='
			return tok, nil
		} else {
			var tok = Token{token.GTR, fmt.Sprintf("%c", s.currRune), Pos(start)}
			s.readRune()
			return tok, nil
		}
	case '<':
		// this can be '<' or '<='
		if s.nextRune == '=' {
			var tok = Token{token.LEQ, fmt.Sprintf("%c%c", s.currRune, s.nextRune), Pos(start)}
			s.readRune()
			s.readRune() // 2nd time for '=='
			return tok, nil
		} else {
			var tok = Token{token.LSS, fmt.Sprintf("%c", s.currRune), Pos(start)}
			s.readRune()
			return tok, nil
		}
	case '!':
		// this can be '<' or '<='
		if s.nextRune == '=' {
			var tok = Token{token.NEQ, fmt.Sprintf("%c%c", s.currRune, s.nextRune), Pos(start)}
			s.readRune()
			s.readRune() // 2nd time for '=='
			return tok, nil
		} else {
			var tok = Token{token.NOT, fmt.Sprintf("%c", s.currRune), Pos(start)}
			s.readRune()
			return tok, nil
		}
	case '?':
		// this can be '??'
		if s.nextRune == '?' {
			var tok = Token{token.COALESCE, fmt.Sprintf("%c%c", s.currRune, s.nextRune), Pos(start)}
			s.readRune()
			s.readRune() // 2nd time for '=='
			return tok, nil
		} else {
			return Token{
				token.ILLEGAL,
				fmt.Sprintf("%c", s.currRune),
				Pos(start),
			}, fmt.Errorf("unexpected character %c", s.currRune)
		}
	case 'f', 's':
		// this can be '<' or '<='
		if s.nextRune == '"' {
			return s.readString('"')
		} else {
			return s.readIdentifier()
		}
	case '@':
		return s.readDateOrTimeOrDatetime()
	case '#':
		return s.readComment()
	case '"':
		return s.readString('"')
	case '\'':
		return s.readString('\'')
	case ' ', '\t':
		if s.skipWhitespace {
			var ret, err = s.readWhitespace()
			if err != nil {
				return ret, err
			}
			return s.nextToken()
		} else {
			return s.readWhitespace()
		}
	case '`':
		return s.readIdentifierQuoted()
	case EndOfInput:
		var tok = Token{token.EOF, "", Pos(start)}
		s.readRune()
		return tok, nil
	default:
		return s.readIdentifier()
	}
}

func (s *Scanner) readIdentifier() (Token, error) {
	if s.isIdentifierPartFirst() {
		return s.readIdentifierUnquoted()
	} else {
		return Token{
			token.ILLEGAL,
			fmt.Sprintf("%c", s.currRune),
			Pos(s.position),
		}, fmt.Errorf("unexpected identifier character %c", s.currRune)
	}
}

func (s *Scanner) readString(ender rune) (Token, error) {
	var position = s.position

	if s.currRune == 'f' || s.currRune == 's' {
		s.readRune() // skip f and f in f"..." and s"..."
	}
	s.readRune() // skip opener, i.e. ' or "

	for {
		switch s.currRune {
		case '\\':
			s.readRune() // skip \
			s.readRune() // skip the char after \
		case ender:
			s.readRune()
			return Token{
				token.STRING,
				string(s.src[position:s.position]),
				Pos(position),
			}, nil
		default:
			s.readRune()
		}
	}
}

func (s *Scanner) readComment() (Token, error) {
	var position = s.position
	for s.currRune != '\n' {
		s.readRune()
	}
	return Token{
		token.COMMENT,
		string(s.src[position:s.position]),
		Pos(position),
	}, nil
}

func (s *Scanner) readWhitespace() (Token, error) {
	var position = s.position
	for s.isWhitespace() {
		s.readRune()
	}
	return Token{
		token.WHITESPACE,
		string(s.src[position:s.position]),
		Pos(position),
	}, nil
}

func (s *Scanner) isWhitespace() bool {
	var ch = s.currRune
	return ' ' == ch || '\t' == ch
}

func (s *Scanner) readIdentifierUnquoted() (Token, error) {
	var position = s.position
	s.readRune()
	for s.isIdentifierMiddle() {
		s.readRune()
	}
	var lit = string(s.src[position:s.position])
	return Token{token.IDENTIFIER, lit, Pos(position)}, nil
}

func (s *Scanner) isIdentifierPartFirst() bool {
	var ch = s.currRune
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

func (s *Scanner) isIdentifierMiddle() bool {
	var ch = s.currRune
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || (ch >= '0' && ch <= '9')
}

func isNumerical(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) readNumber() (Token, error) {
	var position = s.position

	if s.currRune == '-' || s.currRune == '+' {
		s.readRune() // skip sign
	}

	var isFloat = false
	for {
		switch s.currRune {
		case '.':
			isFloat = true
			s.readRune()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.readRune()
		default:
			if s.nextRune >= 'a' && s.nextRune <= 'z' {
				var ident, err = s.readIdentifier()
				if err != nil {
					return ident, err
				}
				switch ident.Lit {
				case "microseconds", "milliseconds", "seconds", "minutes", "hours", "days", "weeks", "months", "years":
					return Token{
						token.INTERVAL,
						string(s.src[position:s.position]),
						Pos(position),
					}, nil
				default:
					return ident, fmt.Errorf("expected an interval, got %s", ident)
				}
			}

			if isFloat {
				return Token{
					token.FLOAT,
					string(s.src[position:s.position]),
					Pos(position),
				}, nil
			} else {
				return Token{
					token.INTEGER,
					string(s.src[position:s.position]),
					Pos(position),
				}, nil
			}
		}
	}
}

func (s *Scanner) readDateOrTimeOrDatetime() (Token, error) {
	var position = s.position
	for {
		s.readRune() // skip '@'
		if s.isEndOfExpression() {
			var lit = string(s.src[position:s.position])
			var typ token.Token

			if strings.Index(lit, "T") != -1 {
				typ = token.TIMESTAMP
			} else if strings.Index(lit, "-") != -1 {
				typ = token.DATE
			} else if strings.Index(lit, ":") != -1 {
				typ = token.TIME
			}

			return Token{
				typ,
				lit,
				Pos(position),
			}, nil
		}
	}
}

func (s *Scanner) readIdentifierQuoted() (Token, error) {
	var position = s.position
	s.readRune() // skip
	for {
		switch s.currRune {
		case '`':
			s.readRune()
			var end = s.position
			return Token{
				token.IDENTIFIER,
				string(s.src[position:end]),
				Pos(position),
			}, nil
		case EndOfInput:
			return Token{
				token.ILLEGAL,
				fmt.Sprintf("%c", s.currRune),
				Pos(position),
			}, fmt.Errorf("missing close `")
		default:
			s.readRune()
		}
	}

}

func (s *Scanner) PrintCursor(layout string, args ...interface{}) {
	var lines = strings.Split(string(s.src), "\n")
	var b strings.Builder
	var line, column = s.getCurrPosition()

	var ch string
	if s.currRune == EndOfInput {
		ch = "EndOfInput"
	} else {
		ch = fmt.Sprintf("%q", s.currRune)
	}

	var prefix = fmt.Sprintf(layout, args...)
	b.WriteString(fmt.Sprintf("%s  %s\n", prefix, lines[line]))
	b.WriteString(fmt.Sprintf("%s  %s▲ [%d]=%s token=%s\n", prefix, strings.Repeat(" ", column), s.position, ch, s.currToken))
	fmt.Print(b.String())
	// fmt.Println("[debug] ", prefix, string(s.src), len(s.src), fmt.Sprintf("[%4d]=%c", s.position, s.currRune))
}

func (s *Scanner) getCurrPosition() (int, int) {
	var line = 0
	var column = 0
	for i := 0; i < s.position && i < len(s.src); i++ {
		if s.src[i] == '\n' {
			line += 1
			column = 0
		} else {
			column += 1
		}
	}
	return line, column
}
