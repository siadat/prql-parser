package parser

import (
	"fmt"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/siadat/prql-parser/ast"
	"github.com/siadat/prql-parser/scanner"
	"github.com/siadat/prql-parser/token"
)

type Parser struct {
	src     io.Reader
	root    *ast.Root
	scanner *scanner.Scanner
	debug   bool
}

type ParseError struct {
	err error
}

func (e ParseError) Error() string {
	return e.err.Error()
}

func (p *Parser) proceed() scanner.Token {
	var t, err = p.scanner.NextToken()
	p.checkErr(err)
	return t
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) init(src io.Reader) error {
	p.scanner = scanner.NewScanner(src)
	p.scanner.SetSkipWhitespace(true)
	p.scanner.SetSkipComment(true)
	p.scanner.SetDebug(p.debug)

	var _, err = p.scanner.NextToken()
	return err
}

func (p *Parser) SetDebug(debug bool) {
	p.debug = debug
	if p.scanner != nil {
		p.scanner.SetDebug(debug)
	}
}

func (p *Parser) Parse(src io.Reader) (retRoot *ast.Root, retErr error) {
	if err := p.init(src); err != nil {
		return nil, err
	}

	defer func() {
		var err = recover()
		if err, ok := err.(ParseError); ok {
			retErr = err
		}
	}()

	retRoot = &ast.Root{Transforms: p.parseTransforms()}
	return
}

func (p *Parser) ParseExpr(src io.Reader) (retExpr ast.Expr, retErr error) {
	if err := p.init(src); err != nil {
		return nil, err
	}

	defer func() {
		var err = recover()
		if err, ok := err.(ParseError); ok {
			retErr = err
		}
	}()
	retExpr = p.parseExpr(nil, token.LowestPrecedence)
	return
}

func (p *Parser) parseTransforms() []ast.Node {
	var nodes []ast.Node
	for {
		var node = p.parseTransform(0)
		if node != nil {
			nodes = append(nodes, node)
		}
		if p.scanner.Eof() {
			return nodes
		}
	}
}

func (p *Parser) parseTransform(indent int) ast.Node {
	var t = p.scanner.CurrToken()
	if t.Lit == "from" {
		return p.parseFromTransform()
	} else if t.Lit == "select" {
		return p.parseSelectTransform()
	} else if t.Lit == "derive" {
		return p.parseDeriveTransform()
	} else if t.Typ == token.NEWLINE {
		p.proceed()
		return nil
	} else if t.Typ == token.EOF {
		return nil
	} else {
		panic(ParseError{fmt.Errorf("failed to parse a transform, unexpected %s", t)})
	}
}

func (p *Parser) parseFromTransform() ast.Node {
	// [x] from table1
	// [x] from e1 = table1
	var ident1 scanner.Token
	var ident2 scanner.Token

	p.expect(token.IDENTIFIER, "from")
	p.proceed()

	ident1 = p.expectType(token.IDENTIFIER)
	p.proceed()

	switch t := p.scanner.CurrToken(); t.Typ {
	case token.ASSIGN:
		p.proceed()
		ident2 = p.expectType(token.IDENTIFIER)
		p.proceed()
	case token.NEWLINE:
		p.proceed()
	case token.EOF:
	}

	if ident1 != (scanner.Token{}) && ident2 != (scanner.Token{}) {
		return ast.FromTransform{
			Alias: &ast.Ident{Name: ident1.Lit, Pos: ident1.Pos},
			Table: ast.Ident{Name: ident2.Lit, Pos: ident2.Pos},
		}
	}
	if ident1 != (scanner.Token{}) {
		return ast.FromTransform{
			Alias: nil,
			Table: ast.Ident{Name: ident1.Lit, Pos: ident1.Pos},
		}
	}
	panic(ParseError{fmt.Errorf("parse failed")})
}

func (p *Parser) parsePrimaryExpr() ast.Expr {
	switch t := p.scanner.CurrToken(); t.Typ {
	case token.STRING:
		p.proceed()

		return ast.String{Value: t.Lit}
	case token.IDENTIFIER:
		p.proceed()

		return ast.Column{
			Name: ast.Ident{Name: t.Lit, Pos: t.Pos},
		}
	case token.DATE:
		p.proceed()

		var t, err = time.Parse("@2006-01-02", t.Lit)
		p.checkErr(err)

		return ast.Date{
			Year:  t.Year(),
			Month: int(t.Month()),
			Day:   t.Day(),
		}
	case token.TIME:
		p.proceed()

		var t, err = time.Parse("@15:04:05", t.Lit)
		p.checkErr(err)

		return ast.Time{
			Hour:   t.Hour(),
			Minute: t.Minute(),
			Second: t.Second(),
		}
	case token.TIMESTAMP:
		p.proceed()
		var t, err = time.Parse("@2006-01-02T15:04:05", t.Lit)
		p.checkErr(err)

		return ast.Timestamp{
			Year:   t.Year(),
			Month:  int(t.Month()),
			Day:    t.Day(),
			Hour:   t.Hour(),
			Minute: t.Minute(),
			Second: t.Second(),
		}
	case token.INTERVAL:
		p.proceed()
		for _, unit := range token.Units {
			var idx = strings.Index(t.Lit, unit)
			if idx != -1 {
				var d, err = strconv.ParseInt(t.Lit[:idx], 10, 64)
				p.checkErr(err)
				return ast.Interval{Count: int(d), Unit: unit}
			}
		}
		panic(ParseError{fmt.Errorf("bad interval format %s", t)})
	case token.ADD, token.SUB:
		var op = t.Typ
		// signed expression, e.g. -1 or +value
		p.proceed()

		switch t := p.scanner.CurrToken(); t.Typ {
		case token.INTEGER,
			token.FLOAT,
			token.IDENTIFIER,
			token.LPAREN:
			return ast.UnaryExpr{
				X:  p.parsePrimaryExpr(),
				Op: op,
			}
		default:
			panic(ParseError{fmt.Errorf("expected integer or float, got %s", t)})
		}
	case token.INTEGER:
		p.proceed()

		var d, err = strconv.ParseInt(t.Lit, 10, 64)
		p.checkErr(err)
		return ast.Integer{Value: int(d)}
	case token.FLOAT:
		p.proceed()

		var f, err = strconv.ParseFloat(t.Lit, 64)
		p.checkErr(err)
		return ast.Float{Value: f}
	case token.LPAREN:
		return p.parseParenExpr()
	default:
		panic(ParseError{fmt.Errorf("failed to parse primary expression, got %s", t)})
	}
}

func (p *Parser) parseParenExpr() ast.Expr {
	p.expect(token.LPAREN, "(")
	p.proceed()
	var expr = p.parseExpr(nil, token.LowestPrecedence)

	p.expect(token.RPAREN, ")")
	p.proceed()

	return ast.ParenExpr{X: expr}
}

func (p *Parser) checkErr(err error) {
	if err != nil {
		if p.debug {
			debug.PrintStack()
		}
		panic(ParseError{err})
	}
}

func (p *Parser) parseExpr(lhs ast.Expr, minPrec token.Precedence) ast.Expr {
	if lhs == nil {
		lhs = p.parsePrimaryExpr()
	}

	for {
		var tk = p.scanner.CurrToken()
		var prec, isOp = token.Precedences[tk.Typ]
		if !isOp {
			return lhs
		}
		if tk.Typ == token.EOF {
			return lhs
		}
		if prec < minPrec {
			return lhs
		}

		p.proceed()

		var rhs = p.parseExpr(nil, prec)
		lhs = ast.BinaryExpr{
			X:  lhs,
			Y:  rhs,
			Op: tk.Typ,
		}
	}
}

// parseAssignExpr returns an expr that might be an AssignExpr
func (p *Parser) parseAssignExpr() ast.Expr {
	switch firstToken := p.scanner.CurrToken(); firstToken.Typ {
	case token.IDENTIFIER:
		var firstIdent = ast.Ident{Name: firstToken.Lit, Pos: firstToken.Pos}
		p.proceed()
		if t := p.scanner.CurrToken(); t.Typ == token.ASSIGN && t.Lit == "=" {
			p.proceed()
			return ast.AssignExpr{
				Name: firstIdent.Name,
				Expr: p.parseExpr(nil, token.LowestPrecedence),
			}
		} else {
			return p.parseExpr(ast.Column{Name: firstIdent}, token.LowestPrecedence)
		}

	default:
		return p.parseExpr(nil, token.LowestPrecedence)
	}
}

func (p *Parser) skipOptionalNewlines() error {
	for {
		var t = p.scanner.CurrToken()
		if t.Typ == token.NEWLINE {
			p.proceed()
		} else {
			return nil
		}
	}
}

func (p *Parser) expect(typ token.Token, lit string) {
	var t = p.scanner.CurrToken()
	if t.Typ == typ && t.Lit == lit {
		return
	}
	panic(ParseError{fmt.Errorf("expected %q, got %s", lit, t)})
}

func (p *Parser) expectType(typ token.Token) scanner.Token {
	var t = p.scanner.CurrToken()
	if t.Typ == typ {
		return t
	}
	panic(ParseError{fmt.Errorf("expected %s, got %s", typ, t)})
}

func (p *Parser) parseExprList() ast.ExprList {
	var list ast.ExprList

	switch t1 := p.scanner.CurrToken(); t1.Typ {
	case token.LBRACK:
		p.proceed()
		for {
			p.checkErr(p.skipOptionalNewlines())
			switch tk := p.scanner.CurrToken(); tk.Typ {
			case token.RBRACK:
				p.proceed()
				return list
			case token.EOF:
				return list
			default:
				var assign = p.parseAssignExpr()
				list.Items = append(list.Items, assign)

				switch tk := p.scanner.CurrToken(); tk.Typ {
				case token.COMMA:
					p.proceed()
				case token.NEWLINE:
					p.checkErr(p.skipOptionalNewlines())
				case token.RBRACK:
					p.proceed()
				default:
					panic(ParseError{fmt.Errorf("unexpected token %s", tk)})
				}
			}
		}
	default:
		var assign = p.parseAssignExpr()
		list.Items = append(list.Items, assign)
		return list
	}
}

func (p *Parser) parseDeriveTransform() ast.Node {
	var list = ast.ExprList{
		Items: nil,
	}

	p.expect(token.IDENTIFIER, "derive")
	p.proceed()

	var items = p.parseExprList()
	list.Items = items.Items

	return ast.DeriveTransform{List: list}
}

func (p *Parser) parseSelectTransform() ast.Node {
	var list = ast.ExprList{
		Items: nil,
	}

	p.expect(token.IDENTIFIER, "select")
	p.proceed()

	var items = p.parseExprList()
	list.Items = items.Items

	return ast.SelectTransform{List: list}
}
